package mesh

import (
	"image/color"
	"math/rand/v2"
	"testing"
)

func BenchmarkTexel(b *testing.B) {
	b.Run("array test", func(b *testing.B) {
		cl := [1024 * 1024]color.RGBA{}

		for b.Loop() {

			for range 13 {
				for range 1920 * 1080 {
					u := rand.Float32()
					v := rand.Float32()

					i := int(u*float32(1024) + v)
					c := cl[i]
					c.B = 200
				}
			}
		}
	})
}
