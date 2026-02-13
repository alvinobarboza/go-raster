package main

import (
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
