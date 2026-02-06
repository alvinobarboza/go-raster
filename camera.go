package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	x, y float32
}

type Camera struct {
	canvas       []rl.Color
	viewDistance float32
	width        int
	height       int
	aspectRation float32
}

func NewCamera(w, h int, d float32) Camera {
	return Camera{
		viewDistance: d,
		width:        w,
		height:       h,
		canvas:       make([]rl.Color, w*h),
		aspectRation: float32(w) / float32(h),
	}
}

func (c *Camera) UpdateCanvasSize(w, h int) {
	c.width = w
	c.height = h
	c.canvas = make([]rl.Color, w*h)
}

func (c Camera) ClearCanvas() {
	for i := range len(c.canvas) {
		c.canvas[i] = rl.RayWhite
	}
}

func (c Camera) ProjectVertex(v rl.Vector3) Point {

	// pvx = px / pz
	// pvy = py * ar / pz

	x := v.X / v.Z
	y := v.Y * c.aspectRation / v.Z

	return Point{
		x: x,
		y: y,
	}
}

func (c Camera) PutPixel(p Point) {
	// psx = (pvx + 1) * (0.5*ScreenWidth - 0.5)
	// psy = -py * (0.5*ScreenHeight - 0.5) + (0.5*ScreenHeight + 0.5)

	x := int((p.x + 1) * (float32(c.width)*0.5 - .5))
	y := int(-p.y*(.5*float32(c.height)-.5) + (0.5*float32(c.height) + .5))

	fmt.Println(p, x, y)

	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}

	c.canvas[y*c.width+x] = rl.Black
}
