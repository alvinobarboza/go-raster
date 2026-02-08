package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 900
	ScreenHeight = 900

	NearPlane = 0.2
)

func main() {

	camera := NewCamera(ScreenWidth/4, ScreenHeight/4, NearPlane, 90)
	scene := NewScene(&camera)

	cube := NewMesh([]Vec3{
		NewVec3(1.0, 1.0, 2.0),   // 0 front top right
		NewVec3(-1.0, 1.0, 2.0),  // 1 front top left
		NewVec3(-1.0, -1.0, 2.0), // 2 front bottom left
		NewVec3(1.0, -1.0, 2.0),  // 3 front bottom rigth
		NewVec3(1.0, 1.0, 3.0),   // 4 back top right
		NewVec3(-1.0, 1.0, 3.0),  // 5 back top left
		NewVec3(-1.0, -1.0, 3.0), // 6 back bottom left
		NewVec3(1.0, -1.0, 3.0),  // 7 back bottom right
	}, rl.Black)

	scene.AddMesh(
		cube,
	)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(ScreenWidth, ScreenHeight, "go-raster")
	defer rl.CloseWindow()

	img := rl.GenImageColor(camera.width, camera.height, rl.RayWhite)
	defer rl.UnloadImage(img)

	renderTexture := rl.LoadTextureFromImage(img)
	defer rl.UnloadTexture(renderTexture)

	backForward := 0.0
	leftRight := 0.0
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		if rl.IsWindowResized() {
			camera.UpdateCanvasSize(rl.GetScreenWidth()/4, rl.GetScreenHeight()/4)
			rl.ImageResize(img, int32(camera.width), int32(camera.height))
			rl.UnloadTexture(renderTexture)
			renderTexture = rl.LoadTextureFromImage(img)
		}

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

		for i := range cube.vertices {
			cube.vertices[i].X += float32(leftRight)
			cube.vertices[i].Z += float32(backForward)
		}
		scene.Render()

		rl.UpdateTexture(renderTexture, scene.activeCam.canvas)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexturePro(
			renderTexture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(camera.width), Height: float32(camera.height)},
			rl.Rectangle{X: 0, Y: 0, Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())},
			rl.Vector2Zero(),
			0,
			rl.White,
		)

		rl.DrawText("raster", 10, 15, 20, rl.Black)
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
	}
}
