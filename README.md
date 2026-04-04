# Simulate

## Project Name, Motive, and Goal

### Project Name

Simulate

### Motive and Goal

I want this project to become an artificial ecosystem where simple local rules create complex emergent behavior, and where I can test biology ideas or try to replicate real life phenomenon through repeatable simulation experiments.

## Expected Functionality

### Phase A: Closed Ecosystem Baseline

- Grid-based world with predators, prey, and resources.
- Energy economy with movement and metabolism costs.
- Predation and feeding mechanics that affect survival.
- Vision and movement rules (including perception constraints).
- Stable simulation loop with reproducible seeded runs.

### Phase B: Heredity and Evolution

- DNA representation and trait mapping.
- Sexual reproduction with crossover and mutation.
- Observable trait inheritance and long-run distribution shift.

### Phase C: Social and Ecological Behaviors

- Warning calls and food-sharing behavior.
- Herding, stalking, and anti-predator strategies.
- Early tests for Hamilton's Rule and niche partitioning.

### Phase D: Experimental Mechanics

- Pheromone/scent systems.
- Viral gene transfer and epigenetic-like effects.
- Environmental cycles, disasters, and dormant states.
- Optional high-chaos mechanics toggled by config.

## Used Skills and Programming Language

### Programming Language

- Go

### Technical Skills Used

- Simulation modeling (agent-based systems)
- Evolutionary systems design
- Data modeling for traits and genetics
- Performance-aware systems programming
- Experiment design with deterministic seeds

## Time Table

### Suggested 8-Week Plan

1. Week 7: Review core behavior and set clear goals for the next steps.
2. Week 8: Improve how agents sense nearby changes and respond.
3. Week 9: Add basic heredity features and adjust random variation.
4. Week 10: Develop group behavior patterns and simple strategy rules.
5. Week 11: Track key results and do broad checks for stability.
6. Week 12: Add optional experiment features with safe on/off controls.
7. Week 13: Connect all parts, fix rough areas, and improve reliability.
8. Week 14: Wrap up results, clean up the project, and finalize the report.

### Deliverables Per Phase

- Phase A deliverable: A stable base simulation where predators and prey can keep interacting over time.
- Phase B deliverable: Basic heredity is working, and trait changes can be seen across generations.
- Phase C deliverable: Social and group behaviors are visible and have a clear effect on outcomes.
- Phase D deliverable: At least two optional experiment features can be turned on and off safely.

## Data Structures and Algorithms That Might Be Used

### Data Structures

- 2D array or matrix for the map/grid state.
- Slice of creature interfaces or typed entities for active agents.
- Struct-based entities for predator, prey, resource, and trait data.
- Bit-packed integer or fixed-size bit array for DNA genome.
- Ring buffer/queue for short-term memory and local event signaling.
- Sparse maps for optional overlays (scent map, heatmap, occupancy stats).

### Algorithms

- Tick-based simulation update loop.
- Local-neighborhood search for perception and interactions.
- Distance-based movement heuristics and steering.
- Probabilistic mutation and crossover for reproduction.
- Selection-pressure evaluation through survival/reproduction outcomes.
- Statistical aggregation for trend analysis (distribution and time-series metrics).

## Notes

- Detailed roadmap and acceptance criteria are documented in [GOAL.md](GOAL.md).