package renderer

import (
	"fmt"
	"image/color"
	"math"
	"sync"

	"github.com/alvinobarboza/go-raster/internal/camera"
	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/scene"
	"github.com/alvinobarboza/go-raster/internal/shapes"
)

type Renderer struct {
	scene   *scene.Scene
	wg      sync.WaitGroup
	mt      sync.Mutex
	indexes chan int

	sHoutputList    []mesh.ClippedVertex
	sHinputList     []mesh.ClippedVertex
	trianglesBuffer []mesh.FullTriangle

	tiles []*ScreenTile

	RenderTileBoundaries     bool
	RenderTriangleBoundaries bool
	RenderMultithreaded      bool
}

// init after loading models,
// otherwise triangle buffer will have 0 size
// TODO -> somehow get around this
func NewRenderer(threads uint) *Renderer {
	r := &Renderer{
		sHoutputList:             make([]mesh.ClippedVertex, 9),
		sHinputList:              make([]mesh.ClippedVertex, 9),
		indexes:                  make(chan int, threads),
		RenderTileBoundaries:     false,
		RenderTriangleBoundaries: false,
		RenderMultithreaded:      true,
	}

	for i := range threads {
		go r.renderTriangleParallel(i)
	}

	return r
}

func (r *Renderer) AddActiveScene(s *scene.Scene) {
	r.scene = s
	biggestTriCount := 0
	for _, o := range s.Objects {
		if len(o.Mesh.Tris) > biggestTriCount {
			biggestTriCount = len(o.Mesh.Tris) * 9
		}
	}

	r.tiles = NewTileSet(
		float32(r.scene.ActiveCam.Width),
		float32(r.scene.ActiveCam.Height),
		float32(r.scene.ActiveCam.Width)/4,
		biggestTriCount,
	)

	r.trianglesBuffer = make([]mesh.FullTriangle, biggestTriCount)
}

func (r *Renderer) DrawLine(a, b camera.ScreenPoint, cl color.RGBA) {
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
			r.scene.ActiveCam.PutPixel(uint(x), uint(ys), cl, 100)
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
		r.scene.ActiveCam.PutPixel(uint(xs), uint(y), cl, 100)
		xs += abX
	}
}

func (r *Renderer) DrawWireframeTriangle(v1, v2, v3 mesh.ClippedVertex) {
	va := r.scene.ActiveCam.ProjectVertexToNDC(v1.V)
	vb := r.scene.ActiveCam.ProjectVertexToNDC(v2.V)
	vc := r.scene.ActiveCam.ProjectVertexToNDC(v3.V)

	a := r.scene.ActiveCam.NDCtoScreen(va)
	b := r.scene.ActiveCam.NDCtoScreen(vb)
	c := r.scene.ActiveCam.NDCtoScreen(vc)

	r.DrawLine(a, b, shapes.Black)
	r.DrawLine(b, c, shapes.Black)
	r.DrawLine(c, a, shapes.Black)
}

func (r *Renderer) drawTileBoundaries() {
	for i := range r.tiles {
		triangles := r.tiles[i].Triangles()

		if len(triangles) > 0 {
			mnX, mnY, mxX, mxY := r.tiles[i].Bounduries()
			r.DrawLine(camera.ScreenPoint{X: mnX, Y: mnY}, camera.ScreenPoint{X: mnX, Y: mxY}, shapes.DarkPurple)
			r.DrawLine(camera.ScreenPoint{X: mnX, Y: mxY}, camera.ScreenPoint{X: mxX, Y: mxY}, shapes.DarkPurple)
			r.DrawLine(camera.ScreenPoint{X: mxX, Y: mxY}, camera.ScreenPoint{X: mxX, Y: mnY}, shapes.DarkPurple)
			r.DrawLine(camera.ScreenPoint{X: mxX, Y: mnY}, camera.ScreenPoint{X: mnX, Y: mnY}, shapes.DarkPurple)
		}
	}
}

func (r *Renderer) renderTriangleParallel(id uint) {
	fmt.Println("Thread ID:", id)
	for i := range r.indexes {
		triangles := r.tiles[i].Triangles()
		mnX, mnY, mxX, mxY := r.tiles[i].Bounduries()

		for _, tri := range triangles {
			va := r.scene.ActiveCam.ProjectVertexToNDC(tri.V1.V)
			vb := r.scene.ActiveCam.ProjectVertexToNDC(tri.V2.V)
			vc := r.scene.ActiveCam.ProjectVertexToNDC(tri.V3.V)

			depthA := 1 / tri.V1.V.Z
			depthB := 1 / tri.V2.V.Z
			depthC := 1 / tri.V3.V.Z

			uv1z := tri.V1.U.Scale(depthA)
			uv2z := tri.V2.U.Scale(depthB)
			uv3z := tri.V3.U.Scale(depthC)

			v0 := r.scene.ActiveCam.NDCtoScreen(va)
			v1 := r.scene.ActiveCam.NDCtoScreen(vb)
			v2 := r.scene.ActiveCam.NDCtoScreen(vc)

			minX := maths.Floor32(maths.Minf(v0.X, maths.Minf(v1.X, v2.X)))
			minY := maths.Floor32(maths.Minf(v0.Y, maths.Minf(v1.Y, v2.Y)))
			maxX := maths.Ceil32(maths.Maxf(v0.X, maths.Maxf(v1.X, v2.X)))
			maxY := maths.Ceil32(maths.Maxf(v0.Y, maths.Maxf(v1.Y, v2.Y)))

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

			area := camera.EdgeCross(v0, v1, v2)
			area = 1 / area

			// check to only run in tile bounds
			cMinY := maths.Maxf(minY, mnY)
			cMaxY := maths.Minf(maxY, mxY)

			cMinX := maths.Maxf(minX, mnX)
			cMaxX := maths.Minf(maxX, mxX)

			if r.RenderTriangleBoundaries {
				r.DrawLine(camera.ScreenPoint{X: cMinX, Y: cMinY}, camera.ScreenPoint{X: cMinX, Y: cMaxY}, shapes.DarkGreen)
				r.DrawLine(camera.ScreenPoint{X: cMinX, Y: cMaxY}, camera.ScreenPoint{X: cMaxX, Y: cMaxY}, shapes.DarkGreen)
				r.DrawLine(camera.ScreenPoint{X: cMaxX, Y: cMaxY}, camera.ScreenPoint{X: cMaxX, Y: cMinY}, shapes.DarkGreen)
				r.DrawLine(camera.ScreenPoint{X: cMaxX, Y: cMinY}, camera.ScreenPoint{X: cMinX, Y: cMinY}, shapes.DarkGreen)
			}

			// pixel's center
			p := camera.ScreenPoint{X: cMinX + 0.5, Y: cMinY + 0.5}

			w0Row := camera.EdgeCross(v0, v1, p) + bias0
			w1Row := camera.EdgeCross(v1, v2, p) + bias1
			w2Row := camera.EdgeCross(v2, v0, p) + bias2

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

			for y := cMinY; y <= cMaxY; y++ {
				w0 := w0Row
				w1 := w1Row
				w2 := w2Row

				for x := cMinX; x <= cMaxX; x++ {
					if w0 >= 0 && w1 >= 0 && w2 >= 0 {
						alpha := w1 * area
						beta := w2 * area
						gama := w0 * area

						depth := depthA*alpha + depthB*beta + depthC*gama

						xx, yy := uint(x), uint(y)
						if r.scene.ActiveCam.DepthPass(xx, yy, depth) {
							uv1 := uv1z.Scale(alpha)
							uv2 := uv2z.Scale(beta)
							uv3 := uv3z.Scale(gama)

							uvCoord := uv1.Add(uv2).Add(uv3).Divide(depth)
							pColor := tri.Texture.TexelColor(uvCoord)

							r.scene.ActiveCam.PutPixel(xx, yy, pColor, depth)
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

		r.wg.Done()
	}
}

func (r *Renderer) assignTrianglesToTiles() {
	for _, t := range r.trianglesBuffer {
		va := r.scene.ActiveCam.ProjectVertexToNDC(t.V1.V)
		vb := r.scene.ActiveCam.ProjectVertexToNDC(t.V2.V)
		vc := r.scene.ActiveCam.ProjectVertexToNDC(t.V3.V)

		v0 := r.scene.ActiveCam.NDCtoScreen(va)
		v1 := r.scene.ActiveCam.NDCtoScreen(vb)
		v2 := r.scene.ActiveCam.NDCtoScreen(vc)

		minX := maths.Floor32(maths.Minf(v0.X, maths.Minf(v1.X, v2.X)))
		minY := maths.Floor32(maths.Minf(v0.Y, maths.Minf(v1.Y, v2.Y)))
		maxX := maths.Ceil32(maths.Maxf(v0.X, maths.Maxf(v1.X, v2.X)))
		maxY := maths.Ceil32(maths.Maxf(v0.Y, maths.Maxf(v1.Y, v2.Y)))

		for _, st := range r.tiles {
			if st.TileTriangleCollision(minX, minY, maxX, maxY) {
				st.AddTriangle(t)
			}
		}
	}
}

func (r *Renderer) RenderTriangle(vert1, vert2, vert3 mesh.ClippedVertex, t *mesh.Texture) {
	va := r.scene.ActiveCam.ProjectVertexToNDC(vert1.V)
	vb := r.scene.ActiveCam.ProjectVertexToNDC(vert2.V)
	vc := r.scene.ActiveCam.ProjectVertexToNDC(vert3.V)

	depthA := 1 / vert1.V.Z
	depthB := 1 / vert2.V.Z
	depthC := 1 / vert3.V.Z

	uv1z := vert1.U.Scale(depthA)
	uv2z := vert2.U.Scale(depthB)
	uv3z := vert3.U.Scale(depthC)

	v0 := r.scene.ActiveCam.NDCtoScreen(va)
	v1 := r.scene.ActiveCam.NDCtoScreen(vb)
	v2 := r.scene.ActiveCam.NDCtoScreen(vc)

	minX := maths.Floor32(maths.Minf(v0.X, maths.Minf(v1.X, v2.X)))
	minY := maths.Floor32(maths.Minf(v0.Y, maths.Minf(v1.Y, v2.Y)))
	maxX := maths.Ceil32(maths.Maxf(v0.X, maths.Maxf(v1.X, v2.X)))
	maxY := maths.Ceil32(maths.Maxf(v0.Y, maths.Maxf(v1.Y, v2.Y)))

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

	area := camera.EdgeCross(v0, v1, v2)
	area = 1 / area

	// pixel's center
	p := camera.ScreenPoint{X: minX + 0.5, Y: minY + 0.5}

	w0Row := camera.EdgeCross(v0, v1, p) + bias0
	w1Row := camera.EdgeCross(v1, v2, p) + bias1
	w2Row := camera.EdgeCross(v2, v0, p) + bias2

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
				if r.scene.ActiveCam.DepthPass(xx, yy, depth) {
					uv1 := uv1z.Scale(alpha)
					uv2 := uv2z.Scale(beta)
					uv3 := uv3z.Scale(gama)

					uvCoord := uv1.Add(uv2).Add(uv3).Divide(depth)
					pColor := t.TexelColor(uvCoord)

					r.scene.ActiveCam.PutPixel(xx, yy, pColor, depth)
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
	for _, o := range r.scene.Objects {
		matTransform := r.scene.ActiveCam.Transforms.MatrixTransforms.MultiplyByMatrix(o.Transforms.MatrixTransforms)
		matRoation := r.scene.ActiveCam.Transforms.RotationMat.MultiplyByMatrix(o.Transforms.RotationMat)

		o.BoundingSphere.CenterWord = matTransform.MultiplyByVec3(o.BoundingSphere.Center)

		if !r.scene.ActiveCam.Frustum.IsBoundsInsideFrustum(&o.BoundingSphere) {
			continue
		}

		r.trianglesBuffer = r.trianglesBuffer[:0]
		for _, t := range o.Mesh.Tris {
			o.Mesh.VertsWorld[t.V1] = matTransform.MultiplyByVec3(o.Mesh.Verts[t.V1])
			o.Mesh.VertsWorld[t.V2] = matTransform.MultiplyByVec3(o.Mesh.Verts[t.V2])
			o.Mesh.VertsWorld[t.V3] = matTransform.MultiplyByVec3(o.Mesh.Verts[t.V3])

			o.Mesh.NormalsWorld[t.N1] = matRoation.MultiplyByVec3(o.Mesh.Normals[t.N1])
			o.Mesh.NormalsWorld[t.N2] = matRoation.MultiplyByVec3(o.Mesh.Normals[t.N2])
			o.Mesh.NormalsWorld[t.N3] = matRoation.MultiplyByVec3(o.Mesh.Normals[t.N3])

			if !t.BackFaceCulling(o.Mesh.VertsWorld, o.Mesh.NormalsWorld) {
				continue
			}

			v1 := mesh.ClippedVertex{
				V: o.Mesh.VertsWorld[t.V1],
				N: o.Mesh.NormalsWorld[t.N1],
				U: o.Mesh.UV[t.U1],
			}

			v2 := mesh.ClippedVertex{
				V: o.Mesh.VertsWorld[t.V2],
				N: o.Mesh.NormalsWorld[t.N2],
				U: o.Mesh.UV[t.U2],
			}

			v3 := mesh.ClippedVertex{
				V: o.Mesh.VertsWorld[t.V3],
				N: o.Mesh.NormalsWorld[t.N3],
				U: o.Mesh.UV[t.U3],
			}

			r.sHoutputList = r.sHoutputList[:0]
			r.sHinputList = r.sHinputList[:3]

			r.sHoutputList = append(r.sHoutputList, v1)
			r.sHoutputList = append(r.sHoutputList, v2)
			r.sHoutputList = append(r.sHoutputList, v3)

			// Sutherland–Hodgman algorithm
			for _, plane := range r.scene.ActiveCam.Frustum.Planes {
				r.sHinputList, r.sHoutputList = r.sHoutputList, r.sHinputList[:0]

				prevI := 0
				for i := range len(r.sHinputList) {
					prevI = i - 1
					if prevI < 0 {
						prevI = len(r.sHinputList) - 1
					}
					currentPoint := r.sHinputList[i]
					prevPoint := r.sHinputList[prevI]

					cp := plane.SignedDistanceToPoint(currentPoint.V)
					pp := plane.SignedDistanceToPoint(prevPoint.V)

					if cp > 0 {
						if pp <= 0 {
							ratio := cp / (cp - pp)
							intersection := mesh.ClippedVertex{
								V: currentPoint.V.LerpTo(prevPoint.V, ratio),
								N: currentPoint.N.LerpTo(prevPoint.N, ratio),
								U: currentPoint.U.LerpTo(prevPoint.U, ratio),
							}
							r.sHoutputList = append(r.sHoutputList, intersection)
						}
						r.sHoutputList = append(r.sHoutputList, currentPoint)
					} else if pp > 0 {
						ratio := cp / (cp - pp)
						intersection := mesh.ClippedVertex{
							V: currentPoint.V.LerpTo(prevPoint.V, ratio),
							N: currentPoint.N.LerpTo(prevPoint.N, ratio),
							U: currentPoint.U.LerpTo(prevPoint.U, ratio),
						}
						r.sHoutputList = append(r.sHoutputList, intersection)
					}
				}
			}

			if len(r.sHoutputList) > 2 {
				for i := 1; i < len(r.sHoutputList)-1; i++ {
					if r.RenderMultithreaded {
						r.trianglesBuffer = append(
							r.trianglesBuffer,
							mesh.NewFullTriangle(
								r.sHoutputList[0],
								r.sHoutputList[i],
								r.sHoutputList[i+1],
								o.Mesh.Texture,
							),
						)
					} else {
						r.RenderTriangle(
							r.sHoutputList[0],
							r.sHoutputList[i],
							r.sHoutputList[i+1],
							o.Mesh.Texture,
						)
					}

					if r.scene.ActiveCam.RenderWire {
						r.DrawWireframeTriangle(
							r.sHoutputList[0],
							r.sHoutputList[i],
							r.sHoutputList[i+1],
						)
					}
				}
			}
		}

		if r.RenderMultithreaded {
			if len(r.trianglesBuffer) > 0 {
				for i := range r.tiles {
					r.tiles[i].ResetBuff()
				}

				r.assignTrianglesToTiles()
				if r.RenderTileBoundaries {
					r.drawTileBoundaries()
				}

				r.wg.Add(len(r.tiles))
				for i := range r.tiles {
					r.indexes <- i
				}
				r.wg.Wait()
			}
		}
		// break
	}
}

func (r *Renderer) Render() {
	r.scene.ActiveCam.ClearCanvas()
	r.renderMeshs()
}

func (r *Renderer) ToggleMultithreaded() {
	r.RenderMultithreaded = !r.RenderMultithreaded
}

func (r *Renderer) ToggleTileBoundaryRender() {
	r.RenderTileBoundaries = !r.RenderTileBoundaries
}

func (r *Renderer) ToggleTriangleBoundaryRender() {
	r.RenderTriangleBoundaries = !r.RenderTriangleBoundaries
}
