package transforms

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X, Y float32
}

func NewVec2(x, y float32) Vec2 {
	return Vec2{x, y}
}

func (v Vec2) DotByVec2(v1 Vec2) float32 {
	return v.X*v1.X + v.Y*v1.Y
}

func (v Vec2) Dot() float32 {
	return v.DotByVec2(v)
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.Dot())))
}

// vector * n
func (v Vec2) Scale(n float32) Vec2 {
	return Vec2{
		X: v.X * n,
		Y: v.Y * n,
	}
}

func (v Vec2) Divide(n float32) Vec2 {
	if n == 0 {
		return Vec2{}
	}
	return Vec2{
		X: v.X / n,
		Y: v.Y / n,
	}
}

func (v Vec2) Normalized() Vec2 {
	return v.Divide(v.Length())
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vec2) Subtract(v2 Vec2) Vec2 {
	return Vec2{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v Vec2) LerpTo(b Vec2, ratio float32) Vec2 {
	if ratio > 1 {
		return b
	}
	if ratio < 0 {
		return v
	}

	return b.Subtract(v).Scale(ratio).Add(v)
}

func (v Vec2) Print(name string) {
	fmt.Printf("%s = %+v\n", name, v)
}
