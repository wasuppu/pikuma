package physics

type Particle struct {
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	SumForces    Vec2
	Mass         float64
	InvMass      float64
	Radius       float64
}

func NewParticle(x, y, mass float64) *Particle {
	invMass := 0.0
	if mass != 0 {
		invMass = 1.0 / mass
	}
	return &Particle{Position: Vec2{x, y}, Mass: mass, InvMass: invMass}
}

func (p *Particle) Integrate(dt float64) {
	// Find the acceleration based on the forces that are being applied and the mass
	p.Acceleration = p.SumForces.Muln(p.InvMass)

	// Integrate the acceleration to find the new velocity
	p.Velocity = p.Velocity.Add(p.Acceleration.Muln(dt))

	// Integrate the velocity to find the new position
	p.Position = p.Position.Add(p.Velocity.Muln(dt))

	// Clear all the forces acting on the object before the next physics step
	p.ClearFore()
}

func (p *Particle) AddFore(force Vec2) {
	p.SumForces = p.SumForces.Add(force)
}

func (p *Particle) ClearFore() {
	p.SumForces = Vec2{0, 0}
}

func (p *Particle) GenerateDragForce(k float64) Vec2 {
	dragForce := Vec2{0, 0}
	magnitude := p.Velocity.Dot(p.Velocity)
	if magnitude > 0 {
		// Calculate the drag direction (inverse of velocity unit vector)
		dragDirection := p.Velocity.Normalize().Muln(-1)

		// Calculate the drag magnitude, k * |v|^2
		dragMagnitude := k * magnitude

		// Generate the final drag force with direction and magnitude
		dragForce = dragDirection.Muln(dragMagnitude)
	}
	return dragForce
}

func (p *Particle) GenerateFrictionForce(k float64) Vec2 {
	// Calculate the friction direction (inverse of velocity unit vector)
	frictionDirection := p.Velocity.Normalize().Muln(-1)

	// Calculate the friction magnitude
	frictionMagnitude := k

	// Calculate the final friction force
	frictionForce := frictionDirection.Muln(frictionMagnitude)
	return frictionForce
}
