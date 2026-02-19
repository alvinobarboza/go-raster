package main

import (
	"image/color"
	"math"
)

type NDCPoint struct {
	X, Y float32
}

type ScreenPoint struct {
	X, Y float32
}

func (sp *ScreenPoint) IsTopOrLeft(sp2 ScreenPoint) bool {
	edge := ScreenPoint{X: sp2.X - sp.X, Y: sp2.Y - sp.Y}
	isTopEdge := edge.Y == 0 && edge.X > 0
	isLeftEdge := edge.Y < 0
	return isTopEdge || isLeftEdge
}

type Plane struct {
	normal   Vec3
	distance float32
}

func NewPlane(p1, normal Vec3) Plane {
	normal = normal.Normalized()
	return Plane{
		normal:   normal,
		distance: normal.DotByVec3(p1),
	}
}

func (p *Plane) SignedDistanceToPoint(point Vec3) float32 {
	return p.normal.DotByVec3(point) - p.distance
}

type Frustum struct {
	topPlane, bottomPlane Plane

	rightPlane, leftPlane Plane

	farPlane, nearPlane Plane
}

func (f *Frustum) IsVertexInsideFrustum(p Vec3) bool {
	np := f.nearPlane.SignedDistanceToPoint(p)
	fp := f.farPlane.SignedDistanceToPoint(p)
	rp := f.rightPlane.SignedDistanceToPoint(p)
	lp := f.leftPlane.SignedDistanceToPoint(p)
	tp := f.topPlane.SignedDistanceToPoint(p)
	bp := f.bottomPlane.SignedDistanceToPoint(p)

	if np > 0 && fp > 0 && rp > 0 && lp > 0 && tp > 0 && bp > 0 {
		return true
	}
	return false
}

func (f *Frustum) IsBoundsInsideFrustum(b *BoundingSphere) bool {
	np := f.nearPlane.SignedDistanceToPoint(b.centerWord)
	fp := f.farPlane.SignedDistanceToPoint(b.centerWord)
	rp := f.rightPlane.SignedDistanceToPoint(b.centerWord)
	lp := f.leftPlane.SignedDistanceToPoint(b.centerWord)
	tp := f.topPlane.SignedDistanceToPoint(b.centerWord)
	bp := f.bottomPlane.SignedDistanceToPoint(b.centerWord)

	r := b.radius

	if np < -r || fp < -r || rp < -r || lp < -r || tp < -r || bp < -r {
		return false
	}

	return true
}

type Camera struct {
	canvas      []color.RGBA
	depthBuffer []float32

	fovAngle    float32
	aspectRatio float32
	fovScaling  float32
	zNear       float32
	zFar        float32
	sensitivity float32

	updateView  bool
	renderDepth bool
	renderWire  bool

	width, height         uint
	halfWidth, halfHeight float32

	transforms Transforms

	frustum Frustum
}

func NewCamera(w, h uint, sensitivity, zNear, zFar, fovAngle float32, pos, rot Vec3) Camera {
	c := Camera{
		fovAngle:    fovAngle,
		fovScaling:  FovScaling(fovAngle),
		zNear:       zNear,
		zFar:        zFar,
		sensitivity: sensitivity,
		updateView:  true,
		renderDepth: false,
		renderWire:  false,
		transforms: Transforms{
			scale:            NewVec3(1, 1, 1),
			rotation:         rot,
			position:         pos,
			forwardDirection: NewVec3(0, 0, 1),
		},
	}

	c.transforms.UpdateCameraTransforms()
	c.UpdateCanvasSize(w, h)
	return c
}

func (c *Camera) UpdateCanvasSize(w, h uint) {
	c.width = w
	c.height = h
	c.halfWidth = float32(w) / 2
	c.halfHeight = float32(h) / 2
	c.aspectRatio = float32(w) / float32(h)
	c.canvas = make([]color.RGBA, w*h)
	c.depthBuffer = make([]float32, w*h)

	c.CalculateFrustum()
}

func (c *Camera) ClearCanvas() {
	for i := range len(c.canvas) {
		c.canvas[i] = Gray
		c.depthBuffer[i] = 0
	}
}

func (c *Camera) ProjectVertexToNDC(v Vec3) NDCPoint {
	zXInverse := 1 / (v.Z * c.aspectRatio)
	zYInverse := 1 / v.Z
	return NDCPoint{
		X: (v.X * c.fovScaling) * zXInverse,
		Y: (v.Y * c.fovScaling) * zYInverse,
	}
}

func (c *Camera) NDCtoScreen(p NDCPoint) ScreenPoint {
	x := (p.X + 1) * c.halfWidth
	y := (1 - p.Y) * c.halfHeight

	return ScreenPoint{
		X: x,
		Y: y,
	}
}

func (c *Camera) DepthPass(x, y uint, depth float32) bool {
	if x >= c.width || y >= c.height {
		return false
	}
	i := y*c.width + x
	if c.depthBuffer[i] > depth {
		return false
	}
	c.depthBuffer[i] = depth
	return true
}

func (c *Camera) PutPixel(x, y uint, cl color.RGBA, depth float32) {
	i := y*c.width + x

	if c.renderDepth {
		c.canvas[i].A = 255
		c.canvas[i].R = uint8(255 * depth)
		c.canvas[i].G = uint8(255 * depth)
		c.canvas[i].B = uint8(255 * depth)
		return
	}

	c.canvas[i] = cl
}

func (c *Camera) MoveBackForwad(unit float32) {
	if unit == 0 {
		return
	}

	rotMat := NewRotationMatrix(c.transforms.rotation)

	direction := rotMat.MultiplyByVec3(c.transforms.forwardDirection)
	normalDirection := direction.Normalized()

	c.transforms.position = c.transforms.position.Add(normalDirection.Scale(unit))

	c.transforms.UpdateCameraTransforms()
}

func (c *Camera) MoveSideways(unit float32) {
	if unit == 0 {
		return
	}

	rotMat := NewRotationMatrix(c.transforms.rotation)

	direction := rotMat.MultiplyByVec3(c.transforms.forwardDirection)
	sideDirection := direction.Cross(NewVec3(0, 1, 0))
	normalDirection := sideDirection.Normalized()

	c.transforms.position = c.transforms.position.Add(normalDirection.Scale(unit))

	c.transforms.UpdateCameraTransforms()
}

func (c *Camera) UpdateRotation(x float32, y float32) {
	if !c.updateView {
		return
	}

	if x == 0 && y == 0 {
		return
	}

	c.transforms.rotation.X -= y * c.sensitivity
	c.transforms.rotation.Y -= x * c.sensitivity

	if c.transforms.rotation.X > 89 {
		c.transforms.rotation.X = 89
	}

	if c.transforms.rotation.X < -89 {
		c.transforms.rotation.X = -89
	}

	c.transforms.UpdateCameraTransforms()
}

func (c *Camera) ToggleViewLock() {
	c.updateView = !c.updateView
}

func (c *Camera) ToggleDepthRender() {
	c.renderDepth = !c.renderDepth
}

func (c *Camera) ToggleWireRender() {
	c.renderWire = !c.renderWire
}

func (c *Camera) MoveVetically(unit float32) {
	c.transforms.position.Y += unit
	c.transforms.UpdateCameraTransforms()
}

func (c *Camera) CalculateFrustum() {
	camFront := c.transforms.forwardDirection
	camRight := NewVec3(1, 0, 0)
	camUp := NewVec3(0, 1, 0)
	camPos := NewVec3(0, 0, 0)

	halfVSide := c.zFar * float32(math.Tan(float64(c.fovAngle*DegToRad)*.5))
	halfHSide := halfVSide * c.aspectRatio
	frontMultFar := camFront.Scale(c.zFar)

	c.frustum.nearPlane = NewPlane(camPos.Add(camFront.Scale(c.zNear)), camFront)
	c.frustum.farPlane = NewPlane(camPos.Add(frontMultFar), camFront.Scale(-1))
	c.frustum.rightPlane = NewPlane(camPos, frontMultFar.Add(camRight.Scale(halfHSide)).Cross(camUp))
	c.frustum.leftPlane = NewPlane(camPos, camUp.Cross(frontMultFar.Subtract(camRight.Scale(halfHSide))))
	c.frustum.topPlane = NewPlane(camPos, frontMultFar.Subtract(camUp.Scale(halfVSide)).Cross(camRight))
	c.frustum.bottomPlane = NewPlane(camPos, camRight.Cross(frontMultFar.Add(camUp.Scale(halfVSide))))
}
