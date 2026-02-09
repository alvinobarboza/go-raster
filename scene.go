package main

import (
	"image/color"
	"math"
)

type Scene struct {
	activeCam *Camera
	objects   []*Model
}

func NewScene(c *Camera) Scene {
	return Scene{
		activeCam: c,
		objects:   make([]*Model, 0),
	}
}

func (s *Scene) AddMesh(o *Model) {
	s.objects = append(s.objects, o)
}

func (s *Scene) DrawLine(a, b ScreenPoint) {
	dx := b.X - a.X
	dy := b.Y - a.Y

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if dx < 0 {
			tmp := a
			a = b
			b = tmp
		}

		abY := float32(b.Y-a.Y) / float32(b.X-a.X)
		ys := float32(a.Y)
		for x := a.X; x <= b.X; x++ {
			s.activeCam.PutPixel(ScreenPoint{X: x, Y: int(ys), color: a.color})
			ys += abY
		}
		return
	}

	if dy < 0 {
		tmp := a
		a = b
		b = tmp
	}

	abX := float32(b.X-a.X) / float32(b.Y-a.Y)
	xs := float32(a.X)

	for y := a.Y; y <= b.Y; y++ {
		s.activeCam.PutPixel(ScreenPoint{X: int(xs), Y: y, color: a.color})
		xs += abX
	}
}

func (s *Scene) DrawWireframeTriangle(v1, v2, v3 Vec3, cl color.RGBA) {
	va := s.activeCam.ProjectVertexToNDC(v1, cl)
	vb := s.activeCam.ProjectVertexToNDC(v2, cl)
	vc := s.activeCam.ProjectVertexToNDC(v3, cl)

	a := s.activeCam.NDCtoScreen(va)
	b := s.activeCam.NDCtoScreen(vb)
	c := s.activeCam.NDCtoScreen(vc)

	s.DrawLine(a, b)
	s.DrawLine(b, c)
	s.DrawLine(c, a)
}

func (s *Scene) Render() {
	for _, o := range s.objects {
		for _, t := range o.mesh.tris {
			v1 := o.mesh.verts[t.v1]
			v2 := o.mesh.verts[t.v2]
			v3 := o.mesh.verts[t.v3]

			v1 = o.transforms.matrixTransforms.MultiplyByVec3(v1)
			v2 = o.transforms.matrixTransforms.MultiplyByVec3(v2)
			v3 = o.transforms.matrixTransforms.MultiplyByVec3(v3)

			s.DrawWireframeTriangle(v1, v2, v3, t.color)
		}
	}
}
