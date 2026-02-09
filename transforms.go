package main

import "math"

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

func (v Vec3) MultiplyByVec3(v1 Vec3) float32 {
	return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
}

func (v Vec3) Dot() float32 {
	return v.MultiplyByVec3(v)
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.Dot())))
}

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

type Matrix [M4x4]float32

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
			transposed[col*MatLength+row] = transposed[row*MatLength+col]
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
	result := Matrix{}
	for row := range MatLength {
		for col := range MatLength {
			result[row*MatLength+col] = 0.0
			for k := range MatLength {
				result[row*MatLength+col] += m[row*MatLength+k] * m2[k*MatLength+col]
			}
		}
	}

	return result
}

type Tranforms struct {
	scale, rotation, position Vec3

	matrixTransforms Matrix
}

func FovScaling(angle float32) float32 {
	return float32(1 / math.Tan(float64(angle*DegToRad/2)))
}
