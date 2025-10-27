package physics

import "math"

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

func NewCircleShape(radius float64) *CircleShape {
	return &CircleShape{Radius: radius}
}

func (s *CircleShape) GetType() ShapeType { return CIRCLE_SHAPE }

func (s *CircleShape) GetMomentOfInertia() float64 {
	// For solid circles, the moment of inertia is 1/2 * r^2
	// But this still needs to be multiplied by the rigidbody's mass
	return 0.5 * s.Radius * s.Radius
}

type PolygonShape interface {
	GetType() ShapeType
	GetMomentOfInertia() float64
	UpdateVertices(angle float64, position Vec2)
	EdgeAt(i int) Vec2
	FindMinSeparation(PolygonShape) float64
	GetWorldVertices() []Vec2
	GetWorldVertice(i int) Vec2
}

type PolygonBase struct {
	LocalVertices []Vec2
	WorldVertices []Vec2
}

func (s *PolygonBase) GetType() ShapeType { return POLYGON_SHAPE }

func (s *PolygonBase) UpdateVertices(angle float64, position Vec2) {
	// Loop all the vertices, transforming from local to world space
	for i := range s.LocalVertices {
		// First rotate, then we translate
		s.WorldVertices[i] = s.LocalVertices[i].Rotate(angle)
		s.WorldVertices[i] = s.WorldVertices[i].Add(position)
	}
}

func (s *PolygonBase) EdgeAt(i int) Vec2 {
	currVertex := i
	nextVertex := (i + 1) % len(s.WorldVertices)
	return s.WorldVertices[nextVertex].Sub(s.WorldVertices[currVertex])
}

func (s *PolygonBase) FindMinSeparation(s2 PolygonShape) float64 {
	separation := -math.MaxFloat64

	// Loop all the vertices of current polygon
	for i := range len(s.GetWorldVertices()) {
		v := s.GetWorldVertice(i)
		normal := s.EdgeAt(i).Normal()

		minSep := math.MaxFloat64

		// Loop all the vertices of other polygon
		for j := range len(s2.GetWorldVertices()) {
			v2 := s2.GetWorldVertice(j)
			minSep = math.Min(minSep, v2.Sub(v).Dot(normal))
		}
		separation = math.Max(separation, minSep)
	}

	return separation
}

func (s *PolygonBase) GetWorldVertices() []Vec2 {
	return s.WorldVertices
}

func (s *PolygonBase) GetWorldVertice(i int) Vec2 {
	return s.WorldVertices[i]
}

type BoxShape struct {
	PolygonBase
	Width  float64
	Height float64
}

func NewBoxShape(width, height float64) *BoxShape {
	box := BoxShape{Width: width, Height: height}

	// Load the vertices of the box polygon
	box.LocalVertices = append(box.LocalVertices, Vec2{-width / 2, -height / 2})
	box.LocalVertices = append(box.LocalVertices, Vec2{width / 2, -height / 2})
	box.LocalVertices = append(box.LocalVertices, Vec2{width / 2, height / 2})
	box.LocalVertices = append(box.LocalVertices, Vec2{-width / 2, height / 2})

	box.WorldVertices = append(box.WorldVertices, Vec2{-width / 2, -height / 2})
	box.WorldVertices = append(box.WorldVertices, Vec2{width / 2, -height / 2})
	box.WorldVertices = append(box.WorldVertices, Vec2{width / 2, height / 2})
	box.WorldVertices = append(box.WorldVertices, Vec2{-width / 2, height / 2})

	return &box
}

func (s *BoxShape) GetType() ShapeType { return BOX_SHAPE }

func (s *BoxShape) GetMomentOfInertia() float64 {
	// For a rectangle, the moment of inertia is 1/12 * (w^2 + h^2)
	// But this still needs to be multiplied by the rigidbody's mass
	return 0.083333 * (s.Width*s.Width + s.Height + s.Height)
}
