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

func (s *Scene) DrawLine(a, b ScreenPoint, cl color.RGBA) {
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
			s.activeCam.PutPixel(uint(x), uint(ys), cl, 100)
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
		s.activeCam.PutPixel(uint(xs), uint(y), cl, 100)
		xs += abX
	}
}

func (s *Scene) DrawWireframeTriangle(verts []Vec3, tri Triangle) {
	va := s.activeCam.ProjectVertexToNDC(verts[tri.v1])
	vb := s.activeCam.ProjectVertexToNDC(verts[tri.v2])
	vc := s.activeCam.ProjectVertexToNDC(verts[tri.v3])

	a := s.activeCam.NDCtoScreen(va)
	b := s.activeCam.NDCtoScreen(vb)
	c := s.activeCam.NDCtoScreen(vc)

	s.DrawLine(a, b, Black)
	s.DrawLine(b, c, Black)
	s.DrawLine(c, a, Black)
}

func (s *Scene) RenderTriangle(verts, uv []Vec3, tri Triangle, t *Texture) {
	va := s.activeCam.ProjectVertexToNDC(verts[tri.v1])
	vb := s.activeCam.ProjectVertexToNDC(verts[tri.v2])
	vc := s.activeCam.ProjectVertexToNDC(verts[tri.v3])

	depthA := 1 / verts[tri.v1].Z
	depthB := 1 / verts[tri.v2].Z
	depthC := 1 / verts[tri.v3].Z

	uv1z := uv[tri.u1].Scale(depthA)
	uv2z := uv[tri.u2].Scale(depthB)
	uv3z := uv[tri.u3].Scale(depthC)

	v0 := s.activeCam.NDCtoScreen(va)
	v1 := s.activeCam.NDCtoScreen(vb)
	v2 := s.activeCam.NDCtoScreen(vc)

	minX := float32(math.Floor(float64(Minf(v0.X, Minf(v1.X, v2.X)))))
	minY := float32(math.Floor(float64(Minf(v0.Y, Minf(v1.Y, v2.Y)))))
	maxX := float32(math.Ceil(float64(Maxf(v0.X, Maxf(v1.X, v2.X)))))
	maxY := float32(math.Ceil(float64(Maxf(v0.Y, Maxf(v1.Y, v2.Y)))))

	deltaW0Col := v0.Y - v1.Y
	deltaW1Col := v1.Y - v2.Y
	deltaW2Col := v2.Y - v0.Y

	deltaW0Row := v1.X - v0.X
	deltaW1Row := v2.X - v1.X
	deltaW2Row := v0.X - v2.X

	bias0 := float32(0)
	bias1 := float32(0)
	bias2 := float32(0)

	if v0.IsTopOrLeft(v1) {
		bias0 = -0.0001
	}

	if v1.IsTopOrLeft(v2) {
		bias1 = -0.0001
	}

	if v2.IsTopOrLeft(v0) {
		bias2 = -0.0001
	}

	area := float32(EdgeCross(v0, v1, v2))
	area = 1 / area

	// pixel's center
	p := ScreenPoint{X: minX + 0.5, Y: minY + 0.5}

	w0Row := EdgeCross(v0, v1, p) + bias0
	w1Row := EdgeCross(v1, v2, p) + bias1
	w2Row := EdgeCross(v2, v0, p) + bias2

	for y := minY; y <= maxY; y++ {
		w0 := w0Row
		w1 := w1Row
		w2 := w2Row
		for x := minX; x <= maxX; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				/*
					      v0 (Top)
					      /\
					     /  \
					    /    \    <-- The distance from this edge (v0-v1)
					   /      \       towards v2 is w0.
					  /   P    \
					 /    |     \
					v1 ---|------v2
					      ^
					      |
					The distance from this edge (v1-v2)
					towards v0 is w1.
					w1 = v1 -> v2 distance to v0 = a = tri.v1
					w2 = v2 -> v0 distance to v1 = b = tri.v2
					w0 = v0 -> v1 distance to v2 = c = tri.v3
				*/
				alpha := w1 * area
				beta := w2 * area
				gama := w0 * area

				depth := depthA*alpha + depthB*beta + depthC*gama

				uv1 := uv1z.Scale(alpha)
				uv2 := uv2z.Scale(beta)
				uv3 := uv3z.Scale(gama)

				uvCoord := uv1.Add(uv2).Add(uv3).Divide(depth)

				s.activeCam.PutPixel(uint(x), uint(y), t.TexelColor(uvCoord), depth)
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

func (s *Scene) Render() {
	for _, o := range s.objects {
		matTransform := s.activeCam.transforms.matrixTransforms.MultiplyByMatrix(o.transforms.matrixTransforms)
		matRoation := s.activeCam.transforms.rotationMat.MultiplyByMatrix(o.transforms.rotationMat)

		o.boundingSphere.centerWord = matTransform.MultiplyByVec3(o.boundingSphere.center)

		if !s.activeCam.frustum.IsBoundsInsideFrustum(&o.boundingSphere) {
			continue
		}

		for _, t := range o.mesh.tris {
			o.mesh.vertsWorld[t.v1] = matTransform.MultiplyByVec3(o.mesh.verts[t.v1])
			o.mesh.vertsWorld[t.v2] = matTransform.MultiplyByVec3(o.mesh.verts[t.v2])
			o.mesh.vertsWorld[t.v3] = matTransform.MultiplyByVec3(o.mesh.verts[t.v3])

			o.mesh.normalsWorld[t.n1] = matRoation.MultiplyByVec3(o.mesh.normals[t.n1])
			o.mesh.normalsWorld[t.n2] = matRoation.MultiplyByVec3(o.mesh.normals[t.n2])
			o.mesh.normalsWorld[t.n3] = matRoation.MultiplyByVec3(o.mesh.normals[t.n3])
		}

		// TODO: generate new tris on frustum plane intersections
		for _, t := range o.mesh.tris {
			if !t.backFaceCulling(o.mesh.vertsWorld, o.mesh.normalsWorld) {
				continue
			}

			if !s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v1]) ||
				!s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v2]) ||
				!s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v3]) {
				continue
			}

			s.RenderTriangle(o.mesh.vertsWorld, o.mesh.uv, t, o.mesh.texture)
		}

		if s.activeCam.renderWire {
			for _, t := range o.mesh.tris {
				// if !t.backFaceCulling(o.mesh.vertsWorld, o.mesh.normalsWorld) {
				// 	continue
				// }

				if !s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v1]) ||
					!s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v2]) ||
					!s.activeCam.frustum.IsVertexInsideFrustum(o.mesh.vertsWorld[t.v3]) {
					continue
				}

				s.DrawWireframeTriangle(o.mesh.vertsWorld, t)
			}
		}
		// break
	}
}
