package camera

import (
	"math"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type ShadowCamera struct {
	DepthBuffer []float32

	FovAngle   float32
	FovScaling float32
	ZNear      float32
	ZFar       float32

	Width, Height         uint
	HalfWidth, HalfHeight float32

	Transforms transforms.Transforms

	Frustum Frustum
}

func NewShadowCamera(w, h uint, zNear, zFar, fovAngle float32, pos, rot transforms.Vec3) *ShadowCamera {
	c := ShadowCamera{
		FovAngle:   fovAngle,
		FovScaling: FovScaling(fovAngle),
		ZNear:      zNear,
		ZFar:       zFar,
		Transforms: transforms.Transforms{
			Scale:            transforms.NewVec3(1, 1, 1),
			Rotation:         rot,
			Position:         pos,
			ForwardDirection: transforms.NewVec3(0, 0, 1),
		},
	}

	c.Transforms.UpdateCameraTransforms()
	c.UpdateCanvasSize(w, h)
	return &c
}

func (c *ShadowCamera) UpdateCanvasSize(w, h uint) {
	c.Width = w
	c.Height = h
	c.HalfWidth = float32(w) / 2
	c.HalfHeight = float32(h) / 2
	c.DepthBuffer = make([]float32, w*h)

	c.CalculateFrustum()
}

func (c *ShadowCamera) ClearDepth() {
	for i := range len(c.DepthBuffer) {
		c.DepthBuffer[i] = 0
	}
}

func (c *ShadowCamera) ProjectVertexToNDC(v transforms.Vec3) transforms.Vec2 {
	zInverse := 1 / v.Z
	return transforms.Vec2{
		X: (v.X * c.FovScaling) * zInverse,
		Y: (v.Y * c.FovScaling) * zInverse,
	}
}

func (c *ShadowCamera) NDCtoScreen(p transforms.Vec2) transforms.Vec2 {
	x := (p.X + 1) * c.HalfWidth
	y := (1 - p.Y) * c.HalfHeight

	return transforms.Vec2{
		X: x,
		Y: y,
	}
}

func (c *ShadowCamera) ScreenToNDC(x, y float32) transforms.Vec2 {
	return transforms.Vec2{
		X: (x / c.HalfWidth) - 1,
		Y: 1 - (y / c.HalfHeight),
	}
}

func (c *ShadowCamera) NDCToVertexRay(p transforms.Vec2) transforms.Vec3 {
	return transforms.NewVec3(p.X/c.FovScaling, p.Y/c.FovScaling, 1).Normalized()
}

func (c *ShadowCamera) ProjectTriangle(v1, v2, v3 mesh.ClippedVertex, t, n, s *mesh.Texture) mesh.FullTriangle {
	triangle := mesh.NewFullTriangle(v1, v2, v3, t, n, s)

	va := c.ProjectVertexToNDC(v1.V)
	vb := c.ProjectVertexToNDC(v2.V)
	vc := c.ProjectVertexToNDC(v3.V)

	triangle.SPV0 = c.NDCtoScreen(va)
	triangle.SPV1 = c.NDCtoScreen(vb)
	triangle.SPV2 = c.NDCtoScreen(vc)

	triangle.DepthZ1 = 1 / v1.V.Z
	triangle.DepthZ2 = 1 / v2.V.Z
	triangle.DepthZ3 = 1 / v3.V.Z

	triangle.V1z = v1.V.Scale(triangle.DepthZ1)
	triangle.V2z = v2.V.Scale(triangle.DepthZ2)
	triangle.V3z = v3.V.Scale(triangle.DepthZ3)

	triangle.UV1z = v1.U.Scale(triangle.DepthZ1)
	triangle.UV2z = v2.U.Scale(triangle.DepthZ2)
	triangle.UV3z = v3.U.Scale(triangle.DepthZ3)

	triangle.N1z = v1.N.Scale(triangle.DepthZ1)
	triangle.N2z = v2.N.Scale(triangle.DepthZ2)
	triangle.N3z = v3.N.Scale(triangle.DepthZ3)

	triangle.Aabb2 = mesh.NewAABB2(
		maths.Floor32(maths.Minf(triangle.SPV0.X, maths.Minf(triangle.SPV1.X, triangle.SPV2.X))),
		maths.Floor32(maths.Minf(triangle.SPV0.Y, maths.Minf(triangle.SPV1.Y, triangle.SPV2.Y))),
		maths.Ceil32(maths.Maxf(triangle.SPV0.X, maths.Maxf(triangle.SPV1.X, triangle.SPV2.X))),
		maths.Ceil32(maths.Maxf(triangle.SPV0.Y, maths.Maxf(triangle.SPV1.Y, triangle.SPV2.Y))),
	)

	return triangle
}

func (c *ShadowCamera) DepthPass(x, y uint, depth float32) {
	if x >= c.Width || y >= c.Height {
		return
	}
	i := y*c.Width + x
	if c.DepthBuffer[i] <= depth {
		c.DepthBuffer[i] = depth
	}
}

func (c *ShadowCamera) CalculateFrustum() {
	camFront := c.Transforms.ForwardDirection
	camRight := transforms.NewVec3(1, 0, 0)
	camUp := transforms.NewVec3(0, 1, 0)
	camPos := transforms.NewVec3(0, 0, 0)

	halfVSide := c.ZFar * float32(math.Tan(float64(c.FovAngle*transforms.DegToRad)*.5))
	frontMultFar := camFront.Scale(c.ZFar)

	c.Frustum.Planes[NearPn] = NewPlane(camPos.Add(camFront.Scale(c.ZNear)), camFront)
	c.Frustum.Planes[FarPn] = NewPlane(camPos.Add(frontMultFar), camFront.Scale(-1))
	c.Frustum.Planes[RightPn] = NewPlane(camPos, frontMultFar.Add(camRight.Scale(halfVSide)).Cross(camUp))
	c.Frustum.Planes[LeftPn] = NewPlane(camPos, camUp.Cross(frontMultFar.Subtract(camRight.Scale(halfVSide))))
	c.Frustum.Planes[TopPn] = NewPlane(camPos, frontMultFar.Subtract(camUp.Scale(halfVSide)).Cross(camRight))
	c.Frustum.Planes[BottomPn] = NewPlane(camPos, camRight.Cross(frontMultFar.Add(camUp.Scale(halfVSide))))
}
