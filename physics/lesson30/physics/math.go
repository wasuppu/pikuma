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

func (v Vec2) Normal() Vec2 {
	return Vec2{v.Y, -v.X}.Normalize()
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

type Vec []float64

func (v Vec) Add(o Vec) Vec {
	u := make(Vec, len(v))
	for i := range v {
		u[i] = v[i] + o[i]
	}
	return u
}

func (v Vec) Sub(o Vec) Vec {
	u := make(Vec, len(v))
	for i := range v {
		u[i] = v[i] - o[i]
	}
	return u
}

func (v Vec) Muln(n float64) Vec {
	u := make(Vec, len(v))
	for i := range v {
		u[i] = v[i] * n
	}
	return u
}

func (v Vec) Dot(o Vec) float64 {
	s := 0.0
	for i := range v {
		s += float64(v[i] * o[i])
	}
	return s
}

type Mat []Vec

func NewMat(rows, cols int) Mat {
	m := make(Mat, rows)
	for i := range m {
		m[i] = make(Vec, cols)
	}
	return m
}

func (m Mat) Nrows() int {
	return len(m)
}

func (m Mat) Transpose() Mat {
	r, c := m.Nrows(), m.Ncols()
	n := NewMat(c, r)
	for i := range r {
		for j := range c {
			n[j][i] = m[i][j]
		}
	}
	return n
}

func (m Mat) Ncols() int {
	if len(m) == 0 {
		return 0
	} else {
		return len(m[0])
	}
}

func (m Mat) Mulv(v Vec) Vec {
	c := m.Ncols()
	if c != len(v) {
		return v
	}

	r := m.Nrows()
	u := make(Vec, r)
	for i := range r {
		u[i] = v.Dot(m[i])
	}
	return u
}

func (m Mat) Mul(n Mat) Mat {
	mr, mc := m.Nrows(), m.Ncols()
	nr, nc := n.Nrows(), n.Ncols()
	if mc != nr {
		return n
	}
	transposed := n.Transpose()
	a := NewMat(mr, nc)
	for i := range mr {
		for j := range nc {
			a[i][j] = m[i].Dot(transposed[j])
		}
	}
	return a
}

func (m Mat) SolveGaussSeidel(v Vec) Vec {
	n := len(v)
	u := make(Vec, n)

	for range n {
		for i := range n {
			dx := (v[i] / m[i][i]) - (m[i].Dot(u) / m[i][i])
			if dx == dx {
				u[i] += dx
			}
		}
	}
	return u
}
