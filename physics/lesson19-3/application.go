package main

import (
	"physics/lesson19-3/physics"

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

	boxA := physics.NewBody(physics.NewBoxShape(200, 200), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2, 1.0)
	boxB := physics.NewBody(physics.NewBoxShape(200, 200), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2, 1.0)
	// boxA.AngularVelocity = 0.4
	// boxB.AngularVelocity = 0.1
	boxB.Rotation = 2.3
	app.bodies = append(app.bodies, boxA)
	app.bodies = append(app.bodies, boxB)
}

// Input processing
func (app *Application) Input() {
	// event := sdl.PollEvent()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
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
		case *sdl.MouseMotionEvent:
			x, y, _ := sdl.GetMouseState()
			app.bodies[0].Position.X = float64(x)
			app.bodies[0].Position.Y = float64(y)
		}
	}
}

// Update function (called several times per second to update objects)
func (app *Application) Update() {
	app.graphic.ClearScreen(0xFF0F0721)

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
	// for _, body := range app.bodies {
	// 	// Apply weight force
	// 	weight := physics.Vec2{X: 0, Y: body.Mass * 9.8 * PIXELS_PER_METER}
	// 	body.AddFore(weight)

	// 	// Apply the wind force
	// 	wind := physics.Vec2{X: 2.0 * PIXELS_PER_METER, Y: 0.0}
	// 	body.AddFore(wind)
	// }

	// Integrate the acceleration and velocity to estimate the new position
	for _, body := range app.bodies {
		body.Update(deltaTime)
	}

	// Check all the rigidbodies with the other rigidbodies for collision
	for i := range len(app.bodies) {
		for j := i + 1; j < len(app.bodies); j++ {
			a := app.bodies[i]
			b := app.bodies[j]
			a.IsColliding = false
			b.IsColliding = false
			if isColliding, contact := physics.IsColliding(a, b); isColliding {
				// Resolve the collision using the impulse method
				// contact.ResolveCollision()

				// Draw debug contact information
				app.graphic.DrawFillCircle(int32(contact.Start.X), int32(contact.Start.Y), 3, 0xFFFF00FF)
				app.graphic.DrawFillCircle(int32(contact.End.X), int32(contact.End.Y), 3, 0xFFFF00FF)
				app.graphic.DrawLine(int32(contact.Start.X), int32(contact.Start.Y), int32(contact.Start.X+contact.Normal.X*15), int32(contact.Start.Y+contact.Normal.Y*15), 0xFFFF00FF)
				a.IsColliding = true
				b.IsColliding = true
			}
		}
	}

	// Check the boundaries of the window applying a hardcoded bounce flip in velocity
	for _, body := range app.bodies {
		if body.Shape.GetType() == physics.CIRCLE_SHAPE {
			circleShape := body.Shape.(*physics.CircleShape)
			if body.Position.X-circleShape.Radius <= 0 {
				body.Position.X = circleShape.Radius
				body.Velocity.X *= -0.9
			} else if body.Position.X+circleShape.Radius >= float64(app.graphic.windowWidth) {
				body.Position.X = float64(app.graphic.windowWidth) - circleShape.Radius
				body.Velocity.X *= -0.9
			}

			if body.Position.Y-circleShape.Radius <= 0 {
				body.Position.Y = circleShape.Radius
				body.Velocity.Y *= -0.9
			} else if body.Position.Y+circleShape.Radius >= float64(app.graphic.windowHeight) {
				body.Position.Y = float64(app.graphic.windowHeight) - circleShape.Radius
				body.Velocity.Y *= -0.9
			}
		}
	}
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	// Draw all bodies
	for _, body := range app.bodies {
		var color uint32 = 0xFFFFFFFF
		if body.IsColliding {
			color = 0xFF0000FF
		}

		switch body.Shape.GetType() {
		case physics.CIRCLE_SHAPE:
			circleShape := body.Shape.(*physics.CircleShape)
			app.graphic.DrawFillCircle(int32(body.Position.X), int32(body.Position.Y), int32(circleShape.Radius), color)
		case physics.BOX_SHAPE:
			boxShape := body.Shape.(*physics.BoxShape)
			app.graphic.DrawPolygon(int32(body.Position.X), int32(body.Position.Y), boxShape.WorldVertices, color)
		}
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
