package render

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *raylibRenderer) drawGraph() {
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())

	panel := rl.Rectangle{X: 24, Y: 60, Width: float32(screenWidth - 48), Height: float32(screenHeight - 84)}
	rl.DrawRectangleRounded(panel, 0.02, 8, rl.Color{R: 250, G: 252, B: 255, A: 255})
	rl.DrawRectangleRoundedLines(panel, 0.02, 8, rl.Color{R: 210, G: 220, B: 234, A: 255})

	rl.DrawText("Average SelfishHerdChance by Generation", int32(panel.X)+18, int32(panel.Y)+14, 22, rl.Color{R: 41, G: 52, B: 74, A: 255})

	if len(r.history) == 0 {
		rl.DrawText("Waiting for generation data...", int32(panel.X)+18, int32(panel.Y)+48, 18, rl.Color{R: 90, G: 100, B: 120, A: 255})
		return
	}

	plot := rl.Rectangle{X: panel.X + 64, Y: panel.Y + 66, Width: panel.Width - 92, Height: panel.Height - 118}
	rl.DrawRectangleLinesEx(plot, 1, rl.Color{R: 190, G: 200, B: 220, A: 255})

	for i := 0; i <= 4; i++ {
		fraction := float32(i) / 4
		y := plot.Y + plot.Height*(1-fraction)
		rl.DrawLineEx(rl.Vector2{X: plot.X, Y: y}, rl.Vector2{X: plot.X + plot.Width, Y: y}, 1, rl.Color{R: 231, G: 236, B: 246, A: 255})
		label := fmt.Sprintf("%.2f", fraction)
		rl.DrawText(label, int32(plot.X)-48, int32(y)-8, 14, rl.Color{R: 90, G: 100, B: 120, A: 255})
	}

	minGen := r.history[0].generation
	maxGen := r.history[len(r.history)-1].generation
	if maxGen <= minGen {
		maxGen = minGen + 1
	}

	for i := 0; i <= 5; i++ {
		fraction := float32(i) / 5
		x := plot.X + plot.Width*fraction
		rl.DrawLineEx(rl.Vector2{X: x, Y: plot.Y + plot.Height}, rl.Vector2{X: x, Y: plot.Y + plot.Height + 6}, 1, rl.Color{R: 120, G: 130, B: 150, A: 255})
		genLabel := minGen + int(float64(maxGen-minGen)*float64(fraction))
		rl.DrawText(fmt.Sprintf("%d", genLabel), int32(x)-10, int32(plot.Y+plot.Height)+10, 14, rl.Color{R: 90, G: 100, B: 120, A: 255})
	}

	lineColor := rl.Color{R: 44, G: 112, B: 255, A: 255}
	for i := 1; i < len(r.history); i++ {
		start := r.history[i-1]
		end := r.history[i]

		x1 := plot.X + (float32(start.generation-minGen)/float32(maxGen-minGen))*plot.Width
		x2 := plot.X + (float32(end.generation-minGen)/float32(maxGen-minGen))*plot.Width

		y1 := plot.Y + plot.Height*(1-float32(clamp01(start.avgSelfishHerdChance)))
		y2 := plot.Y + plot.Height*(1-float32(clamp01(end.avgSelfishHerdChance)))

		rl.DrawLineEx(rl.Vector2{X: x1, Y: y1}, rl.Vector2{X: x2, Y: y2}, 2.6, lineColor)
	}

	for i := range r.history {
		entry := r.history[i]
		x := plot.X + (float32(entry.generation-minGen)/float32(maxGen-minGen))*plot.Width
		y := plot.Y + plot.Height*(1-float32(clamp01(entry.avgSelfishHerdChance)))
		rl.DrawCircleV(rl.Vector2{X: x, Y: y}, 3.3, lineColor)
	}
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
