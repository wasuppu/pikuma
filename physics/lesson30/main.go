package main

import "log"

func main() {
	app := Application{}
	if err := app.Setup(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer app.Destory()

	for app.IsRunning() {
		app.Input()
		app.Update()
		app.Render()
	}
}
