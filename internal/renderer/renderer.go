package renderer

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/scene"
	"github.com/alvinobarboza/go-raster/internal/shapes"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

const MinimumTileSize = 60

type Renderer struct {
	scene   *scene.Scene
	wg      sync.WaitGroup
	mt      sync.Mutex
	indexes chan int

	sHoutputList    []mesh.ClippedVertex
	sHinputList     []mesh.ClippedVertex
	trianglesBuffer []mesh.FullTriangle

	tiles    []*ScreenTile
	tileSize uint

	biggestTriCount int

	RenderTileBoundaries     bool
	RenderTriangleBoundaries bool
	RenderMultithreaded      bool
}

// init after loading models,
// otherwise triangle buffer will have 0 size
// TODO -> somehow get around this
func NewRenderer(threads, tileSize uint) *Renderer {
	if tileSize < MinimumTileSize {
		tileSize = MinimumTileSize
	}

	r := &Renderer{
		sHoutputList:             make([]mesh.ClippedVertex, 9),
		sHinputList:              make([]mesh.ClippedVertex, 9),
		indexes:                  make(chan int, threads),
		RenderTileBoundaries:     false,
		RenderTriangleBoundaries: false,
		RenderMultithreaded:      true,
		tileSize:                 tileSize,
	}

	for i := range threads {
		go r.renderTriangleParallel(i)
	}

	return r
}

func (r *Renderer) UpdateTiles() {
	if r.tileSize > r.scene.ActiveCam.Width {
		r.tileSize = r.scene.ActiveCam.Width / 2
	}

	r.tiles = NewTileSet(
		float32(r.scene.ActiveCam.Width),
		float32(r.scene.ActiveCam.Height),
		float32(r.tileSize),
		r.biggestTriCount,
	)
}

func (r *Renderer) AddActiveScene(s *scene.Scene) {
	r.scene = s
	for _, o := range s.Objects {
		if len(o.Mesh.Tris) > r.biggestTriCount {
			r.biggestTriCount = len(o.Mesh.Tris) * 9
		}
	}

	r.UpdateTiles()

	r.trianglesBuffer = make([]mesh.FullTriangle, r.biggestTriCount)
}

func (r *Renderer) DrawLine(a, b transforms.Vec2, cl color.RGBA) {
	dx := b.X - a.X
	dy := b.Y - a.Y

	if maths.Abs(dx) > maths.Abs(dy) {
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

func (r *Renderer) DrawWireframeTriangle(t mesh.FullTriangle) {
	r.DrawLine(t.SPV0, t.SPV1, shapes.Black)
	r.DrawLine(t.SPV1, t.SPV2, shapes.Black)
	r.DrawLine(t.SPV2, t.SPV0, shapes.Black)
}

func (r *Renderer) DrawWireframeTriangleFromBuff() {
	for _, tri := range r.trianglesBuffer {
		r.DrawWireframeTriangle(tri)
	}
}

func (r *Renderer) drawAABB(aabb mesh.AABB2, c color.RGBA) {
	r.DrawLine(transforms.Vec2{X: aabb.Min.X, Y: aabb.Min.Y}, transforms.Vec2{X: aabb.Min.X, Y: aabb.Max.Y}, c)
	r.DrawLine(transforms.Vec2{X: aabb.Min.X, Y: aabb.Max.Y}, transforms.Vec2{X: aabb.Max.X, Y: aabb.Max.Y}, c)
	r.DrawLine(transforms.Vec2{X: aabb.Max.X, Y: aabb.Max.Y}, transforms.Vec2{X: aabb.Max.X, Y: aabb.Min.Y}, c)
	r.DrawLine(transforms.Vec2{X: aabb.Max.X, Y: aabb.Min.Y}, transforms.Vec2{X: aabb.Min.X, Y: aabb.Min.Y}, c)
}

func (r *Renderer) DrawTriangleBoundary(tri mesh.FullTriangle) {
	r.drawAABB(tri.Aabb2, shapes.Green)
}

func (r *Renderer) DrawTriangleBoundaryFromBuff() {
	for _, tri := range r.trianglesBuffer {
		r.DrawTriangleBoundary(tri)
	}
}

func (r *Renderer) drawTileBoundaries() {
	for i := range r.tiles {
		if r.tiles[i].IsActive {
			r.drawAABB(r.tiles[i].Aabb, shapes.White)
		}
	}
}

func (r *Renderer) renderTriangleParallel(id uint) {
	fmt.Println("Thread ID:", id)
	for i := range r.indexes {
		triangles := r.tiles[i].Triangles()
		tileAabb := r.tiles[i].Aabb

		for _, i := range triangles {
			tri := r.trianglesBuffer[i]

			// check to only run in tile bounds
			tri.Aabb2.Min.Y = maths.Maxf(tri.Aabb2.Min.Y, tileAabb.Min.Y)
			tri.Aabb2.Max.Y = maths.Minf(tri.Aabb2.Max.Y, tileAabb.Max.Y)

			tri.Aabb2.Min.X = maths.Maxf(tri.Aabb2.Min.X, tileAabb.Min.X)
			tri.Aabb2.Max.X = maths.Minf(tri.Aabb2.Max.X, tileAabb.Max.X)

			r.RenderTriangle(tri)
		}

		r.wg.Done()
	}
}

func (r *Renderer) assignTrianglesToTiles() {
	for i, t := range r.trianglesBuffer {

		for j := range r.tiles {
			if r.tiles[j].TileTriangleCollision(t.SPV0, t.SPV1, t.SPV2, t.Aabb2) {
				r.tiles[j].AddTriangle(i)
				r.tiles[j].IsActive = true
			}
		}
	}
}

func (r *Renderer) RenderTriangle(triangle mesh.FullTriangle) {

	deltaW0Col := triangle.SPV0.Y - triangle.SPV1.Y
	deltaW1Col := triangle.SPV1.Y - triangle.SPV2.Y
	deltaW2Col := triangle.SPV2.Y - triangle.SPV0.Y

	deltaW0Row := triangle.SPV1.X - triangle.SPV0.X
	deltaW1Row := triangle.SPV2.X - triangle.SPV1.X
	deltaW2Row := triangle.SPV0.X - triangle.SPV2.X

	bias0 := float32(0)
	bias1 := float32(0)
	bias2 := float32(0)

	if mesh.IsEdgeTopOrLeft(triangle.SPV0, triangle.SPV1) {
		bias0 = -0.0001
	}

	if mesh.IsEdgeTopOrLeft(triangle.SPV1, triangle.SPV2) {
		bias1 = -0.0001
	}

	if mesh.IsEdgeTopOrLeft(triangle.SPV2, triangle.SPV0) {
		bias2 = -0.0001
	}

	area := mesh.EdgeCross(triangle.SPV0, triangle.SPV1, triangle.SPV2)
	area = 1 / area

	// pixel's center
	p := transforms.Vec2{X: triangle.Aabb2.Min.X + 0.5, Y: triangle.Aabb2.Min.Y + 0.5}

	w0Row := mesh.EdgeCross(triangle.SPV0, triangle.SPV1, p) + bias0
	w1Row := mesh.EdgeCross(triangle.SPV1, triangle.SPV2, p) + bias1
	w2Row := mesh.EdgeCross(triangle.SPV2, triangle.SPV0, p) + bias2

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
	for y := triangle.Aabb2.Min.Y; y <= triangle.Aabb2.Max.Y; y++ {
		w0 := w0Row
		w1 := w1Row
		w2 := w2Row

		for x := triangle.Aabb2.Min.X; x <= triangle.Aabb2.Max.X; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				alpha := w1 * area
				beta := w2 * area
				gama := w0 * area

				depth := triangle.DepthZ1*alpha + triangle.DepthZ2*beta + triangle.DepthZ3*gama

				xx, yy := uint(x), uint(y)
				if r.scene.ActiveCam.DepthPass(xx, yy, depth) {
					uv1 := triangle.UV1z.Scale(alpha)
					uv2 := triangle.UV2z.Scale(beta)
					uv3 := triangle.UV3z.Scale(gama)

					uvCoord := uv1.Add(uv2).Add(uv3).Divide(depth)
					pColor := triangle.Texture.TexelColor(uvCoord)

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

					triangle := r.scene.ActiveCam.ProjectTriangle(
						r.sHoutputList[0],
						r.sHoutputList[i],
						r.sHoutputList[i+1],
						o.Mesh.Texture)

					if r.RenderMultithreaded {
						r.trianglesBuffer = append(r.trianglesBuffer, triangle)
					} else {
						r.RenderTriangle(triangle)

						if r.scene.ActiveCam.RenderWire {
							r.DrawWireframeTriangle(triangle)
						}
						if r.RenderTriangleBoundaries {
							r.DrawTriangleBoundary(triangle)
						}
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

				r.wg.Add(len(r.tiles))
				for i := range r.tiles {
					r.indexes <- i
				}
				r.wg.Wait()

				if r.scene.ActiveCam.RenderWire {
					r.DrawWireframeTriangleFromBuff()
				}

				if r.RenderTriangleBoundaries {
					r.DrawTriangleBoundaryFromBuff()
				}
			}
		}
		// break
	}
}

func (r *Renderer) Render() {
	r.scene.ActiveCam.ClearCanvas()
	for i := range r.tiles {
		r.tiles[i].IsActive = false
	}
	r.renderMeshs()
	if r.RenderTileBoundaries {
		r.drawTileBoundaries()
	}
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
