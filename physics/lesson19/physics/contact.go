package physics

import "math"

type Contact struct {
	A      *Body
	B      *Body
	Start  Vec2
	End    Vec2
	Normal Vec2
	Depth  float64
}

func (c *Contact) ResolvePenetration() {
	if c.A.IsStatic() && c.B.IsStatic() {
		return
	}

	da := c.Depth / (c.A.InvMass + c.B.InvMass) * c.A.InvMass
	db := c.Depth / (c.A.InvMass + c.B.InvMass) * c.B.InvMass

	c.A.Position = c.A.Position.Sub(c.Normal.Muln(da))
	c.B.Position = c.B.Position.Add(c.Normal.Muln(db))
}

func (c *Contact) ResolveCollision() {
	// Apply positional correction using the projection method
	c.ResolvePenetration()

	// Define elasticity (coefficient of restitution e)
	e := math.Min(c.A.Restitution, c.B.Restitution)

	// Calculate the relative velocity between the two objects
	vrel := c.A.Velocity.Sub(c.B.Velocity)

	// Calculate the relative velocity along the normal collision vector
	vrelDotNormal := vrel.Dot(c.Normal)

	// Now we proceed to calculate the collision impulse
	impulseDirection := c.Normal
	impulseMagnitude := -(1 + e) * vrelDotNormal / (c.A.InvMass + c.B.InvMass)

	jn := impulseDirection.Muln(impulseMagnitude)

	// Apply the impulse vector to both objects in opposite direction
	c.A.ApplyImpulse(jn)
	c.B.ApplyImpulse(jn.Muln(-1))
}
