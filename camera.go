package main

import (
	"image/color"
)

type Point struct {
	x, y  float32
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

func (c Camera) ProjectVertex(v Vec3) Point {
	zXInverse := 1 / (v.Z * c.aspectRatio)
	zYInverse := 1 / v.Z
	return Point{
		x: (v.X * c.fovScaling) * zXInverse,
		y: (v.Y * c.fovScaling) * zYInverse,
	}
}

func (c Camera) NDCtoScreen(p Point) (x int, y int) {
	x = int((p.x + 1) * c.halfWidth)
	y = int((1 - p.y) * c.halfHeight)

	return x, y
}

func (c Camera) PutPixel(p Point) {
	x, y := c.NDCtoScreen(p)

	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}

	c.canvas[y*c.width+x] = p.color
}
