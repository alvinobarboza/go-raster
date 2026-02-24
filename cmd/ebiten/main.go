package main

import (
	"fmt"
	"log"

	// "net/http"
	// _ "net/http/pprof"
	"unsafe"

	"github.com/alvinobarboza/go-raster/internal/camera"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/renderer"
	"github.com/alvinobarboza/go-raster/internal/scene"
	"github.com/alvinobarboza/go-raster/internal/transforms"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	NearPlane    = 0.2
	FarPlane     = 25
)

type Game struct {
	camera        *camera.Camera
	scene         *scene.Scene
	renderer      *renderer.Renderer
	models        []mesh.Model
	factor        int
	cursorEnabled bool
	prevMouseX    int
	prevMouseY    int
	op            *ebiten.DrawImageOptions
	canvasImage   *ebiten.Image
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.camera.ToggleViewLock()
		g.cursorEnabled = !g.cursorEnabled
		if g.cursorEnabled {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.models[1].Transforms.Scale = g.models[1].Transforms.Scale.Scale(0.99)
		g.models[1].UpdateTransforms()
		g.models[1].Transforms.Scale.Print("scale down")
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			g.models[1].Transforms.Rotation.Y -= .5
		} else {
			g.models[1].Transforms.Rotation.Y += .5
		}
		g.models[1].UpdateTransforms()
		g.models[1].Transforms.Rotation.Print("rotate")
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.models[1].Transforms.Scale = g.models[1].Transforms.Scale.Scale(1.01)
		g.models[1].UpdateTransforms()
		g.models[1].Transforms.Scale.Print("scale up")
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.models[1].Transforms.Position.X -= 0.01
		g.models[1].UpdateTransforms()
		g.models[1].Transforms.Position.Print("position")
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.models[1].Transforms.Position.X += 0.01
		g.models[1].UpdateTransforms()
		g.models[1].Transforms.Position.Print("position")
	}

	var backForwardCam, leftRightCam, upDownCam float32

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		backForwardCam = .02
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		backForwardCam = -.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		leftRightCam = .02
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		leftRightCam = -.02
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		upDownCam = .02
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		upDownCam = -.02
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.camera.ToggleDepthRender()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.camera.ToggleWireRender()
	}

	cx, cy := ebiten.CursorPosition()
	dx := float32(cx - g.prevMouseX)
	dy := float32(cy - g.prevMouseY)
	g.prevMouseX, g.prevMouseY = cx, cy

	if !g.cursorEnabled {
		dt := float32(1.0 / 60.0) // Ebiten updates at a fixed 60 TPS by default
		g.camera.UpdateRotation(dx*dt*0.4, dy*dt*0.4)
	}

	if backForwardCam != 0 || leftRightCam != 0 || upDownCam != 0 {
		g.camera.MoveBackForwad(backForwardCam)
		g.camera.MoveSideways(leftRightCam)
		g.camera.MoveVetically(upDownCam)
	}

	g.renderer.Render()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	canvas := g.scene.ActiveCam.Canvas

	if len(canvas) > 0 {
		// Cast the slice directly to avoid GC overhead
		byteData := unsafe.Slice((*byte)(unsafe.Pointer(&canvas[0])), len(canvas)*4)
		g.canvasImage.WritePixels(byteData)
	}

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	cw, ch := g.canvasImage.Bounds().Dx(), g.canvasImage.Bounds().Dy()

	if g.factor == 1 {
		screen.DrawImage(g.canvasImage, nil)
	} else {
		g.op.GeoM.Scale(float64(sw)/float64(cw), float64(sh)/float64(ch))
		screen.DrawImage(g.canvasImage, g.op)
	}

	ebitenutil.DebugPrintAt(screen, "raster", sw-70, sh-20)

	debugInfo := fmt.Sprintf(
		"FPS: %0.2f\n"+
			"Cam:\nX:%01f\nY:%01f\nZ:%01f\n\n"+
			"Move: A/W/S/D\n"+
			"Mouse view movement:\nTab to Lock/Unlock\n"+
			"Depth toggle: Z\n"+
			"Wireframe toggle: X",
		ebiten.ActualFPS(),
		g.camera.Transforms.Position.X,
		g.camera.Transforms.Position.Y,
		g.camera.Transforms.Position.Z,
	)

	ebitenutil.DebugPrintAt(screen, debugInfo, 10, 10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Handled natively by Ebiten, replaces your rl.IsWindowResized() logic.
	newW := uint(outsideWidth / g.factor)
	newH := uint(outsideHeight / g.factor)

	if newW != g.camera.Width || newH != g.camera.Height {
		g.camera.UpdateCanvasSize(newW, newH)
		g.canvasImage = ebiten.NewImage(int(newW), int(newH))
	}

	return outsideWidth, outsideHeight
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	factor := 1
	cam := camera.NewCamera(
		uint(ScreenWidth/factor),
		uint(ScreenHeight/factor),
		float32(10),
		NearPlane,
		FarPlane,
		float32(53),
		transforms.NewVec3(0, 2, -0.7),
		transforms.NewVec3(-30, 0, 0),
	)

	s := scene.NewScene(cam)
	wp := renderer.NewWorkerPool(4)
	rndr := renderer.NewRenderer(wp)
	rndr.AddActiveScene(s)

	models, err := scene.LoadSceneFromJSON("./scene.json")
	if err != nil {
		panic(err)
	}

	for i := range models {
		s.AddMesh(&models[i])
	}

	g := &Game{
		camera:        cam,
		scene:         s,
		renderer:      rndr,
		models:        models,
		factor:        factor,
		cursorEnabled: false,
		op:            &ebiten.DrawImageOptions{},
		canvasImage:   ebiten.NewImage(int(cam.Width), int(cam.Height)),
	}

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(false)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("go-raster")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
