package physics

import "math"

type Constraint interface {
	Solve()
	PreSolve(dt float64)
	PostSolve()

	A() *Body
	B() *Body
	APoint() Vec2
	BPoint() Vec2
}

type ConstraintBase struct {
	a *Body
	b *Body

	aPoint Vec2 // The anchor point in A's local space
	bPoint Vec2 // The anchor point in B's local space
}

func (c ConstraintBase) A() *Body {
	return c.a
}

func (c ConstraintBase) B() *Body {
	return c.b
}

func (c ConstraintBase) APoint() Vec2 {
	return c.aPoint
}

func (c ConstraintBase) BPoint() Vec2 {
	return c.bPoint
}

// Mat6x6 with the all inverse mass and inverse I of bodies "a" and "b"
func (c ConstraintBase) GetInvM() Mat {
	invM := NewMat(6, 6)
	invM[0][0] = c.a.InvMass
	invM[1][1] = c.a.InvMass
	invM[2][2] = c.a.InvI
	invM[3][3] = c.b.InvMass
	invM[4][4] = c.b.InvMass
	invM[5][5] = c.b.InvI
	return invM
}

// VecN with the all linear and angular velocities of bodies "a" and "b"
func (c ConstraintBase) GetVelocities() Vec {
	v := make(Vec, 6)
	v[0] = c.a.Velocity.X
	v[1] = c.a.Velocity.Y
	v[2] = c.a.AngularVelocity
	v[3] = c.b.Velocity.X
	v[4] = c.b.Velocity.Y
	v[5] = c.b.AngularVelocity
	return v
}

type JointConstraint struct {
	ConstraintBase
	Jacobian     Mat
	cachedLambda Vec
	bias         float64
}

func NewJointConstraint(a, b *Body, anchorPoint Vec2) *JointConstraint {
	return &JointConstraint{ConstraintBase{a: a, b: b, aPoint: a.WorldSpaceToLocalSpace(anchorPoint), bPoint: b.WorldSpaceToLocalSpace(anchorPoint)}, NewMat(1, 6), make(Vec, 1), 0}
}

func (c *JointConstraint) PreSolve(dt float64) {
	// Get the anchor point position in world space
	pa := c.a.LocalSpaceToWorldSpace(c.aPoint)
	pb := c.b.LocalSpaceToWorldSpace(c.bPoint)

	ra := pa.Sub(c.a.Position)
	rb := pb.Sub(c.b.Position)

	c.Jacobian = NewMat(1, 6)

	J1 := pa.Sub(pb).Muln(2.0)
	c.Jacobian[0][0] = J1.X // A linear velocity.x
	c.Jacobian[0][1] = J1.Y // A linear velocity.y

	J2 := ra.Cross(pa.Sub(pb)) * 2.0
	c.Jacobian[0][2] = J2 // A angular velocity

	J3 := pb.Sub(pa).Muln(2.0)
	c.Jacobian[0][3] = J3.X // B linear velocity.x
	c.Jacobian[0][4] = J3.Y // B linear velocity.y

	J4 := rb.Cross(pb.Sub(pa)) * 2.0
	c.Jacobian[0][5] = J4 // B angular velocity

	// Warm starting (apply cached lambda)
	Jt := c.Jacobian.Transpose()
	impulses := Jt.Mulv(c.cachedLambda)

	// Apply the impulses to both bodies
	c.a.ApplyImpulseLinear(Vec2{impulses[0], impulses[1]}) // A linear impulse
	c.a.ApplyImpulseAngular(impulses[2])                   // A angular impulse
	c.b.ApplyImpulseLinear(Vec2{impulses[3], impulses[4]}) // B linear impulse
	c.b.ApplyImpulseAngular(impulses[5])                   // B angular impulse

	// Compute the bias term (baumgarte stabilization)
	beta := 0.2
	C := pb.Sub(pa).Dot(pb.Sub(pa))
	C = math.Max(0, C-0.01)
	c.bias = beta / dt * C
}

func (c *JointConstraint) Solve() {
	V := c.GetVelocities()
	invM := c.GetInvM()

	J := c.Jacobian
	Jt := c.Jacobian.Transpose()

	// Compute lambda using Ax=b (Gauss-Seidel method)
	lhs := J.Mul(invM).Mul(Jt) // A
	rhs := J.Mulv(V).Muln(-1)  // b
	rhs[0] -= c.bias

	lambda := lhs.SolveGaussSeidel(rhs)
	c.cachedLambda = c.cachedLambda.Add(lambda)

	// Compute the impulses with both direction and magnitude
	impulses := Jt.Mulv(lambda)

	// Apply the impulses to both bodies
	c.a.ApplyImpulseLinear(Vec2{impulses[0], impulses[1]}) // A linear impulse
	c.a.ApplyImpulseAngular(impulses[2])                   // A angular impulse
	c.b.ApplyImpulseLinear(Vec2{impulses[3], impulses[4]}) // B linear impulse
	c.b.ApplyImpulseAngular(impulses[5])                   // B angular impulse
}

func (c *JointConstraint) PostSolve() {

}

type PenetrationConstraint struct {
	ConstraintBase
	Jacobian     Mat
	cachedLambda Vec
	bias         float64
	normal       Vec2    // Normal direction of the penetration in A's local space
	friction     float64 // Friction coefficient between the two penetrating bodies
}

func NewPenetrationConstraint(a, b *Body, aCollisionPoint, bCollisionPoint, normal Vec2) *PenetrationConstraint {
	return &PenetrationConstraint{
		ConstraintBase{a: a, b: b, aPoint: a.WorldSpaceToLocalSpace(aCollisionPoint), bPoint: b.WorldSpaceToLocalSpace(bCollisionPoint)},
		NewMat(2, 6),
		make(Vec, 2),
		0,
		a.WorldSpaceToLocalSpace(normal),
		0,
	}
}

func (c *PenetrationConstraint) PreSolve(dt float64) {
	// Get the collision points and normal in world space
	pa := c.a.LocalSpaceToWorldSpace(c.aPoint)
	pb := c.b.LocalSpaceToWorldSpace(c.bPoint)
	n := c.a.LocalSpaceToWorldSpace(c.normal)

	ra := pa.Sub(c.a.Position)
	rb := pb.Sub(c.b.Position)

	c.Jacobian = NewMat(2, 6)

	// Populate the first row of the Jacobian (normal vector)
	c.Jacobian[0][0] = -n.X         // A linear velocity.x
	c.Jacobian[0][1] = -n.Y         // A linear velocity.y
	c.Jacobian[0][2] = -ra.Cross(n) // A angular velocity
	c.Jacobian[0][3] = n.X          // B linear velocity.x
	c.Jacobian[0][4] = n.Y          // B linear velocity.y
	c.Jacobian[0][5] = rb.Cross(n)  // B angular velocity

	// Populate the second row of the Jacobian (tangent vector)
	friction := math.Max(c.a.Friction, c.b.Friction)
	if friction > 0 {
		t := n.Normal()                 // The tangent is the vector perpendicular to the normal
		c.Jacobian[1][0] = -t.X         // A linear velocity.x
		c.Jacobian[1][1] = -t.Y         // A linear velocity.y
		c.Jacobian[1][2] = -ra.Cross(t) // A angular velocity
		c.Jacobian[1][3] = t.X          // B linear velocity.x
		c.Jacobian[1][4] = t.Y          // B linear velocity.y
		c.Jacobian[1][5] = rb.Cross(t)  // B angukar velocity
	}

	// Warm starting (apply cached lambda)
	Jt := c.Jacobian.Transpose()
	impulses := Jt.Mulv(c.cachedLambda)

	// Apply the impulses to both bodies
	c.a.ApplyImpulseLinear(Vec2{impulses[0], impulses[1]}) // A linear impulse
	c.a.ApplyImpulseAngular(impulses[2])                   // A angular impulse
	c.b.ApplyImpulseLinear(Vec2{impulses[3], impulses[4]}) // B linear impulse
	c.b.ApplyImpulseAngular(impulses[5])                   // B angular impulse

	// Compute the bias term (baumgarte stabilization)
	beta := 0.2
	C := pb.Sub(pa).Dot(n.Muln(-1))
	C = math.Min(0, C+0.01)
	c.bias = beta / dt * C
}

func (c *PenetrationConstraint) Solve() {
	V := c.GetVelocities()
	invM := c.GetInvM()

	J := c.Jacobian
	Jt := c.Jacobian.Transpose()

	// Compute lambda using Ax=b (Gauss-Seidel method)
	lhs := J.Mul(invM).Mul(Jt) // A
	rhs := J.Mulv(V).Muln(-1)  // b
	rhs[0] -= c.bias
	lambda := lhs.SolveGaussSeidel(rhs)

	// Accumulate impulses and clamp it within constraint limits
	oldLambda := c.cachedLambda
	c.cachedLambda = c.cachedLambda.Add(lambda)
	if c.cachedLambda[0] < 0 {
		c.cachedLambda[0] = 0
	}

	// Keep friction values between  -(λn*µ) and (λn*µ)
	if c.friction > 0 {
		maxFriction := c.cachedLambda[0] * c.friction
		c.cachedLambda[1] = Clamp(c.cachedLambda[1], -maxFriction, maxFriction)
	}

	lambda = c.cachedLambda.Sub(oldLambda)

	// Compute the impulses with both direction and magnitude
	impulses := Jt.Mulv(lambda)

	// Apply the impulses to both bodies
	c.a.ApplyImpulseLinear(Vec2{impulses[0], impulses[1]}) // A linear impulse
	c.a.ApplyImpulseAngular(impulses[2])                   // A angular impulse
	c.b.ApplyImpulseLinear(Vec2{impulses[3], impulses[4]}) // B linear impulse
	c.b.ApplyImpulseAngular(impulses[5])                   // B angular impulse
}

func (c *PenetrationConstraint) PostSolve() {

}
