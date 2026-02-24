package mesh

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type Triangle struct {
	V1, V2, V3 int
	U1, U2, U3 int
	N1, N2, N3 int
	Color      color.RGBA
}

func (t *Triangle) BackFaceCulling(verts, normals []transforms.Vec3) bool {
	angleA := normals[t.N1].DotByVec3(verts[t.V1].Scale(-1))
	angleB := normals[t.N2].DotByVec3(verts[t.V2].Scale(-1))
	angleC := normals[t.N3].DotByVec3(verts[t.V3].Scale(-1))
	return angleA >= 0 || angleB >= 0 || angleC >= 0
}

type ClippedVertex struct {
	V transforms.Vec3
	N transforms.Vec3
	U transforms.Vec2
}

type Texture struct {
	width, height int
	pixels        []color.RGBA
}

func (t *Texture) TexelColor(uv transforms.Vec2) color.RGBA {
	u := uv.X - maths.Floor32(uv.X)
	v := uv.Y - maths.Floor32(uv.Y)

	w := int(u * float32(t.width))
	h := int(v * float32(t.height))

	i := h*t.width + w

	if uint(i) < uint(len(t.pixels)) {
		return t.pixels[i]
	}

	return t.pixels[0]
}

type MeshData struct {
	Tris                  []Triangle
	Verts, VertsWorld     []transforms.Vec3
	Normals, NormalsWorld []transforms.Vec3
	UV                    []transforms.Vec2
	Texture               *Texture
}

func NewMesh(verts, normals []transforms.Vec3, uvs []transforms.Vec2, tris []Triangle, texture *Texture) MeshData {
	vertsWord := make([]transforms.Vec3, len(verts))
	normalsWord := make([]transforms.Vec3, len(normals))
	return MeshData{
		Verts:        verts,
		Normals:      normals,
		NormalsWorld: normalsWord,
		UV:           uvs,
		VertsWorld:   vertsWord,
		Tris:         tris,
		Texture:      texture,
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
	Center, CenterWord transforms.Vec3
	Radius             float32
}

func NewBoundingSphere() BoundingSphere {
	return BoundingSphere{
		Center:     transforms.NewVec3(0, 0, 0),
		CenterWord: transforms.NewVec3(0, 0, 0),
		Radius:     0,
	}
}

func (s *BoundingSphere) CalculateBoundaries(verts []transforms.Vec3, scale transforms.Matrix) {
	*s = NewBoundingSphere()

	for _, v := range verts {
		s.Center = s.Center.Add(v)
	}

	s.Center = s.Center.Divide(float32(len(verts)))

	for _, v := range verts {
		scaled := scale.MultiplyByVec3(v)
		scaled = scaled.Subtract(s.Center)

		r := scaled.Length()

		if s.Radius < r {
			s.Radius = r
		}
	}
}

type Model struct {
	Transforms     transforms.Transforms
	BoundingSphere BoundingSphere
	Mesh           *MeshData
}

func NewModel(mesh *MeshData, transforms transforms.Transforms) Model {
	m := Model{
		Mesh:       mesh,
		Transforms: transforms,
	}

	m.UpdateTransforms()

	return m
}

func (m *Model) UpdateTransforms() {
	m.Transforms.UpdateModelTransforms()
	m.BoundingSphere.CalculateBoundaries(m.Mesh.Verts, m.Transforms.ScaleMat)
	m.BoundingSphere.CenterWord = m.Transforms.MatrixTransforms.MultiplyByVec3(m.BoundingSphere.Center)
}
