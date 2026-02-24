package main

import (
	"log"
	"math"
	"testing"
)

type Vec3_64 struct {
	X, Y, Z float64
}

func (v Vec3_64) DotByVec3(v1 Vec3_64) float64 {
	return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
}

func (v Vec3_64) Dot() float64 {
	return v.DotByVec3(v)
}

func (v Vec3_64) Length() float64 {
	return math.Sqrt(v.Dot())
}

func (v Vec3_64) Scale(n float64) Vec3_64 {
	return Vec3_64{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func (v Vec3_64) Divide(n float64) Vec3_64 {
	if n == 0 {
		return Vec3_64{}
	}
	return v.Scale(1 / n)
}

func (v Vec3_64) Normalized() Vec3_64 {
	return v.Divide(v.Length())
}

func BenchmarkVecFloat32(t *testing.B) {
	v := Vec3{
		X: 2,
		Y: 3,
		Z: 5,
	}

	for t.Loop() {
		v.Normalized()
	}
}

func BenchmarkVecFloat64(t *testing.B) {
	v := Vec3_64{
		X: 2,
		Y: 3,
		Z: 5,
	}

	for t.Loop() {
		v.Normalized()
	}
}

// Maintain brench
func DivideVec3(v Vec3, n float32) Vec3 {
	if n == 0 {
		return NewVec3(0, 0, 0)
	}

	return Vec3{
		X: v.X / n,
		Y: v.Y / n,
		Z: v.Z / n,
	}
}

func BenchmarkDivision(b *testing.B) {
	v := NewVec3(2, 4, 5)
	for b.Loop() {
		v = DivideVec3(v, 5)
	}
}

func BenchmarkReciprocal(b *testing.B) {
	v := NewVec3(2, 4, 5)
	for b.Loop() {
		v = v.Divide(5)
	}
}

type Matrix64 [M4x4]float64

func (m Matrix64) MultiplyByMatrix64(m2 Matrix64) Matrix64 {
	result := Matrix64{}
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

func BenchmarkMat32(b *testing.B) {
	m := Matrix{
		3, 4, 5, 6,
		4, 5, 6, 7,
		1, 3, 4, 5,
		4, 5, 6, 7,
	}

	for b.Loop() {
		m.MultiplyByMatrix(m)
	}
}

func BenchmarkMat64(b *testing.B) {
	m := Matrix64{
		3, 4, 5, 6,
		4, 5, 6, 7,
		1, 3, 4, 5,
		4, 5, 6, 7,
	}

	for b.Loop() {
		m.MultiplyByMatrix64(m)
	}
}

func TestVec3(t *testing.T) {
	v := Vec3{
		X: 20,
	}

	v = v.Normalized()

	if v.X != 1 {
		t.Errorf("Expected {1 0 0} got %v", v)
	}
}

func TestVec3Lerp(t *testing.T) {
	a := NewVec3(0, 0, 0)
	b := NewVec3(10, 0, 0)

	t.Run("Half", func(t *testing.T) {
		v := a.LerpTo(b, 0.5)
		if v.X != 5 {
			t.Errorf("Wanted {5 0 0} got %v\n", v)
		}
	})

	t.Run("Above", func(t *testing.T) {
		v := a.LerpTo(b, 2)
		if v.X != 10 {
			t.Errorf("Wanted {10 0 0} got %v\n", v)
		}
	})

	t.Run("Below", func(t *testing.T) {
		v := a.LerpTo(b, -1)
		if v.X != 0 {
			t.Errorf("Wanted {0 0 0} got %v\n", v)
		}
	})
}

func TestTransposeMat(t *testing.T) {
	want := Matrix{
		1, 0, 0, 1,
		0, 1, 0, 0,
		1, 0, 1, 0,
		0, 0, 0, 1,
	}

	matTest := Matrix{
		1, 0, 1, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		1, 0, 0, 1,
	}

	got := matTest.Transposed()

	incorrect := true
	for i := range M4x4 {
		if want[i] != got[i] {
			incorrect = false
		}
	}

	if !incorrect {
		t.Errorf("Want %v, got %v", want, got)
	}
}

func IsEqualMatrix(mat1, mat2 Matrix) bool {

	for i := range M4x4 {
		if mat1[i] != mat2[i] {
			log.Println(mat1[i], mat2[i])
			return false
		}
	}

	return true
}

func TestMatrices(t *testing.T) {

	t.Run("scale matrix", func(t *testing.T) {
		want := Matrix{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}

		got := NewScaleMatrix(NewVec3(1, 1, 1))

		if !IsEqualMatrix(want, got) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			got.Print("Got")
		}
	})

	t.Run("rotation matrix", func(t *testing.T) {
		want := Matrix{
			-0.00000004371138828674, +0.00, -1.00, +0.00,
			+0.00, +1.00, +0.00, +0.00,
			+1.00, +0.00, -0.00000004371138828674, +0.00,
			+0.00, +0.00, +0.00, +1.00,
		}

		got := NewRotationMatrix(NewVec3(0, 90, 0))

		if !IsEqualMatrix(want, got) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			got.Print("Got")
		}
	})

	t.Run("translation matrix", func(t *testing.T) {
		want := Matrix{
			+1.00, +0.00, +0.00, +1.00,
			+0.00, +1.00, +0.00, +1.00,
			+0.00, +0.00, +1.00, +0.00,
			+0.00, +0.00, +0.00, +1.00,
		}

		got := NewTranslationMatrix(NewVec3(1, 1, 0))

		if !IsEqualMatrix(want, got) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			got.Print("Got")
		}
	})

	t.Run("mat x mat", func(t *testing.T) {
		mat := Matrix{
			2, 2, 2, 2,
			2, 2, 2, 2,
			2, 2, 2, 2,
			2, 2, 2, 2,
		}

		want := Matrix{
			16, 16, 16, 16,
			16, 16, 16, 16,
			16, 16, 16, 16,
			16, 16, 16, 16,
		}

		got := mat.MultiplyByMatrix(mat)

		if !IsEqualMatrix(got, want) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			got.Print("Got")
		}
	})

	t.Run("matrices multiplication order", func(t *testing.T) {
		scale := Matrix{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}

		rotation := Matrix{
			2, 3, 4, 1,
			3, 2, 3, 2,
			2, 3, 4, 1,
			3, 2, 3, 2,
		}

		translation := Matrix{
			1, 0, 0, 2,
			0, 1, 0, 4,
			0, 0, 1, 5,
			0, 0, 0, 1,
		}

		result := rotation.MultiplyByMatrix(scale)
		got := translation.MultiplyByMatrix(result)

		want := Matrix{
			+8.00, +7.00, +10.00, +5.00,
			+15.00, +10.00, +15.00, +10.00,
			+17.00, +13.00, +19.00, +11.00,
			+3.00, +2.00, +3.00, +2.00,
		}

		if !IsEqualMatrix(want, got) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			got.Print("Got")
		}
	})

	t.Run("calculate corret matrices", func(t *testing.T) {
		want := Matrix{
			-0.00000004371138828674, +0.00, -1.00, +1.00,
			+0.00, +1.00, +0.00, +1.00,
			+1.00, +0.00, -0.00000004371138828674, +0.00,
			+0.00, +0.00, +0.00, +1.00,
		}

		transform := Transforms{
			scale:    NewVec3(1, 1, 1),
			rotation: NewVec3(0, 90, 0),
			position: NewVec3(1, 1, 0),
		}

		transform.UpdateModelTransforms()

		if !IsEqualMatrix(want, transform.matrixTransforms) {
			t.Errorf("want")
			want.Print("Want")
			t.Error("got")
			transform.matrixTransforms.Print("Got")
		}
	})
}

func BenchmarkUnrolled(b *testing.B) {
	mat1 := NewRotationMatrix(NewVec3(30, 20, 20))
	mat2 := NewRotationMatrix(NewVec3(30, 20, 20))

	b.Run("for", func(b *testing.B) {
		for b.Loop() {
			m := Matrix{}
			for row := range MatLength {
				for col := range MatLength {
					for k := range MatLength {
						m[row*MatLength+col] += mat1[row*MatLength+k] * mat2[k*MatLength+col]
					}
				}
			}
		}
	})

	b.Run("unr", func(b *testing.B) {
		for b.Loop() {
			m := Matrix{}
			m[X0Y0] = (mat1[X0Y0] * mat2[X0Y0]) + (mat1[X1Y0] * mat2[X0Y1]) + (mat1[X2Y0] * mat2[X0Y2]) + (mat1[X3Y0] * mat2[X0Y3])
			m[X1Y0] = (mat1[X0Y0] * mat2[X1Y0]) + (mat1[X1Y0] * mat2[X1Y1]) + (mat1[X2Y0] * mat2[X1Y2]) + (mat1[X3Y0] * mat2[X1Y3])
			m[X2Y0] = (mat1[X0Y0] * mat2[X2Y0]) + (mat1[X1Y0] * mat2[X2Y1]) + (mat1[X2Y0] * mat2[X2Y2]) + (mat1[X3Y0] * mat2[X2Y3])
			m[X3Y0] = (mat1[X0Y0] * mat2[X3Y0]) + (mat1[X1Y0] * mat2[X3Y1]) + (mat1[X2Y0] * mat2[X3Y2]) + (mat1[X3Y0] * mat2[X3Y3])

			m[X0Y1] = (mat1[X0Y1] * mat2[X0Y0]) + (mat1[X1Y1] * mat2[X0Y1]) + (mat1[X2Y1] * mat2[X0Y2]) + (mat1[X3Y1] * mat2[X0Y3])
			m[X1Y1] = (mat1[X0Y1] * mat2[X1Y0]) + (mat1[X1Y1] * mat2[X1Y1]) + (mat1[X2Y1] * mat2[X1Y2]) + (mat1[X3Y1] * mat2[X1Y3])
			m[X2Y1] = (mat1[X0Y1] * mat2[X2Y0]) + (mat1[X1Y1] * mat2[X2Y1]) + (mat1[X2Y1] * mat2[X2Y2]) + (mat1[X3Y1] * mat2[X2Y3])
			m[X3Y1] = (mat1[X0Y1] * mat2[X3Y0]) + (mat1[X1Y1] * mat2[X3Y1]) + (mat1[X2Y1] * mat2[X3Y2]) + (mat1[X3Y1] * mat2[X3Y3])

			m[X0Y2] = (mat1[X0Y2] * mat2[X0Y0]) + (mat1[X1Y2] * mat2[X0Y1]) + (mat1[X2Y2] * mat2[X0Y2]) + (mat1[X3Y2] * mat2[X0Y3])
			m[X1Y2] = (mat1[X0Y2] * mat2[X1Y0]) + (mat1[X1Y2] * mat2[X1Y1]) + (mat1[X2Y2] * mat2[X1Y2]) + (mat1[X3Y2] * mat2[X1Y3])
			m[X2Y2] = (mat1[X0Y2] * mat2[X2Y0]) + (mat1[X1Y2] * mat2[X2Y1]) + (mat1[X2Y2] * mat2[X2Y2]) + (mat1[X3Y2] * mat2[X2Y3])
			m[X3Y2] = (mat1[X0Y2] * mat2[X3Y0]) + (mat1[X1Y2] * mat2[X3Y1]) + (mat1[X2Y2] * mat2[X3Y2]) + (mat1[X3Y2] * mat2[X3Y3])

			m[X0Y3] = (mat1[X0Y3] * mat2[X0Y0]) + (mat1[X1Y3] * mat2[X0Y1]) + (mat1[X2Y3] * mat2[X0Y2]) + (mat1[X3Y3] * mat2[X0Y3])
			m[X1Y3] = (mat1[X0Y3] * mat2[X1Y0]) + (mat1[X1Y3] * mat2[X1Y1]) + (mat1[X2Y3] * mat2[X1Y2]) + (mat1[X3Y3] * mat2[X1Y3])
			m[X2Y3] = (mat1[X0Y3] * mat2[X2Y0]) + (mat1[X1Y3] * mat2[X2Y1]) + (mat1[X2Y3] * mat2[X2Y2]) + (mat1[X3Y3] * mat2[X2Y3])
			m[X3Y3] = (mat1[X0Y3] * mat2[X3Y0]) + (mat1[X1Y3] * mat2[X3Y1]) + (mat1[X2Y3] * mat2[X3Y2]) + (mat1[X3Y3] * mat2[X3Y3])
		}
	})

	b.Run("unrOp", func(b *testing.B) {

		// AI suggestion - trying to avoid cache miss

		for b.Loop() {
			// Linearly load mat2 rows into registers to prevent column-jumping
			r0x0, r0x1, r0x2, r0x3 := mat2[X0Y0], mat2[X1Y0], mat2[X2Y0], mat2[X3Y0]
			r1x0, r1x1, r1x2, r1x3 := mat2[X0Y1], mat2[X1Y1], mat2[X2Y1], mat2[X3Y1]
			r2x0, r2x1, r2x2, r2x3 := mat2[X0Y2], mat2[X1Y2], mat2[X2Y2], mat2[X3Y2]
			r3x0, r3x1, r3x2, r3x3 := mat2[X0Y3], mat2[X1Y3], mat2[X2Y3], mat2[X3Y3]

			m := Matrix{}

			// Row Y0
			m[X0Y0] = (mat1[X0Y0] * r0x0) + (mat1[X1Y0] * r1x0) + (mat1[X2Y0] * r2x0) + (mat1[X3Y0] * r3x0)
			m[X1Y0] = (mat1[X0Y0] * r0x1) + (mat1[X1Y0] * r1x1) + (mat1[X2Y0] * r2x1) + (mat1[X3Y0] * r3x1)
			m[X2Y0] = (mat1[X0Y0] * r0x2) + (mat1[X1Y0] * r1x2) + (mat1[X2Y0] * r2x2) + (mat1[X3Y0] * r3x2)
			m[X3Y0] = (mat1[X0Y0] * r0x3) + (mat1[X1Y0] * r1x3) + (mat1[X2Y0] * r2x3) + (mat1[X3Y0] * r3x3)

			// Row Y1
			m[X0Y1] = (mat1[X0Y1] * r0x0) + (mat1[X1Y1] * r1x0) + (mat1[X2Y1] * r2x0) + (mat1[X3Y1] * r3x0)
			m[X1Y1] = (mat1[X0Y1] * r0x1) + (mat1[X1Y1] * r1x1) + (mat1[X2Y1] * r2x1) + (mat1[X3Y1] * r3x1)
			m[X2Y1] = (mat1[X0Y1] * r0x2) + (mat1[X1Y1] * r1x2) + (mat1[X2Y1] * r2x2) + (mat1[X3Y1] * r3x2)
			m[X3Y1] = (mat1[X0Y1] * r0x3) + (mat1[X1Y1] * r1x3) + (mat1[X2Y1] * r2x3) + (mat1[X3Y1] * r3x3)

			// Row Y2
			m[X0Y2] = (mat1[X0Y2] * r0x0) + (mat1[X1Y2] * r1x0) + (mat1[X2Y2] * r2x0) + (mat1[X3Y2] * r3x0)
			m[X1Y2] = (mat1[X0Y2] * r0x1) + (mat1[X1Y2] * r1x1) + (mat1[X2Y2] * r2x1) + (mat1[X3Y2] * r3x1)
			m[X2Y2] = (mat1[X0Y2] * r0x2) + (mat1[X1Y2] * r1x2) + (mat1[X2Y2] * r2x2) + (mat1[X3Y2] * r3x2)
			m[X3Y2] = (mat1[X0Y2] * r0x3) + (mat1[X1Y2] * r1x3) + (mat1[X2Y2] * r2x3) + (mat1[X3Y2] * r3x3)

			// Row Y3
			m[X0Y3] = (mat1[X0Y3] * r0x0) + (mat1[X1Y3] * r1x0) + (mat1[X2Y3] * r2x0) + (mat1[X3Y3] * r3x0)
			m[X1Y3] = (mat1[X0Y3] * r0x1) + (mat1[X1Y3] * r1x1) + (mat1[X2Y3] * r2x1) + (mat1[X3Y3] * r3x1)
			m[X2Y3] = (mat1[X0Y3] * r0x2) + (mat1[X1Y3] * r1x2) + (mat1[X2Y3] * r2x2) + (mat1[X3Y3] * r3x2)
			m[X3Y3] = (mat1[X0Y3] * r0x3) + (mat1[X1Y3] * r1x3) + (mat1[X2Y3] * r2x3) + (mat1[X3Y3] * r3x3)
		}
	})

	b.Run("imp", func(b *testing.B) {
		for b.Loop() {
			mat1.MultiplyByMatrix(mat2)
		}
	})
}

func BenchmarkUnrolledVec(b *testing.B) {
	mat1 := Matrix{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	vec := NewVec3(11, 15, 19)

	b.Run("vFor", func(b *testing.B) {
		for b.Loop() {
			v4 := [MatLength]float32{vec.X, vec.Y, vec.Z, 1.0}
			result := [MatLength]float32{0.0, 0.0, 0.0, 0.0}

			for row := range MatLength {
				for col := range MatLength {
					result[row] += v4[col] * mat1[row*MatLength+col]
				}
			}
		}
	})

	b.Run("vUnr", func(b *testing.B) {
		for b.Loop() {
			NewVec3(
				(vec.X*mat1[X0Y0])+(vec.Y*mat1[X1Y0])+(vec.Z*mat1[X2Y0])+(1*mat1[X3Y0]),
				(vec.X*mat1[X0Y1])+(vec.Y*mat1[X1Y1])+(vec.Z*mat1[X2Y1])+(1*mat1[X3Y1]),
				(vec.X*mat1[X0Y2])+(vec.Y*mat1[X1Y2])+(vec.Z*mat1[X2Y2])+(1*mat1[X3Y2]),
			)
		}
	})
}

func TestMatrixMult(t *testing.T) {
	mat1 := Matrix{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	mat2 := Matrix{
		11, 12, 13, 14,
		15, 16, 17, 18,
		19, 110, 111, 112,
		113, 114, 115, 116,
	}

	vec := NewVec3(11, 15, 19)

	wantMat := Matrix{
		550, 830, 840, 850,
		1182, 1838, 1864, 1890,
		1814, 2846, 2888, 2930,
		2446, 3854, 3912, 3970,
	}

	wantVec := NewVec3(102, 286, 470)

	t.Run("unrolled", func(t *testing.T) {
		m := Matrix{}
		m[X0Y0] = (mat1[X0Y0] * mat2[X0Y0]) + (mat1[X1Y0] * mat2[X0Y1]) + (mat1[X2Y0] * mat2[X0Y2]) + (mat1[X3Y0] * mat2[X0Y3])
		m[X1Y0] = (mat1[X0Y0] * mat2[X1Y0]) + (mat1[X1Y0] * mat2[X1Y1]) + (mat1[X2Y0] * mat2[X1Y2]) + (mat1[X3Y0] * mat2[X1Y3])
		m[X2Y0] = (mat1[X0Y0] * mat2[X2Y0]) + (mat1[X1Y0] * mat2[X2Y1]) + (mat1[X2Y0] * mat2[X2Y2]) + (mat1[X3Y0] * mat2[X2Y3])
		m[X3Y0] = (mat1[X0Y0] * mat2[X3Y0]) + (mat1[X1Y0] * mat2[X3Y1]) + (mat1[X2Y0] * mat2[X3Y2]) + (mat1[X3Y0] * mat2[X3Y3])

		m[X0Y1] = (mat1[X0Y1] * mat2[X0Y0]) + (mat1[X1Y1] * mat2[X0Y1]) + (mat1[X2Y1] * mat2[X0Y2]) + (mat1[X3Y1] * mat2[X0Y3])
		m[X1Y1] = (mat1[X0Y1] * mat2[X1Y0]) + (mat1[X1Y1] * mat2[X1Y1]) + (mat1[X2Y1] * mat2[X1Y2]) + (mat1[X3Y1] * mat2[X1Y3])
		m[X2Y1] = (mat1[X0Y1] * mat2[X2Y0]) + (mat1[X1Y1] * mat2[X2Y1]) + (mat1[X2Y1] * mat2[X2Y2]) + (mat1[X3Y1] * mat2[X2Y3])
		m[X3Y1] = (mat1[X0Y1] * mat2[X3Y0]) + (mat1[X1Y1] * mat2[X3Y1]) + (mat1[X2Y1] * mat2[X3Y2]) + (mat1[X3Y1] * mat2[X3Y3])

		m[X0Y2] = (mat1[X0Y2] * mat2[X0Y0]) + (mat1[X1Y2] * mat2[X0Y1]) + (mat1[X2Y2] * mat2[X0Y2]) + (mat1[X3Y2] * mat2[X0Y3])
		m[X1Y2] = (mat1[X0Y2] * mat2[X1Y0]) + (mat1[X1Y2] * mat2[X1Y1]) + (mat1[X2Y2] * mat2[X1Y2]) + (mat1[X3Y2] * mat2[X1Y3])
		m[X2Y2] = (mat1[X0Y2] * mat2[X2Y0]) + (mat1[X1Y2] * mat2[X2Y1]) + (mat1[X2Y2] * mat2[X2Y2]) + (mat1[X3Y2] * mat2[X2Y3])
		m[X3Y2] = (mat1[X0Y2] * mat2[X3Y0]) + (mat1[X1Y2] * mat2[X3Y1]) + (mat1[X2Y2] * mat2[X3Y2]) + (mat1[X3Y2] * mat2[X3Y3])

		m[X0Y3] = (mat1[X0Y3] * mat2[X0Y0]) + (mat1[X1Y3] * mat2[X0Y1]) + (mat1[X2Y3] * mat2[X0Y2]) + (mat1[X3Y3] * mat2[X0Y3])
		m[X1Y3] = (mat1[X0Y3] * mat2[X1Y0]) + (mat1[X1Y3] * mat2[X1Y1]) + (mat1[X2Y3] * mat2[X1Y2]) + (mat1[X3Y3] * mat2[X1Y3])
		m[X2Y3] = (mat1[X0Y3] * mat2[X2Y0]) + (mat1[X1Y3] * mat2[X2Y1]) + (mat1[X2Y3] * mat2[X2Y2]) + (mat1[X3Y3] * mat2[X2Y3])
		m[X3Y3] = (mat1[X0Y3] * mat2[X3Y0]) + (mat1[X1Y3] * mat2[X3Y1]) + (mat1[X2Y3] * mat2[X3Y2]) + (mat1[X3Y3] * mat2[X3Y3])

		m.Print("unrolled")

		for i := range M4x4 {
			if wantMat[i] != m[i] {
				t.Errorf("Mismatch want %v got %v", wantMat, m)
			}
		}
	})

	t.Run("for loop", func(t *testing.T) {
		m := Matrix{}
		for row := range MatLength {
			for col := range MatLength {
				for k := range MatLength {
					m[row*MatLength+col] += mat1[row*MatLength+k] * mat2[k*MatLength+col]
				}
			}
		}
		m.Print("For")
		for i := range M4x4 {
			if wantMat[i] != m[i] {
				t.Errorf("Mismatch want %v got %v", wantMat, m)
			}
		}
	})

	t.Run("byVec3For", func(t *testing.T) {
		v4 := [MatLength]float32{vec.X, vec.Y, vec.Z, 1.0}
		result := [MatLength]float32{0.0, 0.0, 0.0, 0.0}

		for row := range MatLength {
			for col := range MatLength {
				result[row] += v4[col] * mat1[row*MatLength+col]
			}
		}

		if result[0] != wantVec.X || result[1] != wantVec.Y || result[2] != wantVec.Z {
			t.Errorf("want %v, got %v", wantVec, result)
		}
	})

	t.Run("byVec3Unr", func(t *testing.T) {
		got := NewVec3(
			(vec.X*mat1[X0Y0])+(vec.Y*mat1[X1Y0])+(vec.Z*mat1[X2Y0])+(1*mat1[X3Y0]),
			(vec.X*mat1[X0Y1])+(vec.Y*mat1[X1Y1])+(vec.Z*mat1[X2Y1])+(1*mat1[X3Y1]),
			(vec.X*mat1[X0Y2])+(vec.Y*mat1[X1Y2])+(vec.Z*mat1[X2Y2])+(1*mat1[X3Y2]),
		)
		if got.X != wantVec.X || got.Y != wantVec.Y || got.Z != wantVec.Z {
			t.Errorf("want %v, got %v", wantVec, got)
		}
	})
}
