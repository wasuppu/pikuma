package main

import (
	"path/filepath"
	"physics/lesson28-4/physics"
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

	// Add ragdoll parts (rigid bodies)
	bob := physics.NewBody(physics.NewCircleShape(5), float64(app.graphic.windowWidth)/2, float64(app.graphic.windowHeight)/2-200, 0.0)
	head := physics.NewBody(physics.NewCircleShape(25), bob.Position.X, bob.Position.Y+70, 5.0)
	torso := physics.NewBody(physics.NewBoxShape(50, 100), head.Position.X, head.Position.Y+80, 3.0)

	leftArm := physics.NewBody(physics.NewBoxShape(15, 70), torso.Position.X-32, torso.Position.Y-10, 1.0)
	rightArm := physics.NewBody(physics.NewBoxShape(15, 70), torso.Position.X+32, torso.Position.Y-10, 1.0)

	leftLeg := physics.NewBody(physics.NewBoxShape(20, 90), torso.Position.X-20, torso.Position.Y+97, 1.0)
	rightLeg := physics.NewBody(physics.NewBoxShape(20, 90), torso.Position.X+20, torso.Position.Y+97, 1.0)
	bob.SetTexture(filepath.Join(parentpath, "assets/ragdoll/bob.png"), app.graphic.renderer)
	head.SetTexture(filepath.Join(parentpath, "assets/ragdoll/head.png"), app.graphic.renderer)
	torso.SetTexture(filepath.Join(parentpath, "assets/ragdoll/torso.png"), app.graphic.renderer)
	leftArm.SetTexture(filepath.Join(parentpath, "assets/ragdoll/leftArm.png"), app.graphic.renderer)
	rightArm.SetTexture(filepath.Join(parentpath, "assets/ragdoll/rightArm.png"), app.graphic.renderer)
	leftLeg.SetTexture(filepath.Join(parentpath, "assets/ragdoll/leftLeg.png"), app.graphic.renderer)
	rightLeg.SetTexture(filepath.Join(parentpath, "assets/ragdoll/rightLeg.png"), app.graphic.renderer)
	app.world.AddBody(bob)
	app.world.AddBody(head)
	app.world.AddBody(torso)
	app.world.AddBody(leftArm)
	app.world.AddBody(rightArm)
	app.world.AddBody(leftLeg)
	app.world.AddBody(rightLeg)

	// Add joints between ragdoll parts (distance constraints with one anchor point)
	str := physics.NewJointConstraint(bob, head, bob.Position)
	neck := physics.NewJointConstraint(head, torso, head.Position.Add(physics.Vec2{X: 0, Y: 25}))

	leftShoulder := physics.NewJointConstraint(torso, leftArm, torso.Position.Add(physics.Vec2{X: -28, Y: -45}))
	rightShoulder := physics.NewJointConstraint(torso, rightArm, torso.Position.Add(physics.Vec2{X: 28, Y: -45}))

	leftHip := physics.NewJointConstraint(torso, leftLeg, torso.Position.Add(physics.Vec2{X: -20, Y: 50}))
	rightHip := physics.NewJointConstraint(torso, rightLeg, torso.Position.Add(physics.Vec2{X: 20, Y: 50}))
	app.world.AddConstraint(str)
	app.world.AddConstraint(neck)
	app.world.AddConstraint(leftShoulder)
	app.world.AddConstraint(rightShoulder)
	app.world.AddConstraint(leftHip)
	app.world.AddConstraint(rightHip)

	// Add a floor and walls to contain objects
	floor := physics.NewBody(physics.NewBoxShape(float64(app.graphic.windowWidth)-50, 50), float64(app.graphic.windowWidth)/2.0, float64(app.graphic.windowHeight)-50, 0.0)
	leftWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), 50, float64(app.graphic.windowHeight)/2-25, 0.0)
	rightWall := physics.NewBody(physics.NewBoxShape(50, float64(app.graphic.windowHeight)-100), float64(app.graphic.windowWidth)-50, float64(app.graphic.windowHeight)/2-25, 0.0)
	floor.Restitution = 0.7
	leftWall.Restitution = 0.2
	rightWall.Restitution = 0.2
	app.world.AddBody(floor)
	app.world.AddBody(leftWall)
	app.world.AddBody(rightWall)

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
		case *sdl.MouseMotionEvent:
			x, y, _ := sdl.GetMouseState()
			mouse := physics.Vec2{X: float64(x), Y: float64(y)}
			bob := app.world.GetBodies()[0]
			direction := mouse.Sub(bob.Position).Normalize()
			speed := 1.0
			bob.Position = bob.Position.Add(direction.Muln(speed))
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
	// Draw a line between the bob and the ragdoll head
	bob := app.world.GetBodies()[0]
	head := app.world.GetBodies()[1]
	app.graphic.DrawLine(int32(bob.Position.X), int32(bob.Position.Y), int32(head.Position.X), int32(head.Position.Y), 0xFF555555)

	// Draw all joints anchor points
	for _, joint := range app.world.GetConstraints() {
		if app.debug {
			anchorPoint := joint.A().LocalSpaceToWorldSpace(joint.APoint())
			app.graphic.DrawFillCircle(int32(anchorPoint.X), int32(anchorPoint.Y), 3, 0xFF0000FF)
		}
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
