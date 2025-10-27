package main

import (
	"path/filepath"
	"physics/lesson30-2/physics"
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

	// Add a floor and walls to contain objects
	floor := physics.NewBody(physics.NewBoxShape(float64(app.graphic.windowWidth)-50, 50), float64(app.graphic.windowWidth)/2.0, float64(app.graphic.windowHeight)-50, 0.0)
	leftWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), 50, float64(app.graphic.windowHeight)/2-25, 0.0)
	rightWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), float64(app.graphic.windowWidth)-50, float64(app.graphic.windowHeight)/2-25, 0.0)
	app.world.AddBody(floor)
	app.world.AddBody(leftWall)
	app.world.AddBody(rightWall)

	// Add rigid bodies to the scene
	a := physics.NewBody(physics.NewBoxShape(200, 200), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2, 0.0)
	b := physics.NewBody(physics.NewBoxShape(150, 150), 300, 0.0, 0.0)
	// b.Rotation = 0.1
	app.world.AddBody(a)
	app.world.AddBody(b)

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
					ball := physics.NewBody(physics.NewCircleShape(64), float64(x), float64(y), 1.0)
					ball.SetTexture(filepath.Join(parentpath, "assets/basketball.png"), app.graphic.renderer)
					ball.Restitution = 0.7
					app.world.AddBody(ball)
				}
				if t.Button == sdl.BUTTON_RIGHT {
					x, y, _ := sdl.GetMouseState()
					box := physics.NewBody(physics.NewBoxShape(140, 140), float64(x), float64(y), 1.0)
					box.SetTexture(filepath.Join(parentpath, "assets/crate.png"), app.graphic.renderer)
					box.Restitution = 0.2
					app.world.AddBody(box)
				}
			}
		case *sdl.MouseMotionEvent:
			x, y, _ := sdl.GetMouseState()
			box := app.world.GetBodies()[4]
			box.Position.X = float64(x)
			box.Position.Y = float64(y)
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
	app.world.Update(deltaTime, app.graphic.renderer)
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
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
