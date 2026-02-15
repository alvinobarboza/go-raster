package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720

	NearPlane = 0.2
	FarPlane  = 25
)

func main() {

	// Profiling
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	factor := 1
	sensitivity := float32(10)
	fov := float32(53)
	camera := NewCamera(
		uint(ScreenWidth/factor),
		uint(ScreenHeight/factor),
		sensitivity,
		NearPlane,
		FarPlane,
		fov,
		NewVec3(0, 0.03, -4),
		NewVec3(0, 0, 0),
	)

	scene := NewScene(&camera)

	cubeMesh, err := LoadMeshFromFile("./assets/cube.obj", "./assets/default.jpg")
	if err != nil {
		panic(err)
	}
	cube := NewModel(
		&cubeMesh, NewTransforms(NewVec3(0, 0, 4), NewVec3(16/8, 1, 9/8), NewVec3(0, 0, 0)))

	utahTeapotMesh, err := LoadMeshFromFile("./assets/utah-assets/utah_teapot.obj", "./assets/default.jpg")
	if err != nil {
		panic(err)
	}

	utahTeapot := NewModel(
		&utahTeapotMesh, NewTransforms(NewVec3(5, 0, 4), NewVec3(1, 1, 1), NewVec3(0, 0, 0)))

	triangle := NewTriangle(NewVec3(0, 0, 20.9), NewVec3(1, 1, 1), NewVec3(0, 0, 0))
	triangle.mesh.texture = LoadDefaultTexture()

	scene.AddMesh(&triangle)
	scene.AddMesh(&cube)
	scene.AddMesh(&utahTeapot)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(ScreenWidth, ScreenHeight, "go-raster")
	defer rl.CloseWindow()

	img := rl.GenImageColor(int(camera.width), int(camera.height), rl.RayWhite)
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
			camera.UpdateCanvasSize(uint(rl.GetScreenWidth()/factor), uint(rl.GetScreenHeight()/factor))
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
			triangle.transforms.position.X += leftRight
			triangle.transforms.position.Z += backForward
			triangle.UpdateTransforms()
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

		rl.DrawText("RASTER", int32(rl.GetScreenWidth()-140), int32(rl.GetScreenHeight()-20), 20, rl.Gray)

		rl.DrawRectangle(2, 2, 305, 230, rl.Fade(rl.DarkGray, 0.6))
		rl.DrawRectangleLines(2, 2, 305, 230, rl.Gray)

		rl.DrawFPS(10, 10)
		rl.DrawText(
			fmt.Sprintf(
				"Cam: \nX:%01f \nY:%01f \nZ:%01f",
				camera.transforms.position.X,
				camera.transforms.position.Y,
				camera.transforms.position.Z),
			10, 50, 20, rl.White)
		rl.DrawText("Move: A/W/S/D",
			10, 140, 20, rl.White)
		rl.DrawText("Mouse view moviment: \nTab to Lock/Unlock",
			10, 162, 20, rl.White)

		rl.EndDrawing()
	}
}
