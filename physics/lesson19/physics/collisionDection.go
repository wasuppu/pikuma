package physics

import "math"

func IsColliding(a *Body, b *Body) (bool, Contact) {
	aIsCircle := a.Shape.GetType() == CIRCLE_SHAPE
	bIsCircle := b.Shape.GetType() == CIRCLE_SHAPE
	aIsPolygon := a.Shape.GetType() == BOX_SHAPE
	bIsPolygon := b.Shape.GetType() == BOX_SHAPE

	if aIsCircle && bIsCircle {
		return IsCollidingCircleCircle(a, b)
	}

	if aIsPolygon && bIsPolygon {
		return IsCollidingPolygonPolygon(a, b)
	}

	return false, Contact{}
}

func IsCollidingCircleCircle(a *Body, b *Body) (bool, Contact) {
	aCircleShape := a.Shape.(*CircleShape)
	bCircleShape := b.Shape.(*CircleShape)

	ab := b.Position.Sub(a.Position)
	radiusSum := aCircleShape.Radius + bCircleShape.Radius

	isColliding := ab.Dot(ab) <= radiusSum*radiusSum

	if !isColliding {
		return false, Contact{}
	}

	contact := Contact{A: a, B: b, Normal: ab.Normalize()}

	contact.Start = b.Position.Sub(contact.Normal.Muln(bCircleShape.Radius))
	contact.End = a.Position.Add(contact.Normal.Muln(aCircleShape.Radius))
	contact.Depth = contact.End.Sub(contact.Start).Length()

	return true, contact
}

func IsCollidingPolygonPolygon(a *Body, b *Body) (bool, Contact) {
	aPolygonShape := a.Shape.(PolygonShape)
	bPolygonShape := b.Shape.(PolygonShape)
	if FindMinSeparation(aPolygonShape, bPolygonShape) >= 0 {
		return false, Contact{}
	}
	if FindMinSeparation(bPolygonShape, aPolygonShape) >= 0 {
		return false, Contact{}
	}
	return true, Contact{}
}

func FindMinSeparation(a PolygonShape, b PolygonShape) float64 {
	separation := -math.MaxFloat64

	// Loop all the vertices of current polygon
	for i := range len(a.GetWorldVertices()) {
		va := a.GetWorldVertice(i)
		normal := a.EdgeAt(i).Normal()

		minSep := math.MaxFloat64

		// Loop all the vertices of other polygon
		for j := range len(b.GetWorldVertices()) {
			vb := b.GetWorldVertice(j)
			minSep = math.Min(minSep, vb.Sub(va).Dot(normal))
		}
		separation = math.Max(separation, minSep)
	}

	return separation
}
