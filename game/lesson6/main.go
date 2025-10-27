package main

import "engine/lesson6/engine"

const (
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 600
)

func main() {
	game := engine.Game{}
	if err := game.Initialize(WINDOW_WIDTH, WINDOW_HEIGHT); err != nil {
		panic(err)
	}

	for game.IsRunning() {
		game.ProcessInput()
		game.Update()
		game.Render()
	}

	defer game.Destory()
}
