package physics

func IsColliding(a *Body, b *Body) bool {
	aIsCircle := a.Shape.GetType() == CIRCLE_SHAPE
	bIsCircle := b.Shape.GetType() == CIRCLE_SHAPE

	if aIsCircle && bIsCircle {
		return IsCollidingCircleCircle(a, b)
	}

	return false
}

func IsCollidingCircleCircle(a *Body, b *Body) bool {
	aCircleShape := a.Shape.(*CircleShape)
	bCircleShape := b.Shape.(*CircleShape)

	ab := b.Position.Sub(a.Position)
	radiusSum := aCircleShape.Radius + bCircleShape.Radius

	isColliding := ab.Dot(ab) <= radiusSum*radiusSum

	return isColliding
}
