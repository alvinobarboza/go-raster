package renderer

import (
	"math"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func (r *Renderer) RenderSimple() {
	r.scene.ActiveCam.ClearCanvas()
	r.scene.UpdateLights()

	for _, o := range r.scene.Objects {
		r.renderMeshSimple(o)
	}

	state := r.RenderLight
	r.RenderLight = false
	if r.scene.SkyBox != nil {
		r.renderMeshSimple(r.scene.SkyBox)
	}
	r.RenderLight = state
}

func (r *Renderer) renderMeshSimple(o *mesh.Model) {
	matTransform := r.scene.ActiveCam.Transforms.MatrixTransforms.MultiplyByMatrix(o.Transforms.MatrixTransforms)
	matRotation := r.scene.ActiveCam.Transforms.RotationMat.MultiplyByMatrix(o.Transforms.RotationMat)

	o.BoundingSphere.CenterWord = matTransform.MultiplyByVec3(o.BoundingSphere.Center)

	if !r.scene.ActiveCam.Frustum.IsBoundsInsideFrustum(&o.BoundingSphere) {
		return
	}

	for i := range len(o.Mesh.Verts) {
		o.Mesh.VertsWorld[i] = matTransform.MultiplyByVec3(o.Mesh.Verts[i])
	}

	for i := range len(o.Mesh.Normals) {
		o.Mesh.NormalsWorld[i] = matRotation.MultiplyByVec3(o.Mesh.Normals[i])
	}

	for _, t := range o.Mesh.Tris {
		if !t.BackFaceCulling(o.Mesh.VertsWorld, o.Mesh.NormalsWorld) {
			continue
		}

		v1 := mesh.ClippedVertex{V: o.Mesh.VertsWorld[t.V1], N: o.Mesh.NormalsWorld[t.N1], U: o.Mesh.UV[t.U1]}
		v2 := mesh.ClippedVertex{V: o.Mesh.VertsWorld[t.V2], N: o.Mesh.NormalsWorld[t.N2], U: o.Mesh.UV[t.U2]}
		v3 := mesh.ClippedVertex{V: o.Mesh.VertsWorld[t.V3], N: o.Mesh.NormalsWorld[t.N3], U: o.Mesh.UV[t.U3]}

		r.sHoutputList = r.sHoutputList[:0]
		r.sHinputList = r.sHinputList[:3]

		r.sHoutputList = append(r.sHoutputList, v1)
		r.sHoutputList = append(r.sHoutputList, v2)
		r.sHoutputList = append(r.sHoutputList, v3)

		r.clipTriangles()

		if len(r.sHoutputList) > 2 {
			for i := 1; i < len(r.sHoutputList)-1; i++ {
				triangle := r.scene.ActiveCam.ProjectTriangle(
					r.sHoutputList[0],
					r.sHoutputList[i],
					r.sHoutputList[i+1],
					o.Mesh.Texture, o.Mesh.Normal, o.Mesh.Specular)

				triangle.ShaderSmooth = t.ShaderSmooth
				r.renderTriangleSimple(&triangle)
			}
		}
	}
}

func (r *Renderer) renderTriangleSimple(triangle *mesh.FullTriangle) {
	// Clamp to screen boundaries instead of tiles
	minY := maths.Maxf(triangle.Aabb2.Min.Y, 0)
	maxY := maths.Minf(triangle.Aabb2.Max.Y, float32(r.scene.ActiveCam.Height))
	minX := maths.Maxf(triangle.Aabb2.Min.X, 0)
	maxX := maths.Minf(triangle.Aabb2.Max.X, float32(r.scene.ActiveCam.Width))

	deltaW0Col := triangle.SPV0.Y - triangle.SPV1.Y
	deltaW1Col := triangle.SPV1.Y - triangle.SPV2.Y
	deltaW2Col := triangle.SPV2.Y - triangle.SPV0.Y

	deltaW0Row := triangle.SPV1.X - triangle.SPV0.X
	deltaW1Row := triangle.SPV2.X - triangle.SPV1.X
	deltaW2Row := triangle.SPV0.X - triangle.SPV2.X

	bias0, bias1, bias2 := float32(0), float32(0), float32(0)
	if mesh.IsEdgeTopOrLeft(triangle.SPV0, triangle.SPV1) {
		bias0 = -0.0001
	}
	if mesh.IsEdgeTopOrLeft(triangle.SPV1, triangle.SPV2) {
		bias1 = -0.0001
	}
	if mesh.IsEdgeTopOrLeft(triangle.SPV2, triangle.SPV0) {
		bias2 = -0.0001
	}

	area := 1.0 / mesh.EdgeCross(triangle.SPV0, triangle.SPV1, triangle.SPV2)

	p := transforms.Vec2{X: minX + 0.5, Y: minY + 0.5}

	w0Row := mesh.EdgeCross(triangle.SPV0, triangle.SPV1, p) + bias0
	w1Row := mesh.EdgeCross(triangle.SPV1, triangle.SPV2, p) + bias1
	w2Row := mesh.EdgeCross(triangle.SPV2, triangle.SPV0, p) + bias2

	edge1 := triangle.V2z.Subtract(triangle.V1z)
	edge2 := triangle.V3z.Subtract(triangle.V1z)
	deltaUV1 := triangle.V2.U.Subtract(triangle.V1.U)
	deltaUV2 := triangle.V3.U.Subtract(triangle.V1.U)

	// Prevent divide by zero in UV math
	det := deltaUV1.X*deltaUV2.Y - deltaUV2.X*deltaUV1.Y
	tangent := transforms.NewVec3(1, 0, 0) // Default safe value
	if maths.Abs(det) > 0.000001 {
		f := 1.0 / det
		tangent = transforms.NewVec3(
			f*(deltaUV2.Y*edge1.X-deltaUV1.Y*edge2.X),
			f*(deltaUV2.Y*edge1.Y-deltaUV1.Y*edge2.Y),
			f*(deltaUV2.Y*edge1.Z-deltaUV1.Y*edge2.Z),
		).Normalized()
	}

	for y := minY; y < maxY; y++ {
		w0, w1, w2 := w0Row, w1Row, w2Row

		for x := minX; x < maxX; x++ {
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				alpha := w1 * area
				beta := w2 * area
				gama := w0 * area

				depth := triangle.DepthZ1*alpha + triangle.DepthZ2*beta + triangle.DepthZ3*gama
				xx, yy := uint(x), uint(y)

				if r.scene.ActiveCam.DepthPass(xx, yy, depth) {
					uvCoord := triangle.UV1z.Scale(alpha).Add(triangle.UV2z.Scale(beta)).Add(triangle.UV3z.Scale(gama)).Divide(depth)
					pColor := triangle.Texture.TexelColor(uvCoord)

					if r.RenderLight {
						nCoord := triangle.N1z
						if triangle.ShaderSmooth {
							nCoord = triangle.N1z.Scale(alpha).Add(triangle.N2z.Scale(beta)).Add(triangle.N3z.Scale(gama)).Divide(depth).Normalized()
						}

						if triangle.Normal != nil {
							nMapColor := triangle.Normal.TexelColor(uvCoord)
							nMap := transforms.Vec3{
								X: (float32(nMapColor.R) / 127.5) - 1.0,
								Y: (float32(nMapColor.G) / 127.5) - 1.0,
								Z: (float32(nMapColor.B) / 127.5) - 1.0,
							}

							dotNT := nCoord.DotByVec3(tangent)
							t := tangent.Subtract(nCoord.Scale(dotNT)).Normalized()
							b := nCoord.Cross(t)
							nCoord = t.Scale(nMap.X).Add(b.Scale(nMap.Y)).Add(nCoord.Scale(nMap.Z)).Normalized()
						}

						viewDir := transforms.Vec3{}
						specularStrength := float32(0)
						if triangle.Specular != nil {
							fragPos := triangle.V1z.Scale(alpha).Add(triangle.V2z.Scale(beta)).Add(triangle.V3z.Scale(gama)).Divide(depth)
							viewDir = fragPos.Normalized().Scale(-1)
							specularStrength = triangle.Specular.TexelIntensity(uvCoord)
						}

						result := transforms.NewVec3(float32(pColor.R)/255, float32(pColor.G)/255, float32(pColor.B)/255)

						for _, l := range r.scene.Lights {
							ambient := l.Color.Scale(r.scene.AmbientLightStrength)
							lightIntensity := maths.Maxf(nCoord.DotByVec3(l.DirectionWorld), 0)
							diff := l.Color.Scale(lightIntensity * l.Intensity)

							result = result.Multiply(diff.Add(ambient))

							if specularStrength > 0 {
								halfwayDir := l.DirectionWorld.Add(viewDir).Normalized()
								dot := maths.Maxf(nCoord.DotByVec3(halfwayDir), 0.0)
								spec := float32(math.Pow(float64(dot), 200))
								specular := l.Color.Scale(specularStrength * spec)
								result = result.Add(specular)
							}
						}

						pColor.R = uint8(float32(255) * maths.Minf(result.X, 1))
						pColor.G = uint8(float32(255) * maths.Minf(result.Y, 1))
						pColor.B = uint8(float32(255) * maths.Minf(result.Z, 1))
					}

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
