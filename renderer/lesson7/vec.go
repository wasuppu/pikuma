package main

import "math"

const fovFactor = 640

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

func (v Vec3) project() Vec2 {
	return Vec2{fovFactor * v.x / v.z, fovFactor * v.y / v.z}
}
