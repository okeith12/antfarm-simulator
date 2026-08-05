# AntFarm Architecture

---

## Package graph

```
                    ┌──────────────────────────────────┐
                    │              gui/                │  tcell
                    │  antfarm · renderer · stats      │  ONLY here
                    │  controls · colors               │
                    └───────────────┬──────────────────┘
                                    │ reads
                    ┌───────────────▼──────────────────┐
                    │          simulation/             │  package logic
                    │  updateWorld · antsBehavior      │
                    │  antPlacement · spawn            │
                    │  matureLarvaeToAnt               │
                    └───────┬──────────────────┬───────┘
                            │                  │
              ┌─────────────▼──────┐   ┌───────▼────────┐
              │    pathfinder/     │   │     util/      │
              │  pathfinder        │   │  Abs           │
              │  workerpathfinder  │   └────────────────┘
              │  nursepathfinder   │
              └─────────┬──────────┘
                        │
              ┌─────────▼──────────┐   ┌────────────────┐
              │       types/       │──▶│      rng/      │
              │  world · cell      │   │  xorshift32    │
              │  colony · ant      │   └────────────────┘
              │  queen · nurse     │
              │  worker · soldier  │
              │  larvae · log      │
              └────────────────────┘
```

| Zone | tcell |
|---|---|
| `gui/` | yes |
| `types/` `simulation/` `pathfinder/` `util/` `rng/` | no |

```bash
grep -rn "tcell" types simulation pathfinder util rng main.go   # empty
```

---

## Data model

```
World
├── Width, Height   int
├── Cells           []Cell        flat, row-major
├── Colonies        []*Colony
├── Ticks           int
└── Rng             *rng.Rng

Cell                              Colony
├── Soil     Soil                 ├── Queen         *QueenAnt    reigning, 1 or nil
├── IsTunnel bool                 ├── Queens        []*QueenAnt  heirs, frozen
├── Occupant AntInterface         ├── HeadNurse     *NurseAnt
└── Food     int                  ├── Nurses        []*NurseAnt
                                  ├── Workers       []*WorkerAnt
Ant  (embedded by all five)       ├── Soldiers      []*SoldierAnt
├── ID, Role, Position            ├── Larvae        []*LarvaeAnt
├── Health, MaxHealth             ├── Food, Eggs    int
├── Age, MaxAge                   ├── NextAntID     int
├── ColonyID                      ├── Color         ColonyColor
└── CurrentAction                 └── QueenPosition Position
```

```
AntInterface
├── GetAnt()  *Ant
├── GetIcon() rune
└── GetRole() Role
        ▲
        ├── QueenAnt    ♛   + EggLayingCooldown, TotalEggsLaid, Declining
        ├── NurseAnt    ○   + CurrentlyNursing, LarvaeNursed
        ├── WorkerAnt   ●   + CarryingFood, FoodAmount, DiggingPower, direction
        ├── SoldierAnt  ⚔
        └── LarvaeAnt   ◦   + HasNurseCare, GrowthProgress
```

---

## Grid memory layout

```
Cells []Cell            Index(x, y) = y*Width + x

W = 5
         x=0   1    2    3    4
  y=0  [  0    1    2    3    4  ]
  y=1  [  5    6    7    8    9  ]      Index(2,1) = 1*5+2 = 7
  y=2  [ 10   11   12   13   14  ]

  ├──── row 0 ────┼──── row 1 ────┼──── row 2 ────┤   contiguous
```

| | |
|---|---|
| 3D form | `(z*Height + y)*Width + x` |
| `GetCell(x,y)` | `&Cells[Index(x,y)]`, nil out of bounds |
| returns | pointer into the slice, not a copy |

---

## One tick

```
UpdateWorld(world)
│
├─ Ticks++
│
└─ for each colony ──▶ updateColony
   │
   ├─ 1. queen action := "resting"
   │
   ├─ 2. if Queen.Declining && Ticks % 30 == 0 ──▶ Queen.Health--
   │
   ├─ 3. processDeaths
   │      ├─ queen dead? ──▶ crown Queens[0] ──▶ demote the rest
   │      ├─ head nurse · nurses · workers · soldiers · larvae
   │      └─ removals are backward-iterated
   │
   ├─ 4. Ticks % 50 == 0 && Queen != nil && Food >= 100 ──▶ Eggs++, Food -= 1
   │
   ├─ 5. Ticks % 30 == 0 && Eggs > 0 ──▶ Eggs--, spawn 1 larva near queen
   │
   ├─ 6. all larvae ──▶ Age++
   │
   ├─ 7. larvae with HasNurseCare && Age >= 50 ──▶ matureLarvaeToAnt
   │
   ├─ 8. head nurse ─┐
   │    nurses ──────┤
   │    workers ─────┼──▶ per-role behaviour
   │    soldiers ────┘
   │
   └─ 9. larvae action := "getting care" | "waiting for care"
```

---

## Render loop

```
gui.Antfarm.Run()
│
├─ renderTicker      30 FPS, fixed
├─ simulationTicker  speedPresets[i] Hz
│
└─ loop
   ├─ handleEvents ──▶ Q/ESC quit · L log · P pause · +/- speed
   ├─ simulationTicker fires && !paused ──▶ logic.UpdateWorld
   └─ renderTicker fires && needsRender ──▶ renderer.Render
                                              ├─ terrain   SoilColor(Soil, IsTunnel)
                                              ├─ ants      ColonyColor(Colony.Color)
                                              ├─ stats bar
                                              └─ activity log
```

Simulation writes. Render only reads.

---

## Ant lifecycle

```
   Food >= 100          Ticks%30            nursed
   Ticks%50                │                Age>=50
      │                    │                   │
      ▼                    ▼                   ▼
    Egg ──────────────▶ Larva ─────────────▶ roll 0-99
   -1 food                 │                   │
                           │             ┌─────┼─────┬────────┬────────┐
                    Age>=200             │     │     │        │        │
                           │             0    1-20  21-35   36-99      │
                           ▼             │     │     │        │        │
                         dead          Queen Nurse Soldier Worker      │
                                         1%   20%    15%     64%       │
                                                                       │
   Age >= MaxAge  or  Health <= 0 ──────────────────────────────▶ dead
```

| Role | MaxAge | Health | Drain |
|---|---|---|---|
| Worker | 500 | 100 | -1 per dig |
| Soldier | 600 | 150 | none |
| Nurse | 700 | 80 | none |
| Larva | 200 | 50 | none |
| Queen | 20000, unreached | 200 | -1 per 30 ticks, only when Declining |

---

## Queen succession

```
        larva rolls 0
              │
    ┌─────────┴──────────┐
    │                    │
Queen == nil        Queen != nil
    │                    │
    ▼                    ▼
 CROWNED           append Queens
 QueenPosition     Queen.Declining = true
 = her position          │
    │                    │
    ▼                    ▼
┌────────────┐     ┌───────────┐
│  REIGNING  │     │   HEIR    │
│            │     │           │
│ no ageing  │     │ no ageing │
│ immortal   │     │ no death  │
│ until an   │     │ frozen    │
│ heir       │     └─────┬─────┘
└─────┬──────┘           │
      │ Declining        │
      ▼                  │
┌────────────┐           │
│  FADING    │           │
│ -1hp/30t   │           │
└─────┬──────┘           │
      │ Health <= 0      │
      ▼                  │
    dead ────────────────┤
                         │
              ┌──────────┴──────────┐
              │                     │
         Queens[0]            Queens[1:]
              │                     │
              ▼                     ▼
          CROWNED           demoteHeir
                            20% nurse / 80% worker
                            same ID, same position
```

**Invariant:** `Declining` is set only when an heir exists, and heirs cannot die.
A colony therefore never goes queenless.

---

## Determinism

```
rng.New(seed) ──▶ World.Rng ──┬──▶ world gen      food scatter, 10%
                              ├──▶ egg laying     (none, fixed at 1)
                              ├──▶ caste roll     Below(100)
                              ├──▶ heir demotion  Below(100)
                              └──▶ worker walk    Shuffle, Below(4)

math/rand ──▶ nowhere
```

| | |
|---|---|
| Algorithm | xorshift32 |
| `Below(n)` | `Next() % n`, fixed draw count |
| App seed | clock |
| Test seed | fixed |
| Guarantee | same seed, same colony |

---

## Food

```
stored as tenths          FoodScale = 10

  map pellet    50 units    5 food
  grass         50 units    5 food     one-shot, Food = -1 after
  worker carry  50 or 100
  colony start 500 units   50 food
  egg cost       1 unit   0.1 food
  lay gate     100 units    10 food

display: Food / FoodScale
```

Integer only. No floats anywhere in the simulation.

---

## Constants

`simulation/updateWorld.go`

| | |
|---|---|
| `eggLayingInterval` | 50 |
| `eggHatchTime` | 30 |
| `larvaeGrowTime` | 50 |
| `foodCost` | 1 |
| `layingThreshold` | 100 |
| `queenDeclineInterval` | 30 |

`simulation/matureLarvaeToAnt.go`

| Roll | Caste |
|---|---|
| 0 | Queen |
| 1-20 | Nurse |
| 21-35 | Soldier |
| 36-99 | Worker |

`gui/antfarm.go`

| | |
|---|---|
| `speedPresets` | 0.25, 0.5, 1, 2, 5, 10 |
| `renderFPS` | 30 |

---

## Gotchas

| | |
|---|---|
| `simulation/` declares `package logic` | import as `logic "antfarm/simulation"` |
| `types/solider.go` | filename is misspelled |
| `Cell.Food == -1` | harvested grass, not "no food" |
| `Colony.Queens` | heirs, not all queens |

---

## See also

| | |
|---|---|
| `README.md` | what it is, how to run |
| `WHATNEXT.md` | roadmap |
| `LIFECYCLE-AUDIT.md` | known behavioural problems |
