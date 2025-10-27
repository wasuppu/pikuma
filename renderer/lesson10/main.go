package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS               = 60
	FRAME_TARGET_TIME = 1000 / FPS
)

var (
	trianglesToRender [N_MESH_FACES]Triangle
	cameraPosition    = Vec3{0, 0, -5}
	cubeRotation      = Vec3{0, 0, 0}
)

var (
	isRunning          bool
	window             *sdl.Window
	renderer           *sdl.Renderer
	colorBufferTexture *sdl.Texture
	colorBuffer        []uint32
	windowWidth        int32  = 800
	windowHeight       int32  = 600
	previousFrameTime  uint64 = 0
	TicksPassed               = func(a, b uint64) bool { return int32(b-a) <= 0 }
)

func initial() bool {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize SDL: %s\n", err)
		return false
	}

	// Use SDL to query what is the fullscreen max width and height
	displayMode, _ := sdl.GetCurrentDisplayMode(0)
	windowWidth = displayMode.W
	windowHeight = displayMode.H

	// Create a SDL window
	window, err = sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		windowWidth, windowHeight, sdl.WINDOW_BORDERLESS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create SDL window: %s\n", err)
		return false
	}

	// Create a SDL renderer
	renderer, err = sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create renderer: %s\n", err)
		return false
	}

	// window.SetFullscreen(sdl.WINDOW_FULLSCREEN)

	return true
}

func setup() (err error) {
	// Allocate the required memory in bytes to hold the color buffer
	colorBuffer = make([]uint32, windowWidth*windowHeight)

	// Creating a SDL texture that is used to display the color buffer
	colorBufferTexture, err = renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		windowWidth,
		windowHeight)

	if err != nil {
		return fmt.Errorf("failed to creat texture: %s", err)
	}

	return nil
}

func processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			isRunning = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				isRunning = false
			}
		}
	}
}

func update() {
	// Wait some time until the reach the target frame time in milliseconds
	timeToWait := FRAME_TARGET_TIME - (sdl.GetTicks64() - previousFrameTime)

	// Only delay execution if we are running too fast
	if timeToWait > 0 && timeToWait <= FRAME_TARGET_TIME {
		sdl.Delay(uint32(timeToWait))
	}

	previousFrameTime = sdl.GetTicks64()

	cubeRotation = cubeRotation.addn(0.01)

	// Loop all triangle faces of our mesh
	for i := range N_MESH_FACES {
		meshFace := meshFaces[i]

		var faceVertices [3]Vec3
		faceVertices[0] = meshVertices[meshFace.a-1]
		faceVertices[1] = meshVertices[meshFace.b-1]
		faceVertices[2] = meshVertices[meshFace.c-1]

		var projectedTriangle Triangle

		// Loop all three vertices of this current face and apply transformations
		for j := range 3 {
			transformedVertex := faceVertices[j]

			transformedVertex = transformedVertex.rotateX(cubeRotation.x)
			transformedVertex = transformedVertex.rotateY(cubeRotation.y)
			transformedVertex = transformedVertex.rotateZ(cubeRotation.z)

			// Translate the vertex away from the camera
			transformedVertex.z -= cameraPosition.z

			// Project the current vertex
			projectedPoint := transformedVertex.project()

			// Scale and translate the projected points to the middle of the screen
			projectedPoint.x += float64(windowWidth) / 2
			projectedPoint.y += float64(windowHeight) / 2

			projectedTriangle.points[j] = projectedPoint
		}

		// Save the projected triangle in the array of triangles to render
		trianglesToRender[i] = projectedTriangle
	}
}

func render() {
	drawGrid()

	// Loop all projected triangles and render them
	for i := range N_MESH_FACES {
		triangle := trianglesToRender[i]

		// Draw vertex points
		drawRect(int32(triangle.points[0].x), int32(triangle.points[0].y), 3, 3, 0xFFFFFF00)
		drawRect(int32(triangle.points[1].x), int32(triangle.points[1].y), 3, 3, 0xFFFFFF00)
		drawRect(int32(triangle.points[2].x), int32(triangle.points[2].y), 3, 3, 0xFFFFFF00)

		// Draw unfilled triangle
		drawTriangle(
			int32(triangle.points[0].x),
			int32(triangle.points[0].y),
			int32(triangle.points[1].x),
			int32(triangle.points[1].y),
			int32(triangle.points[2].x),
			int32(triangle.points[2].y),
			0xFF00FF00,
		)
	}

	// Finally draw the color buffer to the SDL window
	renderColorBuffer()

	// Clear all the arrays to get ready for the next frame
	clearColorBuffer(0xFF000000)

	renderer.Present()
}

func destory() {
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()
}

func main() {
	isRunning = initial()

	setup()

	for isRunning {
		processInput()
		update()
		render()
	}

	defer destory()
}
