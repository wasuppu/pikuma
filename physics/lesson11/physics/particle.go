package physics

type Particle struct {
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	SumForces    Vec2
	Mass         float64
	Radius       float64
}

func NewParticle(x, y, mass float64) Particle {
	return Particle{Position: Vec2{x, y}, Mass: mass}
}

func (p *Particle) Integrate(dt float64) {
	// Find the acceleration based on the forces that are being applied and the mass
	p.Acceleration = p.SumForces.Divn(p.Mass)

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
