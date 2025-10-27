package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS               = 60
	FRAME_TARGET_TIME = 1000 / FPS
)

var (
	projectilePos = Vec2{0.0, 0.0}
	projectileVel = Vec2{20.0, 30.0}
)

type Game struct {
	ticksLastFrame uint64
	running        bool
	window         *sdl.Window
	renderer       *sdl.Renderer
}

func (g *Game) Initialize(width, height int) error {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("failed to initialize SDL: %s", err)
	}

	g.window, err = sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height), sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return fmt.Errorf("failed to create window: %s", err)
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %s", err)
	}

	g.running = true

	return nil
}

func (g *Game) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			g.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				g.running = false
			}
		}
	}
}

func (g *Game) Update() {
	// Sleep the execution until we reach the target frame time in milliseconds
	timeToWait := FRAME_TARGET_TIME - (sdl.GetTicks64() - g.ticksLastFrame)

	// Only call delay if we are too fast to process this frame
	if timeToWait > 0 && timeToWait <= FRAME_TARGET_TIME {
		sdl.Delay(uint32(timeToWait))
	}

	// Delta time is the difference in ticks from last frame converted to seconds
	var deltaTime float64 = float64(sdl.GetTicks64()-g.ticksLastFrame) / 1000.0

	// Clamp deltaTime to a maximum value
	if deltaTime > 0.05 {
		deltaTime = 0.05
	}

	// Sets the new ticks for the current frame to be used in the next pass
	g.ticksLastFrame = sdl.GetTicks64()

	// Use deltaTime to update game objects
	projectilePos = Vec2{projectilePos.X() + projectileVel.X()*deltaTime, projectilePos.Y() + projectileVel.Y()*deltaTime}
}

func (g *Game) Render() {
	g.renderer.SetDrawColor(21, 21, 21, 255)
	g.renderer.Clear()

	projectile := sdl.Rect{X: int32(projectilePos.X()), Y: int32(projectilePos.Y()), W: 10, H: 10}

	g.renderer.SetDrawColor(255, 255, 255, 255)
	g.renderer.FillRect(&projectile)

	g.renderer.Present()
}

func (g *Game) Destory() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

func (g Game) IsRunning() bool {
	return g.running
}
