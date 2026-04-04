# Simulate Project Goals

This document is the conceptual blueprint for the simulate project, moving from basic survival mechanics to complex evolutionary modeling.

I want this project to become an artificial ecosystem where simple local rules create complex emergent behavior, and where I can test biology ideas or try to replicate real life phenomenon through repeatable simulation experiments.

## Roadmap Overview

- Phase A: Closed ecosystem baseline (MVP foundation)
- Phase B: Heredity and evolution engine
- Phase C: Social and ecological strategy systems
- Phase D: Experimental and high-chaos mechanics

## Phase A: Closed Ecosystem Baseline

Goal: Build a stable, self-contained predator-prey-resource loop with measurable outcomes.

### Core Mechanics

- Energy economy (metabolism):
	- Every creature has an energy float.
	- Moving costs energy proportional to $speed^2$.
	- Standing still has a base metabolic cost.
- Food chain:
	- Plants/resources spawn on grid tiles.
	- Animals gain energy when eating resources.
	- Predators gain a larger energy boost by occupying prey tile and removing prey from population.
- Perception and movement:
	- Circular vision range (radius from Vision trait).
	- Field of view (FOV) with heading so creatures cannot see behind them unless turning.
- Terrain interaction:
	- Water raises movement cost.
	- Bushes reduce predator detection of prey inside them.

### MVP Acceptance Criteria

- The simulation runs 1000 ticks without crashes or deadlock.
- Predator and prey populations both change over time (not flat line), and extinction is not immediate in most seeded runs.
- Energy is conserved by clear sources/sinks (resource intake, movement cost, base metabolism, predation gain).
- At least 20 seeded runs produce comparable macro patterns (basic reproducibility).

## Phase B: Heredity and Evolution Engine

Goal: Introduce a genetic architecture that allows selection, mutation, and adaptation across generations.

### Genetics and Reproduction

- Genome model:
	- 64-bit DNA string/slice.
	- 2-bit base encoding: 00=A, 01=T, 10=C, 11=G.
- Trait mapping:
	- DNA maps to speed, vision, fear, metabolism, and sociality.
	- Define explicit bit ranges per trait for interpretability.
- Reproduction:
	- Sexual reproduction with male/female assignment.
	- Mating requires proximity and sufficient energy.
	- Offspring generated via crossover and mutation.
- Selection pressures:
	- Natural selection, mutation, genetic drift, and gene flow between sub-grids.

### Phase B Acceptance Criteria

- Trait values are inherited with measurable parent-offspring correlation.
- Mutation rate is configurable and visible in run logs.
- Directional/stabilizing/disruptive selection can each be induced by environment presets.
- At least one trait distribution shifts significantly over long runs (statistical evidence, not anecdotal).

## Phase C: Social and Ecological Strategy Systems

Goal: Add realistic group and behavior strategies that affect fitness and survival.

### Social and Tactical Behaviors

- Sociality and altruism:
	- Warning calls to nearby allies at energy cost.
	- Food sharing to assist starving kin.
- Hunting and anti-predator strategy:
	- Predator stalking mode (lower speed, lower detectability).
	- Herd behavior and dilution effect among genetically similar animals.
- Sexual selection:
	- Runaway selection on costly display traits.

### Theoretical Models to Test

- Hamilton's Rule: $rb > c$ for emergence of helping behavior.
- Niche partitioning: species diverge to consume different resources.

### Phase C Acceptance Criteria

- Helping behavior frequency increases under high-relatedness conditions.
- Grouping behavior reduces average individual predation risk.
- Distinct ecological roles emerge in at least one scenario (resource specialization or habitat preference).

## Phase D: Experimental and High-Chaos Mechanics

Goal: Explore advanced emergent effects beyond core biological realism.

### Experimental Features

- Pheromone trails and scent markers:
	- Decaying scent layers for territory, alarm, and mating.
- Viral evolution and horizontal gene transfer:
	- Viral entities carry DNA snippets and modify host genomes.
- Environmental cycles:
	- Day/night and seasonal modifiers on metabolism and vision.
	- Disasters such as floods and migration bottlenecks.
- Mimicry and deception:
	- Batesian and aggressive mimicry.
- Social learning and culture:
	- Short action-memory buffers copied from successful neighbors.
- Epigenetic-like expression tags:
	- Survival experience biases offspring expression.
- Endosymbiosis events:
	- Rare species-merging outcomes with combined traits.
- Hibernation/spore states:
	- Dormancy during resource crashes with potential reactivation.

### Phase D Acceptance Criteria

- At least two experimental systems can be toggled on/off independently.
- Experimental features do not break baseline simulation invariants.
- Emergent outcomes are logged and replayable from seed + config.

## Data and Validation Plan (Applies Across All Phases)

- Population dynamics:
	- Predator/prey curves for Lotka-Volterra-like oscillations.
- Trait histograms:
	- Distribution tracking for vision, speed, fear, and sociality.
- Energy heatmaps:
	- Spatial concentration of feeding, deaths, and movement costs.
- Fitness landscape:
	- DNA-space vs survival/reproduction outcomes.

## Scope Guardrails

- Must-have for MVP:
	- Entire Phase A.
- Should-have before broad experiments:
	- Phase B core genetics and inheritance observability.
- Not in MVP:
	- Most Phase D features unless they directly support a current hypothesis test.

## Design Consistency Notes

- Prefer deterministic, seed-based runs for debugging and scientific comparison.