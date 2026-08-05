# AntFarm Simulator

[![Platform](https://img.shields.io/badge/Platform-Terminal-black?style=flat&logo=gnometerminal)](https://github.com/gdamore/tcell)
[![CI](https://github.com/okeith12/antfarm-simulator/actions/workflows/ci.yml/badge.svg?cache=false)](https://github.com/okeith12/antfarm-simulator/actions/workflows/ci.yml)
[![Version](https://img.shields.io/github/v/tag/okeith12/antfarm-simulator?label=version)](https://github.com/okeith12/antfarm-simulator/releases)

> **I am interested in becoming a digital Top G and this is just the start of my bionic ecosystem idea.** This ant simulator is v0, the foundation for something much bigger. Next up: a small square panel you prop on a desk, where one pixel is one ant and the colony digs its own maze. After that, **Hardware-in-the-Loop (HWIL) ants** where physical MCUs run the ant brains while this simulator runs their world. Then the 3D printed ants.

---

## What Is This?

A terminal ant colony simulator in Go. The queen lays eggs, nurses tend larvae,
workers dig tunnels and forage, and heirs eventually take the throne. All of it
renders in ASCII through [tcell](https://github.com/gdamore/tcell).

```
🌱🌱🌾🌱🌱🌱🌱🌱🌱🌱🌱🌱   <- Surface (food spawns here)
░░░░░░░░░░░░░░░░░░░░░░░░   <- Sand layer
░░░░░▢▢▢░░░░░░░░░░░░░░░░   <- Queen's chamber
░░░░░▢♛○░░░░░░░░░░░░░░░░   <- ♛ Queen  ○ Nurse  ● Worker
░░░░░▢○▢░░░░░░░░░░░░░░░░
░░░░░░●░░░░░░░░░░░░░░░░░   <- Tunnels being dug
░░░░░░░●░░░░░░░░░░░░░░░░
```

---

## Getting Started

Requires Go 1.24+ and a terminal with UTF-8 and 256-colour support.

```bash
git clone https://github.com/okeith12/antfarm-simulator.git
cd antfarm-simulator
go mod tidy
go run main.go
```

Build a binary instead:

```bash
go build -o antfarm
./antfarm
```

---

## Controls

All of these work.

| Key | Action |
|---|---|
| `Q` / `ESC` | Quit |
| `L` | Toggle the activity log |
| `P` | Pause and resume |
| `+` / `=` | Speed up |
| `-` | Slow down |

Six speed presets, from 0.25x to 10x: `0.25, 0.5, 1, 2, 5, 10` ticks per second.
Rendering stays at 30 FPS independently of simulation speed.

---

## Ant Roles

| Role | Icon | Behaviour |
|---|---|---|
| **Queen** | ♛ | Stays in her chamber and lays one egg every 50 ticks, costing 0.1 food. Does not age. |
| **Nurse** | ○ | Guards the nursery, moves to larvae and tends them until they mature. |
| **Worker** | ● | Wanders with directional momentum, digs tunnels, forages the surface and carries food back. |
| **Soldier** | ⚔ | Patrols. Combat is not implemented. |
| **Larva** | ◦ | Waits for a nurse. With care and 50 ticks of age it matures into an adult. |

When a larva matures it rolls for a caste: **1% queen, 20% nurse, 15% soldier,
64% worker.**

### Succession

A colony holds exactly one reigning queen, and **can never go queenless.** That
is structural rather than lucky:

- A queen does not age and is immortal
- She only begins losing health once she has borne an heir
- Heirs do not age and cannot die while they wait

So a queen is only mortal while an heir exists to replace her. When she dies the
longest-waiting heir is crowned where she stands, the colony centre moves with
her, and any other heirs give up the claim and become workers or nurses.

---

## Architecture

Full detail in **[ARCHITECTURE.md](ARCHITECTURE.md)**.

```
gui/  ───▶  simulation/  ───▶  pathfinder/  ───▶  types/  ───▶  random/
 │                                                               ▲
 └── tcell                          nothing right of gui/ ────────┘
                                    imports tcell
```

```bash
grep -rn "tcell" types simulation pathfinder util random main.go   # empty
```

### Project structure

```
antfarm/
├── main.go              # Entry point
│
├── types/               # The nouns
│   ├── world.go         # World, flat row-major grid
│   ├── cell.go          # Cell, Soil, FoodScale
│   ├── colony.go        # Colony, ColonyColor
│   ├── ant.go           # Base Ant + AntInterface, lifespans, health
│   ├── queen.go         # QueenAnt, including the Declining flag
│   ├── nurse.go  worker.go  solider.go  larvae.go
│   └── log.go
│
├── simulation/          # The verbs (declares `package logic`)
│   ├── updateWorld.go       # One tick, deaths, laying, hatching, succession
│   ├── antsBehavior.go      # Per-role behaviour
│   ├── antPlacement.go      # AddColony, PlaceAnt, RemoveAnt, MoveWorldAnt
│   ├── spawn.go             # Spawning, removal, heir demotion
│   └── matureLarvaeToAnt.go # The caste roll
│
├── pathfinder/          # Movement
│   ├── pathfinder.go        # Directions, CanMoveTo, CanDigTo, Move, DigAndMove
│   ├── workerpathfinder.go  # Random walk with momentum, food delivery
│   └── nursepathfinder.go   # Guard nursery, move to larvae, queen swap
│
├── gui/                 # Terminal rendering, the only package that sees tcell
│   ├── antfarm.go       # Game loop, input, speed and pause
│   └── renderer.go  stats.go  controls.go  colors.go
│
├── random/                 # Deterministic xorshift32
└── util/                # Abs()
```

**Naming trap:** the `simulation/` directory declares `package logic`, left over
from a rename. Import it as `logic "antfarm/simulation"`.

### Two decisions that shape everything

| | |
|---|---|
| Food | scaled integer, `FoodScale = 10`. No floats anywhere in the simulation. |
| Randomness | injected xorshift32 on the World. No `math/rand`. Same seed, same colony. |

```go
world := types.NewWorld(120, 35, random.New(seed))   // clock seed in the app, fixed in tests
```

---

## Configuration

Simulation speed and frame rate, `gui/antfarm.go`:

```go
var speedPresets = []float64{0.25, 0.5, 1, 2, 5, 10}
const renderFPS = 30
```

Timing and cost, `simulation/updateWorld.go`:

```go
eggLayingInterval    = 50   // ticks between eggs
eggHatchTime         = 30   // ticks per hatch, one egg at a time
larvaeGrowTime       = 50   // nursed growth before maturing
foodCost             = 1    // food units per egg, so 0.1 food
layingThreshold      = 100  // store needed to lay, so 10 food
queenDeclineInterval = 30   // ticks per health point once declining
```

---

## Testing

```bash
go test ./...     # 111 tests across 6 packages
```

Everything is deterministic, so a flaky test here means a real bug rather than
bad luck. Two properties matter most: `random` proves that one seed always gives one
sequence, and `TestSameSeedProducesSameColony` runs 500 ticks twice on a single
seed and compares the whole colony.

---

## Known Problems

Detail in **[LIFECYCLE-AUDIT.md](LIFECYCLE-AUDIT.md)**. In short: foraging is a
blind random walk so workers spend most of their lives wandering, the world's
food is finite and never regenerates, and eggs hatch one per 30 ticks no matter
how many are queued. Known and accepted at v0.

---

## Roadmap

`WHATNEXT.md` has the long list. The direction:

- **v0, here now.** Terminal. Queen cycle, nursing, digging, foraging, succession.
- **v0.1.** A small square panel on the desk. One pixel per ant, the colony digs
  its own maze, and you watch it like a toy.
- **v1.** Hardware in the loop. Physical MCUs run ant brains, the simulator runs
  their world.
- **v2.** 3D printed ant bodies.

---

## Contributing

This project is the foundation for a larger bionic ecosystem experiment. If you're interested...
...dont be...justkidding...
...pull requests are welcome! Please open an issue first to discuss major changes.

---

## License

I mean I just put my go code together that anyone can do so do whatever you want with it, just don't blame me for whatever dumpster fire is created.

---

## Acknowledgments

- [tcell](https://github.com/gdamore/tcell), terminal cell library for Go
- Real ants, for being fascinating little engineers
- Myself, for completing a projct

---

<p align="center">
  <i>"The colony is the organism. The ant is just a cell."</i>
</p>
