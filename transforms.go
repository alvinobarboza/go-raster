package main

import (
	"fmt"
	"math"
)

const (
	Pi       = 3.14159265358979323846
	DegToRad = Pi / 180

	MatLength = 4
	M4x4      = MatLength * MatLength
)

type Vec3 struct {
	X, Y, Z float32
}

func NewVec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (v Vec3) DotByVec3(v1 Vec3) float32 {
	return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
}

func (v Vec3) Dot() float32 {
	return v.DotByVec3(v)
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.Dot())))
}

// vector * n
func (v Vec3) Scale(n float32) Vec3 {
	return Vec3{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func (v Vec3) Divide(n float32) Vec3 {
	if n == 0 {
		return Vec3{}
	}
	return v.Scale(1 / n)
}

func (v Vec3) Normalized() Vec3 {
	return v.Divide(v.Length())
}

func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		X: v.Y*v2.Z - v.Z*v2.Y,
		Y: v.Z*v2.X - v.X*v2.Z,
		Z: v.X*v2.Y - v.Y*v2.X,
	}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vec3) Subtract(v2 Vec3) Vec3 {
	return Vec3{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vec3) LerpTo(b Vec3, ratio float32) Vec3 {
	if ratio > 1 {
		return b
	}
	if ratio < 0 {
		return v
	}

	return b.Subtract(v).Scale(ratio).Add(v)
}

func (v Vec3) Print(name string) {
	fmt.Printf("%s = %+v\n", name, v)
}

type Matrix [M4x4]float32

func NewZeroMatrix() Matrix {
	return Matrix{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
}

func NewIdentityMatrix() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func NewScaleMatrix(scale Vec3) Matrix {
	return Matrix{
		scale.X, 0.0, 0.0, 0.0,
		0.0, scale.Y, 0.0, 0.0,
		0.0, 0.0, scale.Z, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func NewRotationMatrix(angle Vec3) Matrix {
	cosa := float32(math.Cos(float64(angle.X * -DegToRad)))
	sina := float32(math.Sin(float64(angle.X * -DegToRad)))

	cosb := float32(math.Cos(float64(angle.Y * -DegToRad)))
	sinb := float32(math.Sin(float64(angle.Y * -DegToRad)))

	cosga := float32(math.Cos(float64(angle.Z * -DegToRad)))
	singa := float32(math.Sin(float64(angle.Z * -DegToRad)))

	// Formula for general 3D roation using matrix
	return Matrix{
		cosb * cosga, sina*sinb*cosga - cosa*singa, cosa*sinb*cosga + sina*singa, 0.0,
		cosb * singa, sina*sinb*singa + cosa*cosga, cosa*sinb*singa - sina*cosga, 0.0,
		-sinb, sina * cosb, cosa * cosb, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

func NewTranslationMatrix(position Vec3) Matrix {
	return Matrix{
		1.0, 0.0, 0.0, position.X,
		0.0, 1.0, 0.0, position.Y,
		0.0, 0.0, 1.0, position.Z,
		0.0, 0.0, 0.0, 1.0,
	}
}

func (m Matrix) Transposed() Matrix {
	transposed := Matrix{}
	for row := range MatLength {
		for col := range MatLength {
			transposed[col*MatLength+row] = m[row*MatLength+col]
		}
	}

	return transposed
}

func (m Matrix) MultiplyByVec3(v Vec3) Vec3 {
	v4 := [MatLength]float32{v.X, v.Y, v.Z, 1.0}
	result := [MatLength]float32{0.0, 0.0, 0.0, 0.0}

	for row := range MatLength {
		for col := range MatLength {
			result[row] += v4[col] * m[row*MatLength+col]
		}
	}

	return Vec3{X: result[0], Y: result[1], Z: result[2]}
}

func (m Matrix) MultiplyByMatrix(m2 Matrix) Matrix {
	result := NewZeroMatrix()
	for row := range MatLength {
		for col := range MatLength {
			for k := range MatLength {
				result[row*MatLength+col] += m[row*MatLength+k] * m2[k*MatLength+col]
			}
		}
	}

	return result
}

type Transforms struct {
	scale            Vec3
	rotation         Vec3
	position         Vec3
	forwardDirection Vec3

	scaleMat         Matrix
	rotationMat      Matrix
	translationMat   Matrix
	matrixTransforms Matrix
}

func NewTransforms(pos, scale, rot Vec3) Transforms {
	return Transforms{
		position: pos,
		scale:    scale,
		rotation: rot,
	}
}

func (t *Transforms) UpdateTransforms(transposeRot bool, invertPosition bool) {
	t.rotationMat = NewRotationMatrix(t.rotation)
	t.scaleMat = NewScaleMatrix(t.scale)
	pos := t.position

	if invertPosition {
		pos = t.position.Scale(-1)
	}

	t.translationMat = NewTranslationMatrix(pos)

	if transposeRot {
		t.rotationMat = t.rotationMat.Transposed()
	}

	// TODO: Fix multiply order
	t.matrixTransforms = t.rotationMat.MultiplyByMatrix(t.scaleMat)
	t.matrixTransforms = t.matrixTransforms.MultiplyByMatrix(t.translationMat)
}

func FovScaling(angle float32) float32 {
	return float32(1 / math.Tan(float64(angle*DegToRad/2)))
}

func (m Matrix) Print(name string) {
	fmt.Println(name, ":")
	for row := range MatLength {
		for col := range MatLength {
			fmt.Printf(" %6.2f", m[col*MatLength+row])
		}
		fmt.Print("\n")
	}
}

func Minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinIn(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxIn(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func EdgeCross(a, b, p ScreenPoint) float32 {
	abX := b.X - a.X
	abY := b.Y - a.Y

	apX := p.X - a.X
	apY := p.Y - a.Y

	return (abX * apY) - (abY * apX)
}
