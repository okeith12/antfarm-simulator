# AntFarm Architecture

---

## Packages

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
├── Soil     Soil                 ├── Queen         *QueenAnt    reigning
├── IsTunnel bool                 ├── Queens        []*QueenAnt  heirs
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

## Grid

```
Cells []Cell            Index(x, y) = y*Width + x

W = 5
         x=0   1    2    3    4
  y=0  [  0    1    2    3    4  ]
  y=1  [  5    6    7    8    9  ]      Index(2,1) = 1*5+2 = 7
  y=2  [ 10   11   12   13   14  ]

  ├──── row 0 ────┼──── row 1 ────┼──── row 2 ────┤   contiguous

GetCell(x, y) ──▶ &Cells[Index(x, y)]
```

---

## Tick

```
UpdateWorld(world)
│
├─ Ticks++
│
└─ for each colony ──▶ updateColony
   │
   ├─ queen decline
   ├─ deaths ──▶ succession
   ├─ lay egg
   ├─ hatch egg ──▶ larva
   ├─ age larvae
   ├─ mature larvae ──▶ caste roll
   ├─ behaviour: head nurse · nurses · workers · soldiers
   └─ larvae care state
```

---

## Render

```
gui.Antfarm.Run()
│
├─ renderTicker      fixed
├─ simulationTicker  speed preset
│
└─ loop
   ├─ handleEvents ──▶ quit · log · pause · speed
   ├─ simulationTicker ──▶ logic.UpdateWorld        writes
   └─ renderTicker ──────▶ renderer.Render          reads
                             ├─ terrain   SoilColor
                             ├─ ants      ColonyColor
                             ├─ stats bar
                             └─ activity log
```

---

## Ant lifecycle

```
                                            ┌──▶ Queen
                                            │
    Egg ──────────▶ Larva ──────────▶ roll ─┼──▶ Nurse
                      │                     │
                      │                     ├──▶ Soldier
                      │                     │
                      ▼                     └──▶ Worker
                    dead                            │
                                                    ▼
                                                  dead
```

---

## Queen succession

```
        larva rolls queen
              │
    ┌─────────┴──────────┐
    │                    │
Queen == nil        Queen != nil
    │                    │
    ▼                    ▼
 CROWNED           append Queens
                   Declining = true
    │                    │
    ▼                    ▼
┌────────────┐     ┌───────────┐
│  REIGNING  │     │   HEIR    │
│            │     │           │
│ no ageing  │     │ no ageing │
│ immortal   │     │ no death  │
└─────┬──────┘     └─────┬─────┘
      │ Declining        │
      ▼                  │
┌────────────┐           │
│  FADING    │           │
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
          CROWNED             demoteHeir
                            nurse or worker
```

---

## Determinism

```
rng.New(seed) ──▶ World.Rng ──┬──▶ world gen
                              ├──▶ caste roll
                              ├──▶ heir demotion
                              └──▶ worker walk

math/rand ──▶ nowhere
```
