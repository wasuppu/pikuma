package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AssetManager struct {
	renderer *sdl.Renderer
	textures map[string]*sdl.Texture
	fonts    map[string]*ttf.Font
}

func (m *AssetManager) ClearData() {
	for k := range m.textures {
		delete(m.textures, k)
	}
	for k := range m.fonts {
		delete(m.fonts, k)
	}
}

func (m *AssetManager) AddTexture(textureId string, filename string) error {
	texture, err := LoadTexture(filename, m.renderer)
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

func (m *AssetManager) AddFont(fontId string, filename string, filesize int) error {
	font, err := LoadFont(filename, filesize)
	if err != nil {
		panic(err)
	}
	m.fonts[fontId] = font
	return nil
}

func (m AssetManager) GetFont(fontId string) *ttf.Font {
	return m.fonts[fontId]
}

func LoadFont(filename string, fontsize int) (*ttf.Font, error) {
	return ttf.OpenFont(filename, fontsize)
}

func DrawFont(texture *sdl.Texture, position sdl.Rect, renderer *sdl.Renderer) {
	renderer.Copy(texture, nil, &position)
}
