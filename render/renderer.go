package render

import "simulate/logic"

type Config struct {
	Title        string
	FPS          int
	WindowHeight int
	WindowWidth  int
	// LinuxDisplayBackend supports: "auto", "wayland", "x11".
	LinuxDisplayBackend string
}

type Renderer interface {
	IsRunning() bool
	Draw(creatures []logic.CreatureSnapshot)
	Close()
}

func normalizeConfig(cfg Config) Config {
	if cfg.Title == "" {
		cfg.Title = "Simulate"
	}
	if cfg.FPS <= 0 {
		cfg.FPS = 10
	}
	if cfg.WindowHeight <= 0 {
		cfg.WindowHeight = 800
	}
	if cfg.WindowWidth <= 0 {
		cfg.WindowWidth = 1200
	}
	if cfg.LinuxDisplayBackend == "" {
		cfg.LinuxDisplayBackend = "auto"
	}
	return cfg
}
