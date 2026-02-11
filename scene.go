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
	va := s.activeCam.ProjectVertexToNDC(verts[tri.v1], tri.color)
	vb := s.activeCam.ProjectVertexToNDC(verts[tri.v2], tri.color)
	vc := s.activeCam.ProjectVertexToNDC(verts[tri.v3], tri.color)

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

// temporary as this must be using frustum calculation
// signed distance from view distance to point
func (s *Scene) signedDistanceToPoint(planeNormal, point Vec3) float32 {
	d := planeNormal.MultiplyByVec3(point)
	d += s.activeCam.zNear
	return d
}

func backFaceCulling(tri *Triangle, verts, normals []Vec3) bool {
	angleA := normals[tri.n1].MultiplyByVec3(verts[tri.v1].Scale(-1))
	angleB := normals[tri.n2].MultiplyByVec3(verts[tri.v2].Scale(-1))
	angleC := normals[tri.n3].MultiplyByVec3(verts[tri.v3].Scale(-1))
	return angleA >= 0 || angleB >= 0 || angleC >= 0
}

func (s *Scene) Render() {
	for _, o := range s.objects {
		matTransform := s.activeCam.transforms.matrixTransforms.MultiplyByMatrix(o.transforms.matrixTransforms)
		matRoation := s.activeCam.transforms.rotationMat.MultiplyByMatrix(o.transforms.rotationMat)

		o.boundingSphere.centerWord = matTransform.MultiplyByVec3(o.boundingSphere.center)
		planeNormal := s.activeCam.transforms.forwardDirection

		for i, v := range o.mesh.verts {
			o.mesh.vertsWorld[i] = matTransform.MultiplyByVec3(v)
		}

		for i, n := range o.mesh.normals {
			o.mesh.normalsWorld[i] = matRoation.MultiplyByVec3(n)
		}

		for _, t := range o.mesh.tris {
			if !backFaceCulling(&t, o.mesh.vertsWorld, o.mesh.normalsWorld) {
				continue
			}

			d1 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v1])
			d2 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v2])
			d3 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v3])

			if d1 <= 0.1 || d2 <= 0.1 || d3 <= 0.1 {
				continue
			}

			s.RenderTriangle(o.mesh.vertsWorld, t)
		}

		for _, t := range o.mesh.tris {
			if !backFaceCulling(&t, o.mesh.vertsWorld, o.mesh.normalsWorld) {
				continue
			}
			d1 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v1])
			d2 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v2])
			d3 := s.signedDistanceToPoint(planeNormal, o.mesh.vertsWorld[t.v3])

			if d1 <= 0.1 || d2 <= 0.1 || d3 <= 0.1 {
				continue
			}

			s.DrawWireframeTriangle(o.mesh.vertsWorld, t)
		}
		break
	}
}
