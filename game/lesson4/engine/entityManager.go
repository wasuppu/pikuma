package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type EntityManager struct {
	entities []*Entity
}

func (m *EntityManager) ClearData() {
	for i := range m.entities {
		m.entities[i].Destroy()
	}
}

func (m *EntityManager) Update(deltaTime float64) {
	for i := range m.entities {
		m.entities[i].Update(deltaTime)
	}
}

func (m *EntityManager) Render(renderer *sdl.Renderer) {
	for i := range m.entities {
		m.entities[i].Render(renderer)
	}
}

func (m EntityManager) HasNoEntities() bool {
	return len(m.entities) == 0
}

func (m *EntityManager) AddEntity(entityName string) *Entity {
	entity := Entity{manager: m, name: entityName, isActive: true}
	m.entities = append(m.entities, &entity)
	return &entity
}

func (m EntityManager) GetEntities() []*Entity {
	return m.entities
}

func (m EntityManager) GetEntityCount() int {
	return len(m.entities)
}
