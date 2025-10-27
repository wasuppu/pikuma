package main

func main() {
	app := Application{}
	app.Setup()
	defer app.Destory()

	for app.IsRunning() {
		app.Input()
		app.Update()
		app.Render()
	}
}
