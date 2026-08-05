# Lifecycle Audit: Colony Collapse and Deadlock

Status: partially fixed (see 4.5). Last updated 2026-08-05.
Companion to `WHATNEXT.md`. Blocks `FIRMWARE-SPEC.md` milestones M3 and M5.

---

## 1. Summary

The colony does not stall. It **dies, and then hard-locks forever**. Once the
last worker dies the simulation enters an absorbing state that no amount of
elapsed time can escape. The queen remains alive, the map still holds hundreds
of units of food, and nothing can ever change again.

Observed over 4000 ticks (`NewWorld(120, 35)`, colony at `(30, 11)`):

```
tick   food   eggs   larvae  workers  grassLeft
400    0      3      2       3        115
800    0      2      1       5        102
1200   0      1      1       4        95
1600   0      0      1       8        83
2000   0      1      1       3        80
2400   0      0      0       2        79
2800   0      0      0       1        79
3200   5      0      0       0        77   <- last worker dies
3600   5      0      0       0        77   <- frozen
4000   5      0      0       0        77   <- frozen
```

Final state: queen alive, `TotalEggsLaid = 51`, 77 unharvested grass cells
(~385 food) sitting on the surface, permanently unreachable.

---

## 2. Root cause: the deadlock

Egg laying is gated on `colony.Food >= 10` (`simulation/updateWorld.go:43`).
Food ends the run at **5**.

Only workers gather food. Nurses only tend larvae and guard the nursery;
soldiers only ever set `CurrentAction = "patrolling"` and do nothing else.
So the dependency cycle has no entry point:

```
no workers -> no food -> food < 10 -> no eggs -> no larvae -> no workers
```

This is an absorbing state, not slow balance. It is unrecoverable by design.

---

## 3. Why workers cannot keep up

Action histogram, 1200 ticks, peak 12 workers:

| Worker action | Ticks |
|---|---|
| `exploring` | 6076 |
| `bringing food to queen` | 1310 |
| `deposited food` | **39** |
| `foraged grass` | 38 |
| `picked up food` | 5 |
| `stuck with food` | 46 |

Three compounding failures:

1. **Foraging is a blind random walk.** 6076 ticks of wandering produced 43 food
   finds. `MoveRandomly` (`pathfinder/workerpathfinder.go:18`) has directional
   momentum but zero food-seeking: no scent, no gradient, no memory.
2. **The return trip burns the profit.** 1310 ticks spent carrying food yielded
   only 39 successful deposits. `BringFoodToQueen` is greedy-directional with no
   obstacle memory, so workers grind against soil walls (`stuck with food`).
3. **Workers die mid-trip and the cargo dies with them.** `WorkerMaxTick = 500`
   (`types/ant.go`), and `DigAndMove` costs 1 HP per dig
   (`pathfinder/pathfinder.go:124`). Workers were observed at hp=3 and hp=5. A
   queen-to-surface round trip via random walk routinely exceeds a worker's
   entire lifespan. `processDeaths` acknowledges the loss directly:
   *"If worker was carrying food, it's lost."*

Net throughput: supply ~0.2 food/tick against demand of up to 0.5 food/tick.

---

## 4. Secondary defects

### 4.1 Eggs queue faster than they can hatch

Laying produces up to 5 eggs per 50 ticks. Hatching is a single `colony.Eggs--`
guarded by `world.Ticks%eggHatchTime == 0` (`simulation/updateWorld.go:61`):
exactly **one egg per 30 ticks**, regardless of backlog depth.

Fill rate 0.1 eggs/tick vs drain rate 0.033 eggs/tick. The backlog was observed
climbing to 14 and still rising. Food is charged at lay time (5 each), so the
colony pays immediately for ants that may never hatch, precisely when food is
scarcest.

### 4.2 Food supply is finite and never regenerates

`types/world.go:NewWorld` is the only food spawner in the codebase. Surface
pellets are one-shot; grass is one-shot and marked spent with `Food = -1`
(`simulation/antsBehavior.go:81`). Total world budget is roughly 660 food for
the entire run. Even a perfectly efficient colony eventually starves. The only
question is when.

### 4.3 Pellet value mismatch

`NewWorld` sets `grid[1][x].Food = 5`, but a worker picking it up takes
`worker.FoodAmount = 10` (`simulation/antsBehavior.go:68`). Free food from
nothing.

### 4.4 Emptied pellet cells re-qualify as grass

The grass check is `currentCell.Food >= 0` (`simulation/antsBehavior.go:78`). A
pellet cell that was just harvested is set to `Food = 0`, which satisfies `>= 0`
and yields another 5 food as "grass". Double-dip.

### 4.5 Outcomes are non-deterministic and wildly divergent (FIXED)

`math/rand` was called directly from `simulation/updateWorld.go`,
`types/world.go` and `pathfinder/workerpathfinder.go`. Two runs of identical
code diverged from "plateaus at 10 ants" to "collapses to 0", so there was no
way to tell whether a balance change helped or the dice simply landed
differently.

Resolved: the `random` package now provides a seeded xorshift32 generator, injected
through `World.Rng`. Nothing in the simulation calls the global pool. Re-measure
any balance change against a fixed seed.

---

## 5. Impact on the hardware port

This audit is a hard gate on two milestones in `FIRMWARE-SPEC.md`:

- **M3 acceptance** requires "a self-sustaining colony grows and stabilises over
  ~30 minutes." At 1 Hz that is 1800 ticks. The Go reference implementation is
  already collapsing by then and dead by ~3000. **A faithful port cannot pass
  M3, because the thing being ported does not satisfy it.**
- **M5 parity gate** diffs Go against C++ on a shared seed. Any economy change
  made after the port invalidates the parity run and forces a re-diff. The fix
  belongs in Go first, which is also what `PHYSICAL-PLAN.md` §8 step 1 assumes.
- §4.1 blocker 4 (inject a seeded PRNG) is upgraded from "nice for parity" to
  **prerequisite for fixing anything**, per 4.5 above.

---

## 6. Candidate fixes

Not yet chosen. Listed cheapest-first; they are complementary, not exclusive.

| # | Fix | Scope | Notes |
|---|---|---|---|
| 1 | **Deadlock guard** | Small | Make 0 workers recoverable: emergency worker spawn, or align the laying gate with `foodCost`. Removes the hard-lock; colony still starves. |
| 2 | ~~**Seeded PRNG injection**~~ | Done | Landed as the `random` package, injected via `World.Rng`. Port blocker 4 cleared. |
| 3 | **Drain the egg queue properly** | Small | Hatch proportional to backlog rather than one per 30 ticks, or charge food at hatch time instead of lay time. |
| 4 | **Economy rebalance** | Medium | Egg cost, laying gate, worker lifespan vs round-trip length, HP cost per dig. |
| 5 | **Food scent detection** | Medium | Pull forward from `WHATNEXT.md` §9. Attacks the root inefficiency, the blind random walk, and is the highest gameplay payoff. |
| 6 | **Renewable food** | Medium | Grass regrowth and/or surface respawn, so the world is not a fixed ~660-food budget guaranteeing eventual death. |

Minimum set to clear the M3 gate: **1 + 2 + 3**, then measure before deciding
whether 4, 5 and 6 are still needed.

---

## 7. Reproducing

No harness is committed. To re-measure, write a throwaway `package main` inside
the module that calls `types.NewWorld` + `types.NewColony` + `logic.AddColony`,
then loops `logic.UpdateWorld` and prints colony counters. Delete it afterward.
It must not ship in the repo.

Pass a fixed seed (`types.NewWorld(w, h, random.New(seed))`) so the run can be
repeated exactly.

Note that `simulation` is imported as package name `logic` (the directory was
renamed but the package declaration was not).
