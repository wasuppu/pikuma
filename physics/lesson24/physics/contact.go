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

	c.A.Shape.UpdateVertices(c.A.Rotation, c.A.Position)
	c.B.Shape.UpdateVertices(c.B.Rotation, c.B.Position)
}

func (c *Contact) ResolveCollision() {
	// Apply positional correction using the projection method
	c.ResolvePenetration()

	// Define coefficient of restitution (elasticity) and friction
	e := math.Min(c.A.Restitution, c.B.Restitution)
	f := math.Min(c.A.Friction, c.B.Friction)

	// Calculate the relative velocity between the two objects
	ra := c.End.Sub(c.A.Position)
	rb := c.Start.Sub(c.B.Position)
	va := c.A.Velocity.Add(Vec2{-c.A.AngularVelocity * ra.Y, c.A.AngularVelocity * ra.X})
	vb := c.B.Velocity.Add(Vec2{-c.B.AngularVelocity * rb.Y, c.B.AngularVelocity * rb.X})
	vrel := va.Sub(vb)

	// Now we proceed to calculate the collision impulse along the normal
	vrelDotNormal := vrel.Dot(c.Normal)
	impulseDirectionN := c.Normal
	impulseMagnitudeN := -(1 + e) * vrelDotNormal / ((c.A.InvMass + c.B.InvMass) + ra.Cross(c.Normal)*ra.Cross(c.Normal)*c.A.InvI + rb.Cross(c.Normal)*rb.Cross(c.Normal)*c.B.InvI)
	jN := impulseDirectionN.Muln(impulseMagnitudeN)

	// Now we proceed to calculate the collision impulse along the tangent
	tangent := c.Normal.Normal()
	vrelDotTangent := vrel.Dot(tangent)
	impulseDirectionT := tangent
	impulseMagnitudeT := f * -(1 + e) * vrelDotTangent / ((c.A.InvMass + c.B.InvMass) + ra.Cross(tangent)*ra.Cross(tangent)*c.A.InvI + rb.Cross(tangent)*rb.Cross(tangent)*c.B.InvI)
	jT := impulseDirectionT.Muln(impulseMagnitudeT)

	// Calculate the final impulse j combining normal and tangent impulses
	j := jN.Add(jT)

	// Apply the impulse vector to both objects in opposite direction
	c.A.ApplyImpulse2(j, ra)
	c.B.ApplyImpulse2(j.Muln(-1), rb)
}
