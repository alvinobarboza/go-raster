package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type Triangle struct {
	v1, v2, v3 int
	u1, u2, u3 int
	n1, n2, n3 int
	color      color.RGBA
}

func (t *Triangle) backFaceCulling(verts, normals []Vec3) bool {
	angleA := normals[t.n1].DotByVec3(verts[t.v1].Scale(-1))
	angleB := normals[t.n2].DotByVec3(verts[t.v2].Scale(-1))
	angleC := normals[t.n3].DotByVec3(verts[t.v3].Scale(-1))
	return angleA >= 0 || angleB >= 0 || angleC >= 0
}

type ClippedVertex struct {
	v Vec3
	n Vec3
	u Vec2
}

type Texture struct {
	width, height int
	pixels        []color.RGBA
}

func (t *Texture) TexelColor(uv Vec2) color.RGBA {
	uv.X = uv.X - Floor32(uv.X)
	uv.Y = uv.Y - Floor32(uv.Y)

	w, h := int(uv.X*float32(t.width)), int(uv.Y*float32(t.height))

	i := h*t.width + w

	if uint(i) < uint(len(t.pixels)) {
		return t.pixels[i]
	}

	return t.pixels[0]
}

type MeshData struct {
	tris                  []Triangle
	verts, vertsWorld     []Vec3
	normals, normalsWorld []Vec3
	uv                    []Vec2
	texture               *Texture
}

func NewMesh(verts, normals []Vec3, uvs []Vec2, tris []Triangle, texture *Texture) MeshData {
	vertsWord := make([]Vec3, len(verts))
	normalsWord := make([]Vec3, len(normals))
	return MeshData{
		verts:        verts,
		normals:      normals,
		normalsWorld: normalsWord,
		uv:           uvs,
		vertsWorld:   vertsWord,
		tris:         tris,
		texture:      texture,
	}
}

func getPixels(file io.Reader) ([]color.RGBA, int, int, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([]color.RGBA, 0)
	// Upside down, since render is upside down
	for y := range height {
		yu := height - y - 1
		for x := range width {
			r, g, b, a := img.At(x, yu).RGBA()

			// From alpha pre-multiplied values
			// 0xFF00 > 0x00FF > 0xFF
			pixels = append(pixels, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	return pixels, width, height, nil
}

type BoundingSphere struct {
	center, centerWord Vec3
	radius             float32
}

func NewBoundingSphere() BoundingSphere {
	return BoundingSphere{
		center:     NewVec3(0, 0, 0),
		centerWord: NewVec3(0, 0, 0),
		radius:     0,
	}
}

func (s *BoundingSphere) CalculateBoundaries(verts []Vec3, scale Matrix) {
	*s = NewBoundingSphere()

	for _, v := range verts {
		s.center = s.center.Add(v)
	}

	s.center = s.center.Divide(float32(len(verts)))

	for _, v := range verts {
		scaled := scale.MultiplyByVec3(v)
		scaled = scaled.Subtract(s.center)

		r := scaled.Length()

		if s.radius < r {
			s.radius = r
		}
	}
}

type Model struct {
	transforms     Transforms
	boundingSphere BoundingSphere
	mesh           *MeshData
}

func NewModel(mesh *MeshData, transforms Transforms) Model {
	m := Model{
		mesh:       mesh,
		transforms: transforms,
	}

	m.UpdateTransforms()

	return m
}

func (m *Model) UpdateTransforms() {
	m.transforms.UpdateModelTransforms()
	m.boundingSphere.CalculateBoundaries(m.mesh.verts, m.transforms.scaleMat)
	m.boundingSphere.centerWord = m.transforms.matrixTransforms.MultiplyByVec3(m.boundingSphere.center)
}
