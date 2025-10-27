package physics

type ShapeType int

const (
	CIRCLE_SHAPE ShapeType = iota
	POLYGON_SHAPE
	BOX_SHAPE
)

type Shape interface {
	GetType() ShapeType
	GetMomentOfInertia() float64
}

type CircleShape struct {
	Radius float64
}

func (s CircleShape) GetType() ShapeType { return CIRCLE_SHAPE }
func (s CircleShape) GetMomentOfInertia() float64 {
	// For solid circles, the moment of inertia is 1/2 * r^2
	// But this still needs to be multiplied by the rigidbody's mass
	return 0.5 * s.Radius * s.Radius
}

type PolygonShape struct {
	Vertices []Vec2
}

func (s PolygonShape) GetType() ShapeType { return POLYGON_SHAPE }
func (s PolygonShape) GetMomentOfInertia() float64 {
	return 0.0
}

type BoxShape struct {
	Width  float64
	Height float64
}

func (s BoxShape) GetType() ShapeType { return BOX_SHAPE }
func (s BoxShape) GetMomentOfInertia() float64 {
	// For a rectangle, the moment of inertia is 1/12 * (w^2 + h^2)
	// But this still needs to be multiplied by the rigidbody's mass
	return 0.083333 * (s.Width*s.Width + s.Height + s.Height)
}
