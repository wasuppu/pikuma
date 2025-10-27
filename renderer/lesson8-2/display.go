package main

import (
	"unsafe"
)

const (
	uint32Size = int(unsafe.Sizeof(uint32(0)))
)

func drawGrid() {
	for y := int32(0); y < windowHeight; y += 10 {
		for x := int32(0); x < windowWidth; x += 10 {
			colorBuffer[windowWidth*y+x] = 0xFF444444
		}
	}
}

func drawPixel(x, y int32, color uint32) {
	if x >= 0 && x < windowWidth && y >= 0 && y < windowHeight {
		colorBuffer[windowWidth*y+x] = color
	}
}

func drawRect(x, y, width, height int32, color uint32) {
	for i := range width {
		for j := range height {
			currentX := x + i
			currentY := y + j
			drawPixel(currentX, currentY, color)
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
