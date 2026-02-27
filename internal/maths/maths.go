package maths

import "math"

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

func Floor32(x float32) float32 {
	i := int(x)
	if x < float32(i) {
		return float32(i - 1)
	}
	return float32(i)
}

func Ceil32(x float32) float32 {
	return -Floor32(-x)
}

func Abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}
