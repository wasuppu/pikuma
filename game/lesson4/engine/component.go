package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Component interface {
	SetOwner(*Entity)
	Initialize()
	Update(float64)
	Render(*sdl.Renderer)
}

type TransformComponent struct {
	owner    *Entity
	position Vec2
	velocity Vec2
	width    int
	height   int
	scale    int
}

func NewTransformComponent(position, velocity Vec2, width, height, scale int) *TransformComponent {
	return &TransformComponent{position: position, velocity: velocity, width: width, height: height, scale: scale}
}

func (c *TransformComponent) SetOwner(e *Entity) {
	c.owner = e
}

func (c TransformComponent) Initialize() {}

func (c *TransformComponent) Update(deltaTime float64) {
	c.position[0] += c.velocity[0] * deltaTime
	c.position[1] += c.velocity[1] * deltaTime
}

func (c *TransformComponent) Render(renderer *sdl.Renderer) {
	transformRectangle := sdl.Rect{
		X: int32(c.position.X()),
		Y: int32(c.position.Y()),
		W: int32(c.width),
		H: int32(c.height),
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(&transformRectangle)
}
