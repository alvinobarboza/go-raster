package main

import (
	"image/color"
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

func TestTexture(t *testing.T) {
	t.Run("(1.5,0)", func(t *testing.T) {
		tex := Texture{width: 100, height: 100}
		uv := Vec3{X: 1.5, Y: 0}

		w, h := tex.UVToWH(uv)

		wantW, wantH := 50, 0

		if w != wantW || h != wantH {
			t.Errorf("want w:%d h:%d, got w:%d h:%d", wantW, wantH, w, h)
		}
	})

	t.Run("(0.5,0.8)", func(t *testing.T) {
		tex := Texture{width: 100, height: 100}
		uv := Vec3{X: 0.5, Y: 0.8}

		w, h := tex.UVToWH(uv)

		wantW, wantH := 50, 80

		if w != wantW || h != wantH {
			t.Errorf("want w:%d h:%d, got w:%d h:%d", wantW, wantH, w, h)
		}
	})

	t.Run("(0.5,0.8) and texel color", func(t *testing.T) {
		tex := Texture{width: 10, height: 10, pixels: make([]color.RGBA, 0)}

		for range tex.width * tex.height {
			tex.pixels = append(tex.pixels, Black)
		}

		wantW, wantH := 5, 8
		wantColor := Red

		tex.pixels[wantH*tex.width+wantW] = wantColor

		uv := Vec3{X: 0.5, Y: 0.8}

		w, h := tex.UVToWH(uv)

		if w != wantW || h != wantH {
			t.Errorf("want w:%d h:%d, got w:%d h:%d", wantW, wantH, w, h)
		}

		tColor := tex.TexelColor(uv)

		if wantColor != tColor {
			t.Errorf("Want: %v, got %v", wantColor, tColor)
		}
	})
}
