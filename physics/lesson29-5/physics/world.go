package physics

const (
	PIXELS_PER_METER = 50
)

type World struct {
	G           float64
	bodies      []*Body
	constraints []Constraint
	forces      []Vec2
	torques     []float64
}

func NewWorld(gravity float64) *World {
	world := World{}
	world.G = -gravity
	return &world
}

func (w *World) AddBody(body *Body) {
	w.bodies = append(w.bodies, body)
}

func (w *World) GetBodies() []*Body {
	return w.bodies
}

func (w *World) AddConstraint(constraint Constraint) {
	w.constraints = append(w.constraints, constraint)
}

func (w *World) GetConstraints() []Constraint {
	return w.constraints
}

func (w *World) AddForce(force Vec2) {
	w.forces = append(w.forces, force)
}

func (w *World) Update(dt float64) {
	// Create a vector of penetration constraints that will be solved frame per frame
	penetrations := []*PenetrationConstraint{}

	// Loop all bodies of the world applying forces
	for _, body := range w.bodies {
		// Apply the weight force to all bodies
		weight := Vec2{X: 0.0, Y: body.Mass * w.G * PIXELS_PER_METER}
		body.AddFore(weight)

		// Apply forces to all bodies
		for _, force := range w.forces {
			body.AddFore(force)
		}

		// Apply torque to all bodies
		for _, torque := range w.torques {
			body.AddTorque(torque)
		}
	}

	// Integrate all the forces
	for _, body := range w.bodies {
		body.IntegrateForces(dt)
	}

	// Check all the bodies with all other bodies detecting collisions
	for i := range w.bodies {
		for j := i + 1; j < len(w.bodies); j++ {
			a := w.bodies[i]
			b := w.bodies[j]
			if isColliding, contact := IsColliding(a, b); isColliding {
				// Create a new penetration constraint
				penetration := NewPenetrationConstraint(contact.A, contact.B, contact.Start, contact.End, contact.Normal)
				penetrations = append(penetrations, penetration)
			}
		}
	}

	// Solve all constraints
	for _, constraint := range w.constraints {
		constraint.PreSolve(dt)
	}
	for _, constraint := range penetrations {
		constraint.PreSolve(dt)
	}
	for range 5 {
		for _, constraint := range w.constraints {
			constraint.Solve()
		}
		for _, constraint := range penetrations {
			constraint.Solve()
		}
	}
	for _, constraint := range w.constraints {
		constraint.PostSolve()
	}
	for _, constraint := range penetrations {
		constraint.PostSolve()
	}

	// Integrate all the velocities
	for _, body := range w.bodies {
		body.IntegrateVelocities(dt)
	}
}
