package engine

import "github.com/veandco/go-sdl2/sdl"

func CheckRectangleCollision(rectangleA, rectangleB sdl.Rect) bool {
	return (rectangleA.X+rectangleA.W >= rectangleB.X &&
		rectangleB.X+rectangleB.W >= rectangleA.X &&
		rectangleA.Y+rectangleA.H >= rectangleB.Y &&
		rectangleB.Y+rectangleB.H >= rectangleA.Y)
}
