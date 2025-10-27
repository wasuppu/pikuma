package main

import (
	"physics/lesson11-2/physics"

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
	particles         []*physics.Particle
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() {
	app.running = app.graphic.OpenWindow()

	smallBall := physics.NewParticle(50, 100, 1.0)
	smallBall.Radius = 4
	app.particles = append(app.particles, smallBall)

	bigBall := physics.NewParticle(200, 100, 3.0)
	bigBall.Radius = 12
	app.particles = append(app.particles, bigBall)
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
	for _, particle := range app.particles {
		particle.AddFore(wind)
	}

	// Apply a "weight" force to the particles
	// weight := physics.Vec2{X: 0.0, Y: 9.8 * PIXELS_PER_METER}
	for _, particle := range app.particles {
		weight := physics.Vec2{X: 0.0, Y: particle.Mass * 9.8 * PIXELS_PER_METER}
		particle.AddFore(weight)
	}

	// Integrate the acceleration and velocity to estimate the new position
	for _, particle := range app.particles {
		particle.Integrate(deltaTime)
	}

	// Check the boundaries of the window
	for _, particle := range app.particles {
		// Nasty hardcoded flip in velocity if it touches the limits of the screen window
		if particle.Position.X-particle.Radius <= 0 {
			particle.Position.X = particle.Radius
			particle.Velocity.X *= -0.9
		} else if particle.Position.X+particle.Radius >= float64(app.graphic.windowWidth) {
			particle.Position.X = float64(app.graphic.windowWidth) - particle.Radius
			particle.Velocity.X *= -0.9
		}

		if particle.Position.Y-particle.Radius <= 0 {
			particle.Position.Y = particle.Radius
			particle.Velocity.Y *= -0.9
		} else if particle.Position.Y+particle.Radius >= float64(app.graphic.windowHeight) {
			particle.Position.Y = float64(app.graphic.windowHeight) - particle.Radius
			particle.Velocity.Y *= -0.9
		}
	}

}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	app.graphic.ClearScreen(0xFF056264)
	for _, particle := range app.particles {
		app.graphic.DrawFillCircle(int32(particle.Position.X), int32(particle.Position.Y), int32(particle.Radius), 0xFFFFFFFF)
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
