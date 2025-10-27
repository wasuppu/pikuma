package physics

func IsColliding(a *Body, b *Body) (bool, Contact) {
	aIsCircle := a.Shape.GetType() == CIRCLE_SHAPE
	bIsCircle := b.Shape.GetType() == CIRCLE_SHAPE
	aIsPolygon := a.Shape.GetType() == POLYGON_SHAPE || a.Shape.GetType() == BOX_SHAPE
	bIsPolygon := b.Shape.GetType() == POLYGON_SHAPE || b.Shape.GetType() == BOX_SHAPE

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
	abSeparation, aAxis, aPoint := aPolygonShape.FindMinSeparation(bPolygonShape)
	if abSeparation >= 0 {
		return false, Contact{}
	}
	baSeparation, bAxis, bPoint := bPolygonShape.FindMinSeparation(aPolygonShape)
	if baSeparation >= 0 {
		return false, Contact{}
	}

	contact := Contact{A: a, B: b}
	if abSeparation > baSeparation {
		contact.Depth = -abSeparation
		contact.Normal = aAxis.Normal()
		contact.Start = aPoint
		contact.End = aPoint.Add(contact.Normal.Muln(contact.Depth))
	} else {
		contact.Depth = -baSeparation
		contact.Normal = bAxis.Normal().Muln(-1)
		contact.Start = bPoint.Sub(contact.Normal.Muln(contact.Depth))
		contact.End = bPoint
	}

	return true, contact
}
