package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	NPOINTS = 9 * 9 * 9
)

var (
	isRunning          bool
	window             *sdl.Window
	renderer           *sdl.Renderer
	colorBufferTexture *sdl.Texture
	colorBuffer        []uint32
	windowWidth        int32         = 800
	windowHeight       int32         = 600
	cubePoints         [NPOINTS]Vec3 // 9x9x9 cube
	projectedPoints    [NPOINTS]Vec2
	cameraPosition     = Vec3{0, 0, -5}
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

	pointCount := 0
	// Start loading an array of Vectors
	// From -1 to 1 (in this 9x9x9 cube)
	for x := -1.0; x <= 1; x += 0.25 {
		for y := -1.0; y <= 1; y += 0.25 {
			for z := -1.0; z <= 1; z += 0.25 {
				cubePoints[pointCount] = Vec3{x, y, z}
				pointCount++
			}
		}
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
	for i := range NPOINTS {
		point := cubePoints[i]
		point.z -= cameraPosition.z
		// Save the projected 2D Vector in the array of projected points
		projectedPoints[i] = point.project()
	}
}

func render() {
	drawGrid()
	for i := range NPOINTS {
		drawRect(
			int32(projectedPoints[i].x)+(windowWidth/2),
			int32(projectedPoints[i].y)+(windowHeight/2),
			4,
			4,
			0xFFFFFF00,
		)
	}

	renderColorBuffer()
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
