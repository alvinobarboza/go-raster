package renderer

import "testing"

func TestCollision(t *testing.T) {
	t.Run("collide", func(t *testing.T) {
		minX := float32(0)
		minY := float32(0)
		maxX := float32(2)
		maxY := float32(2)

		tile1 := ScreenTile{
			minX: 1,
			minY: 1,
			maxX: 3,
			maxY: 3,
		}

		if !tile1.TileTriangleCollision(minX, minY, maxX, maxY) {
			t.Errorf("Expected collision!")
		}

		tile2 := ScreenTile{
			minX: 2.5,
			minY: 2.5,
			maxX: 3,
			maxY: 3,
		}

		if tile2.TileTriangleCollision(minX, minY, maxX, maxY) {
			t.Errorf("Expected no collision!")
		}
	})

	t.Run("collide offset init", func(t *testing.T) {
		minX := float32(15)
		minY := float32(5)
		maxX := float32(19)
		maxY := float32(12)

		tile1 := NewScreenTile(10, 10, 100, 100, 10, 10, 0)

		if !tile1.TileTriangleCollision(minX, minY, maxX, maxY) {
			t.Errorf("Expected collision!")
		}

		tile2 := NewScreenTile(10, 10, 100, 100, 0, 10, 0)

		if tile2.TileTriangleCollision(minX, minY, maxX, maxY) {
			t.Errorf("Expected no collision!")
		}
	})
}
