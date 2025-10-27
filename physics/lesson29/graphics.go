package main

import (
	"fmt"
	"math"
	"os"
	"physics/lesson29/physics"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Graphics struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	windowWidth  int32
	windowHeight int32
}

func (g *Graphics) OpenWindow() bool {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize SDL: %s\n", err)
		return false
	}

	// Use SDL to query what is the fullscreen max width and height
	displayMode, _ := sdl.GetCurrentDisplayMode(0)
	g.windowWidth = displayMode.W
	g.windowHeight = displayMode.H

	// Create a SDL window
	g.window, err = sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		g.windowWidth, g.windowHeight, sdl.WINDOW_BORDERLESS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create SDL window: %s\n", err)
		return false
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create SDL renderer: %v", err)
		return false
	}

	return true
}

func (g *Graphics) ClearScreen(color uint32) {
	g.renderer.SetDrawColor(uint8(color>>16), uint8(color>>8), uint8(color), 255)
	g.renderer.Clear()
}

func (g *Graphics) RenderFrame() {
	g.renderer.Present()
}

func (g *Graphics) CloseWindow() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

func (g *Graphics) DrawLine(x0, y0, x1, y1 int32, color uint32) {
	c := uint32ToColor(color)
	gfx.LineColor(g.renderer, x0, y0, x1, y1, c)
}

func (g *Graphics) DrawCircle(x, y, radius int32, angle float64, color uint32) {
	c := uint32ToColor(color)
	gfx.CircleColor(g.renderer, x, y, radius, c)
	gfx.LineColor(g.renderer, x, y, x+int32(math.Cos(angle)*float64(radius)), y+int32(math.Sin(angle)*float64(radius)), c)
}

func (g *Graphics) DrawFillCircle(x, y, radius int32, color uint32) {
	c := uint32ToColor(color)
	gfx.FilledCircleColor(g.renderer, x, y, radius, c)
}

func (g *Graphics) DrawRect(x, y, width, height int, color uint32) {
	c := uint32ToColor(color)
	gfx.LineColor(g.renderer, int32(float64(x)-float64(width)/2), int32(float64(y)-float64(height)/2), int32(float64(x)+float64(width)/2), int32(float64(y)-float64(height)/2), c)
	gfx.LineColor(g.renderer, int32(float64(x)+float64(width)/2), int32(float64(y)-float64(height)/2), int32(float64(x)+float64(width)/2), int32(float64(y)+float64(height)/2), c)
	gfx.LineColor(g.renderer, int32(float64(x)+float64(width)/2), int32(float64(y)+float64(height)/2), int32(float64(x)-float64(width)/2), int32(float64(y)+float64(height)/2), c)
	gfx.LineColor(g.renderer, int32(float64(x)-float64(width)/2), int32(float64(y)+float64(height)/2), int32(float64(x)-float64(width)/2), int32(float64(y)-float64(height)/2), c)
}

func (g *Graphics) DrawFillRect(x, y, width, height int32, color uint32) {
	c := uint32ToColor(color)
	gfx.BoxColor(g.renderer, int32(float64(x)-float64(width)/2), int32(float64(y)-float64(height)/2), int32(float64(x)+float64(width)/2), int32(float64(y)+float64(height)/2), c)
}

func (g *Graphics) DrawPolygon(x, y int32, vertices []physics.Vec2, color uint32) {
	c := uint32ToColor(color)
	for i := range len(vertices) {
		ni := (i + 1) % len(vertices)
		gfx.LineColor(g.renderer, int32(vertices[i].X), int32(vertices[i].Y), int32(vertices[ni].X), int32(vertices[ni].Y), c)
	}
	gfx.FilledCircleColor(g.renderer, x, y, 1, c)
}

func (g *Graphics) DrawFillPolygon(x, y int32, vertices []physics.Vec2, color uint32) {
	c := uint32ToColor(color)
	var vx, vy []int16
	for i := range len(vertices) {
		vx = append(vx, int16(vertices[i].X))
		vy = append(vy, int16(vertices[i].Y))
	}
	gfx.FilledPolygonColor(g.renderer, vx, vy, c)
	gfx.FilledCircleColor(g.renderer, x, y, 1, uint32ToColor(0xFF000000))
}

func (g *Graphics) DrawTexture(x, y, width, height int32, rotation float64, texture *sdl.Texture) {
	dstRect := sdl.Rect{X: int32(float64(x) - float64(width)/2), Y: int32(float64(y) - float64(height)/2), W: width, H: height}
	rotationDeg := rotation * 57.2958
	g.renderer.CopyEx(texture, nil, &dstRect, rotationDeg, nil, sdl.FLIP_NONE)
}

func uint32ToColor(rgba uint32) sdl.Color {
	a := uint8((rgba >> 24) & 0xFF)
	b := uint8((rgba >> 16) & 0xFF)
	g := uint8((rgba >> 8) & 0xFF)
	r := uint8(rgba & 0xFF)
	return sdl.Color{R: r, G: g, B: b, A: a}
}
