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
	isRunning       bool
	cubePoints      [NPOINTS]Vec3 // 9x9x9 cube
	projectedPoints [NPOINTS]Vec2
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
		// Save the projected 2D vector in the array of projected points
		projectedPoints[i] = cubePoints[i].project()
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
