package physics

import "math"

func IsColliding(a *Body, b *Body, contacts []Contact) (bool, []Contact) {
	aIsCircle := a.Shape.GetType() == CIRCLE_SHAPE
	bIsCircle := b.Shape.GetType() == CIRCLE_SHAPE
	aIsPolygon := a.Shape.GetType() == POLYGON_SHAPE || a.Shape.GetType() == BOX_SHAPE
	bIsPolygon := b.Shape.GetType() == POLYGON_SHAPE || b.Shape.GetType() == BOX_SHAPE

	if aIsCircle && bIsCircle {
		return IsCollidingCircleCircle(a, b, contacts)
	}

	if aIsPolygon && bIsPolygon {
		return IsCollidingPolygonPolygon(a, b, contacts)
	}

	if aIsPolygon && bIsCircle {
		return IsCollidingPolygonCircle(a, b, contacts)
	}

	if aIsCircle && bIsPolygon {
		return IsCollidingPolygonCircle(b, a, contacts)
	}

	return false, contacts
}

func IsCollidingCircleCircle(a *Body, b *Body, contacts []Contact) (bool, []Contact) {
	aCircleShape := a.Shape.(*CircleShape)
	bCircleShape := b.Shape.(*CircleShape)

	ab := b.Position.Sub(a.Position)
	radiusSum := aCircleShape.Radius + bCircleShape.Radius

	isColliding := ab.Dot(ab) <= radiusSum*radiusSum

	if !isColliding {
		return false, contacts
	}

	contact := Contact{A: a, B: b, Normal: ab.Normalize()}

	contact.Start = b.Position.Sub(contact.Normal.Muln(bCircleShape.Radius))
	contact.End = a.Position.Add(contact.Normal.Muln(aCircleShape.Radius))
	contact.Depth = contact.End.Sub(contact.Start).Length()

	contacts = append(contacts, contact)

	return true, contacts
}

func IsCollidingPolygonPolygon(a *Body, b *Body, contacts []Contact) (bool, []Contact) {
	aPolygonShape := a.Shape.(PolygonShape)
	bPolygonShape := b.Shape.(PolygonShape)
	abSeparation, aAxis, aPoint := aPolygonShape.FindMinSeparation(bPolygonShape)
	if abSeparation >= 0 {
		return false, contacts
	}
	baSeparation, bAxis, bPoint := bPolygonShape.FindMinSeparation(aPolygonShape)
	if baSeparation >= 0 {
		return false, contacts
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

	contacts = append(contacts, contact)

	return true, contacts
}

func IsCollidingPolygonCircle(polygon *Body, circle *Body, contacts []Contact) (bool, []Contact) {
	contact := Contact{}
	polygonShape := polygon.Shape.(PolygonShape)
	circleShape := circle.Shape.(*CircleShape)
	polygonVertices := polygonShape.GetWorldVertices()

	isOutside := false
	var minCurVertex, minNextVertex Vec2
	distanceCircleEdge := -math.MaxFloat64

	// Loop all the edges of the polygon/box finding the nearest edge to the circle center
	for i := range len(polygonVertices) {
		n := (i + 1) % len(polygonVertices)
		edge := polygonShape.EdgeAt(i)
		normal := edge.Normal()

		// Compare the circle center with the rectangle vertex
		vertexToCircleCenter := circle.Position.Sub(polygonVertices[i])
		projection := vertexToCircleCenter.Dot(normal)

		// If we found a dot product projection that is in the positive side of the normal
		if projection > 0 {
			// Circle center is outside the polygon
			distanceCircleEdge = projection
			minCurVertex = polygonVertices[i]
			minNextVertex = polygonVertices[n]
			isOutside = true
			break
		} else {
			// Circle center is inside the rectangle, find the min edge (the one with the least negative projection)
			if projection > distanceCircleEdge {
				distanceCircleEdge = projection
				minCurVertex = polygonVertices[i]
				minNextVertex = polygonVertices[n]
			}
		}
	}

	if isOutside {
		// Check if we are inside region A:
		v1 := circle.Position.Sub(minCurVertex) // vector from the nearest vertex to the circle center
		v2 := minNextVertex.Sub(minCurVertex)   // the nearest edge (from curr vertex to next vertex)
		if v1.Dot(v2) < 0 {
			// Distance from vertex to circle center is greater than radius... no collision
			if v1.Length() > circleShape.Radius {
				return false, contacts
			} else {
				// Detected collision in region A:
				contact.A = polygon
				contact.B = circle
				contact.Depth = circleShape.Radius - v1.Length()
				contact.Normal = v1.Normalize()
				contact.Start = circle.Position.Add(contact.Normal.Muln(-circleShape.Radius))
				contact.End = contact.Start.Add(contact.Normal.Muln(contact.Depth))
			}
		} else {
			// Check if we are inside region B:
			v1 := circle.Position.Sub(minNextVertex) // vector from the next nearest vertex to the circle center
			v2 := minCurVertex.Sub(minNextVertex)    // the nearest edge
			if v1.Dot(v2) < 0 {
				// Distance from vertex to circle center is greater than radius... no collision
				if v1.Length() > circleShape.Radius {
					return false, contacts
				} else {
					// Detected collision in region B:
					contact.A = polygon
					contact.B = circle
					contact.Depth = circleShape.Radius - v1.Length()
					contact.Normal = v1.Normalize()
					contact.Start = circle.Position.Add(contact.Normal.Muln(-circleShape.Radius))
					contact.End = contact.Start.Add(contact.Normal.Muln(contact.Depth))
				}
			} else {
				// We are inside region C:
				if distanceCircleEdge > circleShape.Radius {
					// No collision... Distance between the closest distance and the circle center is greater than the radius.
					return false, contacts
				} else {
					// Detected collision in region C:
					contact.A = polygon
					contact.B = circle
					contact.Depth = circleShape.Radius - distanceCircleEdge
					contact.Normal = minNextVertex.Sub(minCurVertex).Normal()
					contact.Start = circle.Position.Sub(contact.Normal.Muln(circleShape.Radius))
					contact.End = contact.Start.Add(contact.Normal.Muln(contact.Depth))
				}
			}
		}
	} else {
		// The center of circle is inside the polygon... it is definitely colliding!
		contact.A = polygon
		contact.B = circle
		contact.Depth = circleShape.Radius - distanceCircleEdge
		contact.Normal = minNextVertex.Sub(minCurVertex).Normal()
		contact.Start = circle.Position.Sub(contact.Normal.Muln(circleShape.Radius))
		contact.End = contact.Start.Add(contact.Normal.Muln(contact.Depth))
	}

	contacts = append(contacts, contact)

	return true, contacts
}
