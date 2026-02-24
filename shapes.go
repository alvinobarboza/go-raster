package main

import "image/color"

func NewColor(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

var (
	LightGray  = NewColor(200, 200, 200, 255)
	Gray       = NewColor(130, 130, 130, 255)
	DarkGray   = NewColor(80, 80, 80, 255)
	Yellow     = NewColor(253, 249, 0, 255)
	Gold       = NewColor(255, 203, 0, 255)
	Orange     = NewColor(255, 161, 0, 255)
	Pink       = NewColor(255, 109, 194, 255)
	Red        = NewColor(230, 41, 55, 255)
	Maroon     = NewColor(190, 33, 55, 255)
	Green      = NewColor(0, 228, 48, 255)
	Lime       = NewColor(0, 158, 47, 255)
	DarkGreen  = NewColor(0, 117, 44, 255)
	SkyBlue    = NewColor(102, 191, 255, 255)
	Blue       = NewColor(0, 121, 241, 255)
	DarkBlue   = NewColor(0, 82, 172, 255)
	Purple     = NewColor(200, 122, 255, 255)
	Violet     = NewColor(135, 60, 190, 255)
	DarkPurple = NewColor(112, 31, 126, 255)
	Beige      = NewColor(211, 176, 131, 255)
	Brown      = NewColor(127, 106, 79, 255)
	DarkBrown  = NewColor(76, 63, 47, 255)
	White      = NewColor(255, 255, 255, 255)
	Black      = NewColor(0, 0, 0, 255)
	Blank      = NewColor(0, 0, 0, 0)
	Magenta    = NewColor(255, 0, 255, 255)
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
	m := NewMesh(vertsCube, nil, nil, trisCube, nil)
	return NewModel(&m, NewTransforms(pos, scale, rotation))
}

func NewTriangle(pos, scale, rotation Vec3) Model {
	uv_verts := []Vec2{
		{X: 0, Y: 0}, // 0 front bottom left
		{X: 0, Y: 1}, // 1 front top left
		{X: 1, Y: 1}, // 2 front top right
	}
	tris_verts := []Vec3{
		{X: -1, Y: -1, Z: 0}, // 0 front bottom left
		{X: -1, Y: 1, Z: 0},  // 1 front top left
		{X: 1, Y: 1, Z: 0},   // 2 front top right
	}
	tris_tris := []Triangle{
		{
			v1: 0, v2: 1, v3: 2,
			u1: 0, u2: 1, u3: 2,
			n1: 0, n2: 0, n3: 0,
			color: Red,
		},
	}
	tris_normal := []Vec3{
		{Z: -1},
	}

	m := NewMesh(tris_verts, tris_normal, uv_verts, tris_tris, nil)
	return NewModel(&m, NewTransforms(pos, scale, rotation))
}
