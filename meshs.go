package main

import "image/color"

type Triangle struct {
	v1, v2, v3 int
	color      color.RGBA
}

type MeshData struct {
	tris  []Triangle
	verts []Vec3
}

func NewMesh(verts []Vec3, tris []Triangle) MeshData {
	return MeshData{
		verts: verts,
		tris:  tris,
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
}
