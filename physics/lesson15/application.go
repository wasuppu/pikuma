package main

import (
	"physics/lesson15/physics"

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
	timePreviousFrame uint64
	bodies            []*physics.Body
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() {
	app.running = app.graphic.OpenWindow()

	body := physics.NewBody(float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2, 1)
	body.Radius = 4
	app.bodies = append(app.bodies, body)
}

// Input processing
func (app *Application) Input() {
	event := sdl.PollEvent()
	if event != nil {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			app.running = false
		case *sdl.KeyboardEvent:
			switch t.Type {
			case sdl.KEYDOWN:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					app.running = false
				}
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

	// Apply forces to the bodies
	for _, body := range app.bodies {
		// Apply weight force
		weight := physics.Vec2{X: 0, Y: body.Mass * 9.8 * PIXELS_PER_METER}
		body.AddFore(weight)
	}

	// Integrate the acceleration and velocity to estimate the new position
	for _, body := range app.bodies {
		body.Integrate(deltaTime)
	}

	// Check the boundaries of the window applying a hardcoded bounce flip in velocity
	for _, body := range app.bodies {
		if body.Position.X-body.Radius <= 0 {
			body.Position.X = body.Radius
			body.Velocity.X *= -0.9
		} else if body.Position.X+body.Radius >= float64(app.graphic.windowWidth) {
			body.Position.X = float64(app.graphic.windowWidth) - body.Radius
			body.Velocity.X *= -0.9
		}

		if body.Position.Y-body.Radius <= 0 {
			body.Position.Y = body.Radius
			body.Velocity.Y *= -0.9
		} else if body.Position.Y+body.Radius >= float64(app.graphic.windowHeight) {
			body.Position.Y = float64(app.graphic.windowHeight) - body.Radius
			body.Velocity.Y *= -0.9
		}
	}
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	app.graphic.ClearScreen(0xFF0F0721)

	// Draw all bodies
	for _, body := range app.bodies {
		app.graphic.DrawFillCircle(int32(body.Position.X), int32(body.Position.Y), int32(body.Radius), 0xFFFFFFFF)
	}

	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
