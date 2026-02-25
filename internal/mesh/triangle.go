package mesh

import (
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type Triangle struct {
	V1, V2, V3 int
	U1, U2, U3 int
	N1, N2, N3 int
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

type FullTriangle struct {
	V1      ClippedVertex
	V2      ClippedVertex
	V3      ClippedVertex
	Texture *Texture
}

func NewFullTriangle(v1, v2, v3 ClippedVertex, t *Texture) FullTriangle {
	return FullTriangle{
		V1:      v1,
		V2:      v2,
		V3:      v3,
		Texture: t,
	}
}
