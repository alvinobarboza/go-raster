package main

import (
	"math/rand/v2"
	"testing"
)

const (
	width  = 1920
	height = 1080
)

var keep bool

func checkIf(i int) bool {
	if i > 100 {
		return true
	}

	if i < 10 {
		return true
	}

	if i == 30 {
		return true
	}

	return false
}

func checkSwitch(i int) bool {
	switch {
	case i > 100:
		return true
	case i < 10:
		return true
	case i == 30:
		return true
	default:
		return false
	}
}

func checkUnsigned(i int) bool {
	if uint(i) >= uint(100) || uint(i) >= uint(50) {
		return true
	}
	return false
}

func BenchmarkIfSwitch(b *testing.B) {
	b.Run("if", func(b *testing.B) {
		for b.Loop() {
			checkIf(rand.IntN(200))
		}
	})

	b.Run("switch", func(b *testing.B) {
		for b.Loop() {
			checkSwitch(rand.IntN(200))
		}
	})

	b.Run("unsigned IF", func(b *testing.B) {
		for b.Loop() {
			checkUnsigned(rand.IntN(200))
		}
	})

	b.Run("BenchmarkBounds_Standard", func(b *testing.B) {
		x, y := 960, 540
		w, h := width, height

		for b.Loop() {
			if x < 0 || x >= w || y < 0 || y >= h {
				keep = true
			} else {
				keep = false
			}
		}
	})

	b.Run("BenchmarkBounds_UintOptimization", func(b *testing.B) {
		x, y := 960, 540
		w, h := width, height

		for b.Loop() {
			// Cast to uint to handle < 0 and >= max in one go
			if uint(x) >= uint(w) || uint(y) >= uint(h) {
				keep = true
			} else {
				keep = false
			}
		}
	})
}
