package main

import (
	"cmp"
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
	RENDER_TEXTURED
	RENDER_TEXTURED_WIRE
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

// Function to draw a solid pixel at position (x,y) using depth interpolation
func drawTrianglePixel(x, y int32, color uint32, pointA, pointB, pointC Vec4) {
	// Create three vec2 to find the interpolation
	pointP := Vec2{float64(x), float64(y)}

	// Calculate the barycentric coordinates of our point 'p' inside the triangle
	weights := barycentricWeights(pointA.v2(), pointB.v2(), pointC.v2(), pointP)

	alpha := weights.x()
	beta := weights.y()
	gamma := weights.z()

	// Interpolate the value of 1/w for the current pixel
	interpolatedReciprocalW := (1/pointA.w())*alpha + (1/pointB.w())*beta + (1/pointC.w())*gamma

	// Adjust 1/w so the pixels that are closer to the camera have smaller values
	interpolatedReciprocalW = 1.0 - interpolatedReciprocalW

	x = clamp(x, 0, windowWidth-1)
	y = clamp(y, 0, windowHeight-1)
	// Only draw the pixel if the depth value is less than the one previously stored in the z-buffer
	if interpolatedReciprocalW < zbuffer[windowWidth*y+x] {
		// Draw a pixel at position (x,y) with the color that comes from the mapped texture
		drawPixel(x, y, color)

		// Update the z-buffer value with the 1/w of this current pixel
		zbuffer[windowWidth*y+x] = interpolatedReciprocalW
	}
}

// Draw a filled triangle with the flat-top/flat-bottom method
// We split the original triangle in two, half flat-bottom and half flat-top
func drawFilledTriangle(
	x0, y0 int32, z0, w0 float64,
	x1, y1 int32, z1, w1 float64,
	x2, y2 int32, z2, w2 float64,
	color uint32) {
	// Draw filled triangles using a z-buffer.
	// You can use a similar technique to the one we used when drawing textured triangles.
	// But now instead of textured pixels, we simply need to draw them with a solid color.
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
		swap(&z0, &z1)
		swap(&w0, &w1)
	}

	if y1 > y2 {
		swap(&y1, &y2)
		swap(&x1, &x2)
		swap(&z1, &z2)
		swap(&w1, &w2)
	}

	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
		swap(&z0, &z1)
		swap(&w0, &w1)
	}

	// Create vector points and texture coords after we sort the vertices
	pointA := Vec4{float64(x0), float64(y0), z0, w0}
	pointB := Vec4{float64(x1), float64(y1), z1, w1}
	pointC := Vec4{float64(x2), float64(y2), z2, w2}

	// Render the upper part of the triangle (flat-bottom)
	invSlope1 := 0.0
	invSlope2 := 0.0
	if y1-y0 != 0 {
		invSlope1 = float64(x1-x0) / math.Abs(float64(y1-y0))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y1-y0 != 0 {
		for y := y0; y <= y1; y++ {
			xStart := int32(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int32(float64(x0) + float64(y-y0)*invSlope2)

			if xEnd < xStart {
				// swap if x_start is to the right of x_end
				swap(&xStart, &xEnd)
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				drawTrianglePixel(x, y, color, pointA, pointB, pointC)
			}
		}
	}

	// Render the bottom part of the triangle (flat-top)
	invSlope1 = 0.0
	invSlope2 = 0.0
	if y2-y1 != 0 {
		invSlope1 = float64(x2-x1) / math.Abs(float64(y2-y1))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y2-y1 != 0 {
		for y := y1; y <= y2; y++ {
			xStart := int32(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int32(float64(x0) + float64(y-y0)*invSlope2)

			if xEnd < xStart {
				// swap if x_start is to the right of x_end
				swap(&xStart, &xEnd)
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				drawTrianglePixel(x, y, color, pointA, pointB, pointC)
			}
		}
	}
}

// Return the barycentric weights alpha, beta, and gamma for point p
func barycentricWeights(a, b, c, p Vec2) Vec3 {
	// Find the vectors between the vertices ABC and point p
	ac := c.sub(a)
	ab := b.sub(a)
	ap := p.sub(a)
	pc := c.sub(p)
	pb := b.sub(p)

	// Compute the area of the full parallegram/triangle ABC using 2D cross product
	areaParallelogramABC := ac.x()*ab.y() - ac.y()*ab.x() // || AC x AB ||

	// Alpha is the area of the small parallelogram/triangle PBC divided by the area of the full parallelogram/triangle ABC
	alpha := (pc.x()*pb.y() - pc.y()*pb.x()) / areaParallelogramABC

	// Beta is the area of the small parallelogram/triangle APC divided by the area of the full parallelogram/triangle ABC
	beta := (ac.x()*ap.y() - ac.y()*ap.x()) / areaParallelogramABC

	// Weight gamma is easily found since barycentric coordinates always add up to 1.0
	gamma := 1 - alpha - beta

	weights := Vec3{alpha, beta, gamma}
	return weights
}

// Function to draw the textured pixel at position x and y using interpolation
func drawTexel(
	x, y int32, texture []uint32,
	pointA, pointB, pointC Vec4,
	auv, buv, cuv Tex2,
) {
	pointP := Vec2{float64(x), float64(y)}

	// Calculate the barycentric coordinates of our point 'p' inside the triangle
	weights := barycentricWeights(pointA.v2(), pointB.v2(), pointC.v2(), pointP)

	alpha := weights.x()
	beta := weights.y()
	gamma := weights.z()

	// Variables to store the interpolated values of U, V, and also 1/w for the current pixel
	// Perform the interpolation of all U/w and V/w values using barycentric weights and a factor of 1/w
	interpolatedU := (auv.u/pointA.w())*alpha + (buv.u/pointB.w())*beta + (cuv.u/pointC.w())*gamma
	interpolatedV := (auv.v/pointA.w())*alpha + (buv.v/pointB.w())*beta + (cuv.v/pointC.w())*gamma
	// Also interpolate the value of 1/w for the current pixel
	interpolatedReciprocalW := (1/pointA.w())*alpha + (1/pointB.w())*beta + (1/pointC.w())*gamma

	// Now we can divide back both interpolated values by 1/w
	interpolatedU /= interpolatedReciprocalW
	interpolatedV /= interpolatedReciprocalW

	// Map the UV coordinate to the full texture width and height
	texX := int(math.Abs(interpolatedU*float64(textureWidth))) % textureWidth
	texY := int(math.Abs(interpolatedV*float64(textureHeight))) % textureHeight

	// Adjust 1/w so the pixels that are closer to the camera have smaller values
	interpolatedReciprocalW = 1.0 - interpolatedReciprocalW

	x = clamp(x, 0, windowWidth-1)
	y = clamp(y, 0, windowHeight-1)
	// Only draw the pixel if the depth value is less than the one previously stored in the z-buffer
	if interpolatedReciprocalW < zbuffer[windowWidth*y+x] {
		// Draw a pixel at position (x,y) with the color that comes from the mapped texture
		drawPixel(x, y, texture[textureWidth*texY+texX])

		// Update the z-buffer value with the 1/w of this current pixel
		zbuffer[windowWidth*y+x] = interpolatedReciprocalW
	}
}

// Draw a textured triangle based on a texture array of colors.
// We split the original triangle in two, half flat-bottom and half flat-top.
func drawTexturedTriangle(
	x0, y0 int32, z0, w0, u0, v0 float64,
	x1, y1 int32, z1, w1, u1, v1 float64,
	x2, y2 int32, z2, w2, u2, v2 float64,
	texture []uint32,
) {
	// We need to sort the vertices by y-coordinate ascending (y0 < y1 < y2)
	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
		swap(&z0, &z1)
		swap(&w0, &w1)
		swap(&u0, &u1)
		swap(&v0, &v1)
	}

	if y1 > y2 {
		swap(&y1, &y2)
		swap(&x1, &x2)
		swap(&z1, &z2)
		swap(&w1, &w2)
		swap(&u1, &u2)
		swap(&v1, &v2)
	}

	if y0 > y1 {
		swap(&y0, &y1)
		swap(&x0, &x1)
		swap(&z0, &z1)
		swap(&w0, &w1)
		swap(&u0, &u1)
		swap(&v0, &v1)
	}

	// Create vector points and texture coords after we sort the vertices
	pointA := Vec4{float64(x0), float64(y0), z0, w0}
	pointB := Vec4{float64(x1), float64(y1), z1, w1}
	pointC := Vec4{float64(x2), float64(y2), z2, w2}
	auv := Tex2{u0, v0}
	buv := Tex2{u1, v1}
	cuv := Tex2{u2, v2}

	// Render the upper part of the triangle (flat-bottom)
	invSlope1 := 0.0
	invSlope2 := 0.0
	if y1-y0 != 0 {
		invSlope1 = float64(x1-x0) / math.Abs(float64(y1-y0))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y1-y0 != 0 {
		for y := y0; y <= y1; y++ {
			xStart := int32(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int32(float64(x0) + float64(y-y0)*invSlope2)

			if xEnd < xStart {
				// swap if x_start is to the right of x_end
				swap(&xStart, &xEnd)
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				drawTexel(x, y, texture, pointA, pointB, pointC, auv, buv, cuv)
			}
		}
	}

	// Render the bottom part of the triangle (flat-top)
	invSlope1 = 0.0
	invSlope2 = 0.0
	if y2-y1 != 0 {
		invSlope1 = float64(x2-x1) / math.Abs(float64(y2-y1))
	}
	if y2-y0 != 0 {
		invSlope2 = float64(x2-x0) / math.Abs(float64(y2-y0))
	}

	if y2-y1 != 0 {
		for y := y1; y <= y2; y++ {
			xStart := int32(float64(x1) + float64(y-y1)*invSlope1)
			xEnd := int32(float64(x0) + float64(y-y0)*invSlope2)

			if xEnd < xStart {
				// swap if x_start is to the right of x_end
				swap(&xStart, &xEnd)
			}

			for x := xStart; x < xEnd; x++ {
				// Draw our pixel with the color that comes from the texture
				drawTexel(x, y, texture, pointA, pointB, pointC, auv, buv, cuv)
			}
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

func clearZBuffer() {
	for y := range windowHeight {
		for x := range windowWidth {
			zbuffer[windowWidth*y+x] = 1
		}
	}
}

func swap[T any](v1 *T, v2 *T) {
	*v1, *v2 = *v2, *v1
}

func clamp[T cmp.Ordered](val T, min T, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}
