package main

import (
	"physics/lesson9/physics"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS                 = 60
	MILLISECS_PER_FRAME = 1000 / FPS
)

type Application struct {
	running  bool
	graphic  Graphics
	particle physics.Particle
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() {
	app.running = app.graphic.OpenWindow()
	app.particle = physics.NewParticle(50, 100, 1.0)
}

// Input processing
func (app *Application) Input() {
	event := sdl.PollEvent()
	if event != nil {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			app.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				app.running = false
			}
		}
	}
}

// Update function (called several times per second to update objects)
func (app *Application) Update() {
	app.particle.Velocity = physics.Vec2{X: 2.0, Y: 0.0}
	app.particle.Position = app.particle.Position.Add(app.particle.Velocity)
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	app.graphic.ClearScreen(0xFF056264)
	app.graphic.DrawFillCircle(int32(app.particle.Position.X), int32(app.particle.Position.Y), 4, 0xFFFFFFFF)
	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
