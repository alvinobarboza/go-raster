package transforms

import "math"

const (
	X0Y0 int = iota
	X1Y0
	X2Y0
	X3Y0
	X0Y1
	X1Y1
	X2Y1
	X3Y1
	X0Y2
	X1Y2
	X2Y2
	X3Y2
	X0Y3
	X1Y3
	X2Y3
	X3Y3

	MatLength = 4
	M4x4      = MatLength * MatLength
)

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
	return Vec3{
		X: (v.X * m[X0Y0]) + (v.Y * m[X1Y0]) + (v.Z * m[X2Y0]) + (1 * m[X3Y0]),
		Y: (v.X * m[X0Y1]) + (v.Y * m[X1Y1]) + (v.Z * m[X2Y1]) + (1 * m[X3Y1]),
		Z: (v.X * m[X0Y2]) + (v.Y * m[X1Y2]) + (v.Z * m[X2Y2]) + (1 * m[X3Y2]),
	}
}

func (m Matrix) MultiplyByMatrix(m2 Matrix) Matrix {
	result := Matrix{}
	result[X0Y0] = (m[X0Y0] * m2[X0Y0]) + (m[X1Y0] * m2[X0Y1]) + (m[X2Y0] * m2[X0Y2]) + (m[X3Y0] * m2[X0Y3])
	result[X1Y0] = (m[X0Y0] * m2[X1Y0]) + (m[X1Y0] * m2[X1Y1]) + (m[X2Y0] * m2[X1Y2]) + (m[X3Y0] * m2[X1Y3])
	result[X2Y0] = (m[X0Y0] * m2[X2Y0]) + (m[X1Y0] * m2[X2Y1]) + (m[X2Y0] * m2[X2Y2]) + (m[X3Y0] * m2[X2Y3])
	result[X3Y0] = (m[X0Y0] * m2[X3Y0]) + (m[X1Y0] * m2[X3Y1]) + (m[X2Y0] * m2[X3Y2]) + (m[X3Y0] * m2[X3Y3])

	result[X0Y1] = (m[X0Y1] * m2[X0Y0]) + (m[X1Y1] * m2[X0Y1]) + (m[X2Y1] * m2[X0Y2]) + (m[X3Y1] * m2[X0Y3])
	result[X1Y1] = (m[X0Y1] * m2[X1Y0]) + (m[X1Y1] * m2[X1Y1]) + (m[X2Y1] * m2[X1Y2]) + (m[X3Y1] * m2[X1Y3])
	result[X2Y1] = (m[X0Y1] * m2[X2Y0]) + (m[X1Y1] * m2[X2Y1]) + (m[X2Y1] * m2[X2Y2]) + (m[X3Y1] * m2[X2Y3])
	result[X3Y1] = (m[X0Y1] * m2[X3Y0]) + (m[X1Y1] * m2[X3Y1]) + (m[X2Y1] * m2[X3Y2]) + (m[X3Y1] * m2[X3Y3])

	result[X0Y2] = (m[X0Y2] * m2[X0Y0]) + (m[X1Y2] * m2[X0Y1]) + (m[X2Y2] * m2[X0Y2]) + (m[X3Y2] * m2[X0Y3])
	result[X1Y2] = (m[X0Y2] * m2[X1Y0]) + (m[X1Y2] * m2[X1Y1]) + (m[X2Y2] * m2[X1Y2]) + (m[X3Y2] * m2[X1Y3])
	result[X2Y2] = (m[X0Y2] * m2[X2Y0]) + (m[X1Y2] * m2[X2Y1]) + (m[X2Y2] * m2[X2Y2]) + (m[X3Y2] * m2[X2Y3])
	result[X3Y2] = (m[X0Y2] * m2[X3Y0]) + (m[X1Y2] * m2[X3Y1]) + (m[X2Y2] * m2[X3Y2]) + (m[X3Y2] * m2[X3Y3])

	result[X0Y3] = (m[X0Y3] * m2[X0Y0]) + (m[X1Y3] * m2[X0Y1]) + (m[X2Y3] * m2[X0Y2]) + (m[X3Y3] * m2[X0Y3])
	result[X1Y3] = (m[X0Y3] * m2[X1Y0]) + (m[X1Y3] * m2[X1Y1]) + (m[X2Y3] * m2[X1Y2]) + (m[X3Y3] * m2[X1Y3])
	result[X2Y3] = (m[X0Y3] * m2[X2Y0]) + (m[X1Y3] * m2[X2Y1]) + (m[X2Y3] * m2[X2Y2]) + (m[X3Y3] * m2[X2Y3])
	result[X3Y3] = (m[X0Y3] * m2[X3Y0]) + (m[X1Y3] * m2[X3Y1]) + (m[X2Y3] * m2[X3Y2]) + (m[X3Y3] * m2[X3Y3])

	return result
}
