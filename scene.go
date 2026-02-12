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

	v0 := s.activeCam.NDCtoScreen(va)
	v1 := s.activeCam.NDCtoScreen(vb)
	v2 := s.activeCam.NDCtoScreen(vc)

	minX := MinIn(v0.X, MinIn(v1.X, v2.X))
	minY := MinIn(v0.Y, MinIn(v1.Y, v2.Y))
	maxX := MaxIn(v0.X, MaxIn(v1.X, v2.X))
	maxY := MaxIn(v0.Y, MaxIn(v1.Y, v2.Y))

	deltaW0Col := v1.Y - v2.Y
	deltaW1Col := v2.Y - v0.Y
	deltaW2Col := v0.Y - v1.Y

	deltaW0Row := v2.X - v1.X
	deltaW1Row := v0.X - v2.X
	deltaW2Row := v1.X - v0.X

	bias0 := 0
	bias1 := 0
	bias2 := 0

	if v1.IsTopOrLeft(v2) {
		bias0 = -1
	}

	if v2.IsTopOrLeft(v0) {
		bias1 = -1
	}

	if v0.IsTopOrLeft(v1) {
		bias2 = -1
	}

	area := float32(EdgeCross(v0, v1, v2))

	p := ScreenPoint{X: minX, Y: minY}
	w0Row := EdgeCross(v1, v2, p) + bias0
	w1Row := EdgeCross(v2, v0, p) + bias1
	w2Row := EdgeCross(v0, v1, p) + bias2

	for y := minY; y <= maxY; y++ {
		w0 := w0Row
		w1 := w1Row
		w2 := w2Row
		for x := minX; x <= maxX; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				// TODO: use to interpolate depth and uv coordinates
				alpha := float32(w0) / area
				beta := float32(w1) / area
				gama := float32(w2) / area

				r := 255 * alpha
				g := 255 * beta
				b := 255 * gama

				s.activeCam.PutPixel(ScreenPoint{X: x, Y: y, color: color.RGBA{
					A: 255,
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
				}})
			}
			w0 += deltaW0Col
			w1 += deltaW1Col
			w2 += deltaW2Col
		}
		w0Row += deltaW0Row
		w1Row += deltaW1Row
		w2Row += deltaW2Row
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
