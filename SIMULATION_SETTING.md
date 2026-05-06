# Simulation Settings

## Purpose

This simulation is meant to show how prey grouping changes under predator pressure.

The current code keeps prey movement simple, but not fully random:

- prey keep a persistent heading
- headings change with small random jitter
- predators chase prey, then rest after a catch
- generations restart when prey population falls too low

The result is a repeatable setup for comparing prey movement over time and across generations.

## World Rules

- World size: `100 x 100`
- Boundary behavior: toroidal wrapping on both axes
- Initial population is set when generation `0` is created
- The simulation tracks generation history, including a generation `0` baseline

## Prey Movement

### Heading and Motion

Each prey stores a persistent movement heading.

- On spawn, the heading is random if it has not already been set
- Every tick, the heading is jittered by a bounded random angle
- Movement uses a fixed prey speed of `0.15`
- Position updates wrap around the map edges

This makes prey travel in smooth, continuous paths instead of re-randomizing every frame.

### Movement Modes

Prey switch between two modes based on predator proximity:

- `RoamingMode` when no predator is within vision range
- `EvadingMode` when a predator is within vision range

The active code path currently does this:

- If a predator is found within vision, the prey points away from that predator
- If the prey is in evading mode, it may blend that escape direction toward the nearest prey using `selfishHerdChance`
- The `roaming` blending branch exists in the codebase, but it is currently commented out

So in the current build, the social-blending trait that affects movement is `selfishHerdChance` during evasion.

### Prey Traits

The prey currently initialize these traits:

- `vision`
- `roaming`
- `selfishHerdChance`

In the current movement logic:

- `vision` decides whether the prey switches into evading mode
- `selfishHerdChance` affects how strongly it may herd while escaping
- `roaming` is initialized and preserved for evolution, but its movement branch is not active right now

## Predator Movement

Predators follow a chase-and-rest cycle.

### Chase

- Predator speed is `3x` prey speed
- A predator looks for the nearest prey
- If prey are available, the predator moves toward the nearest one
- Catch radius is `0.5`

### Rest

- After a successful catch, the predator rests for `120` ticks
- If no prey are found, the predator also enters rest mode for `120` ticks
- While resting, the predator does not move

This keeps predators from locking on permanently and makes the pressure on prey more dynamic.

## Generation Cycle

The generation system is based on prey survival.

### When a Generation Ends

A generation restarts when alive prey count drops to `50%` or less of the original initial prey count.

### What Happens on Restart

When a generation restarts:

1. Surviving prey are collected
2. The completed generation stats are recorded
3. Generation counters are advanced
4. All world creatures are cleared
5. The next generation is spawned

### Next Generation Spawn

The next generation is built from survivors:

- If survivors exist, each survivor produces two children
- Children are then filled up to the original prey count if needed
- Predator count is restored to the original starting predator count
- If no survivors remain, the world falls back to a fresh random spawn

## Graph and History

The simulation keeps generation statistics so the graph can update in real time.

Tracked values:

- generation number
- average `selfishHerdChance` of surviving prey
- survivor count

Generation `0` is recorded as a baseline from the current live prey state before the first step completes.

## Practical Notes

- The simulation is currently tuned around simple movement rules, not complex flocking physics
- Prey movement is intentionally lightweight so generation changes are easy to observe
- The `roaming` trait is still part of prey initialization and inheritance, but the active movement logic currently uses `selfishHerdChance` for escape-time social blending