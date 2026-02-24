package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/alvinobarboza/go-raster/internal/camera"
	"github.com/alvinobarboza/go-raster/internal/renderer"
	"github.com/alvinobarboza/go-raster/internal/scene"
	"github.com/alvinobarboza/go-raster/internal/shapes"
	"github.com/alvinobarboza/go-raster/internal/transforms"
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
	camera := camera.NewCamera(
		uint(ScreenWidth/factor),
		uint(ScreenHeight/factor),
		sensitivity,
		NearPlane,
		FarPlane,
		fov,
		transforms.NewVec3(0, 2, -0.7),
		transforms.NewVec3(-30, 0, 0),
	)

	s := scene.NewScene(camera)

	wp := renderer.NewWorkerPool(4)
	renderer := renderer.NewRenderer(wp)

	renderer.AddActiveScene(s)

	models, err := scene.LoadSceneFromJSON("./scene.json")
	if err != nil {
		panic(err)
	}

	for i := range models {
		s.AddMesh(&models[i])
	}

	// triangle := NewTriangle(NewVec3(0, 0, 2), NewVec3(1, 1, 1), NewVec3(1, 1, 1))
	// triangle.mesh.texture = LoadDefaultTexture()
	// scene.AddMesh(&triangle)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(ScreenWidth, ScreenHeight, "go-raster")
	defer rl.CloseWindow()

	img := rl.GenImageColor(int(camera.Width), int(camera.Height), rl.RayWhite)
	defer rl.UnloadImage(img)

	renderTexture := rl.LoadTextureFromImage(img)
	defer rl.UnloadTexture(renderTexture)

	backForwardCam := float32(0.0)
	leftRightCam := float32(0.0)
	upDownCam := float32(0.0)

	rl.SetTargetFPS(60)
	rl.DisableCursor()
	cursorEnabled := false

	for !rl.WindowShouldClose() {

		if rl.IsWindowResized() {
			camera.UpdateCanvasSize(uint(rl.GetScreenWidth()/factor), uint(rl.GetScreenHeight()/factor))
			rl.ImageResize(img, int32(camera.Width), int32(camera.Height))
			rl.UnloadTexture(renderTexture)
			renderTexture = rl.LoadTextureFromImage(img)
		}

		if rl.IsKeyPressed(rl.KeyTab) {
			camera.ToggleViewLock()
			if cursorEnabled {
				rl.DisableCursor()
			} else {
				rl.EnableCursor()
			}
			cursorEnabled = !cursorEnabled
		}

		backForwardCam = 0
		leftRightCam = 0
		upDownCam = 0

		if rl.IsKeyDown(rl.KeyQ) {
			models[1].Transforms.Scale = models[1].Transforms.Scale.Scale(0.99)
			models[1].UpdateTransforms()
			models[1].Transforms.Scale.Print("scale down")
		}

		if rl.IsKeyDown(rl.KeyR) {
			if rl.IsKeyDown(rl.KeyLeftShift) {
				models[1].Transforms.Rotation.Y -= .5
			} else {
				models[1].Transforms.Rotation.Y += .5
			}
			models[1].UpdateTransforms()
			models[1].Transforms.Rotation.Print("rotate")
		}

		if rl.IsKeyDown(rl.KeyE) {
			models[1].Transforms.Scale = models[1].Transforms.Scale.Scale(1.01)
			models[1].UpdateTransforms()
			models[1].Transforms.Scale.Print("scale up")
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			models[1].Transforms.Position.X -= 0.01
			models[1].UpdateTransforms()
			models[1].Transforms.Position.Print("position")
		}

		if rl.IsKeyDown(rl.KeyRight) {
			models[1].Transforms.Position.X += 0.01
			models[1].UpdateTransforms()
			models[1].Transforms.Position.Print("position")
		}

		if rl.IsKeyDown(rl.KeyW) {
			backForwardCam = .02
		}

		if rl.IsKeyDown(rl.KeyS) {
			backForwardCam = -.02
		}

		if rl.IsKeyDown(rl.KeyA) {
			leftRightCam = .02
		}

		if rl.IsKeyDown(rl.KeyD) {
			leftRightCam = -.02
		}

		if rl.IsKeyDown(rl.KeySpace) {
			upDownCam = .02
		}

		if rl.IsKeyDown(rl.KeyLeftControl) {
			upDownCam = -.02
		}

		if rl.IsKeyPressed(rl.KeyZ) {
			camera.ToggleDepthRender()
		}

		if rl.IsKeyPressed(rl.KeyX) {
			camera.ToggleWireRender()
		}

		mouseDelta := rl.GetMouseDelta()
		camera.UpdateRotation(mouseDelta.X*rl.GetFrameTime()*0.4, mouseDelta.Y*rl.GetFrameTime()*0.4)

		if backForwardCam != 0 || leftRightCam != 0 || upDownCam != 0 {
			camera.MoveBackForwad(backForwardCam)
			camera.MoveSideways(leftRightCam)
			camera.MoveVetically(upDownCam)
		}

		renderer.Render()

		rl.UpdateTexture(renderTexture, s.ActiveCam.Canvas)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexturePro(
			renderTexture,
			rl.Rectangle{X: 0, Y: 0, Width: float32(camera.Width), Height: float32(camera.Height)},
			rl.Rectangle{X: 0, Y: 0, Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())},
			rl.Vector2Zero(),
			0,
			rl.White,
		)

		rl.DrawText("raster", int32(rl.GetScreenWidth()-70), int32(rl.GetScreenHeight()-20), 20, shapes.Black)

		rl.DrawRectangle(2, 2, 305, 250, rl.Fade(rl.DarkGray, 0.6))
		rl.DrawRectangleLines(2, 2, 305, 250, rl.Gray)

		rl.DrawFPS(10, 10)
		rl.DrawText(
			fmt.Sprintf(
				"Cam: \nX:%01f \nY:%01f \nZ:%01f",
				camera.Transforms.Position.X,
				camera.Transforms.Position.Y,
				camera.Transforms.Position.Z),
			10, 50, 20, rl.White)
		rl.DrawText("Move: A/W/S/D",
			10, 140, 20, rl.White)
		rl.DrawText("Mouse view moviment: \nTab to Lock/Unlock",
			10, 163, 20, rl.White)
		rl.DrawText("Depth toggle: Z", 10, 204, 20, rl.White)
		rl.DrawText("Wireframe toggle: X", 10, 224, 20, rl.White)

		rl.EndDrawing()
	}
}
