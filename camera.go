package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	x, y float32
}

type Camera struct {
	canvas      []rl.Color
	fovAngle    float32
	aspectRatio float32
	fovScaling  float32
	zNear       float32

	width, height         int
	halfWidth, halfHeight float32
}

func NewCamera(w, h int, zNear, fovAngle float32) Camera {
	return Camera{
		fovAngle:    fovAngle,
		fovScaling:  FovScaling(fovAngle),
		width:       w,
		height:      h,
		halfWidth:   float32(w) / 2,
		halfHeight:  float32(h) / 2,
		canvas:      make([]rl.Color, w*h),
		aspectRatio: float32(w) / float32(h),
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
	return Point{
		x: (v.X * c.fovScaling) / (v.Z * c.aspectRatio),
		y: (v.Y * c.fovScaling) / v.Z,
	}
}

func (c Camera) PutPixel(p Point) {
	x := int((p.x + 1) * c.halfWidth)
	y := int((1 - p.y) * c.halfHeight)

	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}

	c.canvas[y*c.width+x] = rl.Black
}
