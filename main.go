package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

func main() {

	camera := NewCamera(ScreenWidth/4, ScreenHeight/4, 1)

	plane := [8]rl.Vector3{
		rl.NewVector3(1.0, 1.0, 3.0),   // 0 front top right
		rl.NewVector3(-1.0, 1.0, 3.0),  // 1 front top left
		rl.NewVector3(-1.0, -1.0, 3.0), // 2 front bottom left
		rl.NewVector3(1.0, -1.0, 3.0),  // 3 front bottom rigth
		rl.NewVector3(1.0, 1.0, 2.0),   // 4 back top right
		rl.NewVector3(-1.0, 1.0, 2.0),  // 5 back top left
		rl.NewVector3(-1.0, -1.0, 2.0), // 6 back bottom left
		rl.NewVector3(1.0, -1.0, 2.0),  // 7 back bottom right
	}

	// rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(ScreenWidth, ScreenHeight, "go-raster")
	defer rl.CloseWindow()

	img := rl.GenImageColor(ScreenWidth/4, ScreenHeight/4, rl.RayWhite)
	defer rl.UnloadImage(img)

	renderTexture := rl.LoadTextureFromImage(img)
	defer rl.UnloadTexture(renderTexture)

	backForward := 0.0
	leftRight := 0.0
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		camera.ClearCanvas()

		backForward = 0
		leftRight = 0

		if rl.IsKeyDown(rl.KeyUp) {
			backForward = .1
		}

		if rl.IsKeyDown(rl.KeyDown) {
			backForward = -.1
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			leftRight = -.1
		}

		if rl.IsKeyDown(rl.KeyRight) {
			leftRight = .1
		}

		for i := range plane {
			plane[i].X += float32(leftRight)
			plane[i].Z += float32(backForward)
			p := camera.ProjectVertex(plane[i])
			camera.PutPixel(p)
		}

		rl.UpdateTexture(renderTexture, camera.canvas)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexturePro(
			renderTexture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(camera.width), Height: float32(camera.height)},
			rl.Rectangle{X: 0, Y: 0, Width: float32(ScreenWidth), Height: float32(ScreenHeight)},
			rl.Vector2Zero(),
			0,
			rl.White,
		)

		rl.DrawText("raster", 10, 15, 20, rl.Black)
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
	}
}
