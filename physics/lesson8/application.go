package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Application struct {
	running bool
	graphic Graphics
}

// Setup function (executed once in the beginning of the simulation)
func (app *Application) Setup() {
	app.running = app.graphic.OpenWindow()
}

// Input processing
func (app *Application) Input() {
	event := sdl.PollEvent()
	if event != nil {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			app.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				app.running = false
			}
		}
	}
}

// Update function (called several times per second to update objects)
func (app *Application) Update() {}

// Render function (called several times per second to draw objects)
func (app *Application) Render() {
	app.graphic.ClearScreen(0xFF056264)
	app.graphic.DrawFillCircle(200, 200, 40, 0xFFFFFFFF)
	app.graphic.RenderFrame()
}

// Destroy function to delete objects and close the window
func (app *Application) Destory() {
	app.graphic.CloseWindow()
}

func (app Application) IsRunning() bool {
	return app.running
}
