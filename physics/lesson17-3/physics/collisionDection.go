package physics

func IsColliding(a *Body, b *Body) (bool, Contact) {
	aIsCircle := a.Shape.GetType() == CIRCLE_SHAPE
	bIsCircle := b.Shape.GetType() == CIRCLE_SHAPE

	if aIsCircle && bIsCircle {
		return IsCollidingCircleCircle(a, b)
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
