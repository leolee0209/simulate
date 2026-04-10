package main

import (
	"flag"
	"fmt"
	"simulate/logic"
	"simulate/render"
)

func main() {
	gen := flag.Int("gen", 0, "run with graphics until generation N completes, then stop simulation updates and keep graph visible")
	flag.Parse()

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
}
