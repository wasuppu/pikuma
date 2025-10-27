package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type EntityManager struct {
	renderer *sdl.Renderer
	event    *sdl.Event
	camera   *sdl.Rect
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

func (m *EntityManager) Render() {
	for layerNumber := range NUM_LAYERS {
		for _, entity := range m.GetEntitiesByLayer(LayerType(layerNumber)) {
			entity.Render(m.renderer)
		}
	}
}

func (m EntityManager) HasNoEntities() bool {
	return len(m.entities) == 0
}

func (m *EntityManager) AddEntity(entityName string, layer LayerType) *Entity {
	entity := Entity{manager: m, name: entityName, isActive: true, componentTypeMap: make(map[ComponentType]Component)}
	m.entities = append(m.entities, &entity)
	return &entity
}

func (m EntityManager) GetEntities() []*Entity {
	return m.entities
}

func (m EntityManager) GetEntitiesByLayer(layer LayerType) []*Entity {
	selectedEntities := []*Entity{}
	for _, entity := range m.entities {
		if entity.layer == layer {
			selectedEntities = append(selectedEntities, entity)
		}
	}
	return selectedEntities
}

func (m EntityManager) GetEntityCount() int {
	return len(m.entities)
}

func (m EntityManager) CheckEntityCollisions(myEntity *Entity) string {
	myCollider := myEntity.GetComponent(COLLIDER_COMPONENT).(*ColliderComponent)
	for _, entity := range m.entities {
		if entity.name != myEntity.name && entity.name != "Tile" {
			if entity.HasComponent(COLLIDER_COMPONENT) {
				otherCollider := entity.GetComponent(COLLIDER_COMPONENT).(*ColliderComponent)
				if CheckRectangleCollision(myCollider.collider, otherCollider.collider) {
					return otherCollider.colliderTag
				}
			}
		}
	}
	return ""
}
