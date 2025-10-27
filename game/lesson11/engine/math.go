package engine

import "math"

type Vec2 [2]float64

func (v Vec2) X() float64 {
	return v[0]
}

func (v Vec2) Y() float64 {
	return v[1]
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{v.X() - o.X(), v.Y() - o.Y()}
}

func (v Vec2) Dot(o Vec2) float64 {
	return v[0]*o[0] + v[1]*o[1]
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func Radians(angle float64) float64 {
	return angle * math.Pi / 180
}
