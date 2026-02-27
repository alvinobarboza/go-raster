package renderer

import (
	"testing"

	"github.com/alvinobarboza/go-raster/internal/mesh"
)

func TestCollision(t *testing.T) {
	t.Run("collide", func(t *testing.T) {
		aabb := mesh.NewAABB2(0, 0, 2, 2)

		tile1 := ScreenTile{
			Aabb: mesh.NewAABB2(1, 1, 3, 3),
		}

		if !tile1.TileTriangleCollision(aabb) {
			t.Errorf("Expected collision!")
		}

		tile2 := ScreenTile{
			Aabb: mesh.NewAABB2(2.5, 2.5, 3, 3),
		}

		if tile2.TileTriangleCollision(aabb) {
			t.Errorf("Expected no collision!")
		}
	})

	t.Run("collide offset init", func(t *testing.T) {
		aabb := mesh.NewAABB2(15, 5, 19, 12)

		tile1 := NewScreenTile(10, 10, 100, 100, 10, 10, 0)

		if !tile1.TileTriangleCollision(aabb) {
			t.Errorf("Expected collision!")
		}

		tile2 := NewScreenTile(10, 10, 100, 100, 0, 10, 0)

		if tile2.TileTriangleCollision(aabb) {
			t.Errorf("Expected no collision!")
		}
	})
}

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

func BenchmarkTriangleTileCollision(b *testing.B) {
	s := ScreenTile{
		Aabb: mesh.NewAABB2(0, 0, 10, 10),
	}

	aabb := mesh.NewAABB2(5, 5, 15, 15)

	b.Run("AABB", func(b *testing.B) {
		for b.Loop() {
			s.TileTriangleCollision(aabb)
		}
	})
}
