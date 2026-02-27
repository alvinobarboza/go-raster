package renderer

import (
	"math"
	"testing"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func TestTileAutoGen(t *testing.T) {
	t.Run("Tile Calc", func(t *testing.T) {
		const w, h = 14, 12
		tileLength := 4

		tiles := NewTileSet(w, h, float32(tileLength), 0)

		if len(tiles) != 12 {
			t.Errorf("Want len 12, got %d", len(tiles))
		}
	})
}

// AI generated

type Triangle struct {
	V0, V1, V2 transforms.Vec2
}

func CheckAABBTriangle(aabb mesh.AABB2, tri Triangle) bool {
	// 1. Broad-phase AABB test (Axes X and Y)
	triMinX := maths.Minf(tri.V0.X, maths.Minf(tri.V1.X, tri.V2.X))
	triMaxX := maths.Maxf(tri.V0.X, maths.Maxf(tri.V1.X, tri.V2.X))

	if triMinX > aabb.Max.X || triMaxX < aabb.Min.X {
		return false
	}

	triMinY := maths.Minf(tri.V0.Y, maths.Minf(tri.V1.Y, tri.V2.Y))
	triMaxY := maths.Maxf(tri.V0.Y, maths.Maxf(tri.V1.Y, tri.V2.Y))

	if triMinY > aabb.Max.Y || triMaxY < aabb.Min.Y {
		return false
	}

	// 2. Triangle edge axes test
	center := transforms.Vec2{
		X: (aabb.Min.X + aabb.Max.X) * 0.5,
		Y: (aabb.Min.Y + aabb.Max.Y) * 0.5,
	}

	extents := transforms.Vec2{
		X: (aabb.Max.X - aabb.Min.X) * 0.5,
		Y: (aabb.Max.Y - aabb.Min.Y) * 0.5,
	}

	edges := [3]transforms.Vec2{
		{X: tri.V1.X - tri.V0.X, Y: tri.V1.Y - tri.V0.Y},
		{X: tri.V2.X - tri.V1.X, Y: tri.V2.Y - tri.V1.Y},
		{X: tri.V0.X - tri.V2.X, Y: tri.V0.Y - tri.V2.Y},
	}

	for _, edge := range edges {
		// Perpendicular normal: (-y, x)
		nx, ny := -edge.Y, edge.X

		// Project AABB
		r := extents.X*float32(math.Abs(float64(nx))) + extents.Y*float32(math.Abs(float64(ny)))
		pc := center.X*nx + center.Y*ny

		// Project Triangle vertices
		p0 := tri.V0.X*nx + tri.V0.Y*ny
		p1 := tri.V1.X*nx + tri.V1.Y*ny
		p2 := tri.V2.X*nx + tri.V2.Y*ny

		triMin := maths.Minf(p0, maths.Minf(p1, p2))
		triMax := maths.Maxf(p0, maths.Maxf(p1, p2))

		// Check for separation
		if triMin > pc+r || triMax < pc-r {
			return false
		}
	}

	return true
}

func TestCheckAABBTriangle(t *testing.T) {
	tests := []struct {
		name string
		aabb mesh.AABB2
		tri  Triangle
		want bool
	}{
		{
			name: "Triangle fully inside AABB",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 10, Y: 10}},
			tri:  Triangle{V0: transforms.Vec2{X: 2, Y: 2}, V1: transforms.Vec2{X: 8, Y: 2}, V2: transforms.Vec2{X: 5, Y: 8}},
			want: true,
		},
		{
			name: "AABB fully inside Triangle",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 4, Y: 4}, Max: transforms.Vec2{X: 6, Y: 6}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 0}, V1: transforms.Vec2{X: 10, Y: 0}, V2: transforms.Vec2{X: 5, Y: 10}},
			want: true,
		},
		{
			name: "Partial overlap (Vertex inside AABB)",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 4, Y: 4}, V1: transforms.Vec2{X: 10, Y: 4}, V2: transforms.Vec2{X: 4, Y: 10}},
			want: true,
		},
		{
			name: "Edge crossing (No vertices inside AABB)",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 2, Y: 2}, Max: transforms.Vec2{X: 4, Y: 4}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 3}, V1: transforms.Vec2{X: 6, Y: 3}, V2: transforms.Vec2{X: 3, Y: 6}},
			want: true,
		},
		{
			name: "Separated on X axis",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 6, Y: 0}, V1: transforms.Vec2{X: 10, Y: 0}, V2: transforms.Vec2{X: 8, Y: 5}},
			want: false,
		},
		{
			name: "Separated on Y axis",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 6}, V1: transforms.Vec2{X: 5, Y: 6}, V2: transforms.Vec2{X: 2.5, Y: 10}},
			want: false,
		},
		{
			name: "SAT Diagonal Separation (Overlaps X and Y bounds, separated by edge)",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 1, Y: 1}},
			tri:  Triangle{V0: transforms.Vec2{X: 3, Y: 0}, V1: transforms.Vec2{X: 0, Y: 3}, V2: transforms.Vec2{X: 3, Y: 3}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckAABBTriangle(tt.aabb, tt.tri); got != tt.want {
				t.Errorf("CheckAABBTriangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCheckAABBTriangle(b *testing.B) {
	benchmarks := []struct {
		name string
		aabb mesh.AABB2
		tri  Triangle
	}{
		{
			name: "TriInsideAABB",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 10, Y: 10}},
			tri:  Triangle{V0: transforms.Vec2{X: 2, Y: 2}, V1: transforms.Vec2{X: 8, Y: 2}, V2: transforms.Vec2{X: 5, Y: 8}},
		},
		{
			name: "AABBInsideTri",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 4, Y: 4}, Max: transforms.Vec2{X: 6, Y: 6}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 0}, V1: transforms.Vec2{X: 10, Y: 0}, V2: transforms.Vec2{X: 5, Y: 10}},
		},
		{
			name: "PartialOverlap",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 4, Y: 4}, V1: transforms.Vec2{X: 10, Y: 4}, V2: transforms.Vec2{X: 4, Y: 10}},
		},
		{
			name: "EdgeCrossing",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 2, Y: 2}, Max: transforms.Vec2{X: 4, Y: 4}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 3}, V1: transforms.Vec2{X: 6, Y: 3}, V2: transforms.Vec2{X: 3, Y: 6}},
		},
		{
			name: "SeparatedX",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 6, Y: 0}, V1: transforms.Vec2{X: 10, Y: 0}, V2: transforms.Vec2{X: 8, Y: 5}},
		},
		{
			name: "SeparatedY",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 5, Y: 5}},
			tri:  Triangle{V0: transforms.Vec2{X: 0, Y: 6}, V1: transforms.Vec2{X: 5, Y: 6}, V2: transforms.Vec2{X: 2.5, Y: 10}},
		},
		{
			name: "SATDiagonal",
			aabb: mesh.AABB2{Min: transforms.Vec2{X: 0, Y: 0}, Max: transforms.Vec2{X: 1, Y: 1}},
			tri:  Triangle{V0: transforms.Vec2{X: 3, Y: 0}, V1: transforms.Vec2{X: 0, Y: 3}, V2: transforms.Vec2{X: 3, Y: 3}},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				// The result is assigned to a blank identifier to prevent the
				// compiler from optimizing the function call away.
				_ = CheckAABBTriangle(bm.aabb, bm.tri)
			}
		})
	}
}

// AI generated
