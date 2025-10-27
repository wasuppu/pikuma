package main

import (
	"path/filepath"
	"physics/lesson28-3/physics"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS                 = 60
	MILLISECS_PER_FRAME = 1000 / FPS
	NUM_BODIES          = 8
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
	world             *physics.World
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() error {
	app.running = app.graphic.OpenWindow()

	// Create a physics world with gravity of -9.8 m/s2
	app.world = physics.NewWorld(-9.8)

	// Add several bodies
	for i := range NUM_BODIES {
		mass := 1.0
		if i == 0 {
			mass = 0.0
		}
		body := physics.NewBody(physics.NewBoxShape(30, 30), float64(app.graphic.windowWidth)/2-(float64(i)*40), 100, mass)
		body.SetTexture(filepath.Join(parentpath, "assets/crate.png"), app.graphic.renderer)
		app.world.AddBody(body)
	}

	// Add joints to connect them (distance constraints)
	for i := range NUM_BODIES - 1 {
		a := app.world.GetBodies()[i]
		b := app.world.GetBodies()[i+1]
		joint := physics.NewJointConstraint(a, b, a.Position)
		app.world.AddConstraint(joint)
	}

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
					ball.Restitution = 0.7
					app.world.AddBody(ball)
				}
				if t.Button == sdl.BUTTON_RIGHT {
					x, y, _ := sdl.GetMouseState()
					box := physics.NewBody(physics.NewBoxShape(60, 60), float64(x), float64(y), 1.0)
					box.SetTexture(filepath.Join(parentpath, "assets/crate.png"), app.graphic.renderer)
					box.Restitution = 0.2
					app.world.AddBody(box)
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

	// Update world bodies (integration, collision detection, etc.)
	app.world.Update(deltaTime)
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	// Draw a line between joint objects
	for _, joint := range app.world.GetConstraints() {
		pa := joint.A().LocalSpaceToWorldSpace(joint.APoint())
		pb := joint.B().LocalSpaceToWorldSpace(joint.APoint())
		app.graphic.DrawLine(int32(pa.X), int32(pa.Y), int32(pb.X), int32(pb.Y), 0xFF555555)
	}

	// Draw all bodies
	for _, body := range app.world.GetBodies() {
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
	for _, body := range app.world.GetBodies() {
		body.Texture.Destroy()
	}
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
