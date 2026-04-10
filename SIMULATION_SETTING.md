# Goal

The goal of this simulation is to distinguish predator-driven herding from stronger intrinsic prey grouping by using simple but persistent movement rules. Prey should not reset to a fully random heading every frame. Instead, each prey keeps a direction memory, changes it with small random turning noise, and applies mode-specific social influence chances. Predators still provide ecological pressure through chase and rest cycles. The system should remain simple enough to observe how grouping behavior changes across generations.

# Explanation

Each prey keeps a persistent heading direction. At spawn, this heading starts random. Every tick, the heading is perturbed by bounded angular noise so paths are continuous and smooth rather than frame-by-frame random jumps.

Prey mode selection happens first. If a predator is within vision, the prey enters evading mode. Otherwise, it is in roaming mode.

In roaming mode, social influence is controlled by the `roaming` chance trait. With probability `roaming`, prey blends its current base direction with the nearest prey heading direction. If the chance does not trigger, it keeps its own base direction.

In evading mode, social influence is controlled by a separate `selfishHerdChance` trait. The evasion order is explicit: first point straight away from the predator, then with probability `selfishHerdChance`, blend that escape direction toward the nearest prey position. This separates escape tendency from herd-following tendency during danger.

Predators run a chase-and-rest cycle. They chase visible prey at high speed, and after a catch they rest for fixed ticks. This prevents permanent lock-on behavior and keeps pressure dynamic.

Generations continue to reset when prey count reaches half of the original count. Survivors seed the next generation, and both prey chance traits (`roaming` and `selfishHerdChance`) can mutate with small bounded variance. This allows independent evolution of social influence in safe and dangerous contexts.

The result is a repeatable experiment where you can compare whether prey are merely being pushed by predator pressure or showing stronger endogenous grouping tendencies through persistent heading and split chance traits.