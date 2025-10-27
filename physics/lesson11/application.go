package main

import (
	"physics/lesson11/physics"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS                 = 60
	MILLISECS_PER_FRAME = 1000 / FPS
	PIXELS_PER_METER    = 50
)

type Application struct {
	running           bool
	graphic           Graphics
	particle          physics.Particle
	timePreviousFrame uint64
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() {
	app.running = app.graphic.OpenWindow()
	app.particle = physics.NewParticle(50, 100, 1.0)
	app.particle.Radius = 4
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
	// Wait some time until the reach the target frame time in milliseconds
	timeToWait := int(MILLISECS_PER_FRAME - (sdl.GetTicks64() - app.timePreviousFrame))

	// Only call delay if we are too fast to process this frame
	if timeToWait > 0 {
		sdl.Delay(uint32(timeToWait))
	}

	// Calculate the deltatime in seconds
	deltaTime := min(float64(sdl.GetTicks64()-app.timePreviousFrame)/1000.0, 0.016)

	// Set the time of the current frame to be used in the next one
	app.timePreviousFrame = sdl.GetTicks64()

	// Apply a "wind" force to my particle
	wind := physics.Vec2{X: 0.2 * PIXELS_PER_METER, Y: 0.0}
	app.particle.AddFore(wind)

	// Integrate the acceleration and velocity to estimate the new position
	app.particle.Integrate(deltaTime)

	// Nasty hardcoded flip in velocity if it touches the limits of the screen window
	if app.particle.Position.X-app.particle.Radius <= 0 {
		app.particle.Position.X = app.particle.Radius
		app.particle.Velocity.X *= -0.9
	} else if app.particle.Position.X+app.particle.Radius >= float64(app.graphic.windowWidth) {
		app.particle.Position.X = float64(app.graphic.windowWidth) - app.particle.Radius
		app.particle.Velocity.X *= -0.9
	}

	if app.particle.Position.Y-app.particle.Radius <= 0 {
		app.particle.Position.Y = app.particle.Radius
		app.particle.Velocity.Y *= -0.9
	} else if app.particle.Position.Y+app.particle.Radius >= float64(app.graphic.windowHeight) {
		app.particle.Position.Y = float64(app.graphic.windowHeight) - app.particle.Radius
		app.particle.Velocity.Y *= -0.9
	}
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	app.graphic.ClearScreen(0xFF056264)
	app.graphic.DrawFillCircle(int32(app.particle.Position.X), int32(app.particle.Position.Y), int32(app.particle.Radius), 0xFFFFFFFF)
	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
