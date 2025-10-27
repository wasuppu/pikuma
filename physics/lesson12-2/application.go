package main

import (
	"physics/lesson12-2/physics"

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
	pushForce         physics.Vec2
	liquid            sdl.Rect
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

	app.liquid = sdl.Rect{X: 0, Y: app.graphic.windowHeight / 2, W: app.graphic.windowWidth, H: app.graphic.windowHeight / 2}
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
				case sdl.K_UP:
					app.pushForce.Y = -50 * PIXELS_PER_METER
				case sdl.K_RIGHT:
					app.pushForce.X = 50 * PIXELS_PER_METER
				case sdl.K_DOWN:
					app.pushForce.Y = 50 * PIXELS_PER_METER
				case sdl.K_LEFT:
					app.pushForce.X = -50 * PIXELS_PER_METER
				}
			case sdl.KEYUP:
				switch t.Keysym.Sym {
				case sdl.K_UP:
					app.pushForce.Y = 0
				case sdl.K_RIGHT:
					app.pushForce.X = 0
				case sdl.K_DOWN:
					app.pushForce.Y = 0
				case sdl.K_LEFT:
					app.pushForce.X = 0
				}
			}
		case *sdl.MouseButtonEvent:
			switch t.Type {
			case sdl.MOUSEBUTTONDOWN:
				if t.Button == sdl.BUTTON_LEFT {
					x, y, _ := sdl.GetMouseState()
					particle := physics.NewParticle(float64(x), float64(y), 1.0)
					particle.Radius = 5
					app.particles = append(app.particles, particle)
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

	// Apply forces to the particles
	for _, particle := range app.particles {
		// Apply a "weight" force to the particles
		weight := physics.Vec2{X: 0.0, Y: particle.Mass * 9.8 * PIXELS_PER_METER}
		particle.AddFore(weight)

		// Apply a "push" force to the particles
		particle.AddFore(app.pushForce)

		// Apply a drag force if we are inside the liquid...
		if particle.Position.Y >= float64(app.liquid.Y) {
			drag := particle.GenerateDragForce(0.03)
			particle.AddFore(drag)
		}
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

	// Draw the liquid in the screen
	app.graphic.DrawFillRect(app.liquid.X+app.liquid.W/2, app.liquid.Y+app.liquid.H/2, app.liquid.W, app.liquid.H, 0xFF6E3713)

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
