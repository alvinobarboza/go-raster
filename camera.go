package main

import (
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
	// TODO: Calculate d = 1/tan(fovAngleRad/2)
	// Near is another variable can be zNear
	return Point{
		x: v.X / v.Z,
		y: v.Y * c.aspectRation / v.Z,
	}
}

func (c Camera) PutPixel(p Point) {
	x := int((p.x + 1) * 0.5 * float32(c.width))
	y := int((1 - p.y) * 0.5 * float32(c.height))

	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}

	c.canvas[y*c.width+x] = rl.Black
}
