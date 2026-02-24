package mesh

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/transforms"
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
	verts := []transforms.Vec3{
		transforms.NewVec3(1, 0, 0),
		transforms.NewVec3(-1, 0, 0),
	}

	bounds := NewBoundingSphere()

	bounds.CalculateBoundaries(verts, transforms.NewIdentityMatrix())

	if bounds.Radius != 1 {
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
			sum += maths.Floor32(x)
		}
	})
}
