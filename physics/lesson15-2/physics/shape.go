package physics

type ShapeType int

const (
	CIRCLE_SHAPE ShapeType = iota
	POLYGON_SHAPE
	BOX_SHAPE
)

type Shape interface {
	GetType() ShapeType
}

type CircleShape struct {
	Radius float64
}

func (s CircleShape) GetType() ShapeType { return CIRCLE_SHAPE }

type PolygonShape struct {
	Vertices []Vec2
}

func (s PolygonShape) GetType() ShapeType { return POLYGON_SHAPE }

type BoxShape struct {
	Width  float64
	Height float64
}

func (s BoxShape) GetType() ShapeType { return BOX_SHAPE }
