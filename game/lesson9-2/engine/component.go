package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ComponentType int

const (
	TRANSFORM_COMPONENT ComponentType = iota
	SPRITE_COMPONENT
	KEYBOARD_CONTROL_COMPONENT
	TILE_COMPONENT
	COLLIDER_COMPONENT
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
	camera := c.owner.manager.camera

	if c.isAnimated {
		c.sourceRectangle.X = c.sourceRectangle.W * int32(int(float64(sdl.GetTicks64())/float64(c.animationSpeed))%c.numFrames)
	}
	c.sourceRectangle.Y = int32(c.animationIndex) * int32(c.transform.height)

	c.destinationRectangle.X = int32(c.transform.position.X())
	c.destinationRectangle.Y = int32(c.transform.position.Y())
	if !c.isFixed {
		c.destinationRectangle.X -= camera.X
		c.destinationRectangle.Y -= camera.Y
	}
	c.destinationRectangle.W = int32(c.transform.width * c.transform.scale)
	c.destinationRectangle.H = int32(c.transform.height * c.transform.scale)
}

func (c *SpriteComponent) Render(renderer *sdl.Renderer) {
	DrawTexture(c.texture, c.sourceRectangle, c.destinationRectangle, c.spriteFilp, renderer)
}

type KeyboardControlComponent struct {
	owner     *Entity
	upKey     string
	downKey   string
	rightKey  string
	leftKey   string
	shootKey  string
	transform *TransformComponent
	sprite    *SpriteComponent
}

func NewKeyboardControlComponent(upKey, rightKey, downKey, leftKey, shootKey string) *KeyboardControlComponent {
	return &KeyboardControlComponent{upKey: upKey, rightKey: rightKey, downKey: downKey, leftKey: leftKey, shootKey: shootKey}
}

func (c *KeyboardControlComponent) SetOwner(e *Entity) {
	c.owner = e
}

func (c *KeyboardControlComponent) Initialize() {
	c.transform = c.owner.GetComponent(TRANSFORM_COMPONENT).(*TransformComponent)
	c.sprite = c.owner.GetComponent(SPRITE_COMPONENT).(*SpriteComponent)
}

func (c *KeyboardControlComponent) Update(deltaTime float64) {
	event := *c.owner.manager.event
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		key := sdl.GetKeyName(t.Keysym.Sym)
		switch t.Type {
		case sdl.KEYDOWN:
			switch key {
			case c.upKey:
				c.transform.velocity[1] = -25
				c.transform.velocity[0] = 0
				c.sprite.Play("UpAnimation")
			case c.rightKey:
				c.transform.velocity[1] = 0
				c.transform.velocity[0] = 25
				c.sprite.Play("RightAnimation")
			case c.downKey:
				c.transform.velocity[1] = 25
				c.transform.velocity[0] = 0
				c.sprite.Play("DownAnimation")
			case c.leftKey:
				c.transform.velocity[1] = 0
				c.transform.velocity[0] = -25
				c.sprite.Play("LeftAnimation")
			}
		case sdl.KEYUP:
			switch key {
			case c.upKey:
				c.transform.velocity[1] = 0
			case c.rightKey:
				c.transform.velocity[0] = 0
			case c.downKey:
				c.transform.velocity[1] = 0
			case c.leftKey:
				c.transform.velocity[0] = 0
			}
		}
	}
}

func (c *KeyboardControlComponent) Render(renderer *sdl.Renderer) {}

type TileComponent struct {
	owner                *Entity
	texture              *sdl.Texture
	sourceRectangle      sdl.Rect
	destinationRectangle sdl.Rect
	position             Vec2
}

func NewTileComponent(sourceRectX, sourceRectY, x, y, tileSize, tileScale int, assetTexture *sdl.Texture) *TileComponent {
	tile := &TileComponent{texture: assetTexture}

	tile.sourceRectangle.X = int32(sourceRectX)
	tile.sourceRectangle.Y = int32(sourceRectY)
	tile.sourceRectangle.W = int32(tileSize)
	tile.sourceRectangle.H = int32(tileSize)

	tile.destinationRectangle.X = int32(x)
	tile.destinationRectangle.Y = int32(y)
	tile.destinationRectangle.W = int32(tileSize * tileScale)
	tile.destinationRectangle.H = int32(tileSize * tileScale)

	tile.position[0] = float64(x)
	tile.position[1] = float64(y)

	return tile
}

func (c *TileComponent) SetOwner(e *Entity) {
	c.owner = e
}

func (c *TileComponent) Initialize() {}

func (c *TileComponent) Update(deltaTime float64) {
	camera := c.owner.manager.camera
	c.destinationRectangle.X = int32(c.position.X() - float64(camera.X))
	c.destinationRectangle.Y = int32(c.position.Y() - float64(camera.Y))
}

func (c *TileComponent) Render(renderer *sdl.Renderer) {
	DrawTexture(c.texture, c.sourceRectangle, c.destinationRectangle, sdl.FLIP_NONE, renderer)
}

type ColliderComponent struct {
	owner                *Entity
	colliderTag          string
	collider             sdl.Rect
	sourceRectangle      sdl.Rect
	destinationRectangle sdl.Rect
	transform            *TransformComponent
}

func NewColliderComponent(colliderTag string, x, y, width, height int) *ColliderComponent {
	collider := &ColliderComponent{colliderTag: colliderTag}
	collider.collider = sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)}
	return collider
}

func (c *ColliderComponent) SetOwner(e *Entity) {
	c.owner = e
}

func (c *ColliderComponent) Initialize() {
	if c.owner.HasComponent(TRANSFORM_COMPONENT) {
		c.transform = c.owner.GetComponent(TRANSFORM_COMPONENT).(*TransformComponent)
		c.sourceRectangle = sdl.Rect{X: 0, Y: 0, W: int32(c.transform.width), H: int32(c.transform.height)}
		c.destinationRectangle = sdl.Rect{X: c.collider.X, Y: c.collider.Y, W: c.collider.W, H: c.collider.H}
	}
}

func (c *ColliderComponent) Update(deltaTime float64) {
	camera := c.owner.manager.camera

	c.collider.X = int32(c.transform.position.X())
	c.collider.Y = int32(c.transform.position.Y())
	c.collider.W = int32(c.transform.width * c.transform.scale)
	c.collider.H = int32(c.transform.height * c.transform.scale)

	c.destinationRectangle.X = c.collider.X - camera.X
	c.destinationRectangle.Y = c.collider.Y - camera.Y
}

func (c *ColliderComponent) Render(renderer *sdl.Renderer) {}
