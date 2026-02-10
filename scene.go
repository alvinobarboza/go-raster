package main

import (
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

func (s *Scene) DrawWireframeTriangle(verts []Vec3, tri Triangle) {
	va := s.activeCam.ProjectVertexToNDC(verts[tri.v1], Black)
	vb := s.activeCam.ProjectVertexToNDC(verts[tri.v2], Black)
	vc := s.activeCam.ProjectVertexToNDC(verts[tri.v3], Black)

	a := s.activeCam.NDCtoScreen(va)
	b := s.activeCam.NDCtoScreen(vb)
	c := s.activeCam.NDCtoScreen(vc)

	s.DrawLine(a, b)
	s.DrawLine(b, c)
	s.DrawLine(c, a)
}

func (s *Scene) RenderTriangle(verts []Vec3, tri Triangle) {
	va := s.activeCam.ProjectVertexToNDC(verts[tri.v3], tri.color)
	vb := s.activeCam.ProjectVertexToNDC(verts[tri.v2], tri.color)
	vc := s.activeCam.ProjectVertexToNDC(verts[tri.v1], tri.color)

	a := s.activeCam.NDCtoScreen(va)
	b := s.activeCam.NDCtoScreen(vb)
	c := s.activeCam.NDCtoScreen(vc)

	maxX := MaxIn(a.X, MaxIn(b.X, c.X))
	minX := MinIn(a.X, MinIn(b.X, c.X))
	maxY := MaxIn(a.Y, MaxIn(b.Y, c.Y))
	minY := MinIn(a.Y, MinIn(b.Y, c.Y))

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			p := ScreenPoint{X: x, Y: y}
			w1 := EdgeCross(b, c, p)
			w2 := EdgeCross(c, a, p)
			w3 := EdgeCross(a, b, p)

			if w1 >= 0 && w2 >= 0 && w3 >= 0 {
				s.activeCam.PutPixel(ScreenPoint{X: x, Y: y, color: tri.color})
			}
		}
	}
}

func (s *Scene) Render() {
	for _, o := range s.objects {
		matTransform := s.activeCam.transforms.matrixTransforms.MultiplyByMatrix(o.transforms.matrixTransforms)
		o.boundingSphere.centerWord = matTransform.MultiplyByVec3(o.boundingSphere.center)

		for i, v := range o.mesh.verts {
			o.mesh.vertsWorld[i] = matTransform.MultiplyByVec3(v)
		}

		for _, t := range o.mesh.tris {
			s.RenderTriangle(o.mesh.vertsWorld, t)
		}
		for _, t := range o.mesh.tris {
			s.DrawWireframeTriangle(o.mesh.vertsWorld, t)
		}
	}
}
