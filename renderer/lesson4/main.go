package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
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

}

func render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	drawGrid()
	drawRect(300, 200, 300, 150, 0xFFFF00FF)

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
