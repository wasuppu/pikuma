package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	uint32Size = int(unsafe.Sizeof(uint32(0)))
)

var (
	isRunning          bool
	window             *sdl.Window
	renderer           *sdl.Renderer
	colorBufferTexture *sdl.Texture
	colorBuffer        []uint32
	windowWidth        int32 = 800
	windowHeight       int32 = 600
)

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

func drawGrid() {
	for y := int32(0); y < windowHeight; y += 10 {
		for x := int32(0); x < windowWidth; x += 10 {
			colorBuffer[windowWidth*y+x] = 0xFF444444
		}
	}
}

func drawRect(x, y, width, height int32, color uint32) {
	for i := range width {
		for j := range height {
			currentX := x + i
			currentY := y + j
			colorBuffer[windowWidth*currentY+currentX] = color
		}
	}
}

func renderColorBuffer() {
	colorBufferTexture.Update(nil, unsafe.Pointer(&colorBuffer[0]), int(windowWidth)*uint32Size)
	renderer.Copy(colorBufferTexture, nil, nil)
}

func clearColorBuffer(color uint32) {
	for y := range windowHeight {
		for x := range windowWidth {
			colorBuffer[windowWidth*y+x] = color
		}
	}
}
