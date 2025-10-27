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

type Vec2 [2]float64

func (v Vec2) x() float64 {
	return v[0]
}

func (v Vec2) y() float64 {
	return v[1]
}

type Vec4 [4]float64

func (v Vec4) x() float64 {
	return v[0]
}

func (v Vec4) y() float64 {
	return v[1]
}

func (v Vec4) z() float64 {
	return v[2]
}

func (v Vec4) w() float64 {
	return v[3]
}

func (v Vec4) v3() Vec3 {
	return Vec3{v[0], v[1], v[2]}
}

func (v Vec4) dot(o Vec4) float64 {
	s := 0.0
	for i := range v {
		s += float64(v[i] * o[i])
	}
	return s
}

type Vec3 [3]float64

func (v Vec3) x() float64 {
	return v[0]
}

func (v Vec3) y() float64 {
	return v[1]
}

func (v Vec3) z() float64 {
	return v[2]
}

func (v Vec3) v4() Vec4 {
	return Vec4{v[0], v[1], v[2], 1}
}

func (v Vec3) addn(n float64) Vec3 {
	return Vec3{v.x() + n, v.y() + n, v.z() + n}
}

func (v Vec3) muln(n float64) Vec3 {
	return Vec3{v.x() * n, v.y() * n, v.z() * n}
}

func (v Vec3) sub(o Vec3) Vec3 {
	return Vec3{v.x() - o.x(), v.y() - o.y(), v.z() - o.z()}
}

func (v Vec3) dot(o Vec3) float64 {
	s := 0.0
	for i := range v {
		s += float64(v[i] * o[i])
	}
	return s
}

func (v Vec3) cross(o Vec3) Vec3 {
	return Vec3{
		v.y()*o.z() - v.z()*o.y(),
		v.z()*o.x() - v.x()*o.z(),
		v.x()*o.y() - v.y()*o.x(),
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
	return Vec3{v.x(), v.y()*cos - v.z()*sin, v.y()*sin + v.z()*cos}
}

func (v Vec3) rotateY(angle float64) Vec3 {
	sin, cos := math.Sincos(angle)
	return Vec3{v.x()*cos - v.z()*sin, v.y(), v.x()*sin + v.z()*cos}
}

func (v Vec3) rotateZ(angle float64) Vec3 {
	sin, cos := math.Sincos(angle)
	return Vec3{v.x()*cos - v.y()*sin, v.x()*sin + v.y()*cos, v.z()}
}

// Function that receives a 3D Vector and returns a projected 2D point
func (v Vec3) project() Vec2 {
	return Vec2{fovFactor * v.x() / v.z(), fovFactor * v.y() / v.z()}
}

type Mat4 [4]Vec4

func (m Mat4) mulv(v Vec4) Vec4 {
	u := Vec4{}
	for i := range 4 {
		u[i] = m[i].dot(v)
	}
	return u
}

func Identity4() Mat4 {
	m := Mat4{}
	for i := range 4 {
		for j := range 4 {
			if i == j {
				m[i][j] = 1
			} else {
				m[i][j] = 0
			}
		}
	}
	return m
}

func Scale(v Vec3) Mat4 {
	m := Identity4()
	m[0][0] = v[0]
	m[1][1] = v[1]
	m[2][2] = v[2]
	return m
}

func Translation(v Vec3) Mat4 {
	m := Identity4()
	m[0][3] = v[0]
	m[1][3] = v[1]
	m[2][3] = v[2]
	return m
}

func RotateX(angle float64) Mat4 {
	m := Identity4()
	s, c := math.Sincos(angle)
	m[1][1] = c
	m[1][2] = -s
	m[2][1] = s
	m[2][2] = c
	return m
}

func RotateY(angle float64) Mat4 {
	m := Identity4()
	s, c := math.Sincos(angle)
	m[0][0] = c
	m[0][2] = s
	m[2][0] = -s
	m[2][2] = c
	return m
}

func RotateZ(angle float64) Mat4 {
	m := Identity4()
	s, c := math.Sincos(angle)
	m[0][0] = c
	m[0][1] = -s
	m[1][0] = s
	m[1][1] = c
	return m
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
	vertices    []Vec3 // dynamic array of vertices
	faces       []Face // dynamic array of faces
	rotation    Vec3   // rotation with x, y, and z values
	scale       Vec3   // scale with x, y, and z values
	translation Vec3   // translation with x, y, and z values
}

func NewMesh() Mesh {
	return Mesh{rotation: Vec3{0, 0, 0}, scale: Vec3{1, 1, 1}, translation: Vec3{0, 0, 0}}
}

func LoadCubeMeshData() Mesh {
	mesh := NewMesh()
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
		return NewMesh(), err
	}
	defer file.Close()

	mesh := NewMesh()

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
			for i := range 3 {
				v[i] = f[i]
			}

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
