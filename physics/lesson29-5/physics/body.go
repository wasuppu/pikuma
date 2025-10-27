package physics

import (
	"math"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Body struct {
	// Linear motion
	Position     Vec2
	Velocity     Vec2
	Acceleration Vec2

	Rotation            float64
	AngularVelocity     float64
	AngularAcceleration float64

	// Forces and torque
	SumForces Vec2
	sumTorque float64

	// Mass and Moment of Inertia
	Mass    float64
	InvMass float64
	I       float64
	InvI    float64

	// Coefficient of restitution (elasticity)
	Restitution float64

	// Coefficient of friction
	Friction float64

	// shape/geometry of this rigid body
	Shape Shape

	// Pointer to an SDL texture
	Texture *sdl.Texture
}

func NewBody(shape Shape, x, y, mass float64) *Body {
	body := &Body{Position: Vec2{x, y}, Mass: mass, Shape: shape}
	body.Restitution = 1.0
	body.Friction = 0.7

	invMass := 0.0
	if mass != 0 {
		invMass = 1.0 / mass
	}
	I := shape.GetMomentOfInertia() * mass
	invI := 0.0
	if I != 0 {
		invI = 1.0 / I
	}

	body.InvMass = invMass
	body.I = I
	body.InvI = invI

	body.Shape.UpdateVertices(body.Rotation, body.Position)
	return body
}

func (b *Body) SetTexture(textureFileName string, renderer *sdl.Renderer) error {
	surface, err := img.Load(textureFileName)
	if err != nil {
		return err
	}
	defer surface.Free()

	b.Texture, err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	return nil
}

func (b *Body) IsStatic() bool {
	epsilon := 0.005
	return math.Abs(b.InvMass-0.0) < epsilon
}

func (b *Body) LocalSpaceToWorldSpace(point Vec2) Vec2 {
	rotated := point.Rotate(b.Rotation)
	return rotated.Add(b.Position)
}

func (b *Body) WorldSpaceToLocalSpace(point Vec2) Vec2 {
	translated := point.Sub(b.Position)
	sin, cos := math.Sincos(-b.Rotation)
	rotatedX := cos*translated.X - sin*translated.Y
	rotatedY := cos*translated.Y + sin*translated.X
	return Vec2{rotatedX, rotatedY}
}

func (b *Body) ApplyImpulseLinear(j Vec2) {
	if b.IsStatic() {
		return
	}

	b.Velocity = b.Velocity.Add(j.Muln(b.InvMass))
}

func (b *Body) ApplyImpulseAngular(j float64) {
	if b.IsStatic() {
		return
	}

	b.AngularVelocity += j * b.InvI
}

func (b *Body) ApplyImpulseAtPoint(j, r Vec2) {
	if b.IsStatic() {
		return
	}

	b.Velocity = b.Velocity.Add(j.Muln(b.InvMass))
	b.AngularVelocity += r.Cross(j) * b.InvI
}

func (b *Body) IntegrateForces(dt float64) {
	if b.IsStatic() {
		return
	}

	// Find the acceleration based on the forces that are being applied and the mass
	b.Acceleration = b.SumForces.Muln(b.InvMass)

	// Integrate the acceleration to find the new velocity
	b.Velocity = b.Velocity.Add(b.Acceleration.Muln(dt))

	// Find the angular acceleration based on the torque that is being applied and the moment of inertia
	b.AngularAcceleration = b.sumTorque * b.InvI

	// Integrate the angular acceleration to find the new angular velocity
	b.AngularVelocity += b.AngularAcceleration * dt

	// Clear all the forces and torque acting on the object before the next physics step
	b.ClearFore()
	b.ClearTorque()
}

func (b *Body) IntegrateVelocities(dt float64) {
	if b.IsStatic() {
		return
	}

	// Integrate the velocity to find the new position
	b.Position = b.Position.Add(b.Velocity.Muln(dt))

	// Integrate the angular velocity to find the new rotation angle
	b.Rotation += b.AngularVelocity * dt

	// Update the vertices to adjust them to the new position/rotation
	b.Shape.UpdateVertices(b.Rotation, b.Position)
}

func (b *Body) AddFore(force Vec2) {
	b.SumForces = b.SumForces.Add(force)
}

func (b *Body) ClearFore() {
	b.SumForces = Vec2{0, 0}
}

func (b *Body) AddTorque(torque float64) {
	b.sumTorque += torque
}

func (b *Body) ClearTorque() {
	b.sumTorque = 0
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
