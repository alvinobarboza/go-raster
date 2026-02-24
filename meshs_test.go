package main

import (
	"math"
	"math/rand/v2"
	"testing"
)

func BenchmarkWrap(b *testing.B) {
	b.Run("wrap directly", func(b *testing.B) {
		wrap := func(x float32) float32 {
			return x - float32(math.Floor(float64(x)))
		}

		for b.Loop() {
			wrap(rand.Float32() + float32(rand.IntN(10)))
		}
	})

	b.Run("wrap if", func(b *testing.B) {
		wrap := func(x float32) float32 {
			return x - float32(math.Floor(float64(x)))
		}

		for b.Loop() {
			i := float32(rand.IntN(10))
			if i != 1 {
				wrap(rand.Float32() + i)
			}
		}
	})
}

func TestBoundary(t *testing.T) {
	verts := []Vec3{
		NewVec3(1, 0, 0),
		NewVec3(-1, 0, 0),
	}

	bounds := NewBoundingSphere()

	bounds.CalculateBoundaries(verts, NewIdentityMatrix())

	if bounds.radius != 1 {
		t.Errorf("Expercted r=1, got %+v", bounds)
	}
}

func BenchmarkFloor(b *testing.B) {
	x := float32(2.76)

	b.Run("mathFloor", func(b *testing.B) {
		sum := float32(0)
		for b.Loop() {
			sum += float32(math.Floor(float64(x)))
		}
	})

	b.Run("custFloor", func(b *testing.B) {
		sum := float32(0)
		for b.Loop() {
			sum += Floor32(x)
		}
	})
}
