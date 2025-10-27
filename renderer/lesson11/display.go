package main

import (
	"math"
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

func drawLine(x0, y0, x1, y1 int32, color uint32) {
	deltaX := x1 - x0
	deltaY := y1 - y0

	longestSideLength := int32(math.Abs(float64(deltaY)))
	if math.Abs(float64(deltaX)) >= math.Abs(float64(deltaY)) {
		longestSideLength = int32(math.Abs(float64(deltaX)))
	}

	xinc := float64(deltaX) / float64(longestSideLength)
	yinc := float64(deltaY) / float64(longestSideLength)

	currentX := float64(x0)
	currentY := float64(y0)
	for i := 0; i <= int(longestSideLength); i++ {
		drawPixel(int32(math.Round(currentX)), int32(math.Round(currentY)), color)
		currentX += xinc
		currentY += yinc
	}
}

// Draw a triangle using three raw line calls
func drawTriangle(x0, y0, x1, y1, x2, y2 int32, color uint32) {
	drawLine(x0, y0, x1, y1, color)
	drawLine(x1, y1, x2, y2, color)
	drawLine(x2, y2, x0, y0, color)
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
