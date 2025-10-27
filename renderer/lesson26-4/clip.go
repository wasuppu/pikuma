package main

import (
	"math"
)

const (
	NUM_PLANES            = 6
	MAX_NUM_POLY_VERTICES = 10
)

type FrustumPlane int

const (
	LEFT_FRUSTUM_PLANE FrustumPlane = iota
	RIGHT_FRUSTUM_PLANE
	TOP_FRUSTUM_PLANE
	BOTTOM_FRUSTUM_PLANE
	NEAR_FRUSTUM_PLANE
	FAR_FRUSTUM_PLANE
)

type Plane struct {
	point  Vec3
	normal Vec3
}

// Frustum planes are defined by a point and a normal vector
func initFrustumPlanes(fovX, fovY, znear, zfar float64) [NUM_PLANES]Plane {
	sinHalfFovX, cosHalfFovX := math.Sincos(fovX / 2)
	sinHalfFovY, cosHalfFovY := math.Sincos(fovY / 2)

	frustumPlanes := [NUM_PLANES]Plane{}
	frustumPlanes[LEFT_FRUSTUM_PLANE].point = Vec3{0, 0, 0}
	frustumPlanes[LEFT_FRUSTUM_PLANE].normal[0] = cosHalfFovX
	frustumPlanes[LEFT_FRUSTUM_PLANE].normal[1] = 0
	frustumPlanes[LEFT_FRUSTUM_PLANE].normal[2] = sinHalfFovX

	frustumPlanes[RIGHT_FRUSTUM_PLANE].point = Vec3{0, 0, 0}
	frustumPlanes[RIGHT_FRUSTUM_PLANE].normal[0] = -cosHalfFovX
	frustumPlanes[RIGHT_FRUSTUM_PLANE].normal[1] = 0
	frustumPlanes[RIGHT_FRUSTUM_PLANE].normal[2] = sinHalfFovX

	frustumPlanes[TOP_FRUSTUM_PLANE].point = Vec3{0, 0, 0}
	frustumPlanes[TOP_FRUSTUM_PLANE].normal[0] = 0
	frustumPlanes[TOP_FRUSTUM_PLANE].normal[1] = -cosHalfFovY
	frustumPlanes[TOP_FRUSTUM_PLANE].normal[2] = sinHalfFovY

	frustumPlanes[BOTTOM_FRUSTUM_PLANE].point = Vec3{0, 0, 0}
	frustumPlanes[BOTTOM_FRUSTUM_PLANE].normal[0] = 0
	frustumPlanes[BOTTOM_FRUSTUM_PLANE].normal[1] = cosHalfFovY
	frustumPlanes[BOTTOM_FRUSTUM_PLANE].normal[2] = sinHalfFovY

	frustumPlanes[NEAR_FRUSTUM_PLANE].point = Vec3{0, 0, znear}
	frustumPlanes[NEAR_FRUSTUM_PLANE].normal[0] = 0
	frustumPlanes[NEAR_FRUSTUM_PLANE].normal[1] = 0
	frustumPlanes[NEAR_FRUSTUM_PLANE].normal[2] = 1

	frustumPlanes[FAR_FRUSTUM_PLANE].point = Vec3{0, 0, zfar}
	frustumPlanes[FAR_FRUSTUM_PLANE].normal[0] = 0
	frustumPlanes[FAR_FRUSTUM_PLANE].normal[1] = 0
	frustumPlanes[FAR_FRUSTUM_PLANE].normal[2] = -1
	return frustumPlanes
}

type Polygon struct {
	vertices    [MAX_NUM_POLY_VERTICES]Vec3
	texcoords   [MAX_NUM_POLY_VERTICES]Tex2
	numVertices int
}

func createPolygonFromTriangle(v0, v1, v2 Vec3, t0, t1, t2 Tex2) Polygon {
	return Polygon{
		[MAX_NUM_POLY_VERTICES]Vec3{v0, v1, v2},
		[MAX_NUM_POLY_VERTICES]Tex2{t0, t1, t2},
		3,
	}
}

func (polygon *Polygon) clip() {
	polygon.clipAgainstPlane(LEFT_FRUSTUM_PLANE)
	polygon.clipAgainstPlane(RIGHT_FRUSTUM_PLANE)
	polygon.clipAgainstPlane(TOP_FRUSTUM_PLANE)
	polygon.clipAgainstPlane(BOTTOM_FRUSTUM_PLANE)
	polygon.clipAgainstPlane(NEAR_FRUSTUM_PLANE)
	polygon.clipAgainstPlane(FAR_FRUSTUM_PLANE)
}

func (polygon *Polygon) clipAgainstPlane(plane FrustumPlane) {
	planePoint := frustumPlanes[plane].point
	planeNormal := frustumPlanes[plane].normal

	// Declare a static array of inside vertices that will be part of the final polygon returned via parameter
	insideVertices := [MAX_NUM_POLY_VERTICES]Vec3{}
	insideTexcoord := [MAX_NUM_POLY_VERTICES]Tex2{}
	numInsideVertices := 0

	// Start the previous vertex with the last polygon vertex and texture coordinate
	var previousVertex Vec3
	var previousTexcoord Tex2
	if polygon.numVertices >= 3 {
		previousVertex = polygon.vertices[polygon.numVertices-1]
		previousTexcoord = polygon.texcoords[polygon.numVertices-1]
	}

	// Calculate the dot product of the current and previous vertex
	previousDot := previousVertex.sub(planePoint).dot(planeNormal)

	// Loop all the polygon vertices while the current is different than the last one
	for i := range polygon.numVertices {
		// Start the current vertex with the first polygon vertex and texture coordinate
		currentVertex := polygon.vertices[i]
		currentTexcoord := polygon.texcoords[i]

		currentDot := currentVertex.sub(planePoint).dot(planeNormal)

		// If we changed from inside to outside or from outside to inside
		if currentDot*previousDot < 0 {
			// Find the interpolation factor t
			t := previousDot / (previousDot - currentDot)

			// Calculate the intersection point I = Q1 + t(Q2-Q1)
			intersectionPoint := Vec3{
				lerp(previousVertex.x(), currentVertex.x(), t),
				lerp(previousVertex.y(), currentVertex.y(), t),
				lerp(previousVertex.z(), currentVertex.z(), t),
			}

			// Use the lerp formula to get the interpolated U and V texture coordinates
			interpolatedTexcoord := Tex2{
				lerp(previousTexcoord.u, currentTexcoord.u, t),
				lerp(previousTexcoord.v, currentTexcoord.v, t),
			}

			// Insert the intersection point to the list of "inside vertices"
			insideVertices[numInsideVertices] = intersectionPoint
			insideTexcoord[numInsideVertices] = interpolatedTexcoord
			numInsideVertices++
		}

		// Current vertex is inside the plane
		if currentDot > 0 {
			// Insert the current vertex to the list of "inside vertices"
			insideVertices[numInsideVertices] = currentVertex
			insideTexcoord[numInsideVertices] = currentTexcoord
			numInsideVertices++
		}

		// Move to the next vertex
		previousDot = currentDot
		previousVertex = currentVertex
		previousTexcoord = currentTexcoord
	}

	// At the end, copy the list of inside vertices into the destination polygon (out parameter)
	for i := range numInsideVertices {
		polygon.vertices[i] = insideVertices[i]
		polygon.texcoords[i] = insideTexcoord[i]
	}
	polygon.numVertices = numInsideVertices
}

func (polygon *Polygon) triangles() ([MAX_NUM_POLY_VERTICES]Triangle, int) {
	var triangles [MAX_NUM_POLY_VERTICES]Triangle
	for i := range polygon.numVertices - 2 {
		triangles[i].points[0] = polygon.vertices[0].v4()
		triangles[i].points[1] = polygon.vertices[i+1].v4()
		triangles[i].points[2] = polygon.vertices[i+2].v4()
		triangles[i].texcoords[0] = polygon.texcoords[0]
		triangles[i].texcoords[1] = polygon.texcoords[i+1]
		triangles[i].texcoords[2] = polygon.texcoords[i+2]
	}
	numTriangles := polygon.numVertices - 2
	return triangles, numTriangles
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}
