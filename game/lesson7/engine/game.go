package engine

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	FPS               = 60
	FRAME_TARGET_TIME = 1000 / FPS
)

var (
	basepath string
	rootpath string
)

func init() {
	_, exepath, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(exepath)
	rootpath = filepath.Dir(filepath.Dir(basepath))
}

type Game struct {
	ticksLastFrame uint64
	running        bool
	window         *sdl.Window
	renderer       *sdl.Renderer
	event          sdl.Event
	manager        *EntityManager
	assetManager   *AssetManager
}

func (g *Game) Initialize(width, height int) error {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("failed to initialize SDL: %s", err)
	}

	g.window, err = sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height), sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return fmt.Errorf("failed to create window: %s", err)
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %s", err)
	}

	g.manager = &EntityManager{renderer: g.renderer, event: &g.event}
	g.assetManager = &AssetManager{manager: g.manager, textures: make(map[string]*sdl.Texture)}

	if err = g.LoadLevel(0); err != nil {
		return err
	}

	g.running = true

	return nil
}

func (g *Game) LoadLevel(levelNumber int) error {
	/* Start including new assets to the assetmanager list */
	if err := g.assetManager.AddTexture("tank-image", filepath.Join(rootpath, "assets/images/tank-big-right.png")); err != nil {
		return err
	}

	if err := g.assetManager.AddTexture("chopper-image", filepath.Join(rootpath, "assets/images/chopper-spritesheet.png")); err != nil {
		return err
	}

	if err := g.assetManager.AddTexture("radar-image", filepath.Join(rootpath, "assets/images/radar.png")); err != nil {
		return err
	}

	if err := g.assetManager.AddTexture("jungle-tiletexture", filepath.Join(rootpath, "assets/tilemaps/jungle.png")); err != nil {
		return err
	}

	m := Map{g.manager, g.assetManager.GetTexture("jungle-tiletexture"), 2, 32}
	if err := m.LoadMap(filepath.Join(rootpath, "assets/tilemaps/jungle.map"), 25, 20); err != nil {
		return err
	}

	/* Start including entities and also components to them */
	chopperEntity := g.manager.AddEntity("chopper")
	chopperEntity.AddComponent(NewTransformComponent(Vec2{240, 106}, Vec2{0, 0}, 32, 32, 1), TRANSFORM_COMPONENT)
	chopperEntity.AddComponent(NewSpriteComponent2(g.assetManager.GetTexture("chopper-image"), 2, 90, true, false), SPRITE_COMPONENT)
	chopperEntity.AddComponent(NewKeyboardControlComponent("Up", "Right", "Down", "Left", "Space"), KEYBOARD_CONTROL_COMPONENT)

	tankEntity := g.manager.AddEntity("tank")
	tankEntity.AddComponent(NewTransformComponent(Vec2{150, 495}, Vec2{5, 0}, 32, 32, 1), TRANSFORM_COMPONENT)
	tankEntity.AddComponent(NewSpriteComponent(g.assetManager.GetTexture("tank-image")), SPRITE_COMPONENT)

	radarEntity := g.manager.AddEntity("Radar")
	radarEntity.AddComponent(NewTransformComponent(Vec2{720, 15}, Vec2{0, 0}, 64, 64, 1), TRANSFORM_COMPONENT)
	radarEntity.AddComponent(NewSpriteComponent2(g.assetManager.GetTexture("radar-image"), 8, 150, false, true), SPRITE_COMPONENT)

	return nil
}

func (g *Game) ProcessInput() {
	g.event = sdl.PollEvent()
	if g.event != nil {
		switch t := g.event.(type) {
		case *sdl.QuitEvent:
			g.running = false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				g.running = false
			}
		}
	}
}

func (g *Game) Update() {
	// Sleep the execution until we reach the target frame time in milliseconds
	timeToWait := FRAME_TARGET_TIME - (sdl.GetTicks64() - g.ticksLastFrame)

	// Only call delay if we are too fast to process this frame
	if timeToWait > 0 && timeToWait <= FRAME_TARGET_TIME {
		sdl.Delay(uint32(timeToWait))
	}

	// Delta time is the difference in ticks from last frame converted to seconds
	var deltaTime float64 = float64(sdl.GetTicks64()-g.ticksLastFrame) / 1000.0

	// Clamp deltaTime to a maximum value
	if deltaTime > 0.05 {
		deltaTime = 0.05
	}

	// Sets the new ticks for the current frame to be used in the next pass
	g.ticksLastFrame = sdl.GetTicks64()

	g.manager.Update(deltaTime)
}

func (g *Game) Render() {
	g.renderer.SetDrawColor(21, 21, 21, 255)
	g.renderer.Clear()

	if g.manager.HasNoEntities() {
		return
	}

	g.manager.Render()

	g.renderer.Present()
}

func (g *Game) Destory() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

func (g Game) IsRunning() bool {
	return g.running
}
