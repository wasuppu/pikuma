package physics

type Particle struct {
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2
	Mass         float64
}

func NewParticle(x, y, mass float64) Particle {
	return Particle{Position: Vec2{x, y}, Mass: mass}
}
