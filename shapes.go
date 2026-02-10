package main

import "image/color"

var (
	Red    = color.RGBA{R: 230, G: 41, B: 55, A: 255}
	Green  = color.RGBA{R: 0, G: 228, B: 48, A: 255}
	Blue   = color.RGBA{R: 0, G: 121, B: 241, A: 255}
	Yellow = color.RGBA{R: 253, G: 249, B: 0, A: 255}
	Purple = color.RGBA{R: 200, G: 122, B: 255, A: 255}
	Black  = color.RGBA{R: 50, G: 50, B: 50, A: 255}
)

// CUBE === begging
var vertsCube = []Vec3{
	{X: 1.0, Y: 1.0, Z: -1.0},   // 0 front top right
	{X: -1.0, Y: 1.0, Z: -1.0},  // 1 front top left
	{X: -1.0, Y: -1.0, Z: -1.0}, // 2 front bottom left
	{X: 1.0, Y: -1.0, Z: -1.0},  // 3 front bottom rigth
	{X: 1.0, Y: 1.0, Z: 1.0},    // 4 back top right
	{X: -1.0, Y: 1.0, Z: 1.0},   // 5 back top left
	{X: -1.0, Y: -1.0, Z: 1.0},  // 6 back bottom left
	{X: 1.0, Y: -1.0, Z: 1.0},   // 7 back bottom right
}

var trisCube = []Triangle{
	{v1: 0, v2: 1, v3: 2, color: Red},
	{v1: 0, v2: 2, v3: 3, color: Red},
	{v1: 4, v2: 0, v3: 3, color: Green},
	{v1: 4, v2: 3, v3: 7, color: Green},
	{v1: 5, v2: 4, v3: 7, color: Blue},
	{v1: 5, v2: 7, v3: 6, color: Blue},
	{v1: 1, v2: 5, v3: 6, color: Yellow},
	{v1: 1, v2: 6, v3: 2, color: Yellow},
	{v1: 4, v2: 5, v3: 1, color: Purple},
	{v1: 4, v2: 1, v3: 0, color: Purple},
	{v1: 2, v2: 6, v3: 7, color: color.RGBA{A: 255, R: 0, G: 255, B: 255}},
	{v1: 2, v2: 7, v3: 3, color: color.RGBA{A: 255, R: 0, G: 255, B: 255}},
}

// CUBE === end

func NewCube(pos, scale, rotation Vec3) Model {
	m := NewMesh(vertsCube, trisCube, nil)
	return NewModel(&m, Transforms{
		position: pos,
		scale:    scale,
		rotation: rotation,
	})
}
