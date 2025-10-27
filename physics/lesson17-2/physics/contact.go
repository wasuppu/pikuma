package physics

type Contact struct {
	A      *Body
	B      *Body
	Start  Vec2
	End    Vec2
	Normal Vec2
	Depth  float64
}
