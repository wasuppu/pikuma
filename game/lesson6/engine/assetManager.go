package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type AssetManager struct {
	manager  *EntityManager
	textures map[string]*sdl.Texture
}

func (m *AssetManager) ClearData() {
	for k := range m.textures {
		delete(m.textures, k)
	}
}

func (m *AssetManager) AddTexture(textureId string, filename string) error {
	texture, err := LoadTexture(filename, m.manager.renderer)
	if err != nil {
		return err
	}
	m.textures[textureId] = texture
	return nil
}

func (m AssetManager) GetTexture(textureId string) *sdl.Texture {
	return m.textures[textureId]
}

func LoadTexture(filename string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	surface, err := img.Load(filename)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture: %v", err)
	}
	return texture, nil
}

func DrawTexture(texture *sdl.Texture, sourceRectangle, destinationRectangle sdl.Rect, flip sdl.RendererFlip, renderer *sdl.Renderer) {
	renderer.CopyEx(texture, &sourceRectangle, &destinationRectangle, 0.0, nil, flip)
}
