package engine

import (
	"fmt"
	"reflect"

	"github.com/veandco/go-sdl2/sdl"
)

type Entity struct {
	manager    *EntityManager
	isActive   bool
	components []Component
	name       string
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

func (e *Entity) AddComponent(component Component) Component {
	component.SetOwner(e)
	component.Initialize()
	typ := reflect.TypeOf(component)
	fmt.Println(typ)
	e.components = append(e.components, component)
	return component
}
