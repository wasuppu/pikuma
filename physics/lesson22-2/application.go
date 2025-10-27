package main

import (
	"path/filepath"
	"physics/lesson22-2/physics"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS                 = 60
	MILLISECS_PER_FRAME = 1000 / FPS
	PIXELS_PER_METER    = 50
)

var (
	basepath   string
	parentpath string
)

func init() {
	_, exepath, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(exepath)
	parentpath = filepath.Dir(basepath)
}

type Application struct {
	debug             bool
	running           bool
	graphic           Graphics
	timePreviousFrame uint64
	bodies            []*physics.Body
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() error {
	app.running = app.graphic.OpenWindow()

	// Add a floor and walls to contain objects inside the screen
	floor := physics.NewBody(physics.NewBoxShape(float64(app.graphic.windowWidth)-50, 50), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)-50, 0.0)
	leftWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), 50, float64(app.graphic.windowHeight)/2-25, 0.0)
	rightWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), float64(app.graphic.windowWidth)-50, float64(app.graphic.windowHeight)/2-25, 0.0)
	floor.Restitution = 0.5
	leftWall.Restitution = 0.2
	rightWall.Restitution = 0.2
	app.bodies = append(app.bodies, floor)
	app.bodies = append(app.bodies, leftWall)
	app.bodies = append(app.bodies, rightWall)

	// Add a static box so other boxes can collide
	bigBox := physics.NewBody(physics.NewBoxShape(200, 200), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2, 0.0)
	err := bigBox.SetTexture(filepath.Join(parentpath, "assets/crate.png"), app.graphic.renderer)
	if err != nil {
		return err
	}
	bigBox.Restitution = 0.7
	bigBox.Rotation = 1.4
	app.bodies = append(app.bodies, bigBox)

	return nil
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
				case sdl.K_d:
					app.debug = !app.debug
				}
			}
		case *sdl.MouseButtonEvent:
			if t.Type == sdl.MOUSEBUTTONDOWN {
				if t.Button == sdl.BUTTON_LEFT {
					x, y, _ := sdl.GetMouseState()
					ball := physics.NewBody(physics.NewCircleShape(30), float64(x), float64(y), 1.0)
					ball.SetTexture(filepath.Join(parentpath, "assets/basketball.png"), app.graphic.renderer)
					ball.Restitution = 0.5
					app.bodies = append(app.bodies, ball)
				}
				if t.Button == sdl.BUTTON_RIGHT {
					x, y, _ := sdl.GetMouseState()
					box := physics.NewBody(physics.NewBoxShape(60, 60), float64(x), float64(y), 1.0)
					box.SetTexture(filepath.Join(parentpath, "assets/crate.png"), app.graphic.renderer)
					box.Restitution = 0.2
					app.bodies = append(app.bodies, box)
				}
			}
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
	for _, body := range app.bodies {
		// Apply weight force
		weight := physics.Vec2{X: 0, Y: body.Mass * 9.8 * PIXELS_PER_METER}
		body.AddFore(weight)
	}

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
				contact.ResolveCollision()

				// Draw debug contact information
				if app.debug {
					app.graphic.DrawFillCircle(int32(contact.Start.X), int32(contact.Start.Y), 3, 0xFFFF00FF)
					app.graphic.DrawFillCircle(int32(contact.End.X), int32(contact.End.Y), 3, 0xFFFF00FF)
					app.graphic.DrawLine(int32(contact.Start.X), int32(contact.Start.Y), int32(contact.Start.X+contact.Normal.X*15), int32(contact.Start.Y+contact.Normal.Y*15), 0xFFFF00FF)
					a.IsColliding = true
					b.IsColliding = true
				}
			}
		}
	}
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	// Draw all bodies
	for _, body := range app.bodies {
		switch body.Shape.GetType() {
		case physics.CIRCLE_SHAPE:
			circleShape := body.Shape.(*physics.CircleShape)
			if !app.debug && body.Texture != nil {
				app.graphic.DrawTexture(int32(body.Position.X), int32(body.Position.Y), int32(circleShape.Radius*2), int32(circleShape.Radius*2), body.Rotation, body.Texture)
			} else {
				app.graphic.DrawCircle(int32(body.Position.X), int32(body.Position.Y), int32(circleShape.Radius), body.Rotation, 0xFF00FF00)
			}
		case physics.BOX_SHAPE:
			boxShape := body.Shape.(*physics.BoxShape)
			if !app.debug && body.Texture != nil {
				app.graphic.DrawTexture(int32(body.Position.X), int32(body.Position.Y), int32(boxShape.Width), int32(boxShape.Height), body.Rotation, body.Texture)
			} else {
				app.graphic.DrawPolygon(int32(body.Position.X), int32(body.Position.Y), boxShape.WorldVertices, 0xFF00FF00)
			}
		case physics.POLYGON_SHAPE:
			polygonShape := body.Shape.(physics.PolygonShape)
			if !app.debug {
				app.graphic.DrawFillPolygon(int32(body.Position.X), int32(body.Position.Y), polygonShape.GetWorldVertices(), 0xFF444444)
			} else {
				app.graphic.DrawPolygon(int32(body.Position.X), int32(body.Position.Y), polygonShape.GetWorldVertices(), 0xFF00FF00)
			}
		}
	}

	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	for _, body := range app.bodies {
		body.Texture.Destroy()
	}
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
