package physics

import (
	"math"
)

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{v.X + o.X, v.Y + o.Y}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{v.X - o.X, v.Y - o.Y}
}

func (v Vec2) Muln(t float64) Vec2 {
	return Vec2{v.X * t, v.Y * t}
}

func (v Vec2) Divn(t float64) Vec2 {
	return Vec2{v.X / t, v.Y / t}
}

func (v Vec2) Dot(o Vec2) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vec2) Cross(o Vec2) float64 {
	return v.X*o.Y - v.Y*o.X
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l != 0 {
		return v.Divn(l)
	}
	return v
}

func (v Vec2) Rotate(angle float64) Vec2 {
	sin, cos := math.Sincos(angle)
	return Vec2{v.X*cos - v.Y*sin, v.X*sin + v.Y*cos}
}

func Clamp[T int | float64](val T, min T, max T) T {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}
