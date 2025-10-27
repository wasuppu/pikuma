package main

import (
	"math"
	"unsafe"
)

type CullMethod int

const (
	CULL_NONE CullMethod = iota
	CULL_BACKFACE
)

type RenderMethod int

const (
	RENDER_WIRE RenderMethod = iota
	RENDER_WIRE_VERTEX
	RENDER_FILL_TRIANGLE
	RENDER_FILL_TRIANGLE_WIRE
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

// Draw a filled a triangle with a flat bottom
func fillFlatBottomTriangle(x0, y0, x1, y1, x2, y2 int32, color uint32) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(x1-x0) / float64(y1-y0)
	invSlope2 := float64(x2-x0) / float64(y2-y0)

	// Start x_start and x_end from the top vertex (x0,y0)
	xStart := float64(x0)
	xEnd := float64(x0)

	// Loop all the scanlines from top to bottom
	for y := y0; y <= y2; y++ {
		drawLine(int32(xStart), y, int32(xEnd), y, color)
		xStart += invSlope1
		xEnd += invSlope2
	}
}

// Draw a filled a triangle with a flat top
func fillFlatTopTriangle(x0, y0, x1, y1, x2, y2 int32, color uint32) {
	// Find the two slopes (two triangle legs)
	invSlope1 := float64(x2-x0) / float64(y2-y0)
	invSlope2 := float64(x2-x1) / float64(y2-y1)

	// Start x_start and x_end from the bottom vertex (x2,y2)
	xStart := float64(x2)
	xEnd := float64(x2)

	// Loop all the scanlines from bottom to top
	for y := y2; y >= y0; y-- {
		drawLine(int32(xStart), y, int32(xEnd), y, color)
		xStart -= invSlope1
		xEnd -= invSlope2
	}
}

// Draw a filled triangle with the flat-top/flat-bottom method
// We split the original triangle in two, half flat-bottom and half flat-top
func drawFilledTriangle(x0, y0, x1, y1, x2, y2 int32, color uint32) {
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
	}

	if y1 > y2 {
		swap(&y1, &y2)
		swap(&x1, &x2)
	}

	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
	}

	if y1 == y2 {
		// Draw flat-bottom triangle
		fillFlatBottomTriangle(x0, y0, x1, y1, x2, y2, color)
	} else if y0 == y1 {
		// Draw flat-top triangle
		fillFlatTopTriangle(x0, y0, x1, y1, x2, y2, color)
	} else {
		// Calculate the new vertex (mx,my) using triangle similarity
		my := y1
		mx := (x2-x0)*(y1-y0)/(y2-y0) + x0
		// Draw flat-bottom triangle
		fillFlatBottomTriangle(x0, y0, x1, y1, mx, my, color)
		// Draw flat-top triangle
		fillFlatTopTriangle(x1, y1, mx, my, x2, y2, color)
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

func swap[T any](v1 *T, v2 *T) {
	*v1, *v2 = *v2, *v1
}
