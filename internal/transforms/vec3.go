package transforms

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float32
}

func NewVec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (v Vec3) DotByVec3(v1 Vec3) float32 {
	return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
}

func (v Vec3) Dot() float32 {
	return v.DotByVec3(v)
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.Dot())))
}

// vector * n
func (v Vec3) Scale(n float32) Vec3 {
	return Vec3{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func (v Vec3) Divide(n float32) Vec3 {
	if n == 0 {
		return Vec3{}
	}

	return Vec3{
		X: v.X / n,
		Y: v.Y / n,
		Z: v.Z / n,
	}
}

func (v Vec3) Normalized() Vec3 {
	return v.Divide(v.Length())
}

func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		X: v.Y*v2.Z - v.Z*v2.Y,
		Y: v.Z*v2.X - v.X*v2.Z,
		Z: v.X*v2.Y - v.Y*v2.X,
	}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vec3) Subtract(v2 Vec3) Vec3 {
	return Vec3{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vec3) LerpTo(b Vec3, ratio float32) Vec3 {
	if ratio > 1 {
		return b
	}
	if ratio < 0 {
		return v
	}

	return b.Subtract(v).Scale(ratio).Add(v)
}

func (v Vec3) Print(name string) {
	fmt.Printf("%s = %+v\n", name, v)
}

func ReflectRay(ray, normal Vec3) Vec3 {
	r_dot_n := ray.DotByVec3(normal)
	return Vec3{
		X: 2*normal.X*r_dot_n - ray.X,
		Y: 2*normal.Y*r_dot_n - ray.Y,
		Z: 2*normal.Z*r_dot_n - ray.Z,
	}
}
