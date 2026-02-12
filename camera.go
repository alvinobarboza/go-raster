package main

import (
	"image/color"
	"math"
)

type NDCPoint struct {
	X, Y  float32
	color color.RGBA
}

type ScreenPoint struct {
	X, Y  float32
	color color.RGBA
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
		normal: normal,
		distance: normal.MultiplyByVec3(p1),
	}
}

func (p *Plane) SignedDistanceToPoint(point Vec3) float32 {
	return p.normal.MultiplyByVec3(point) - p.distance
}

type Frustum struct {
	topFace, bottomFace Plane

	rightFace, leftFace Plane

	farPlane, nearPlane Plane
}

type Camera struct {
	canvas      []color.RGBA
	fovAngle    float32
	aspectRatio float32
	fovScaling  float32
	zNear       float32
	sensitivity float32

	updateView bool

	width, height         uint
	halfWidth, halfHeight float32

	transforms Transforms

	frustum Frustum
}

func NewCamera(w, h uint, sensitivity, zNear, fovAngle float32, pos, rot Vec3) Camera {
	c := Camera{
		fovAngle:    fovAngle,
		fovScaling:  FovScaling(fovAngle),
		zNear:       zNear,
		sensitivity: sensitivity,
		updateView:  true,
		transforms: Transforms{
			scale:            NewVec3(1, 1, 1),
			rotation:         rot,
			position:         pos,
			forwardDirection: NewVec3(0, 0, 1),
		},
	}

	c.transforms.UpdateTransforms(true, true)
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
}

func (c *Camera) ClearCanvas() {
	for i := range len(c.canvas) {
		c.canvas[i] = Gray
	}
}

func (c *Camera) ProjectVertexToNDC(v Vec3, cl color.RGBA) NDCPoint {
	zXInverse := 1 / (v.Z * c.aspectRatio)
	zYInverse := 1 / v.Z
	return NDCPoint{
		X:     (v.X * c.fovScaling) * zXInverse,
		Y:     (v.Y * c.fovScaling) * zYInverse,
		color: cl,
	}
}

func (c *Camera) NDCtoScreen(p NDCPoint) ScreenPoint {
	x := (p.X + 1) * c.halfWidth
	y := (1 - p.Y) * c.halfHeight

	return ScreenPoint{
		X:     x,
		Y:     y,
		color: p.color,
	}
}

func (c *Camera) PutPixel(p ScreenPoint) {
	x, y := uint(p.X), uint(p.Y)
	if x >= c.width || y >= c.height {
		return
	}
	c.canvas[y*c.width+x] = p.color
}

func (c *Camera) MoveBackForwad(unit float32) {
	if unit == 0 {
		return
	}

	rotMat := NewRotationMatrix(c.transforms.rotation)

	direction := rotMat.MultiplyByVec3(c.transforms.forwardDirection)
	normalDirection := direction.Normalized()

	c.transforms.position = c.transforms.position.Add(normalDirection.Scale(unit))

	c.transforms.UpdateTransforms(true, true)
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

	c.transforms.UpdateTransforms(true, true)
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

	c.transforms.UpdateTransforms(true, true)
}

func (c *Camera) ToggleViewLock() {
	c.updateView = !c.updateView
}

func (c *Camera) MoveVetically(unit float32) {
	c.transforms.position.Y += unit
	c.transforms.UpdateTransforms(true, true)
}

func (c *Camera) CalculateFrustum() {
	zFar := float32(100)
	camFront := c.transforms.forwardDirection
	halfVSide := zFar * math.Tan(c.fovAngle * .5)
	halfHSide := halfVSide * aspect;
	frontMultFar := camFront.Scale(zFar);
	//TODO: translate cpp to go
	c.frustum.nearPlane = NewPlane(c.transforms.position.Add(camFront.Scale(zNear)), camFront );
	c.frustum.farFace = NewPlane( c.transforms.position.Add(frontMultFar) , camFront.Scale(-1) )
	c.frustum.rightFace = NewPlane( c.transforms.position, glm::cross(frontMultFar - cam.Right * halfHSide, cam.Up) )
	c.frustum.leftFace = NewPlane( c.transforms.position, glm::cross(cam.Up, frontMultFar + cam.Right * halfHSide) )
	c.frustum.topFace = NewPlane( c.transforms.position, glm::cross(cam.Right, frontMultFar - cam.Up * halfVSide) )
	c.frustum.bottomFace = NewPlane( c.transforms.position, glm::cross(frontMultFar + cam.Up * halfVSide, cam.Right) )
}
