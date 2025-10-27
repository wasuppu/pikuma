package main

import (
	"path/filepath"
	"physics/lesson30-4/physics"
	"runtime"

	"github.com/veandco/go-sdl2/img"
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
	bgTexture         *sdl.Texture
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() error {
	app.running = app.graphic.OpenWindow()

	// Create a physics world with gravity of -9.8 m/s2
	app.world = physics.NewWorld(-9.8)

	// Load texture for the background image
	bgSurface, err := img.Load(filepath.Join(parentpath, "assets/angrybirds/background.png"))
	if err != nil {
		return err
	}
	defer bgSurface.Free()
	app.bgTexture, err = app.graphic.renderer.CreateTextureFromSurface(bgSurface)
	if err != nil {
		return err
	}

	// Add bird
	bird := physics.NewBody(physics.NewCircleShape(45), 100, float64(app.graphic.windowHeight)/2+220, 3.0)
	bird.SetTexture(filepath.Join(parentpath, "assets/angrybirds/bird-red.png"), app.graphic.renderer)
	app.world.AddBody(bird)

	// Add a floor and walls to contain objects
	floor := physics.NewBody(physics.NewBoxShape(float64(app.graphic.windowWidth)-50, 50), float64(app.graphic.windowWidth)/2.0, float64(app.graphic.windowHeight)/2+340, 0.0)
	leftFence := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-200), 0, float64(app.graphic.windowHeight)/2-35, 0.0)
	rightFence := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-200), float64(app.graphic.windowWidth), float64(app.graphic.windowHeight)/2-35, 0.0)
	app.world.AddBody(floor)
	app.world.AddBody(leftFence)
	app.world.AddBody(rightFence)

	// Add a stack of boxes
	for i := 1; i <= 4; i++ {
		mass := 10.0 / float64(i)
		box := physics.NewBody(physics.NewBoxShape(50, 50), 600, floor.Position.Y-float64(i)*55, mass)
		box.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-box.png"), app.graphic.renderer)
		box.Friction = 0.9
		box.Restitution = 0.1
		app.world.AddBody(box)
	}

	// Add structure with blocks
	plank1 := physics.NewBody(physics.NewBoxShape(50, 150), float64(app.graphic.windowWidth)/2+20, floor.Position.Y-100, 5.0)
	plank2 := physics.NewBody(physics.NewBoxShape(50, 150), float64(app.graphic.windowWidth)/2+180, floor.Position.Y-100, 5.0)
	plank3 := physics.NewBody(physics.NewBoxShape(250, 25), float64(app.graphic.windowWidth)/2+100, floor.Position.Y-200, 2.0)
	plank1.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-plank-solid.png"), app.graphic.renderer)
	plank2.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-plank-solid.png"), app.graphic.renderer)
	plank3.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-plank-cracked.png"), app.graphic.renderer)
	app.world.AddBody(plank1)
	app.world.AddBody(plank2)
	app.world.AddBody(plank3)

	// Add a triangle polygon
	triangleVertices := []physics.Vec2{{X: 30, Y: 30}, {X: -30, Y: 30}, {X: 0, Y: -30}}
	triangle := physics.NewBody(physics.NewPolygoShape(triangleVertices), plank3.Position.X, plank3.Position.Y-50, 0.5)
	triangle.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-triangle.png"), app.graphic.renderer)
	app.world.AddBody(triangle)

	// Add a pyramid of boxes
	numRows := 5
	for col := range numRows {
		for row := range col {
			x := (plank3.Position.X + 200) + float64(col)*50.0 - (float64(row) * 25.0)
			y := (floor.Position.Y - 50) - float64(row)*52
			mass := 5.0 / (float64(row) + 1.0)
			box := physics.NewBody(physics.NewBoxShape(50, 50), x, y, mass)
			box.Friction = 0.9
			box.Restitution = 0.0
			box.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-box.png"), app.graphic.renderer)
			app.world.AddBody(box)
		}
	}

	// Add a bridge of connected steps and joints
	numSteps := 10
	spacing := 33
	startStep := physics.NewBody(physics.NewBoxShape(80, 20), 200, 200, 0.0)
	startStep.SetTexture(filepath.Join(parentpath, "assets/angrybirds/rock-bridge-anchor.png"), app.graphic.renderer)
	app.world.AddBody(startStep)
	last := floor
	for i := 1; i <= numSteps; i++ {
		x := startStep.Position.X + 30 + (float64(i) * float64(spacing))
		y := startStep.Position.Y + 20
		mass := 3.0
		if i == numSteps {
			mass = 0.0
		}
		step := physics.NewBody(physics.NewCircleShape(15), x, y, mass)
		step.SetTexture(filepath.Join(parentpath, "assets/angrybirds/wood-bridge-step.png"), app.graphic.renderer)
		app.world.AddBody(step)
		joint := physics.NewJointConstraint(last, step, step.Position)
		app.world.AddConstraint(joint)
		last = step
	}
	endStep := physics.NewBody(physics.NewBoxShape(80, 20), last.Position.X+60, last.Position.Y-20, 0.0)
	endStep.SetTexture(filepath.Join(parentpath, "assets/angrybirds/rock-bridge-anchor.png"), app.graphic.renderer)
	app.world.AddBody(endStep)

	// Add pigs
	pig1 := physics.NewBody(physics.NewCircleShape(30), plank1.Position.X+80, floor.Position.Y-50, 3.0)
	pig2 := physics.NewBody(physics.NewCircleShape(30), plank2.Position.X+400, floor.Position.Y-50, 3.0)
	pig3 := physics.NewBody(physics.NewCircleShape(30), plank2.Position.X+460, floor.Position.Y-50, 3.0)
	pig4 := physics.NewBody(physics.NewCircleShape(30), 220, 130, 1.0)
	pig1.SetTexture(filepath.Join(parentpath, "assets/angrybirds/pig-1.png"), app.graphic.renderer)
	pig2.SetTexture(filepath.Join(parentpath, "assets/angrybirds/pig-2.png"), app.graphic.renderer)
	pig3.SetTexture(filepath.Join(parentpath, "assets/angrybirds/pig-1.png"), app.graphic.renderer)
	pig4.SetTexture(filepath.Join(parentpath, "assets/angrybirds/pig-2.png"), app.graphic.renderer)
	app.world.AddBody(pig1)
	app.world.AddBody(pig2)
	app.world.AddBody(pig3)
	app.world.AddBody(pig4)

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
				case sdl.K_UP:
					app.world.GetBodies()[0].ApplyImpulseLinear(physics.Vec2{X: 0.0, Y: -600.0})
				case sdl.K_LEFT:
					app.world.GetBodies()[0].ApplyImpulseLinear(physics.Vec2{X: -400.0, Y: 0.0})
				case sdl.K_RIGHT:
					app.world.GetBodies()[0].ApplyImpulseLinear(physics.Vec2{X: 400.0, Y: 0.0})
				}
			}
		case *sdl.MouseButtonEvent:
			if t.Type == sdl.MOUSEBUTTONDOWN {
				if t.Button == sdl.BUTTON_LEFT {
					x, y, _ := sdl.GetMouseState()
					box := physics.NewBody(physics.NewBoxShape(60, 60), float64(x), float64(y), 1.0)
					box.SetTexture(filepath.Join(parentpath, "assets/angrybirds/rock-box.png"), app.graphic.renderer)
					box.Friction = 0.9
					app.world.AddBody(box)
				}
				if t.Button == sdl.BUTTON_RIGHT {
					x, y, _ := sdl.GetMouseState()
					rock := physics.NewBody(physics.NewCircleShape(30), float64(x), float64(y), 1.0)
					rock.SetTexture(filepath.Join(parentpath, "assets/angrybirds/rock-round.png"), app.graphic.renderer)
					rock.Friction = 0.4
					app.world.AddBody(rock)
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
	app.world.Update(deltaTime, app.graphic.renderer)
}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	// Draw background texture
	app.graphic.DrawTexture(app.graphic.windowWidth/2, app.graphic.windowHeight/2, app.graphic.windowWidth, app.graphic.windowHeight, 0.0, app.bgTexture)

	// Draw all bodies
	for _, body := range app.world.GetBodies() {
		switch body.Shape.GetType() {
		case physics.CIRCLE_SHAPE:
			circleShape := body.Shape.(*physics.CircleShape)
			if !app.debug && body.Texture != nil {
				app.graphic.DrawTexture(int32(body.Position.X), int32(body.Position.Y), int32(circleShape.Radius*2), int32(circleShape.Radius*2), body.Rotation, body.Texture)
			} else if app.debug {
				app.graphic.DrawCircle(int32(body.Position.X), int32(body.Position.Y), int32(circleShape.Radius), body.Rotation, 0xFF00FF00)
			}
		case physics.BOX_SHAPE:
			boxShape := body.Shape.(*physics.BoxShape)
			if !app.debug && body.Texture != nil {
				app.graphic.DrawTexture(int32(body.Position.X), int32(body.Position.Y), int32(boxShape.Width()), int32(boxShape.Height()), body.Rotation, body.Texture)
			} else if app.debug {
				app.graphic.DrawPolygon(int32(body.Position.X), int32(body.Position.Y), boxShape.WorldVertices, 0xFF00FF00)
			}
		case physics.POLYGON_SHAPE:
			polygonShape := body.Shape.(physics.PolygonShape)
			if !app.debug && body.Texture != nil {
				app.graphic.DrawTexture(int32(body.Position.X), int32(body.Position.Y), int32(polygonShape.Width()), int32(polygonShape.Height()), body.Rotation, body.Texture)
			} else if app.debug {
				app.graphic.DrawPolygon(int32(body.Position.X), int32(body.Position.Y), polygonShape.GetWorldVertices(), 0xFF00FF00)
			}
		}
	}

	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	app.bgTexture.Destroy()
	for _, body := range app.world.GetBodies() {
		body.Texture.Destroy()
	}
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
