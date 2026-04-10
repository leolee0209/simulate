package render

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"simulate/logic"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type raylibRenderer struct {
	cfg          Config
	worldWidth   float64
	worldHeight  float64
	camera       rl.Camera2D
	backend      string
	lastGen      int
	transition   float32
	inTransition bool
	activeTab    int
	history      []generationMetric
	lastCaptured int
}

type generationMetric struct {
	generation           int
	avgSelfishHerdChance float64
	survivorCount        int
}

func NewRenderer(worldWidth float64, worldHeight float64, cfg Config) (Renderer, error) {
	cfg = normalizeConfig(cfg)
	backend := configureLinuxBackend(cfg.LinuxDisplayBackend)

	windowWidth := int32(cfg.WindowWidth)
	windowHeight := int32(cfg.WindowHeight)
	rl.InitWindow(windowWidth, windowHeight, cfg.Title)
	//rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetTargetFPS(int32(cfg.FPS))

	margins := 48.0
	fitX := (float64(cfg.WindowWidth) - margins) / worldWidth
	fitY := (float64(cfg.WindowHeight) - margins) / worldHeight
	initialZoom := float32(math.Min(fitX, fitY))
	if initialZoom <= 0 {
		initialZoom = 1
	}

	camera := rl.Camera2D{
		Offset: rl.Vector2{X: float32(cfg.WindowWidth) / 2, Y: float32(cfg.WindowHeight) / 2},
		Target: rl.Vector2{X: float32(worldWidth / 2), Y: float32(worldHeight / 2)},
		Zoom:   initialZoom,
	}

	return &raylibRenderer{
		cfg:          cfg,
		worldWidth:   worldWidth,
		worldHeight:  worldHeight,
		camera:       camera,
		backend:      backend,
		lastGen:      -1,
		activeTab:    0,
		lastCaptured: -1,
	}, nil
}

func (r *raylibRenderer) IsRunning() bool {
	return !rl.WindowShouldClose()
}

func (r *raylibRenderer) Draw(creatures []logic.CreatureSnapshot) {
	r.captureGenerationMetrics()
	r.handleTabInput()

	if r.activeTab == 0 {
		r.updateGenerationTransition()
		r.syncCameraOffset()
		r.handleCameraInput()
	}

	rl.BeginDrawing()
	if r.activeTab == 0 && r.inTransition {
		rl.ClearBackground(rl.White)
	} else {
		rl.ClearBackground(rl.Color{R: 245, G: 248, B: 252, A: 255})
	}

	if r.activeTab == 0 {
		rl.BeginMode2D(r.camera)
		rl.DrawRectangle(0, 0, int32(r.worldWidth), int32(r.worldHeight), rl.Color{R: 252, G: 253, B: 255, A: 255})

		radiusScale := r.creatureRadiusScale()
		for _, c := range creatures {
			color := creatureRenderColor(c)
			radius := float32(0.5) * radiusScale
			if radius > 0 {
				rl.DrawCircleV(rl.Vector2{X: float32(c.Pos.X), Y: float32(c.Pos.Y)}, radius, color)
			}
		}
		rl.EndMode2D()
		r.drawWorldBorder()

		rl.DrawRectangleRounded(rl.Rectangle{X: 10, Y: 56, Width: 330, Height: 42}, 0.3, 8, rl.Fade(rl.White, 0.92))
		rl.DrawText(fmt.Sprintf("World %.0f x %.0f", r.worldWidth, r.worldHeight), 22, 67, 18, rl.Color{R: 45, G: 52, B: 64, A: 255})
		rl.DrawText(fmt.Sprintf("Gen %d", logic.Generation()+1), 248, 67, 18, rl.Color{R: 45, G: 52, B: 64, A: 255})
	} else {
		r.drawGraph()
	}

	r.drawTabBar()

	rl.EndDrawing()
}

func (r *raylibRenderer) ShowGraphTab() {
	r.activeTab = 1
}

func (r *raylibRenderer) tabHeight() int32 {
	return 44
}

func (r *raylibRenderer) handleTabInput() {
	if rl.IsKeyPressed(rl.KeyTab) {
		r.activeTab = (r.activeTab + 1) % 2
	}

	if !rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		return
	}

	mouse := rl.GetMousePosition()
	if mouse.Y < 0 || mouse.Y > float32(r.tabHeight()) {
		return
	}

	if rl.CheckCollisionPointRec(mouse, r.simulationTabRect()) {
		r.activeTab = 0
		return
	}
	if rl.CheckCollisionPointRec(mouse, r.graphTabRect()) {
		r.activeTab = 1
	}
}

func (r *raylibRenderer) simulationTabRect() rl.Rectangle {
	return rl.Rectangle{X: 12, Y: 6, Width: 134, Height: 32}
}

func (r *raylibRenderer) graphTabRect() rl.Rectangle {
	return rl.Rectangle{X: 156, Y: 6, Width: 134, Height: 32}
}

func (r *raylibRenderer) drawTabBar() {
	tabH := r.tabHeight()
	screenWidth := int32(rl.GetScreenWidth())
	barColor := rl.Color{R: 233, G: 239, B: 247, A: 255}
	rl.DrawRectangle(0, 0, screenWidth, tabH, barColor)

	simRect := r.simulationTabRect()
	graphRect := r.graphTabRect()
	r.drawSingleTab(int32(simRect.X), int32(simRect.Y), int32(simRect.Width), int32(simRect.Height), "Simulation", r.activeTab == 0)
	r.drawSingleTab(int32(graphRect.X), int32(graphRect.Y), int32(graphRect.Width), int32(graphRect.Height), "Graph", r.activeTab == 1)
}

func (r *raylibRenderer) drawSingleTab(x, y, w, h int32, label string, active bool) {
	bg := rl.Color{R: 215, G: 224, B: 238, A: 255}
	fg := rl.Color{R: 62, G: 74, B: 92, A: 255}
	if active {
		bg = rl.Color{R: 82, G: 116, B: 196, A: 255}
		fg = rl.White
	}
	rl.DrawRectangleRounded(rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}, 0.35, 8, bg)
	textW := rl.MeasureText(label, 18)
	textX := x + (w-int32(textW))/2
	textY := y + (h-18)/2
	rl.DrawText(label, textX, textY, 18, fg)
}

func (r *raylibRenderer) updateGenerationTransition() {
	currentGen := logic.Generation()
	if r.lastGen == -1 {
		r.lastGen = currentGen
		return
	}

	if currentGen != r.lastGen {
		r.lastGen = currentGen
		r.inTransition = true
		r.transition = 0
	}

	if !r.inTransition {
		return
	}

	r.transition += rl.GetFrameTime()
	if r.transition >= 1 {
		r.transition = 1
		r.inTransition = false
	}
}

func (r *raylibRenderer) creatureRadiusScale() float32 {
	if !r.inTransition {
		return 1
	}

	const blankPhase = float32(0.18)
	if r.transition <= blankPhase {
		return 0
	}

	normalized := (r.transition - blankPhase) / (1 - blankPhase)
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// ease-out for smoother emergence
	return 1 - float32(math.Pow(float64(1-normalized), 2))
}

func (r *raylibRenderer) syncCameraOffset() {
	r.camera.Offset = rl.Vector2{
		X: float32(rl.GetScreenWidth()) / 2,
		Y: float32(rl.GetScreenHeight()) / 2,
	}
}

func (r *raylibRenderer) drawWorldBorder() {
	topLeft := rl.GetWorldToScreen2D(rl.Vector2{X: 0, Y: 0}, r.camera)
	topRight := rl.GetWorldToScreen2D(rl.Vector2{X: float32(r.worldWidth), Y: 0}, r.camera)
	bottomRight := rl.GetWorldToScreen2D(rl.Vector2{X: float32(r.worldWidth), Y: float32(r.worldHeight)}, r.camera)
	bottomLeft := rl.GetWorldToScreen2D(rl.Vector2{X: 0, Y: float32(r.worldHeight)}, r.camera)

	borderColor := rl.Color{R: 72, G: 84, B: 104, A: 255}
	thickness := float32(2.5)
	rl.DrawLineEx(topLeft, topRight, thickness, borderColor)
	rl.DrawLineEx(topRight, bottomRight, thickness, borderColor)
	rl.DrawLineEx(bottomRight, bottomLeft, thickness, borderColor)
	rl.DrawLineEx(bottomLeft, topLeft, thickness, borderColor)
}

func (r *raylibRenderer) handleCameraInput() {
	if r.activeTab != 0 {
		return
	}

	if rl.GetMousePosition().Y <= float32(r.tabHeight()) {
		return
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		delta := rl.GetMouseDelta()
		r.camera.Target.X -= delta.X / r.camera.Zoom
		r.camera.Target.Y -= delta.Y / r.camera.Zoom
	}

	wheel := rl.GetMouseWheelMove()
	if wheel == 0 {
		return
	}

	mouse := rl.GetMousePosition()
	before := rl.GetScreenToWorld2D(mouse, r.camera)

	newZoom := r.camera.Zoom * (1 + wheel*0.1)
	if newZoom < 0.05 {
		newZoom = 0.05
	}
	if newZoom > 40 {
		newZoom = 40
	}
	r.camera.Zoom = newZoom

	after := rl.GetScreenToWorld2D(mouse, r.camera)
	r.camera.Target.X += before.X - after.X
	r.camera.Target.Y += before.Y - after.Y
}

func (r *raylibRenderer) Close() {
	rl.CloseWindow()
}

func creatureRenderColor(c logic.CreatureSnapshot) rl.Color {
	base := rl.Color{R: c.BaseColor.R, G: c.BaseColor.G, B: c.BaseColor.B, A: 255}

	switch c.Kind {
	case logic.AnimalKind:
		if c.Mode == logic.EvadingMode {
			return blendColor(base, rl.Color{R: 70, G: 220, B: 90, A: 255}, 0.5)
		}
		return base
	case logic.PredatorKind:
		return base
	default:
		return base
	}
}

func blendColor(a, b rl.Color, ratio float64) rl.Color {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	blend := func(x, y uint8) uint8 {
		return uint8(float64(x)*(1-ratio) + float64(y)*ratio)
	}

	return rl.Color{
		R: blend(a.R, b.R),
		G: blend(a.G, b.G),
		B: blend(a.B, b.B),
		A: 255,
	}
}

func configureLinuxBackend(requested string) string {
	if runtime.GOOS != "linux" {
		return ""
	}

	backend := strings.ToLower(strings.TrimSpace(requested))
	if backend == "" {
		backend = "auto"
	}

	if backend == "auto" {
		sessionType := strings.ToLower(strings.TrimSpace(os.Getenv("XDG_SESSION_TYPE")))
		hasWayland := os.Getenv("WAYLAND_DISPLAY") != "" || sessionType == "wayland"
		if hasWayland {
			backend = "wayland"
		} else {
			backend = "x11"
		}
	}

	if backend != "wayland" && backend != "x11" {
		backend = "x11"
	}

	if os.Getenv("GLFW_PLATFORM") == "" || requested != "" {
		_ = os.Setenv("GLFW_PLATFORM", backend)
	}

	return backend
}
