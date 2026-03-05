package mesh

import (
	"testing"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func TestMeshLoader(t *testing.T) {
	t.Run("loadVec3 Vertex", func(t *testing.T) {
		lineToParse := "v 1.0 1.0 1.0"

		want := transforms.NewVec3(1, 1, 1)

		x, y, z := LoadVec3(lineToParse, 3, true)
		got := transforms.NewVec3(x, y, z)

		if want.X != got.X &&
			want.Y != got.Y &&
			want.Z != got.Z {
			t.Errorf("wanted %+v, got %+v", want, got)
		}
	})

	t.Run("loadVec3 UV", func(t *testing.T) {
		lineToParse := "v 1.0 1.0 0.0"

		want := transforms.NewVec3(1, 1, 0)

		x, y, z := LoadVec3(lineToParse, 3, true)
		got := transforms.NewVec3(x, y, z)

		if want.X != got.X &&
			want.Y != got.Y &&
			want.Z != got.Z {
			t.Errorf("wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex", func(t *testing.T) {
		lineToParse := "f 2 2 2"
		want := Triangle{
			V1: 1,
			V2: 1,
			V3: 1,
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].V1 != want.V1 &&
			got[0].V2 != want.V2 &&
			got[0].V3 != want.V3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex UV", func(t *testing.T) {
		lineToParse := "f 2/2 2/2 2/2\n\r"
		want := Triangle{
			V1: 1,
			V2: 1,
			V3: 1,
			U1: 1,
			U2: 1,
			U3: 1,
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].V1 != want.V1 &&
			got[0].V2 != want.V2 &&
			got[0].V3 != want.V3 &&
			got[0].U1 != want.U1 &&
			got[0].U2 != want.U2 &&
			got[0].U3 != want.U3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Vertex Normal", func(t *testing.T) {
		lineToParse := "f 2//2 2//2 2//2\n\r"
		want := Triangle{
			V1: 1,
			V2: 1,
			V3: 1,
			N1: 1,
			N2: 1,
			N3: 1,
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].V1 != want.V1 &&
			got[0].V2 != want.V2 &&
			got[0].V3 != want.V3 &&
			got[0].N1 != want.N1 &&
			got[0].N2 != want.N2 &&
			got[0].N3 != want.N3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Full 1", func(t *testing.T) {
		lineToParse := "f 2/3/2 2/3/2 2/3/2\n\r"
		want := Triangle{
			V1: 1,
			V2: 1,
			V3: 1,
			U1: 2,
			U2: 2,
			U3: 2,
			N1: 1,
			N2: 1,
			N3: 1,
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 1 {
			t.Error("Didn't parse anything")
		}
		if got[0].V1 != want.V1 &&
			got[0].V2 != want.V2 &&
			got[0].V3 != want.V3 &&
			got[0].N1 != want.N1 &&
			got[0].N2 != want.N2 &&
			got[0].N3 != want.N3 &&
			got[0].U1 != want.U1 &&
			got[0].U2 != want.U2 &&
			got[0].U3 != want.U3 {
			t.Errorf("Wanted %+v, got %+v", want, got)
		}
	})

	t.Run("LoadTriangle Full 2", func(t *testing.T) {
		lineToParse := "f 2/3/4 5/6/7 8/9/10 11/12/13\n\r"
		want := []Triangle{
			{
				V1: 1, U1: 2, N1: 3,
				V2: 4, U2: 5, N2: 6,
				V3: 7, U3: 8, N3: 9,
			},
			{
				V1: 1, U1: 2, N1: 3,
				V2: 7, U2: 8, N2: 9,
				V3: 10, U3: 11, N3: 12,
			},
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 2 {
			t.Error("Didn't parse anything")
		}
		for i := range 2 {
			if got[i].V1 != want[i].V1 &&
				got[i].V2 != want[i].V2 &&
				got[i].V3 != want[i].V3 &&
				got[i].N1 != want[i].N1 &&
				got[i].N2 != want[i].N2 &&
				got[i].N3 != want[i].N3 &&
				got[i].U1 != want[i].U1 &&
				got[i].U2 != want[i].U2 &&
				got[i].U3 != want[i].U3 {
				t.Errorf("Wanted %+v, got %+v", want, got)
			}
		}
	})

	t.Run("LoadTriangle Full 3", func(t *testing.T) {
		lineToParse := "f 2/3/4 5/6/7 8/9/10 11/12/13 14/15/16\n\r"
		want := []Triangle{
			{
				V1: 1, U1: 2, N1: 3,
				V2: 4, U2: 5, N2: 6,
				V3: 7, U3: 8, N3: 9,
			},
			{
				V1: 1, U1: 2, N1: 3,
				V2: 7, U2: 8, N2: 9,
				V3: 10, U3: 11, N3: 12,
			},
			{
				V1: 1, U1: 2, N1: 3,
				V2: 10, U2: 11, N2: 12,
				V3: 13, U3: 14, N3: 15,
			},
		}

		got := LoadTriangle(lineToParse, true)

		if len(got) != 3 {
			t.Error("Didn't parse anything")
		}
		for i := range 3 {
			if got[i].V1 != want[i].V1 &&
				got[i].V2 != want[i].V2 &&
				got[i].V3 != want[i].V3 &&
				got[i].N1 != want[i].N1 &&
				got[i].N2 != want[i].N2 &&
				got[i].N3 != want[i].N3 &&
				got[i].U1 != want[i].U1 &&
				got[i].U2 != want[i].U2 &&
				got[i].U3 != want[i].U3 {
				t.Errorf("Wanted %+v, got %+v", want, got)
			}
		}
	})
}
