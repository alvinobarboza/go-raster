package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"os"
)

type Triangle struct {
	v1, v2, v3 int
	color      color.RGBA
}

type Texture struct {
	width, height int
	pixels        []color.RGBA
}

type MeshData struct {
	tris  []Triangle
	verts []Vec3

	texture *Texture
}

func NewMesh(verts []Vec3, tris []Triangle, texture *Texture) MeshData {
	return MeshData{
		verts:   verts,
		tris:    tris,
		texture: texture,
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

func LoadDefaultTexture() *Texture {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open("./assets/default.jpg")

	if err != nil {
		log.Println("Error: File could not be opened")
		panic(err)
	}

	defer file.Close()

	pixels, w, h, err := getPixels(file)

	if err != nil {
		log.Println("Error: Image could not be decoded")
		panic(err)
	}

	return &Texture{
		width:  w,
		height: h,
		pixels: pixels,
	}
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
	transforms     Tranforms
	boundingSphere BoundingSphere
	mesh           *MeshData
}

func NewModel(mesh *MeshData, transforms Tranforms) Model {
	m := Model{
		mesh:       mesh,
		transforms: transforms,
	}

	m.UpdateTransforms()

	return m
}

func (m *Model) UpdateTransforms() {
	m.transforms.UpdateTransforms()
	m.boundingSphere.CalculateBoundaries(m.mesh.verts, m.transforms.scaleMat)
	m.boundingSphere.centerWord = m.transforms.matrixTransforms.MultiplyByVec3(m.boundingSphere.center)
}
