package camera

import (
	"image/color"
	"math"

	"github.com/alvinobarboza/go-raster/internal/shapes"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type Camera struct {
	Canvas      []color.RGBA
	DepthBuffer []float32

	FovAngle    float32
	AspectRatio float32
	FovScaling  float32
	ZNear       float32
	ZFar        float32
	Sensitivity float32

	UpdateView  bool
	RenderDepth bool
	RenderWire  bool

	Width, Height         uint
	HalfWidth, HalfHeight float32

	Transforms transforms.Transforms

	Frustum Frustum
}

func NewCamera(w, h uint, sensitivity, zNear, zFar, fovAngle float32, pos, rot transforms.Vec3) *Camera {
	c := Camera{
		FovAngle:    fovAngle,
		FovScaling:  FovScaling(fovAngle),
		ZNear:       zNear,
		ZFar:        zFar,
		Sensitivity: sensitivity,
		UpdateView:  true,
		RenderDepth: false,
		RenderWire:  false,
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

func (c *Camera) UpdateCanvasSize(w, h uint) {
	c.Width = w
	c.Height = h
	c.HalfWidth = float32(w) / 2
	c.HalfHeight = float32(h) / 2
	c.AspectRatio = float32(w) / float32(h)
	c.Canvas = make([]color.RGBA, w*h)
	c.DepthBuffer = make([]float32, w*h)

	c.CalculateFrustum()
}

func (c *Camera) ClearCanvas() {
	for i := range len(c.Canvas) {
		c.Canvas[i] = shapes.DarkGray
		c.DepthBuffer[i] = 0
	}
}

func (c *Camera) ProjectVertexToNDC(v transforms.Vec3) NDCPoint {
	zXInverse := 1 / (v.Z * c.AspectRatio)
	zYInverse := 1 / v.Z
	return NDCPoint{
		X: (v.X * c.FovScaling) * zXInverse,
		Y: (v.Y * c.FovScaling) * zYInverse,
	}
}

func (c *Camera) NDCtoScreen(p NDCPoint) ScreenPoint {
	x := (p.X + 1) * c.HalfWidth
	y := (1 - p.Y) * c.HalfHeight

	return ScreenPoint{
		X: x,
		Y: y,
	}
}

func (c *Camera) DepthPass(x, y uint, depth float32) bool {
	if x >= c.Width || y >= c.Height {
		return false
	}
	i := y*c.Width + x
	if c.DepthBuffer[i] > depth {
		return false
	}
	c.DepthBuffer[i] = depth
	return true
}

func (c *Camera) PutPixel(x, y uint, cl color.RGBA, depth float32) {
	i := y*c.Width + x

	if int(i) >= len(c.Canvas) {
		return
	}

	if c.RenderDepth {
		c.Canvas[i].A = 255
		c.Canvas[i].R = uint8(255 * depth)
		c.Canvas[i].G = uint8(255 * depth)
		c.Canvas[i].B = uint8(255 * depth)
		return
	}

	c.Canvas[i] = cl
}

func (c *Camera) MoveBackForwad(unit float32) {
	if unit == 0 {
		return
	}

	rotMat := transforms.NewRotationMatrix(c.Transforms.Rotation)

	direction := rotMat.MultiplyByVec3(c.Transforms.ForwardDirection)
	normalDirection := direction.Normalized()

	c.Transforms.Position = c.Transforms.Position.Add(normalDirection.Scale(unit))

	c.Transforms.UpdateCameraTransforms()
}

func (c *Camera) MoveSideways(unit float32) {
	if unit == 0 {
		return
	}

	rotMat := transforms.NewRotationMatrix(c.Transforms.Rotation)

	direction := rotMat.MultiplyByVec3(c.Transforms.ForwardDirection)
	sideDirection := direction.Cross(transforms.NewVec3(0, 1, 0))
	normalDirection := sideDirection.Normalized()

	c.Transforms.Position = c.Transforms.Position.Add(normalDirection.Scale(unit))

	c.Transforms.UpdateCameraTransforms()
}

func (c *Camera) UpdateRotation(x float32, y float32) {
	if !c.UpdateView {
		return
	}

	if x == 0 && y == 0 {
		return
	}

	c.Transforms.Rotation.X -= y * c.Sensitivity
	c.Transforms.Rotation.Y -= x * c.Sensitivity

	if c.Transforms.Rotation.X > 89 {
		c.Transforms.Rotation.X = 89
	}

	if c.Transforms.Rotation.X < -89 {
		c.Transforms.Rotation.X = -89
	}

	c.Transforms.UpdateCameraTransforms()
}

func (c *Camera) ToggleViewLock() {
	c.UpdateView = !c.UpdateView
}

func (c *Camera) ToggleDepthRender() {
	c.RenderDepth = !c.RenderDepth
}

func (c *Camera) ToggleWireRender() {
	c.RenderWire = !c.RenderWire
}

func (c *Camera) MoveVetically(unit float32) {
	c.Transforms.Position.Y += unit
	c.Transforms.UpdateCameraTransforms()
}

func (c *Camera) CalculateFrustum() {
	camFront := c.Transforms.ForwardDirection
	camRight := transforms.NewVec3(1, 0, 0)
	camUp := transforms.NewVec3(0, 1, 0)
	camPos := transforms.NewVec3(0, 0, 0)

	halfVSide := c.ZFar * float32(math.Tan(float64(c.FovAngle*transforms.DegToRad)*.5))
	halfHSide := halfVSide * c.AspectRatio
	frontMultFar := camFront.Scale(c.ZFar)

	c.Frustum.Planes[NearPn] = NewPlane(camPos.Add(camFront.Scale(c.ZNear)), camFront)
	c.Frustum.Planes[FarPn] = NewPlane(camPos.Add(frontMultFar), camFront.Scale(-1))
	c.Frustum.Planes[RightPn] = NewPlane(camPos, frontMultFar.Add(camRight.Scale(halfHSide)).Cross(camUp))
	c.Frustum.Planes[LeftPn] = NewPlane(camPos, camUp.Cross(frontMultFar.Subtract(camRight.Scale(halfHSide))))
	c.Frustum.Planes[TopPn] = NewPlane(camPos, frontMultFar.Subtract(camUp.Scale(halfVSide)).Cross(camRight))
	c.Frustum.Planes[BottomPn] = NewPlane(camPos, camRight.Cross(frontMultFar.Add(camUp.Scale(halfVSide))))
}

func FovScaling(angle float32) float32 {
	return float32(1 / math.Tan(float64(angle*transforms.DegToRad/2)))
}
