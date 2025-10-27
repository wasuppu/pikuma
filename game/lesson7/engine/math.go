package engine

type Vec2 [2]float64

func (v Vec2) X() float64 {
	return v[0]
}

func (v Vec2) Y() float64 {
	return v[1]
}
