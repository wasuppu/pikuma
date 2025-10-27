package main

import "engine/lesson9/engine"

func main() {
	game := engine.Game{}
	if err := game.Initialize(); err != nil {
		panic(err)
	}

	for game.IsRunning() {
		game.ProcessInput()
		game.Update()
		game.Render()
	}

	defer game.Destory()
}
