package render

import "simulate/logic"

func (r *raylibRenderer) captureGenerationMetrics() {
	history := logic.GenerationHistory()
	if len(history) == 0 {
		return
	}

	if len(history) == len(r.history) {
		last := history[len(history)-1]
		if r.lastCaptured == last.Generation {
			return
		}
	}

	r.history = r.history[:0]
	for i := range history {
		entry := history[i]
		r.history = append(r.history, generationMetric{
			generation:           entry.Generation,
			avgSelfishHerdChance: entry.AvgSelfishHerdChance,
			survivorCount:        entry.SurvivorCount,
		})
	}
	r.lastCaptured = history[len(history)-1].Generation
}
