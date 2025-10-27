package main

import (
	"math"
)

type Vec2 [2]float64

func (v Vec2) X() float64 {
	return v[0]
}

func (v Vec2) Y() float64 {
	return v[1]
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{v[0] + o[0], v[1] + o[1]}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{v[0] - o[0], v[1] - o[1]}
}

func (v Vec2) Muln(t float64) Vec2 {
	return Vec2{v[0] * t, v[1] * t}
}

func (v Vec2) Divn(t float64) Vec2 {
	return Vec2{v[0] / t, v[1] / t}
}

func (v Vec2) Dot(o Vec2) float64 {
	return v[0]*o[0] + v[1]*o[1]
}

func (v Vec2) Cross(o Vec2) float64 {
	return v[0]*o[1] - v[1]*o[0]
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec2) Normalize() Vec2 {
	return v.Divn(v.Length())
}
