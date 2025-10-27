package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ComponentType int

const (
	TRANSFORM_COMPONENT ComponentType = iota
	SPRITE_COMPONENT
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

func (c *TransformComponent) Initialize() {}

func (c *TransformComponent) Update(deltaTime float64) {
	c.position[0] += c.velocity[0] * deltaTime
	c.position[1] += c.velocity[1] * deltaTime
}

func (c *TransformComponent) Render(renderer *sdl.Renderer) {}

type SpriteComponent struct {
	owner                *Entity
	transform            *TransformComponent
	texture              *sdl.Texture
	sourceRectangle      sdl.Rect
	destinationRectangle sdl.Rect
	isAnimated           bool
	numFrames            int
	animationSpeed       int
	isFixed              bool
	animations           map[string]Animation
	currentAnimationName string
	animationIndex       uint
	spriteFilp           sdl.RendererFlip
}

func NewSpriteComponent(texture *sdl.Texture) *SpriteComponent {
	return &SpriteComponent{texture: texture, animations: make(map[string]Animation)}
}

func NewSpriteComponent2(texture *sdl.Texture, numFrames, animationSpeed int, hasDirections, isFixed bool) *SpriteComponent {
	sprite := &SpriteComponent{texture: texture, animations: make(map[string]Animation), isAnimated: true, numFrames: numFrames, animationSpeed: animationSpeed, isFixed: isFixed}

	if hasDirections {
		downAnimation := Animation{0, numFrames, animationSpeed}
		rightAnimation := Animation{1, numFrames, animationSpeed}
		leftAnimation := Animation{2, numFrames, animationSpeed}
		upAnimation := Animation{3, numFrames, animationSpeed}

		sprite.animations["DownAnimation"] = downAnimation
		sprite.animations["RightAnimation"] = rightAnimation
		sprite.animations["LeftAnimation"] = leftAnimation
		sprite.animations["UpAnimation"] = upAnimation

		sprite.animationIndex = 0
		sprite.currentAnimationName = "DownAnimation"
	} else {
		singleAnimation := Animation{0, numFrames, animationSpeed}
		sprite.animations["SingleAnimation"] = singleAnimation
		sprite.animationIndex = 0
		sprite.currentAnimationName = "SingleAnimation"
	}

	sprite.Play(sprite.currentAnimationName)

	return sprite
}

func (c *SpriteComponent) Play(animationName string) {
	animation := c.animations[animationName]
	c.numFrames = animation.numFrames
	c.animationIndex = animation.index
	c.animationSpeed = animation.animationSpeed
	c.currentAnimationName = animationName
}

func (c *SpriteComponent) SetOwner(e *Entity) {
	c.owner = e
}

func (c *SpriteComponent) Initialize() {
	c.transform = c.owner.GetComponent(TRANSFORM_COMPONENT).(*TransformComponent)
	c.sourceRectangle.X = 0
	c.sourceRectangle.Y = 0
	c.sourceRectangle.W = int32(c.transform.width)
	c.sourceRectangle.H = int32(c.transform.height)
}

func (c *SpriteComponent) Update(deltaTime float64) {
	if c.isAnimated {
		c.sourceRectangle.X = c.sourceRectangle.W * int32(int(float64(sdl.GetTicks64())/float64(c.animationSpeed))%c.numFrames)
	}
	c.sourceRectangle.Y = int32(c.animationIndex) * int32(c.transform.height)

	c.destinationRectangle.X = int32(c.transform.position.X())
	c.destinationRectangle.Y = int32(c.transform.position.Y())
	c.destinationRectangle.W = int32(c.transform.width * c.transform.scale)
	c.destinationRectangle.H = int32(c.transform.height * c.transform.scale)
}

func (c *SpriteComponent) Render(renderer *sdl.Renderer) {
	DrawTexture(c.texture, c.sourceRectangle, c.destinationRectangle, c.spriteFilp, renderer)
}
