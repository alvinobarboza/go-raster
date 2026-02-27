package shapes

import (
	"image/color"

	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func NewColor(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

var (
	LighterGray = NewColor(200, 200, 200, 255)
	LightGray   = NewColor(180, 180, 180, 255)
	Gray        = NewColor(130, 130, 130, 255)
	DarkGray    = NewColor(80, 80, 80, 255)
	Yellow      = NewColor(253, 249, 0, 255)
	Gold        = NewColor(255, 203, 0, 255)
	Orange      = NewColor(255, 161, 0, 255)
	Pink        = NewColor(255, 109, 194, 255)
	Red         = NewColor(230, 41, 55, 255)
	Maroon      = NewColor(190, 33, 55, 255)
	Green       = NewColor(0, 228, 48, 255)
	Lime        = NewColor(0, 158, 47, 255)
	DarkGreen   = NewColor(0, 117, 44, 255)
	SkyBlue     = NewColor(102, 191, 255, 255)
	Blue        = NewColor(0, 121, 241, 255)
	DarkBlue    = NewColor(0, 82, 172, 255)
	Purple      = NewColor(200, 122, 255, 255)
	Violet      = NewColor(135, 60, 190, 255)
	DarkPurple  = NewColor(112, 31, 126, 255)
	Beige       = NewColor(211, 176, 131, 255)
	Brown       = NewColor(127, 106, 79, 255)
	DarkBrown   = NewColor(76, 63, 47, 255)
	White       = NewColor(255, 255, 255, 255)
	Black       = NewColor(0, 0, 0, 255)
	Blank       = NewColor(0, 0, 0, 0)
	Magenta     = NewColor(255, 0, 255, 255)
)

// CUBE === begging
var vertsCube = []transforms.Vec3{
	{X: 1.0, Y: 1.0, Z: -1.0},   // 0 front top right
	{X: -1.0, Y: 1.0, Z: -1.0},  // 1 front top left
	{X: -1.0, Y: -1.0, Z: -1.0}, // 2 front bottom left
	{X: 1.0, Y: -1.0, Z: -1.0},  // 3 front bottom rigth
	{X: 1.0, Y: 1.0, Z: 1.0},    // 4 back top right
	{X: -1.0, Y: 1.0, Z: 1.0},   // 5 back top left
	{X: -1.0, Y: -1.0, Z: 1.0},  // 6 back bottom left
	{X: 1.0, Y: -1.0, Z: 1.0},   // 7 back bottom right
}

var trisCube = []mesh.Triangle{
	{V1: 0, V2: 1, V3: 2},
	{V1: 0, V2: 2, V3: 3},
	{V1: 4, V2: 0, V3: 3},
	{V1: 4, V2: 3, V3: 7},
	{V1: 5, V2: 4, V3: 7},
	{V1: 5, V2: 7, V3: 6},
	{V1: 1, V2: 5, V3: 6},
	{V1: 1, V2: 6, V3: 2},
	{V1: 4, V2: 5, V3: 1},
	{V1: 4, V2: 1, V3: 0},
	{V1: 2, V2: 6, V3: 7},
	{V1: 2, V2: 7, V3: 3},
}

// CUBE === end

func NewCube(pos, scale, rotation transforms.Vec3) mesh.Model {
	m := mesh.NewMesh(vertsCube, nil, nil, trisCube, nil)
	return mesh.NewModel(&m, transforms.NewTransforms(pos, scale, rotation))
}

func NewTriangle(pos, scale, rotation transforms.Vec3) mesh.Model {
	uv_verts := []transforms.Vec2{
		{X: 0, Y: 0}, // 0 front bottom left
		{X: 0, Y: 1}, // 1 front top left
		{X: 1, Y: 1}, // 2 front top right
	}
	tris_verts := []transforms.Vec3{
		{X: -1, Y: -1, Z: 0}, // 0 front bottom left
		{X: -1, Y: 1, Z: 0},  // 1 front top left
		{X: 1, Y: 1, Z: 0},   // 2 front top right
	}
	tris_tris := []mesh.Triangle{
		{
			V1: 0, V2: 1, V3: 2,
			U1: 0, U2: 1, U3: 2,
			N1: 0, N2: 0, N3: 0,
		},
	}
	tris_normal := []transforms.Vec3{
		{Z: -1},
	}

	m := mesh.NewMesh(tris_verts, tris_normal, uv_verts, tris_tris, nil)
	return mesh.NewModel(&m, transforms.NewTransforms(pos, scale, rotation))
}
