package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	fovFactor       = 640
	N_CUBE_VERTICES = 8
	N_CUBE_FACES    = 6 * 2 // 6 cube faces, 2 triangles per face
)

var (
	cubeVertices = [N_CUBE_VERTICES]Vec3{
		{-1, -1, -1}, // 1
		{-1, 1, -1},  // 2
		{1, 1, -1},   // 3
		{1, -1, -1},  // 4
		{1, 1, 1},    // 5
		{1, -1, 1},   // 6
		{-1, 1, 1},   // 7
		{-1, -1, 1},  // 8
	}
	cubeFaces = [N_CUBE_FACES]Face{
		// front
		{1, 2, 3, 0xFFFF0000},
		{1, 3, 4, 0xFFFF0000},
		// right
		{4, 3, 5, 0xFF00FF00},
		{4, 5, 6, 0xFF00FF00},
		// back
		{6, 5, 7, 0xFF0000FF},
		{6, 7, 8, 0xFF0000FF},
		// left
		{8, 7, 2, 0xFFFFFF00},
		{8, 2, 1, 0xFFFFFF00},
		// top
		{2, 7, 5, 0xFFFF00FF},
		{2, 5, 3, 0xFFFF00FF},
		// bottom
		{6, 8, 1, 0xFF00FFFF},
		{6, 1, 4, 0xFF00FFFF},
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

func (v Vec3) muln(n float64) Vec3 {
	return Vec3{v.x * n, v.y * n, v.z * n}
}

func (v Vec3) sub(o Vec3) Vec3 {
	return Vec3{v.x - o.x, v.y - o.y, v.z - o.z}
}

func (v Vec3) dot(o Vec3) float64 {
	return v.x*o.x + v.y*o.y + v.z*o.z
}

func (v Vec3) cross(o Vec3) Vec3 {
	return Vec3{
		v.y*o.z - v.z*o.y,
		v.z*o.x - v.x*o.z,
		v.x*o.y - v.y*o.x,
	}
}

func (v Vec3) length() float64 {
	return math.Sqrt(v.dot(v))
}

func (v Vec3) normalize() Vec3 {
	return v.muln(1 / v.length())
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
	color   uint32
}

type Triangle struct {
	points   [3]Vec2
	color    uint32
	avgDepth float64
}

type Mesh struct {
	vertices []Vec3 // dynamic array of vertices
	faces    []Face // dynamic array of faces
	rotation Vec3   // rotation with x, y, and z values
}

func LoadCubeMeshData() Mesh {
	mesh := Mesh{}
	mesh.rotation = Vec3{0, 0, 0}
	for i := range N_CUBE_VERTICES {
		mesh.vertices = append(mesh.vertices, cubeVertices[i])
	}
	for i := range N_CUBE_FACES {
		mesh.faces = append(mesh.faces, cubeFaces[i])
	}
	return mesh
}

func LoadObjFileData(filename string) (Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Mesh{}, err
	}
	defer file.Close()

	mesh := Mesh{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		args := fields[1:]
		switch key {
		case "v": // vertex
			f := parseFloats(args)
			v := Vec3{}
			v.x = f[0]
			v.y = f[1]
			v.z = f[2]

			mesh.vertices = append(mesh.vertices, v)
		case "f": // face
			var vertexIndices, textureIndices, normalIndices [3]int
			for i := range 3 {
				indices := parseInts(strings.Split(args[i], "/"))
				vertexIndices[i] = indices[0]
				textureIndices[i] = indices[1]
				normalIndices[i] = indices[2]
			}
			face := Face{a: vertexIndices[0], b: vertexIndices[1], c: vertexIndices[2]}
			mesh.faces = append(mesh.faces, face)
		}
	}

	return mesh, nil
}

func parseFloats(items []string) []float64 {
	result := make([]float64, len(items))
	for i, item := range items {
		f, _ := strconv.ParseFloat(item, 64)
		result[i] = f
	}
	return result
}

func parseInts(items []string) []int {
	result := make([]int, len(items))
	for i, item := range items {
		n, _ := strconv.Atoi(item)
		result[i] = n
	}
	return result
}
