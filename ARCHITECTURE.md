# AntFarm Architecture

How the Go implementation is put together, and why it is shaped this way.

Status: v0, terminal. Last updated 2026-08-05.

Companion documents: `README.md` (what it is), `WHATNEXT.md` (where it is going),
`LIFECYCLE-AUDIT.md` (known behavioural problems).

---

## 1. The one rule

**The simulation core knows nothing about the display.**

```
gui/  ────uses────>  simulation/  ────uses────>  pathfinder/  ───>  types/
                                                                      ^
                                                                      |
                                              nothing here imports tcell
```

The arrow never reverses. `types/`, `simulation/`, `pathfinder/`, `util/` and
`rng/` contain no reference to tcell, to a terminal, or to any display concept.
You can verify it in one command:

```bash
grep -rn "tcell" types simulation pathfinder util rng main.go   # returns nothing
```

This matters because of where the project is going. The terminal build is v0.
The same simulation has to run on a microcontroller driving an SPI panel, where
no terminal library exists. Anything the core knows about rendering is something
that has to be torn out later.

`gui/` is the only package allowed to import tcell, because rendering is its
entire job.

---

## 2. Packages

| Package | Role | Depends on |
|---|---|---|
| `types/` | The nouns. World, Cell, Colony, and the five ant kinds. | `rng/` |
| `simulation/` | The verbs. One tick of world update, spawning, behaviour, placement. | `types/`, `pathfinder/`, `util/` |
| `pathfinder/` | Movement. Shared helpers plus per-role navigation. | `types/`, `util/` |
| `gui/` | Terminal rendering, input, the game loop. | everything, plus tcell |
| `rng/` | Deterministic xorshift32 generator. | nothing |
| `util/` | `Abs()` for integers. | nothing |

**A naming trap worth knowing.** The `simulation/` directory declares
`package logic`, left over from a rename. Import it as:

```go
import logic "antfarm/simulation"
```

---

## 3. Data model

### World

```go
type World struct {
    Width, Height int
    Cells         []Cell     // flat, row-major
    Colonies      []*Colony
    Ticks         int
    Rng           *rng.Rng
}
```

The grid is a **single flat slice**, not a slice of slices of pointers. The cell
at `(x, y)` lives at `Cells[y*Width + x]`:

```
W = 5
        x=0  1   2   3   4
 y=0  [  0   1   2   3   4 ]
 y=1  [  5   6   7   8   9 ]     Index(2,1) = 1*5 + 2 = 7
 y=2  [ 10  11  12  13  14 ]
```

Row-major, so stepping `x` moves one slot and stepping `y` moves a whole row.
Two reasons for that order: the renderer scans along rows, so consecutive reads
are contiguous in memory, and it is the 2D case of the firmware's 3D formula
`(z*Height + y)*Width + x`, which keeps the port mechanical.

`GetCell(x, y)` returns a pointer **into** the slice, so writes through it are
visible to the next reader. It returns nil out of bounds.

### Colony

One reigning queen, one head nurse, and slices for every other role. Food is
shared, held as scaled integers (see 5.1).

```go
type Colony struct {
    Queen     *QueenAnt      // the reigning queen, nil if the colony has none
    Queens    []*QueenAnt    // heirs in waiting
    HeadNurse *NurseAnt
    Nurses    []*NurseAnt
    Workers   []*WorkerAnt
    Soldiers  []*SoldierAnt
    Larvae    []*LarvaeAnt
    Food, Eggs, NextAntID int
    QueenPosition Position
}
```

Every ant embeds a base `Ant` and satisfies `AntInterface` (`GetAnt`, `GetIcon`,
`GetRole`). A `Cell` holds one `Occupant`.

---

## 4. One tick

`logic.UpdateWorld(world)` increments `Ticks` and runs `updateColony` for each
colony, in this order:

```
UpdateWorld
   |
   +-- Ticks++
   |
   +-- for each colony: updateColony
            |
            +-- queen action defaults to "resting"
            +-- if the queen is Declining, she loses health
            +-- processDeaths          (removals, then succession)
            +-- queen lays one egg     (needs a queen, the interval, and food)
            +-- one egg hatches into a larva
            +-- larvae age
            +-- nursed, grown larvae mature into adults
            +-- update head nurse, nurses, workers, soldiers
            +-- mark each larva as cared for or waiting
```

Deaths are processed first so the rest of the tick never operates on a corpse.

Rendering is driven separately by `gui/`, which reads world state and never
mutates it.

---

## 5. Design decisions worth knowing

### 5.1 Food is a scaled integer, never a float

`types.FoodScale = 10`, so internally food is counted in tenths. One food unit
displays as 0.1 food, and the stats bar divides by `FoodScale` on the way out.

Floats are avoided deliberately. The firmware port has to produce bit-identical
results from the same seed, and float arithmetic does not survive that: compilers
reassociate under optimisation, x86 may use 80-bit internal precision where the
target uses 32-bit, and fused multiply-add changes rounding. One differing bit
compounds into an entirely different colony a thousand ticks later. Integers are
exact on every platform.

### 5.2 Randomness is injected and reproducible

Nothing in the simulation calls `math/rand`. The generator is xorshift32, held on
the World and passed in at construction:

```go
world := types.NewWorld(120, 35, rng.New(seed))
```

The terminal app seeds from the clock, so every run differs. Tests pass fixed
seeds. Same seed, same colony, always.

xorshift32 was chosen over Go's `math/rand` because it is ten lines and
reimplements identically in C. `Rng.Below` uses plain modulo rather than
rejection sampling for the same reason: rejection sampling consumes a variable
number of draws and would desynchronise two implementations.

### 5.3 Colour lives only in the display layer

`gui/colors.go` holds the single mapping from simulation enums to tcell colours:
`SoilColor(soil, isTunnel)` and `ColonyColor(c)`. A colony stores a
`types.ColonyColor` enum, which is a palette slot, not a colour.

On the firmware this file becomes an rgb565 lookup table. Nothing else changes.

### 5.4 A colony can never go queenless

This is a structural guarantee, not a balance tuning:

- A queen does not age and is immortal
- She only loses health once `Declining` is set
- `Declining` is only set when an heir is born
- Heirs do not age and cannot die while waiting

So a queen is mortal only while an heir exists to replace her. When she dies the
longest-waiting heir is crowned where she stands, the colony centre moves with
her, and every other heir is demoted to worker or nurse, because a colony holds
exactly one queen.

---

## 6. Tuning

Caste roll on maturity, `simulation/matureLarvaeToAnt.go`:

| Roll | Becomes |
|---|---|
| 0 | Queen (1%) |
| 1-20 | Nurse (20%) |
| 21-35 | Soldier (15%) |
| 36-99 | Worker (64%) |

Timing and cost, `simulation/updateWorld.go`:

| Constant | Value | Meaning |
|---|---|---|
| `eggLayingInterval` | 50 | Ticks between eggs |
| `eggHatchTime` | 30 | Ticks per hatch, one egg at a time |
| `larvaeGrowTime` | 50 | Ticks of nursed growth before maturing |
| `foodCost` | 1 | Food units per egg, so 0.1 food |
| `layingThreshold` | 100 | Store needed to lay, so 10 food |
| `queenDeclineInterval` | 30 | Ticks per health point once declining |

Lifespans and health, `types/ant.go`:

| Role | Max age | Health |
|---|---|---|
| Worker | 500 | 100 |
| Soldier | 600 | 150 |
| Nurse | 700 | 80 |
| Queen | 20000, unused while immortal | 200 |
| Larva | 200 | 50 |

---

## 7. Testing

```bash
go test ./...     # 111 tests across 6 packages
```

Two properties are worth calling out because the port depends on them:

- `rng` proves same seed produces the same sequence, that a zero seed does not
  collapse xorshift32 to zeroes, and that `Shuffle` is deterministic.
- `TestSameSeedProducesSameColony` runs 500 ticks twice on one seed and compares
  full colony state.

Everything is deterministic. A flaky test here means a real bug, not luck.

---

## 8. Known problems

`LIFECYCLE-AUDIT.md` has the detail. In short: foraging is a blind random walk,
so workers spend most of their life wandering; the world's food is finite and
never regenerates; and eggs hatch one per 30 ticks regardless of how many are
waiting. These are known and accepted at v0.
