# Simulation Settings

## Purpose

This simulation models predator-prey dynamics to study how the **selfishHerdChance** trait evolves under predation pressure. The trait determines whether prey cooperate or compete when escaping threats, with values ranging from -1 (sacrifice nearby prey) to +1 (use nearby prey as cover).

The simulation is designed to be:
- **Repeatable** — generation milestones and statistics are recorded
- **Analyzable** — graphs show trait evolution and survivor counts per generation
- **Configurable** — trait averages can be set and noise is consistently applied

[Wikipedia for Selfish herd theory](https://en.wikipedia.org/wiki/Selfish_herd_theory)

## Conclusion

All tests converge to around 0.1. The average deviation is quite high sits at around 0.07.
I say this simulation successfully recreates a simple situation for the selfish herd befaviour to come along evolutionary. The result's sitting around 0.1 says that a prey would steer a little bit towards their friend when being chased by predator. The fact that it doesn't steer away is very good for this conclusion and experiment.

---

## World Rules

- **World size:** 100 × 100 units
- **Boundary behavior:** Toroidal wrapping (both X and Y axes wrap around edges)
- **Coordinate system:** Position updates use wrapped distance calculations to account for the torus topology
- **Initial population:** Set when generation 0 is initialized; recorded as a baseline
- **Generation tracking:** History records generation number, average selfishHerdChance, and survivor count

---

## Prey Movement

### Heading and Motion

Each prey maintains a persistent movement heading that evolves through:

- **Initialization:** On spawn, heading is set to a random unit vector if not already assigned
- **Per-tick jitter:** Heading is modified by a random rotation within ±20 degrees
- **Fixed speed:** Prey move at 0.15 units per tick
- **Wrap-aware movement:** Positions wrap around map edges using the toroidal system

This creates smooth, continuous paths that appear directed rather than purely random.

### Movement Modes

Prey dynamically switch between two modes:

- **RoamingMode:** No predator detected within vision range → moves with original heading
- **EvadingMode:** Predator detected within vision range → points away from the nearest predator

### SelfishHerdChance Behavior During Evasion

When a prey is in EvadingMode and a nearby prey is found, the following logic applies:

1. **Probability check:** `rand() < |selfishHerdChance|`
   - Allows the trait to influence behavior probabilistically
   - Absolute value ensures negative values still trigger blending

2. **Direction modifier:**
   - **If selfishHerdChance ≥ 0:** Blend escape direction **toward** nearby prey (selfish/herding)
   - **If selfishHerdChance < 0:** Blend escape direction **away from** nearby prey (sacrificial)

3. **Blending:** The final escape direction averages the pure evasion vector with the prey-relative vector

### Prey Traits

Each prey initializes three traits:

| Trait | Range | Purpose |
|-------|-------|---------|
| `vision` | 5–20 | Distance at which predators trigger evasion mode |
| `roaming` | 0–1 | Reserved for future roaming-mode blending (currently inactive) |
| `selfishHerdChance` | -1–1 | Controls blending behavior during evasion |

**Current movement logic uses:** Only `vision` and `selfishHerdChance`; `roaming` is initialized and inherited but not actively used.

### Trait Initialization and Noise

The simulation uses an **average + noise** model for reproducible but varied runs:

```
initialValue = SetAverage(avg) + noise
noise = (random() - 0.5) × 0.5
result = clamp(initialValue, [-1, 1])
```

For each new prey or generation:

1. A **configurable average** is set via `SetInitialSelfishHerdChance(avg)` (default: 0.75)
2. **Noise is applied:** `(random() - 0.5) × spread` where spread = 0.5
3. **Result is clamped:** Clamped to [-1, 1] to stay within valid bounds

**Effect:** Traits cluster around the set average but vary by ~±0.25 per initialization, allowing experimental control while maintaining population diversity.

**Multirun feature:** Using `-multi N` generates N+1 runs with evenly distributed averages:
- Step size = 2 / N
- Example: `-multi 5` runs with averages: 1.0, 0.6, 0.2, -0.2, -0.6, -1.0
- Each run applies the same ±0.25 noise, creating comparable distributions across the trait space

### Trait Inheritance and Mutation

When the next generation spawns from survivors:

Each survivor produces exactly two children, plus additional offspring fill the population to target count:

- **Roaming mutation:** `parent.roaming + (random() × 0.2 - 0.1)`, clamped to [0, 1]
- **SelfishHerdChance mutation:** `parent.selfishHerdChance + (random() × 0.1 - 0.05)`, clamped to [-1, 1]

This creates gradual evolution with small random steps, preserving beneficial traits across generations.

---

## Predator Movement

Predators follow a **chase-and-rest cycle**:

### Chase Phase

- **Speed:** 0.45 units per tick (3× prey speed)
- **Target:** Nearest alive prey
- **Catch radius:** 0.5 units
- **Movement:** Direct vector toward nearest prey, wrapped for toroidal distance

### Rest Phase

- **Duration:** 120 ticks (fixed)
- **Triggered by:**
  - A successful catch
  - No prey found in the world
- **Behavior:** Stationary; no movement during rest

This cycle prevents predators from permanently locking onto prey and creates dynamic pressure fluctuations.

---

## Generation Cycle

### When a Generation Ends

A generation restarts when prey population falls to **50% or less** of the initial prey count:

```
nextGeneration_if (alivePreyCount ≤ 0.5 × initialPreyCount)
```

### What Happens on Restart

1. **Survivors collected:** All alive prey are recorded
2. **Stats recorded:** Average selfishHerdChance and survivor count saved to history
3. **World cleared:** All prey and predators removed
4. **Next generation initialized:** New population spawned

### Next Generation Spawn

- **From survivors:** Each survivor produces exactly 2 children (mutations applied)
- **Population fill:** If survivors exist but don't reach target count, additional offspring are randomly selected from survivors
- **Predators reset:** Predator count restored to initial count, placed at random positions
- **Fallback:** If no survivors remain, spawn a fresh generation with randomized traits

---

## Graph and History

The simulation records generation statistics for visualization:

**Tracked per generation:**
- Generation number
- Average selfishHerdChance of survivors
- Survivor count

**Baseline:** Generation 0 is captured before the first simulation step completes, serving as a reference point.

**Export:** When `--export` is used with `--multi`, the HTML output includes:
- Line graph showing trait evolution across all runs
- Data table with raw values
- Simulation settings and parameters
- Y-axis range: [-1, 1] to show full spectrum from sacrificial to selfish behavior

---

## Practical Notes

- **Lightweight physics:** Simulation prioritizes observable generational effects over complex flocking
- **Toroidal world:** The wrapped distance system ensures no edge effects; prey can escape across map boundaries
- **Configurable runs:** The `-multi N` flag generates N+1 evenly spaced trait averages from 1.0 to -1.0, enabling controlled comparison
- **Noise model:** Set averages allow reproducible experimental conditions while population variance maintains natural selection pressure
- **Visualization:** Exported HTML graphs show the full -1 to 1 range, making both selfish and sacrificial behavior patterns visible
