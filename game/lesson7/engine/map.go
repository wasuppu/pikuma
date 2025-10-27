package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type Map struct {
	manager   *EntityManager
	texture   *sdl.Texture
	scale     int
	titleSize int
}

func (m *Map) LoadMap(filepath string, mapSizeX, mapSizeY int) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %v", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	for y := range mapSizeY {
		for x := range mapSizeX {
			n, err := readDigit(reader)
			if err != nil {
				return fmt.Errorf("failed to read digit %v", err)
			}
			sourceRectY := n * m.titleSize
			n, err = readDigit(reader)
			if err != nil {
				return fmt.Errorf("failed to read digit %v", err)
			}
			sourceRectX := n * m.titleSize
			m.AddTile(sourceRectX, sourceRectY, x*m.scale*m.titleSize, y*m.scale*m.titleSize)

			_, err = reader.ReadByte()
			if err != nil {
				return fmt.Errorf("failed to read separator %v", err)
			}
		}
	}

	return nil
}

func (m *Map) AddTile(sourceRectX, sourceRectY, x, y int) {
	tile := m.manager.AddEntity("tile")
	tile.AddComponent(NewTileComponent(sourceRectX, sourceRectY, x, y, m.titleSize, m.scale, m.texture), TILE_COMPONENT)
}

func readDigit(r *bufio.Reader) (int, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(b))
}
