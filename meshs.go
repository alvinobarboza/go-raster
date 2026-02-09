package main

import "image/color"

type Mesh struct {
	vertices []Vec3
	color    color.RGBA
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

	boundingSphere BoundingSphere
}
