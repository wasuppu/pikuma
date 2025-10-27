package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS               = 60
	FRAME_TARGET_TIME = 1000 / FPS
)

var (
	trianglesToRender []Triangle // Array of triangles that should be rendered frame by frame
	cameraPosition    = Vec3{0, 0, 0}
	mesh              Mesh
	cullMethod        CullMethod
	renderMethod      RenderMethod
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
	basepath           string
	parentpath         string
)

func init() {
	_, exepath, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(exepath)
}

// Initial SDL
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

// Setup function to initialize variables and game objects
func setup() (err error) {
	// Initialize render mode and triangle culling method
	renderMethod = RENDER_WIRE
	cullMethod = CULL_BACKFACE

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

	// Loads the vertex and face values for the mesh data structure
	mesh, err = LoadObjFileData(filepath.Join(parentpath, "assets", "cube.obj"))
	if err != nil {
		return fmt.Errorf("failed to load model: %s", err)
	}

	return nil
}

// Poll system events and handle keyboard input
func processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			isRunning = false
		case *sdl.KeyboardEvent:
			switch t.Type {
			case sdl.KEYDOWN:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					isRunning = false
				case sdl.K_1:
					renderMethod = RENDER_WIRE_VERTEX
				case sdl.K_2:
					renderMethod = RENDER_WIRE
				case sdl.K_3:
					renderMethod = RENDER_FILL_TRIANGLE
				case sdl.K_4:
					renderMethod = RENDER_FILL_TRIANGLE_WIRE
				case sdl.K_c:
					cullMethod = CULL_BACKFACE
				case sdl.K_d:
					cullMethod = CULL_NONE
				}
			}
		}
	}
}

// Update function frame by frame with a fixed time step
func update() {
	// Wait some time until the reach the target frame time in milliseconds
	timeToWait := FRAME_TARGET_TIME - (sdl.GetTicks64() - previousFrameTime)

	// Only delay execution if we are running too fast
	if timeToWait > 0 && timeToWait <= FRAME_TARGET_TIME {
		sdl.Delay(uint32(timeToWait))
	}

	previousFrameTime = sdl.GetTicks64()

	mesh.rotation = mesh.rotation.addn(0.01)

	// Loop all triangle faces of our mesh
	for i := range len(mesh.faces) {
		meshFace := mesh.faces[i]

		var faceVertices [3]Vec3
		faceVertices[0] = mesh.vertices[meshFace.a-1]
		faceVertices[1] = mesh.vertices[meshFace.b-1]
		faceVertices[2] = mesh.vertices[meshFace.c-1]

		var transformedVertices [3]Vec3

		// Loop all three vertices of this current face and apply transformations
		for j := range 3 {
			transformedVertex := faceVertices[j]

			transformedVertex = transformedVertex.rotateX(mesh.rotation.x)
			transformedVertex = transformedVertex.rotateY(mesh.rotation.y)
			transformedVertex = transformedVertex.rotateZ(mesh.rotation.z)

			// Translate the vertex away from the camera
			transformedVertex.z += 5

			// Save transformed vertex in the array of transformed vertices
			transformedVertices[j] = transformedVertex
		}

		// Backface culling test to see if the current face should be projected
		if cullMethod == CULL_BACKFACE {
			// Check backface culling
			vectorA := transformedVertices[0] /*   A   */
			vectorB := transformedVertices[1] /*  / \  */
			vectorC := transformedVertices[2] /* C---B */

			// Get the vector subtraction of B-A and C-A
			vectorAB := vectorB.sub(vectorA)
			vectorAC := vectorC.sub(vectorA)
			vectorAB = vectorAB.normalize()
			vectorAC = vectorAC.normalize()

			// Compute the face normal (using cross product to find perpendicular)
			normal := vectorAB.cross(vectorAC)
			normal = normal.normalize()

			// Find the vector between a point in the triangle and the camera origin
			cameraRay := cameraPosition.sub(vectorA)

			// Calculate how aligned the camera ray is with the face normal (using dot product)
			dotNormalCamera := normal.dot(cameraRay)

			// Bypass the triangles that are looking away from the camera
			if dotNormalCamera < 0 {
				continue
			}
		}

		var projectedTriangle Triangle

		// Loop all three vertices of this current face and apply transformations
		for j := range 3 {

			// Project the current vertex
			projectedPoint := transformedVertices[j].project()

			// Scale and translate the projected points to the middle of the screen
			projectedPoint.x += float64(windowWidth) / 2
			projectedPoint.y += float64(windowHeight) / 2

			projectedTriangle.points[j] = projectedPoint
		}

		// Save the projected triangle in the array of triangles to render
		trianglesToRender = append(trianglesToRender, projectedTriangle)
	}
}

// Render function to draw objects on the display
func render() {
	drawGrid()

	// // Loop all projected triangles and render them
	for i := range trianglesToRender {
		triangle := trianglesToRender[i]

		// Draw filled triangle
		if renderMethod == RENDER_FILL_TRIANGLE || renderMethod == RENDER_FILL_TRIANGLE_WIRE {
			drawFilledTriangle(
				int32(triangle.points[0].x), int32(triangle.points[0].y), // vertex A
				int32(triangle.points[1].x), int32(triangle.points[1].y), // vertex B
				int32(triangle.points[2].x), int32(triangle.points[2].y), // vertex C
				0xFF555555,
			)
		}

		// Draw triangle wireframe
		if renderMethod == RENDER_WIRE || renderMethod == RENDER_WIRE_VERTEX || renderMethod == RENDER_FILL_TRIANGLE_WIRE {
			drawTriangle(
				int32(triangle.points[0].x), int32(triangle.points[0].y), // vertex A
				int32(triangle.points[1].x), int32(triangle.points[1].y), // vertex B
				int32(triangle.points[2].x), int32(triangle.points[2].y), // vertex C
				0xFFFFFFFF,
			)
		}

		// Draw triangle vertex points
		if renderMethod == RENDER_WIRE_VERTEX {
			drawRect(int32(triangle.points[0].x)-3, int32(triangle.points[0].y)-3, 6, 6, 0xFFFF0000) // vertex A
			drawRect(int32(triangle.points[1].x)-3, int32(triangle.points[1].y)-3, 6, 6, 0xFFFF0000) // vertex B
			drawRect(int32(triangle.points[2].x)-3, int32(triangle.points[2].y)-3, 6, 6, 0xFFFF0000) // vertex C
		}
	}

	// Clear the array of triangles to render every frame loop
	trianglesToRender = []Triangle{}

	// Finally draw the color buffer to the SDL window
	renderColorBuffer()

	// Clear all the arrays to get ready for the next frame
	clearColorBuffer(0xFF000000)

	renderer.Present()
}

// Clean up SDL
func destory() {
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()
}

// Main function
func main() {
	isRunning = initial()

	if err := setup(); err != nil {
		log.Fatalf("%#v", err)
	}

	for isRunning {
		processInput()
		update()
		render()
	}

	defer destory()
}
