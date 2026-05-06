package main

import (
	"flag"
	"fmt"
	"os"
	"simulate/logic"
	"simulate/render"
	"time"
)

func main() {
	gen := flag.Int("gen", 0, "run with graphics until generation N completes, then stop simulation updates and keep graph visible")
	export := flag.String("export", "", "after simulation ends, export graph and settings to HTML file (e.g., './results.html')")
	multi := flag.Bool("multi", false, "run multiple headlessly with starting averages: 0.0, 0.2, 0.4, 0.6, 0.8, 1.0 down to the target gen")
	flag.Parse()

	if *multi && *export != "" && *gen > 0 {
		runMultiHeadless(*gen, *export)
		return
	}

	logic.InitSimulation()
	logic.SetCycle(0.2)

	logic.AddPredators(4)
	logic.AddAnimals(80)

	cols, rows := logic.Dimensions()
	renderer, err := render.NewRenderer(cols, rows, render.Config{
		Title:        "Simulate",
		FPS:          60,
		WindowHeight: 1200,
		WindowWidth:  1200,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Close()

	targetGen := *gen
	graphShown := false

	for renderer.IsRunning() {
		if targetGen <= 0 {
			logic.Step()
		} else if !graphShown {
			for logic.LastCompletedGeneration() < targetGen {
				logic.Step()
			}

			if tabSetter, ok := renderer.(interface{ ShowGraphTab() }); ok {
				tabSetter.ShowGraphTab()
			}
			graphShown = true
		}

		renderer.Draw(logic.Snapshot())
	}

	// Export results if requested
	if *export != "" {
		exportResults(*export)
	}
}

func exportResults(filePath string) {
	history := logic.GenerationHistory()

	// Extract history metrics for export
	metrics := make([]render.ExportMetric, len(history))
	for i, h := range history {
		metrics[i] = render.ExportMetric{
			Generation:           h.Generation,
			AvgSelfishHerdChance: h.AvgSelfishHerdChance,
			SurvivorCount:        h.SurvivorCount,
		}
	}

	runs := []render.RunData{
		{
			Label:   "Default Run",
			History: metrics,
		},
	}

	// Read SIMULATION_SETTING.md
	settingsMarkdown := ""
	if data, err := os.ReadFile("SIMULATION_SETTING.md"); err == nil {
		settingsMarkdown = string(data)
	}

	cols, rows := logic.Dimensions()
	settings := render.SimulationSettings{
		WorldWidth:        cols,
		WorldHeight:       rows,
		InitialPreyCount:  80,
		InitialPredators:  4,
		PreySpeed:         0.15,
		PredatorSpeed:     0.45,
		PredatorRestTicks: 120,
		PreyCatchRadius:   0.5,
		GenerationCount:   logic.LastCompletedGeneration() + 1,
		TotalTicks:        logic.Tick(),
		ExportTime:        time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := render.ExportGraphHTML(filePath, runs, settings, settingsMarkdown); err != nil {
		fmt.Printf("Error exporting results: %v\n", err)
		return
	}

	fmt.Printf("Results exported to: %s\n", filePath)
}

func runMultiHeadless(targetGen int, exportPath string) {
	fmt.Printf("Starting multi-run simulation for %d target generations...\n", targetGen)
	starts := []float64{0.0, 0.2, 0.4, 0.6, 0.8, 1.0}
	var runs []render.RunData

	for _, startAvg := range starts {
		fmt.Printf("Running simulation with initial average %.2f...\n", startAvg)
		logic.SetInitialSelfishHerdChance(startAvg)
		logic.InitSimulation()
		// Don't limit the cycle so it runs fast natively
		logic.AddPredators(4)
		logic.AddAnimals(80)

		for logic.LastCompletedGeneration() < targetGen {
			logic.Step()
		}

		history := logic.GenerationHistory()
		metrics := make([]render.ExportMetric, len(history))
		for i, h := range history {
			metrics[i] = render.ExportMetric{
				Generation:           h.Generation,
				AvgSelfishHerdChance: h.AvgSelfishHerdChance,
				SurvivorCount:        h.SurvivorCount,
			}
		}

		runs = append(runs, render.RunData{
			Label:   fmt.Sprintf("Start Avg %.2f", startAvg),
			History: metrics,
		})
	}

	// Read SIMULATION_SETTING.md
	settingsMarkdown := ""
	if data, err := os.ReadFile("SIMULATION_SETTING.md"); err == nil {
		settingsMarkdown = string(data)
	}

	cols, rows := logic.Dimensions()
	settings := render.SimulationSettings{
		WorldWidth:        cols,
		WorldHeight:       rows,
		InitialPreyCount:  80,
		InitialPredators:  4,
		PreySpeed:         0.15,
		PredatorSpeed:     0.45,
		PredatorRestTicks: 120,
		PreyCatchRadius:   0.5,
		GenerationCount:   targetGen,    // They all reached this expected target
		TotalTicks:        logic.Tick(), // Not strictly accurate across all runs, but last one works
		ExportTime:        time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := render.ExportGraphHTML(exportPath, runs, settings, settingsMarkdown); err != nil {
		fmt.Printf("Error exporting results: %v\n", err)
		return
	}

	fmt.Printf("Multi-run results exported to: %s\n", exportPath)
}
