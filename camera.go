package main

import (
	"image/color"
)

type NDCPoint struct {
	X, Y  float32
	color color.RGBA
}

type ScreenPoint struct {
	X, Y  int
	color color.RGBA
}

type Camera struct {
	canvas      []color.RGBA
	fovAngle    float32
	aspectRatio float32
	fovScaling  float32
	zNear       float32

	width, height         int
	halfWidth, halfHeight float32
}

func NewCamera(w, h int, zNear, fovAngle float32) Camera {
	c := Camera{
		fovAngle:   fovAngle,
		fovScaling: FovScaling(fovAngle),
	}

	c.UpdateCanvasSize(w, h)
	return c
}

func (c *Camera) UpdateCanvasSize(w, h int) {
	c.width = w
	c.height = h
	c.halfWidth = float32(w) / 2
	c.halfHeight = float32(h) / 2
	c.aspectRatio = float32(w) / float32(h)
	c.canvas = make([]color.RGBA, w*h)
}

func (c Camera) ClearCanvas() {
	for i := range len(c.canvas) {
		c.canvas[i].R = 240
		c.canvas[i].G = 240
		c.canvas[i].B = 240
		c.canvas[i].A = 255
	}
}

func (c Camera) ProjectVertexToNDC(v Vec3, cl color.RGBA) NDCPoint {
	zXInverse := 1 / (v.Z * c.aspectRatio)
	zYInverse := 1 / v.Z
	return NDCPoint{
		X:     (v.X * c.fovScaling) * zXInverse,
		Y:     (v.Y * c.fovScaling) * zYInverse,
		color: cl,
	}
}

func (c Camera) NDCtoScreen(p NDCPoint) ScreenPoint {
	x := int((p.X + 1) * c.halfWidth)
	y := int((1 - p.Y) * c.halfHeight)

	return ScreenPoint{
		X:     x,
		Y:     y,
		color: p.color,
	}
}

func (c Camera) PutPixel(p ScreenPoint) {
	if p.X < 0 || p.X >= c.width || p.Y < 0 || p.Y >= c.height {
		return
	}
	c.canvas[p.Y*c.width+p.X] = p.color
}
