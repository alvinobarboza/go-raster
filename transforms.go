package main

import "math"

const Pi = 3.14159265358979323846

func DegToRad(angle float32) float32 {
	return (angle * Pi) / 180
}

func FovScaling(angle float32) float32 {
	return 1 / float32(math.Tan(float64(DegToRad(angle)/2)))
}
