package main

import (
	"image/color"
	"math"
	"sync"
)

type Renderer struct {
	scene      *Scene
	wp         *WorkerPool
	wg         sync.WaitGroup
	outputList []ClippedVertex
	inputList  []ClippedVertex
}

func NewRenderer(wp *WorkerPool) *Renderer {
	return &Renderer{
		outputList: make([]ClippedVertex, 9),
		inputList:  make([]ClippedVertex, 9),
		wp:         wp,
	}
}

func (r *Renderer) AddActiveScene(s *Scene) {
	r.scene = s
}

func (r *Renderer) DrawLine(a, b ScreenPoint, cl color.RGBA) {
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
			r.scene.activeCam.PutPixel(uint(x), uint(ys), cl, 100)
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
		r.scene.activeCam.PutPixel(uint(xs), uint(y), cl, 100)
		xs += abX
	}
}

func (r *Renderer) DrawWireframeTriangle(v1, v2, v3 ClippedVertex) {
	va := r.scene.activeCam.ProjectVertexToNDC(v1.v)
	vb := r.scene.activeCam.ProjectVertexToNDC(v2.v)
	vc := r.scene.activeCam.ProjectVertexToNDC(v3.v)

	a := r.scene.activeCam.NDCtoScreen(va)
	b := r.scene.activeCam.NDCtoScreen(vb)
	c := r.scene.activeCam.NDCtoScreen(vc)

	r.DrawLine(a, b, Black)
	r.DrawLine(b, c, Black)
	r.DrawLine(c, a, Black)
}

func (r *Renderer) RenderTriangle(vert1, vert2, vert3 ClippedVertex, t *Texture) {
	va := r.scene.activeCam.ProjectVertexToNDC(vert1.v)
	vb := r.scene.activeCam.ProjectVertexToNDC(vert2.v)
	vc := r.scene.activeCam.ProjectVertexToNDC(vert3.v)

	depthA := 1 / vert1.v.Z
	depthB := 1 / vert2.v.Z
	depthC := 1 / vert3.v.Z

	uv1z := vert1.u.Scale(depthA)
	uv2z := vert2.u.Scale(depthB)
	uv3z := vert3.u.Scale(depthC)

	v0 := r.scene.activeCam.NDCtoScreen(va)
	v1 := r.scene.activeCam.NDCtoScreen(vb)
	v2 := r.scene.activeCam.NDCtoScreen(vc)

	minX := Floor32(Minf(v0.X, Minf(v1.X, v2.X)))
	minY := Floor32(Minf(v0.Y, Minf(v1.Y, v2.Y)))
	maxX := Ceil32(Maxf(v0.X, Maxf(v1.X, v2.X)))
	maxY := Ceil32(Maxf(v0.Y, Maxf(v1.Y, v2.Y)))

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
	for y := minY; y <= maxY; y++ {
		w0 := w0Row
		w1 := w1Row
		w2 := w2Row

		for x := minX; x <= maxX; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				alpha := w1 * area
				beta := w2 * area
				gama := w0 * area

				depth := depthA*alpha + depthB*beta + depthC*gama

				xx, yy := uint(x), uint(y)
				if r.scene.activeCam.DepthPass(xx, yy, depth) {
					uv1 := uv1z.Scale(alpha)
					uv2 := uv2z.Scale(beta)
					uv3 := uv3z.Scale(gama)

					uvCoord := uv1.Add(uv2).Add(uv3).Divide(depth)
					pColor := t.TexelColor(uvCoord)

					r.scene.activeCam.PutPixel(xx, yy, pColor, depth)
				}
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

func (r *Renderer) renderMeshs() {
	for _, o := range r.scene.objects {
		matTransform := r.scene.activeCam.transforms.matrixTransforms.MultiplyByMatrix(o.transforms.matrixTransforms)
		matRoation := r.scene.activeCam.transforms.rotationMat.MultiplyByMatrix(o.transforms.rotationMat)

		o.boundingSphere.centerWord = matTransform.MultiplyByVec3(o.boundingSphere.center)

		if !r.scene.activeCam.frustum.IsBoundsInsideFrustum(&o.boundingSphere) {
			continue
		}

		for _, t := range o.mesh.tris {
			o.mesh.vertsWorld[t.v1] = matTransform.MultiplyByVec3(o.mesh.verts[t.v1])
			o.mesh.vertsWorld[t.v2] = matTransform.MultiplyByVec3(o.mesh.verts[t.v2])
			o.mesh.vertsWorld[t.v3] = matTransform.MultiplyByVec3(o.mesh.verts[t.v3])

			o.mesh.normalsWorld[t.n1] = matRoation.MultiplyByVec3(o.mesh.normals[t.n1])
			o.mesh.normalsWorld[t.n2] = matRoation.MultiplyByVec3(o.mesh.normals[t.n2])
			o.mesh.normalsWorld[t.n3] = matRoation.MultiplyByVec3(o.mesh.normals[t.n3])

			if !t.backFaceCulling(o.mesh.vertsWorld, o.mesh.normalsWorld) {
				continue
			}

			v1 := ClippedVertex{
				v: o.mesh.vertsWorld[t.v1],
				n: o.mesh.normalsWorld[t.n1],
				u: o.mesh.uv[t.u1],
			}

			v2 := ClippedVertex{
				v: o.mesh.vertsWorld[t.v2],
				n: o.mesh.normalsWorld[t.n2],
				u: o.mesh.uv[t.u2],
			}

			v3 := ClippedVertex{
				v: o.mesh.vertsWorld[t.v3],
				n: o.mesh.normalsWorld[t.n3],
				u: o.mesh.uv[t.u3],
			}

			r.outputList = r.outputList[:0]
			r.inputList = r.inputList[:3]

			r.outputList = append(r.outputList, v1)
			r.outputList = append(r.outputList, v2)
			r.outputList = append(r.outputList, v3)

			// Sutherland–Hodgman algorithm
			for _, plane := range r.scene.activeCam.frustum.planes {
				r.inputList, r.outputList = r.outputList, r.inputList[:0]

				prevI := 0
				for i := range len(r.inputList) {
					prevI = i - 1
					if prevI < 0 {
						prevI = len(r.inputList) - 1
					}
					currentPoint := r.inputList[i]
					prevPoint := r.inputList[prevI]

					cp := plane.SignedDistanceToPoint(currentPoint.v)
					pp := plane.SignedDistanceToPoint(prevPoint.v)

					if cp > 0 {
						if pp <= 0 {
							ratio := cp / (cp - pp)
							intersection := ClippedVertex{
								v: currentPoint.v.LerpTo(prevPoint.v, ratio),
								n: currentPoint.n.LerpTo(prevPoint.n, ratio),
								u: currentPoint.u.LerpTo(prevPoint.u, ratio),
							}
							r.outputList = append(r.outputList, intersection)
						}
						r.outputList = append(r.outputList, currentPoint)
					} else if pp > 0 {
						ratio := cp / (cp - pp)
						intersection := ClippedVertex{
							v: currentPoint.v.LerpTo(prevPoint.v, ratio),
							n: currentPoint.n.LerpTo(prevPoint.n, ratio),
							u: currentPoint.u.LerpTo(prevPoint.u, ratio),
						}
						r.outputList = append(r.outputList, intersection)
					}
				}
			}

			if len(r.outputList) > 2 {
				for i := 1; i < len(r.outputList)-1; i++ {
					r.RenderTriangle(
						r.outputList[0],
						r.outputList[i],
						r.outputList[i+1],
						o.mesh.texture,
					)
					if r.scene.activeCam.renderWire {
						r.DrawWireframeTriangle(
							r.outputList[0],
							r.outputList[i],
							r.outputList[i+1],
						)
					}
				}
			}
		}
		// break
	}
}

func (r *Renderer) Render() {
	r.scene.activeCam.ClearCanvas()
	r.renderMeshs()
}
