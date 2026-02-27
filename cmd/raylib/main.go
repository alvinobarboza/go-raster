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
	ScreenWidth  = 1920
	ScreenHeight = 1080

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
		transforms.NewVec3(0.43, 1.4, 0.6),
		transforms.NewVec3(-40, 40, 0),
	)

	s := scene.NewScene(camera)

	models, err := scene.LoadSceneFromJSON("./scene.json")
	if err != nil {
		panic(err)
	}

	for i := range models {
		s.AddMesh(&models[i])
	}

	renderer := renderer.NewRenderer(8, 200)
	renderer.AddActiveScene(s)

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

	rl.SetTargetFPS(40)
	rl.DisableCursor()
	cursorEnabled := false

	for !rl.WindowShouldClose() {

		if rl.IsWindowResized() {
			camera.UpdateCanvasSize(uint(rl.GetScreenWidth()/factor), uint(rl.GetScreenHeight()/factor))
			rl.ImageResize(img, int32(camera.Width), int32(camera.Height))
			rl.UnloadTexture(renderTexture)
			renderTexture = rl.LoadTextureFromImage(img)
			renderer.UpdateTiles()
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

		if rl.IsKeyPressed(rl.KeyT) {
			renderer.ToggleTriangleBoundaryRender()
		}

		if rl.IsKeyPressed(rl.KeyM) {
			renderer.ToggleMultithreaded()
		}

		if rl.IsKeyPressed(rl.KeyB) {
			renderer.ToggleTileBoundaryRender()
		}

		if rl.IsKeyPressed(rl.KeyY) {
			renderer.IncrementTileSize(20)
		}

		if rl.IsKeyPressed(rl.KeyU) {
			renderer.IncrementTileSize(-20)
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

		w := int32(420)
		h := int32(370)
		rl.DrawRectangle(2, 2, w, h, rl.Fade(rl.Black, 0.6))
		rl.DrawRectangleLines(2, 2, w, h, shapes.LightGray)

		textSize := int32(20)
		yOffSet := int32(10)
		gap := int32(30)

		rl.DrawFPS(10, yOffSet)
		yOffSet += gap
		rl.DrawText(
			fmt.Sprintf(
				"Cam:\n X: %02.2f Y: %02.2f Z: %02.2f\n X: %02.2f' Y: %02.2f' Z: %02.2f'",
				camera.Transforms.Position.X,
				camera.Transforms.Position.Y,
				camera.Transforms.Position.Z,
				camera.Transforms.Rotation.X,
				camera.Transforms.Rotation.Y,
				camera.Transforms.Rotation.Z),
			10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap * 3

		rl.DrawText("Move: A/W/S/D", 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Lock/Unlock mouse: Tab = status: %v", !cursorEnabled), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Depth toggle: Z = status: %v", camera.RenderDepth), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Wireframe toggle: X = status: %v", camera.RenderWire), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Triangle AABB: T = status: %v", renderer.RenderTriangleBoundaries), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Render Tiles: B = status: %v", renderer.RenderTileBoundaries), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("Multithread: M = status: %v", renderer.RenderMultithreaded), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap
		rl.DrawText(fmt.Sprintf("(De)Increase tile size: Y/U = size: %v", renderer.TileSize()), 10, yOffSet, textSize, shapes.LightGray)
		yOffSet += gap

		rl.EndDrawing()
	}
}
