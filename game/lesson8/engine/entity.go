package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Entity struct {
	manager          *EntityManager
	isActive         bool
	components       []Component
	componentTypeMap map[ComponentType]Component
	name             string
	layer            LayerType
}

func (e *Entity) Update(deltaTime float64) {
	for i := range e.components {
		e.components[i].Update(deltaTime)
	}
}

func (e *Entity) Render(renderer *sdl.Renderer) {
	for i := range e.components {
		e.components[i].Render(renderer)
	}
}

func (e *Entity) Destroy() {
	e.isActive = false
}

func (e Entity) IsActive() bool {
	return e.isActive
}

func (e *Entity) AddComponent(component Component, typ ComponentType) Component {
	component.SetOwner(e)
	component.Initialize()
	e.componentTypeMap[typ] = component
	e.components = append(e.components, component)
	return component
}

func (e Entity) GetComponent(typ ComponentType) Component {
	return e.componentTypeMap[typ]
}
