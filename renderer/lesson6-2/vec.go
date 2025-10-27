package main

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

func (v Vec3) project() Vec2 {
	return Vec2{fovFactor * v.x / v.z, fovFactor * v.y / v.z}
}
