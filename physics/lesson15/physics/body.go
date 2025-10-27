package physics

type Body struct {
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	SumForces    Vec2
	Mass         float64
	InvMass      float64
	Radius       float64
}

func NewBody(x, y, mass float64) *Body {
	invMass := 0.0
	if mass != 0 {
		invMass = 1.0 / mass
	}
	return &Body{Position: Vec2{x, y}, Mass: mass, InvMass: invMass}
}

func (b *Body) Integrate(dt float64) {
	// Find the acceleration based on the forces that are being applied and the mass
	b.Acceleration = b.SumForces.Muln(b.InvMass)

	// Integrate the acceleration to find the new velocity
	b.Velocity = b.Velocity.Add(b.Acceleration.Muln(dt))

	// Integrate the velocity to find the new position
	b.Position = b.Position.Add(b.Velocity.Muln(dt))

	// Clear all the forces acting on the object before the next physics step
	b.ClearFore()
}

func (b *Body) AddFore(force Vec2) {
	b.SumForces = b.SumForces.Add(force)
}

func (b *Body) ClearFore() {
	b.SumForces = Vec2{0, 0}
}

func (b *Body) GenerateDragForce(k float64) Vec2 {
	dragForce := Vec2{0, 0}
	magnitude := b.Velocity.Dot(b.Velocity)
	if magnitude > 0 {
		// Calculate the drag direction (inverse of velocity unit vector)
		dragDirection := b.Velocity.Normalize().Muln(-1)

		// Calculate the drag magnitude, k * |v|^2
		dragMagnitude := k * magnitude

		// Generate the final drag force with direction and magnitude
		dragForce = dragDirection.Muln(dragMagnitude)
	}
	return dragForce
}

func (b *Body) GenerateFrictionForce(k float64) Vec2 {
	// Calculate the friction direction (inverse of velocity unit vector)
	frictionDirection := b.Velocity.Normalize().Muln(-1)

	// Calculate the friction magnitude (just k, for now)
	frictionMagnitude := k

	// Calculate the final resulting friction force vector
	frictionForce := frictionDirection.Muln(frictionMagnitude)
	return frictionForce
}

func (b *Body) GenerateGravitationalForce(b2 *Body, G, minDistance, maxDistance float64) Vec2 {
	// Calculate the distance between the two objects
	d := b2.Position.Sub(b.Position)
	distanceSquared := d.Dot(d)

	// Clamp the values of the distance (to allow for some insteresting visual effects)
	distanceSquared = Clamp(distanceSquared, minDistance, maxDistance)

	// Calculate the direction of the attraction force
	attractionDirection := d.Normalize()

	// Calculate the strength of the attraction force
	attractionMagnitude := G * (b.Mass * b2.Mass) / distanceSquared

	// Calculate the final resulting attraction force vector
	attractionForce := attractionDirection.Muln(attractionMagnitude)
	return attractionForce
}

func (b *Body) GenerateSpringForce(anchor Vec2, restLength, k float64) Vec2 {
	// Calculate the distance between the anchor and the object
	d := b.Position.Sub(anchor)

	// Find the spring displacement considering the rest length
	displacement := d.Length() - restLength

	// Calculate the direction of the spring force
	springDirection := d.Normalize()

	// Calculate the magnitude of the spring force
	springMagnitude := -k * displacement

	// Calculate the final resulting spring force vector
	springForce := springDirection.Muln(springMagnitude)
	return springForce
}
