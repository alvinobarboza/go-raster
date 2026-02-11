package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
)

type Triangle struct {
	v1, v2, v3 int
	u1, u2, u3 int
	n1, n2, n3 int
	color      color.RGBA
}

type Texture struct {
	width, height int
	pixels        []color.RGBA
}

type MeshData struct {
	tris       []Triangle
	verts      []Vec3
	vertsWorld []Vec3

	normals []Vec3
	uv      []Vec3

	texture *Texture
}

func NewMesh(verts []Vec3, tris []Triangle, texture *Texture) MeshData {
	vertsWord := make([]Vec3, len(verts))
	return MeshData{
		verts:      verts,
		vertsWorld: vertsWord,
		tris:       tris,
		texture:    texture,
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
	for y := range height {
		for x := range width {
			r, g, b, a := img.At(x, y).RGBA()
			pixels = append(pixels, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
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
		scaled = scaled.Sub(s.center)

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
	m.transforms.UpdateTransforms(false, false)
	m.boundingSphere.CalculateBoundaries(m.mesh.verts, m.transforms.scaleMat)
	m.boundingSphere.centerWord = m.transforms.matrixTransforms.MultiplyByVec3(m.boundingSphere.center)
}
