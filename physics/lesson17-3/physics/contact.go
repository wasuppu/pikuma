package physics

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
