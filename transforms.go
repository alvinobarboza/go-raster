package main

import "math"

const Pi = 3.14159265358979323846

type Vec3 struct {
	X, Y, Z float32
}

func NewVec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func DegToRad(angle float32) float32 {
	return (angle * Pi) / 180
}

func FovScaling(angle float32) float32 {
	return 1 / float32(math.Tan(float64(DegToRad(angle)/2)))
}
