package main

import "testing"

func valueTester(v1, v2, v3 Vec3) Vec3 {
	return v1.Divide(2).Add(v2.Divide(2)).Add(v3.Divide(2))
}

func referenceTester(vs []Vec3, i1, i2, i3 int) Vec3 {
	return vs[i1].Divide(2).Add(vs[i2].Divide(2)).Add(vs[i3].Divide(2))
}

func BenchmarkPassByValue(b *testing.B) {
	b.Run("value", func(b *testing.B) {
		v := NewVec3(2, 3, 4)
		for b.Loop() {
			v = valueTester(v, v, v)
		}
	})

	b.Run("reference", func(b *testing.B) {
		vs := make([]Vec3, 0)

		for range 200 {
			vs = append(vs, NewVec3(2, 3, 4))
		}

		for b.Loop() {
			referenceTester(vs, 0, 3, 4)
		}
	})
}
