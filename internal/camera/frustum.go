package camera

import (
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

const (
	TopPn uint = iota
	BottomPn
	LeftPn
	RightPn
	NearPn
	FarPn

	Planes uint = 6
)

type Plane struct {
	normal   transforms.Vec3
	distance float32
}

func NewPlane(p1, normal transforms.Vec3) Plane {
	normal = normal.Normalized()
	return Plane{
		normal:   normal,
		distance: normal.DotByVec3(p1),
	}
}

func (p *Plane) SignedDistanceToPoint(point transforms.Vec3) float32 {
	return p.normal.DotByVec3(point) - p.distance
}

type Frustum struct {
	Planes [Planes]Plane
}

func (f *Frustum) IsVertexInsideFrustum(p transforms.Vec3) bool {
	for i := range Planes {
		if f.Planes[i].SignedDistanceToPoint(p) <= 0 {
			return false
		}
	}

	return true
}

func (f *Frustum) IsBoundsInsideFrustum(b *mesh.BoundingSphere) bool {
	for i := range Planes {
		if f.Planes[i].SignedDistanceToPoint(b.CenterWord) < -b.Radius {
			return false
		}
	}

	return true
}
