package physics

type Particle struct {
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	Mass         float64
	Radius       float64
}

func NewParticle(x, y, mass float64) Particle {
	return Particle{Position: Vec2{x, y}, Mass: mass}
}

func (p *Particle) Integrate(dt float64) {
	p.Velocity = p.Velocity.Add(p.Acceleration.Muln(dt))
	p.Position = p.Position.Add(p.Velocity.Muln(dt))
}
