# WhatNext - Future Work and Roadmap

This document tracks planned features, improvements, and the long-term vision for AntFarm Simulator.

---

## 2 - Core Code Fixes

Foundation work before adding new features. Get the basics solid.

- [X] Separate stats out of renderer.go into its own package
- [X] Move logic out of types folder into logic folder 
- [ ] Update controls to actually work (speed up, slow down, pause)
- [ ] Add test files for existing code, boohooo, AI come in handy here :crossed-fingers:
- [ ] Add GitHub Actions to create releases

---

## 3 - Ant Lifecycle

Make ants feel aallive with real consequences.

- [ ] Worker health decreases by 1 per tunnel dug
- [ ] Lifespan system for all ant types (ants die of old age)
- [ ] Queen dies if colony runs out of food
- [ ] Spawn more nurse ants as colony grows
- [ ] "All for one" mindset - ants prioritize protecting the queen

---

## 4 - Soldier Ants

Defense and combat mechanics.

- [ ] Implement soldier patrol behavior
- [ ] Soldiers respond to threats near queen
- [ ] Combat system (soldier vs soldier, soldier vs worker, soldier vs nurse, soldier eat larvae, solider kill queen)
- [ ] Soldiers spawn more when colony detects enemy activity

---

## 5 - MOre Support

Competition and territory.

- [ ] Multiple colonies in one world
- [ ] Random queen placement
- [ ] Visual indicator for colony base/chamber
- [ ] Colony territory surrounded by harder soil
- [ ] Resource competition between colonies

---

## 6 - World Generation Overhaul

More interesting terrain.

- [ ] Use all soil types throughout world (not just sand)
- [ ] Procedural soil distribution based on depth
- [ ] Underground stability/moisture system
- [ ] Custom world configurations
- [ ] Seed-based world generation for reproducibility

---

## 7 - Colony Customization

Custom colonies.

- [ ] Predefined colony templates (aggressive, defensive, balanced)
- [ ] Custom colony creation
- [ ] Colony traits/perks system
- [ ] Starting resource configuration

---

## 8 - UI and Visualization

Make it pretty and informative.

- [ ] GUI launcher instead of `go run main.go`
- [ ] Smoother tcell rendering
- [ ] Minimap for larger colonies
- [ ] Detailed statistics dashboard
- [ ] Population graphs over time
- [ ] Resource trends visualization
- [ ] Zoom in/out

---

## 9 - Emergent Behaviors

Ants acting like real ants.

- [ ] Pheromone trail system
- [ ] Food scent detection
- [ ] Ants follow successful paths
- [ ] Trail decay over time
- [ ] Predator threats (spiders, beetles)
- [ ] Food competition mechanics

---

## 10 - Simulation Features

Full-featured single-machine simulator.

- [ ] Save/load simulation states
- [ ] Timelapse mode (fast forward with smooth visuals)
- [ ] Wipeout mode (flood, shake terminal)
- [ ] Simulation statistics export (CSV, JSON)
- [ ] Replay system

---

## 20 - Microservices Architecture

Prepare for distributed simulation.

- [ ] Extract terminal behavior into its own microservice
- [ ] Extract world behavior into its own microservice
- [ ] Define clean APIs between services or use messages
- [ ] Message queue for ant actions
- [ ] Database for world state persistence

---

## 30 - Distributed Colonies

Multiple terminals, mutlipe worlds, one ecosystem.

- [ ] Stack terminals and connect farms together
- [ ] Ants can traverse between connected terminals
- [ ] Shared world state across instances
- [ ] Network protocol for terminal sync
- [ ] Handle disconnection gracefully (can't unlink mid-simulation)
- [ ] Cross-terminal colony wars

---

## 40 - Hardware in the Loop (HWIL)

Physical MCUs running ant brains. *Separate README for this phase.*

- [ ] Serial/BLE protocol for external MCU ants
- [ ] Sensor simulation (what the ant "sees")
- [ ] Action parsing (movement commands from MCU)
- [ ] Mixed simulation: software + hardware ants
- [ ] Use my ESp32 HAL code 
- [ ] Latency compensation
- [ ] Hardware ant registration/discovery

---

## 50 - Physical

Real robots in the real world.

- [ ] 3D printable ant robot designs
- [ ] ESP32/nRF firmware templates
- [ ] Real-world to simulation bridge
- [ ] Camera-based position tracking
- [ ] Swarm coordination algorithms, hopefullyyyyy 
- [ ] Physical obstacle detection

---

1. **v0.2** - Clean up code first, makes everything else easier
2. **v0.3** - Lifecycle gives the simulation meaning
3. **v0.4** - Soldiers make multi-colony interesting
4. **v0.5** - Multi-colony is the real game
5. **v0.6** - Better worlds to fight over
6. **v0.8** - UI polish (can do alongside other work)
7. **v0.9** - Emergent behavior is the magic
8. **v1.0** - Save/load before going distributed
9. **v2.0+** - Microservices only when single-machine is solid

---

## Ideas Parking Lot

Random ideas to explore later:

- Ant genetics (traits passed to offspring)
- Seasonal cycles affecting food spawns
- Underground fungus farming
- Ant communication animations
- Sound effects for digging/combat and happy music for bg
- Mobile companion app showing colony stats


---

*Last updated: January 2026*