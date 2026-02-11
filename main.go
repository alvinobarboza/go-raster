package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768

	NearPlane = 0.2
)

func main() {

	sensitivity := float32(20)
	fov := float32(90)
	camera := NewCamera(
		ScreenWidth/4,
		ScreenHeight/4,
		sensitivity,
		NearPlane,
		fov,
		NewVec3(0, 0, -5),
		NewVec3(0, 0, 0),
	)

	cubeMesh, err := LoadMeshFromFile("./assets/cube.obj", "./assets/default.jpg")
	if err != nil {
		panic(err)
	}

	utahTeapotMesh, err := LoadMeshFromFile("./assets/utah-assets/utah_teapot.obj", "./assets/default.jpg")
	if err != nil {
		panic(err)
	}
	cube := NewModel(
		&cubeMesh, NewTransforms(NewVec3(0, 0, 4), NewVec3(1, 1, 1), NewVec3(0, 0, 0)))

	utahTeapot := NewModel(
		&utahTeapotMesh, NewTransforms(NewVec3(5, 0, 4), NewVec3(1, 1, 1), NewVec3(0, 0, 0)))

	scene := NewScene(&camera)
	scene.AddMesh(&cube)
	scene.AddMesh(&utahTeapot)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(ScreenWidth, ScreenHeight, "go-raster")
	defer rl.CloseWindow()

	img := rl.GenImageColor(camera.width, camera.height, rl.RayWhite)
	defer rl.UnloadImage(img)

	renderTexture := rl.LoadTextureFromImage(img)
	defer rl.UnloadTexture(renderTexture)

	backForward := float32(0.0)
	leftRight := float32(0.0)

	backForwardCam := float32(0.0)
	leftRightCam := float32(0.0)
	upDownCam := float32(0.0)

	rl.SetTargetFPS(60)
	rl.DisableCursor()
	cursorEnabled := false

	for !rl.WindowShouldClose() {

		if rl.IsWindowResized() {
			camera.UpdateCanvasSize(rl.GetScreenWidth()/4, rl.GetScreenHeight()/4)
			rl.ImageResize(img, int32(camera.width), int32(camera.height))
			rl.UnloadTexture(renderTexture)
			renderTexture = rl.LoadTextureFromImage(img)
		}

		camera.ClearCanvas()

		if rl.IsKeyPressed(rl.KeyTab) {
			camera.ToggleViewLock()
			if cursorEnabled {
				rl.DisableCursor()
			} else {
				rl.EnableCursor()
			}
			cursorEnabled = !cursorEnabled
		}

		backForward = 0
		leftRight = 0

		backForwardCam = 0
		leftRightCam = 0
		upDownCam = 0

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

		if rl.IsKeyDown(rl.KeyW) {
			backForwardCam = .1
		}

		if rl.IsKeyDown(rl.KeyS) {
			backForwardCam = -.1
		}

		if rl.IsKeyDown(rl.KeyA) {
			leftRightCam = .1
		}

		if rl.IsKeyDown(rl.KeyD) {
			leftRightCam = -.1
		}

		if rl.IsKeyDown(rl.KeySpace) {
			upDownCam = .1
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			upDownCam = -.1
		}

		mouseDelta := rl.GetMouseDelta()
		camera.UpdateRotation(mouseDelta.X*rl.GetFrameTime(), mouseDelta.Y*rl.GetFrameTime())

		if backForward != 0 || leftRight != 0 {
			cube.transforms.position.X += leftRight
			cube.transforms.position.Z += backForward
			cube.UpdateTransforms()
		}

		if backForwardCam != 0 || leftRightCam != 0 || upDownCam != 0 {
			camera.MoveBackForwad(backForwardCam)
			camera.MoveSideways(leftRightCam)
			camera.MoveVetically(upDownCam)
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
		rl.DrawFPS(10, 0)
		rl.EndDrawing()
	}
}
