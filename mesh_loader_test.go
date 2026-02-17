package main

import (
	"testing"
)

func TestMeshLoader(t *testing.T) {
	t.Run("loadVec3 Vertex", func(t *testing.T) {
		lineToParse := "v 1.0 1.0 1.0"

		want := NewVec3(1, 1, 1)

		x, y, z := LoadVec3(lineToParse, 3, true)
		got := NewVec3(x, y, z)

		if want.X != got.X &&
			want.Y != got.Y &&
			want.Z != got.Z {
			t.Errorf("wanted %+v, got %+v", want, got)
		}
	})

	t.Run("loadVec3 UV", func(t *testing.T) {
		lineToParse := "v 1.0 1.0 0.0"

		want := NewVec3(1, 1, 0)

		x, y, z := LoadVec3(lineToParse, 3, true)
		got := NewVec3(x, y, z)

		if want.X != got.X &&
			want.Y != got.Y &&
			want.Z != got.Z {
			t.Errorf("wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex", func(t *testing.T) {
		lineToParse := "f 2 2 2"
		want := Triangle{
			v1: 1,
			v2: 1,
			v3: 1,
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].v1 != want.v1 &&
			got[0].v2 != want.v2 &&
			got[0].v3 != want.v3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex UV", func(t *testing.T) {
		lineToParse := "f 2/2 2/2 2/2\n\r"
		want := Triangle{
			v1: 1,
			v2: 1,
			v3: 1,
			u1: 1,
			u2: 1,
			u3: 1,
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].v1 != want.v1 &&
			got[0].v2 != want.v2 &&
			got[0].v3 != want.v3 &&
			got[0].u1 != want.u1 &&
			got[0].u2 != want.u2 &&
			got[0].u3 != want.u3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex Normal", func(t *testing.T) {
		lineToParse := "f 2//2 2//2 2//2\n\r"
		want := Triangle{
			v1: 1,
			v2: 1,
			v3: 1,
			n1: 1,
			n2: 1,
			n3: 1,
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].v1 != want.v1 &&
			got[0].v2 != want.v2 &&
			got[0].v3 != want.v3 &&
			got[0].n1 != want.n1 &&
			got[0].n2 != want.n2 &&
			got[0].n3 != want.n3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Full 1", func(t *testing.T) {
		lineToParse := "f 2/3/2 2/3/2 2/3/2\n\r"
		want := Triangle{
			v1: 1,
			v2: 1,
			v3: 1,
			u1: 2,
			u2: 2,
			u3: 2,
			n1: 1,
			n2: 1,
			n3: 1,
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].v1 != want.v1 &&
			got[0].v2 != want.v2 &&
			got[0].v3 != want.v3 &&
			got[0].n1 != want.n1 &&
			got[0].n2 != want.n2 &&
			got[0].n3 != want.n3 &&
			got[0].u1 != want.u1 &&
			got[0].u2 != want.u2 &&
			got[0].u3 != want.u3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Full 2", func(t *testing.T) {
		lineToParse := "f 2/3/4 5/6/7 8/9/10 11/12/13\n\r"
		want := []Triangle{
			{
				v1: 1, u1: 2, n1: 3,
				v2: 4, u2: 5, n2: 6,
				v3: 7, u3: 8, n3: 9,
			},
			{
				v1: 1, u1: 2, n1: 3,
				v2: 7, u2: 8, n2: 9,
				v3: 10, u3: 11, n3: 12,
			},
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 2 {
			t.Error("Didn't parse anything")
		}
		for i := range 2 {
			if got[i].v1 != want[i].v1 &&
				got[i].v2 != want[i].v2 &&
				got[i].v3 != want[i].v3 &&
				got[i].n1 != want[i].n1 &&
				got[i].n2 != want[i].n2 &&
				got[i].n3 != want[i].n3 &&
				got[i].u1 != want[i].u1 &&
				got[i].u2 != want[i].u2 &&
				got[i].u3 != want[i].u3 {
				t.Errorf("Wanted %+v, got %+v", want, got)
			}
		}
	})

	t.Run("LoadTriangle Full 3", func(t *testing.T) {
		lineToParse := "f 2/3/4 5/6/7 8/9/10 11/12/13 14/15/16\n\r"
		want := []Triangle{
			{
				v1: 1, u1: 2, n1: 3,
				v2: 4, u2: 5, n2: 6,
				v3: 7, u3: 8, n3: 9,
			},
			{
				v1: 1, u1: 2, n1: 3,
				v2: 7, u2: 8, n2: 9,
				v3: 10, u3: 11, n3: 12,
			},
			{
				v1: 1, u1: 2, n1: 3,
				v2: 10, u2: 11, n2: 12,
				v3: 13, u3: 14, n3: 15,
			},
		}

		got := LoadTriangle(lineToParse)

		if len(got) != 3 {
			t.Error("Didn't parse anything")
		}
		for i := range 3 {
			if got[i].v1 != want[i].v1 &&
				got[i].v2 != want[i].v2 &&
				got[i].v3 != want[i].v3 &&
				got[i].n1 != want[i].n1 &&
				got[i].n2 != want[i].n2 &&
				got[i].n3 != want[i].n3 &&
				got[i].u1 != want[i].u1 &&
				got[i].u2 != want[i].u2 &&
				got[i].u3 != want[i].u3 {
				t.Errorf("Wanted %+v, got %+v", want, got)
			}
		}
	})
}
