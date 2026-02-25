package mesh

import "github.com/alvinobarboza/go-raster/internal/transforms"

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
