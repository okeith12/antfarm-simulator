# AntFarm Simulator

[![Platform](https://img.shields.io/badge/Platform-Terminal-black?style=flat&logo=gnometerminal)](https://github.com/gdamore/tcell)
[![CI](https://github.com/okeith12/antfarm/actions/workflows/ci.yml/badge.svg)](https://github.com/okeith12/antfarm/actions/workflows/ci.yml)
[![Version](https://img.shields.io/github/v/tag/okeith12/antfarm?label=version)](https://github.com/okeith12/antfarm/releases)

> **I am interested in becoming a digital Top G and this is just the start of my bionic ecosystem idea.** This ant simulator is v0 — the foundation for something much bigger. Next up: **Hardware-in-the-Loop (HWIL) ants** where physical MCUs run the ant brains while this simulator runs their world. Then comes the 3D printed ants. Join me on this journey from terminal ants to real-world bionic colonies. 

---

##  What Is This?

A terminal-based ant colony simulator written in Go. Watch your colony as the queen lays eggs, nurses tend to larvae, and workers dig tunnels through procedurally generated terrain — all rendered in ASCII art using [tcell](https://github.com/gdamore/tcell).

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
## Architecture Overview

Here is the antchitecture overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              TERMINAL                                    │
│                         (tcell.Screen)                                   │
└─────────────────────────────────────────────────────────────────────────┘
                                   ▲
                                   │ Render()
                                   │
┌─────────────────────────────────────────────────────────────────────────┐
│                               gui/                                       │
│  ┌─────────────────┐  ┌─────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │   renderer.go   │  │  stats.go   │  │ controls.go  │  │ colors.go │ │
│  │                 │  │             │  │              │  │           │ │
│  │ • Render()      │  │ • render    │  │ • render     │  │ • Color   │ │
│  │ • ToggleLog()   │  │   Stats()   │  │   Controls() │  │   consts  │ │
│  │                 │  │ • render    │  │              │  │           │ │
│  │                 │  │   Activity  │  │              │  │           │ │
│  │                 │  │   Log()     │  │              │  │           │ │
│  └─────────────────┘  └─────────────┘  └──────────────┘  └───────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                                   ▲
                                   │ reads
                                   │
┌─────────────────────────────────────────────────────────────────────────┐
│                              types/                                      │
│                                                                          │
│  ┌─────────┐    ┌──────────────────────────────────────────────┐        │
│  │  World  │───▶│                  Colony                       │        │
│  │         │    │                                                │        │
│  │ • Grid  │    │  ┌───────┐ ┌───────┐ ┌────────┐ ┌──────────┐ │        │
│  │ • Width │    │  │ Queen │ │ Nurse │ │ Worker │ │ Soldier  │ │        │
│  │ • Height│    │  │  ♛    │ │   ○   │ │   ●    │ │    ⚔     │ │        │
│  │ • Ticks │    │  └───┬───┘ └───┬───┘ └───┬────┘ └────┬─────┘ │        │
│  └────┬────┘    │      │         │         │           │        │        │
│       │         │      └─────────┴─────────┴───────────┘        │        │
│       │         │                    │                           │        │
│       ▼         │              implements                        │        │
│  ┌─────────┐    │                    ▼                           │        │
│  │  Cell   │    │           ┌──────────────┐    ┌─────────┐     │        │
│  │         │    │           │ AntInterface │    │ Larvae  │     │        │
│  │ • Soil  │    │           │              │    │   ◦     │     │        │
│  │ • Food  │    │           │ • GetAnt()   │    └─────────┘     │        │
│  │ • Tunnel│    │           │ • GetIcon()  │                     │        │
│  │ • Ant   │    │           │ • GetRole()  │                     │        │
│  └─────────┘    │           └──────────────┘                     │        │
│                 └──────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────────────────────┘
                                   ▲
                                   │ updates
                                   │
┌─────────────────────────────────────────────────────────────────────────┐
│                              logic/                                      │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────────┐              │
│  │   world.go    │  │    ant.go     │  │ world_colony.go│              │
│  │               │  │               │  │                 │              │
│  │ • UpdateWorld │  │ • update      │  │ • AddColony()   │              │
│  │ • updateColony│  │   Worker()    │  │ • PlaceAnt()    │              │
│  │ • egg hatching│  │ • update      │  │ • RemoveAnt()   │              │
│  │ • larvae grow │  │   Nurse()     │  │ • MoveWorldAnt()│              │
│  └───────────────┘  └───────────────┘  └─────────────────┘              │
│                                                                          │
│  ┌─────────────────┐                                                    │
│  │colony_ants.go│                                                    │
│  │                 │                                                    │
│  │ • SpawnWorker() │                                                    │
│  │ • SpawnNurse()  │                                                    │
│  │ • SpawnLarvae() │                                                    │
│  │ • RemoveLarvae()│                                                    │
│  └─────────────────┘                                                    │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
                                   │ uses
                                   ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           pathfinder/                                    │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────────┐      │
│  │  pathfinder.go  │  │workerpathfinder  │  │ nursepathfinder    │      │
│  │                 │  │                  │  │                    │      │
│  │ • Direction     │  │ • MoveRandomly() │  │ • GuardNursery()   │      │
│  │ • CanMoveTo()   │  │ • BringFoodTo    │  │ • MoveTowardLarvae │      │
│  │ • CanDigTo()    │  │   Queen()        │  │ • Queen swap logic │      │
│  │ • MoveAnt()     │  │ • pickNewDir()   │  │                    │      │
│  └─────────────────┘  └──────────────────┘  └────────────────────┘      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
antfarm/
├── main.go                 # Entry point, game loop, event handling
├── go.mod                  # Go module 
├── go.sum                  
│
├── types/                  # the "whats"
│   ├── world.go            # World struct
│   ├── cell.go             # Individual grid cells (soil, tunnel, food)
│   ├── colony.go           # Colony struct
│   ├── ant.go              # Base Ant + AntInterface
│   ├── queen.go            # QueenAnt - lays eggs
│   ├── nurse.go            # NurseAnt - tends larvae
│   ├── worker.go           # WorkerAnt - digs, forages
│   ├── solider.go          # SoldierAnt - defends (WIP)
│   ├── larvae.go           # LarvaeAnt - baby ants
│   └── log.go              # Activity logging helper
│
├── logic/                  # the "how"
│   ├── world.go            # UpdateWorld(), updateColony(), egg/larvae lifecycle
│   ├── ant.go              # Ant behavior dispatchers
│   ├── world_colony.go    # AddColony(), PlaceAnt(), RemoveAnt(), MoveWorldAnt()
│   └── colony_ant.go   # SpawnWorker(), SpawnNurse(), SpawnLarvae(), RemoveLarvae()
│
├── pathfinder/             # Movement and navigation
│   ├── pathfinder.go       # Shared utilities, directions, movement
│   ├── workerpathfinder.go # Worker-specific: random walk, food delivery
│   └── nursepathfinder.go  # Nurse-specific: guard nursery, tend larvae
│
├── gui/                    # Terminal rendering
|   ├── antfarm.go          # Main entrypoint and orchestrator 
│   ├── renderer.go         # Render(), ToggleLog()
│   ├── stats.go            # renderStats(), renderActivityLog()
│   ├── controls.go         # renderControls()
│   └── colors.go           # Color constants
│
└── util/                   # Helpers
    └── abs.go              # Abs() for integers
```

---

## Simulation Flow

```
main.go
   │
   ├──▶ Initialize tcell.Screen
   │
   ├──▶ Create World (width × height grid)
   │         │
   │         └──▶ Generate terrain layers
   │         └──▶ Scatter food on surface
   │
   ├──▶ Create Colony at position
   │         │
   │         ├──▶ Spawn Queen (ID: 0)
   │         └──▶ Spawn Head Nurse (ID: 1)
   │
   └──▶ Game Loop
            │
            ├──▶ Handle Input (ESC/Q to quit, L to toggle log)
            │
            ├──▶ Simulation Tick (1 Hz default)
            │         │
            │         └──▶ logic.UpdateWorld()
            │                   │
            │                   ├──▶ Queen lays eggs (if enough food)
            │                   ├──▶ Eggs hatch into Larvae
            │                   ├──▶ Larvae + NurseCare → Workers
            │                   ├──▶ Update each Nurse behavior
            │                   ├──▶ Update each Worker behavior
            │                   └──▶ Update each Soldier behavior
            │
            └──▶ Render (30 FPS)
                      │
                      └──▶ gui.Renderer.Render()
                                │
                                ├──▶ Draw terrain grid
                                ├──▶ Draw ants with role icons
                                ├──▶ Draw stats bar
                                └──▶ Draw activity log
```

---

##  Ant Roles & Behaviors

| Role | Icon | Behavior |
|------|------|----------|
| **Queen** | ♛ | Stays in chamber. Lays 1-5 eggs every 50 ticks (costs 10 food each). |
| **Nurse** | ○ | Guards nursery near queen. When larvae spawn, moves to them and provides care until they mature into workers. |
| **Worker** | ● | Explores randomly using "tryna be" ant like movement. Digs tunnels through sand. Forages food from surface and delivers to queen. |
| **Soldier** | ⚔ | Patrols (WIP - combat not implemented yet). |
| **Larvae** | ◦ | Waits for nurse care. After receiving care + 50 ticks of age → becomes Worker. |

---

## Controls

| Key | Action |
|-----|--------|
| `Q` / `ESC` | Quit simulation |
| `L` | Toggle activity log |
| `Ctrl+>` | Speed up (not yet implemented) |
| `Ctrl+<` | Slow down (not yet implemented) |
| `P` | Pause (not yet implemented) |

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- A terminal with UTF-8 and 256-color support

### Installation

```bash
# Clone the repository
git clone https://github.com/okeith12/antfarm.git
cd antfarm

# Download dependencies
go mod tidy

# Run the simulation
go run main.go
```

### Build

```bash
# Build binary
go build -o antfarm

# Run
./antfarm
```
When I actually get somewhere with this then you can 
Download the latest release for your platform from the [Releases](https://github.com/okeith12/antfarm/releases) page.

**macOS / Linux:**

```bash
# Make it executable
chmod +x antfarm

# Run it
./antfarm
```

**Windows:**

```
antfarm.exe
```
---

## Configuration

Key constants in `main.go`:

```go
const (
    simulationUpdatesPerSecond = 1   // Simulation speed 
    renderFPS                  = 30  // Frames per second
)
```

Key timing in `logic/world.go`:

```go
var (
    eggLayingInterval = 50  // Ticks between egg batches
    eggHatchTime      = 30  // Ticks for egg → larvae
    larvaeGrowTime    = 50  // Ticks for larvae → worker (with nurse care)
)
```

---

## 🗺️ Roadmap

### v0.1 - Current
- [x] Basic world generation
- [x] Queen egg-laying cycle
- [x] Nurse-larvae care system
- [x] Worker digging & foraging
- [x] Terminal rendering with tcell
- [ ] Activity logging system

- [ ] Study Ant behavior and read up on them to further ehance the simulator

### v0.2 - Planned SUbjected to change, 
- [ ] Multiple colonies with different colors
- [ ] Soldier patrol and combat
- [ ] Pheromone trail system
- [ ] Food scent detection
- [ ] Colony statistics dashboard

### v1.0 - HWIL Integration YASSSSSS
- [ ] Serial/BLE protocol for external MCU ants
- [ ] Sensor simulation (what the ant "sees")
- [ ] Action parsing (movement commands from MCU)
- [ ] Mixed simulation: software + hardware ants

### v2.0 - Physical 
- [ ] 3D printable ant robot designs
- [ ] ESP32/nRF firmware templates
- [ ] Real-world ↔ simulation bridge

---

##  Contributing

This project is the foundation for a larger bionic ecosystem experiment. If you're interested...
...dont be...justkidding...
...pull requests are welcome! Please open an issue first to discuss major changes.

---

## License

I mean I just put my go code together that anyone can do so do whatever you want with it, just don't blame me for whatever dumpster fire is created.

---

## Acknowledgments

- [tcell](https://github.com/gdamore/tcell) - Terminal cell library for Go
- Real ants - for being fascinating little engineers
- Myself - for completing a projct

---

<p align="center">
  <i>"The colony is the organism. The ant is just a cell."</i>
</p>
