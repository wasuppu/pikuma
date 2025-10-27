package main

import "math"

const (
	fovFactor       = 640
	N_MESH_VERTICES = 8
	N_MESH_FACES    = 6 * 2 // 6 cube faces, 2 triangles per face
)

var (
	meshVertices = [N_MESH_VERTICES]Vec3{
		{-1, -1, -1}, // 1
		{-1, 1, -1},  // 2
		{1, 1, -1},   // 3
		{1, -1, -1},  // 4
		{1, 1, 1},    // 5
		{1, -1, 1},   // 6
		{-1, 1, 1},   // 7
		{-1, -1, 1},  // 8
	}
	meshFaces = [N_MESH_FACES]Face{
		// front
		{1, 2, 3},
		{1, 3, 4},
		// right
		{4, 3, 5},
		{4, 5, 6},
		// back
		{6, 5, 7},
		{6, 7, 8},
		// left
		{8, 7, 2},
		{8, 2, 1},
		// top
		{2, 7, 5},
		{2, 5, 3},
		// bottom
		{6, 8, 1},
		{6, 1, 4},
	}
)

type Vec2 struct {
	x float64
	y float64
}

type Vec3 struct {
	x float64
	y float64
	z float64
}

func (v Vec3) addn(n float64) Vec3 {
	return Vec3{v.x + n, v.y + n, v.z + n}
}

func (v Vec3) rotateX(angle float64) Vec3 {
	sin, cos := math.Sincos(angle)
	return Vec3{v.x, v.y*cos - v.z*sin, v.y*sin + v.z*cos}
}

func (v Vec3) rotateY(angle float64) Vec3 {
	sin, cos := math.Sincos(angle)
	return Vec3{v.x*cos - v.z*sin, v.y, v.x*sin + v.z*cos}
}

func (v Vec3) rotateZ(angle float64) Vec3 {
	sin, cos := math.Sincos(angle)
	return Vec3{v.x*cos - v.y*sin, v.x*sin + v.y*cos, v.z}
}

// Function that receives a 3D Vector and returns a projected 2D point
func (v Vec3) project() Vec2 {
	return Vec2{fovFactor * v.x / v.z, fovFactor * v.y / v.z}
}

type Face struct {
	a, b, c int
}

type Triangle struct {
	points [3]Vec2
}

type Mesh struct {
	vertices []Vec3 // dynamic array of vertices
	faces    []Face // dynamic array of faces
	rotation Vec3   // rotation with x, y, and z values
}
