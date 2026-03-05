package mesh

import (
	_ "image/jpeg"
	_ "image/png"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type MeshData struct {
	Tris                  []Triangle
	Verts, VertsWorld     []transforms.Vec3
	Normals, NormalsWorld []transforms.Vec3
	UV                    []transforms.Vec2
	Texture               *Texture
	Normal                *Texture
	Specular              *Texture
}

func NewMesh(
	verts, normals []transforms.Vec3,
	uvs []transforms.Vec2,
	tris []Triangle,
	texture, normal, specular *Texture) MeshData {

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
		Normal:       normal,
		Specular:     specular,
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
