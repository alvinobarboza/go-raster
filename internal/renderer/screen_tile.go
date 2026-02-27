package renderer

import (
	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type ScreenTile struct {
	Width, Height float32
	OffW, OffH    float32

	FullWidth, FullHeight float32

	Aabb mesh.AABB2

	trianglesBuffer []int

	IsActive bool
}

func NewScreenTile(w, h, fw, fh, offx, offy float32, buffSize int) *ScreenTile {
	s := &ScreenTile{}
	s.UpdateTileSize(w, h, fw, fh, offx, offy)
	s.UpdateBufferSize(buffSize)
	s.ResetBuff()
	return s
}

func RecalculateTiles(w, h, tileLength float32, buffSize int, tiles []*ScreenTile) []*ScreenTile {
	wOffSet := float32(0)
	hOffSet := float32(0)
	tW, tH := tileLength, tileLength

	i := 0
	isNew := len(tiles) == 0

	for {
		if !isNew && i < len(tiles) {
			tiles[i].UpdateTileSize(tW, tH, w, h, wOffSet, hOffSet)
			i++
		} else {
			tt := NewScreenTile(tW, tH, w, h, wOffSet, hOffSet, buffSize)
			// fmt.Printf("%+v\n", tt)
			tiles = append(tiles, tt)
		}

		wOffSet += tileLength
		offOffSetW := wOffSet + tileLength

		if offOffSetW > w {
			if offOffSetW-w < tileLength && offOffSetW-w > 0 {
				wOffSet = w - (tileLength - (offOffSetW - w))
				tW = tileLength - (offOffSetW - w)
			} else {
				tW = tileLength
				wOffSet = 0
				hOffSet += tileLength
			}
		}

		offOffSetH := hOffSet + tileLength

		if offOffSetH > h {
			if offOffSetH-h < tileLength && offOffSetH-h > 0 {
				hOffSet = h - (tileLength - (offOffSetH - h))
				tH = tileLength - (offOffSetH - h)
			} else {
				tH = tileLength
			}
		}

		if hOffSet >= h || wOffSet >= w {
			if i > 0 {
				tiles = tiles[:i]
			}
			break
		}
	}

	return tiles
}

func NewTileSet(w, h, tileLength float32, buffSize int) []*ScreenTile {
	tiles := make([]*ScreenTile, 0)
	tiles = RecalculateTiles(w, h, tileLength, buffSize, tiles)
	return tiles
}

func (s *ScreenTile) UpdateTileSize(w, h, fw, fh, offx, offy float32) {
	s.Width = w
	s.Height = h
	s.FullWidth = fw
	s.FullHeight = fh
	s.OffW = offx
	s.OffH = offy
	s.Aabb.Min.X = 0 + offx
	s.Aabb.Min.Y = 0 + offy
	s.Aabb.Max.X = w + offx
	s.Aabb.Max.Y = h + offy
}

func (s *ScreenTile) UpdateBufferSize(size int) {
	s.trianglesBuffer = make([]int, size)
}

// AI helped SAT collision
func (s *ScreenTile) TileTriangleCollision(v0, v1, v2 transforms.Vec2, triAabb mesh.AABB2) bool {
	if s.Aabb.Collide(triAabb) {
		// 1. Broad-phase AABB test (Axes X and Y)
		triMinX := maths.Minf(v0.X, maths.Minf(v1.X, v2.X))
		triMaxX := maths.Maxf(v0.X, maths.Maxf(v1.X, v2.X))

		if triMinX > s.Aabb.Max.X || triMaxX < s.Aabb.Min.X {
			return false
		}

		triMinY := maths.Minf(v0.Y, maths.Minf(v1.Y, v2.Y))
		triMaxY := maths.Maxf(v0.Y, maths.Maxf(v1.Y, v2.Y))

		if triMinY > s.Aabb.Max.Y || triMaxY < s.Aabb.Min.Y {
			return false
		}

		// 2. Triangle edge axes test
		center := transforms.Vec2{
			X: (s.Aabb.Min.X + s.Aabb.Max.X) * 0.5,
			Y: (s.Aabb.Min.Y + s.Aabb.Max.Y) * 0.5,
		}

		extents := transforms.Vec2{
			X: (s.Aabb.Max.X - s.Aabb.Min.X) * 0.5,
			Y: (s.Aabb.Max.Y - s.Aabb.Min.Y) * 0.5,
		}

		edges := [3]transforms.Vec2{
			{X: v1.X - v0.X, Y: v1.Y - v0.Y},
			{X: v2.X - v1.X, Y: v2.Y - v1.Y},
			{X: v0.X - v2.X, Y: v0.Y - v2.Y},
		}

		for _, edge := range edges {
			// Perpendicular normal: (-y, x)
			nx, ny := -edge.Y, edge.X

			// Project AABB
			r := extents.X*maths.Abs(nx) + extents.Y*maths.Abs(ny)
			pc := center.X*nx + center.Y*ny

			// Project Triangle vertices
			p0 := v0.X*nx + v0.Y*ny
			p1 := v1.X*nx + v1.Y*ny
			p2 := v2.X*nx + v2.Y*ny

			triMin := maths.Minf(p0, maths.Minf(p1, p2))
			triMax := maths.Maxf(p0, maths.Maxf(p1, p2))

			// Check for separation
			if triMin > pc+r || triMax < pc-r {
				return false
			}
		}

		return true
	}
	return false
}

func (s *ScreenTile) AddTriangle(index int) {
	s.trianglesBuffer = append(s.trianglesBuffer, index)
}

func (s *ScreenTile) ResetBuff() {
	s.trianglesBuffer = s.trianglesBuffer[:0]
}

func (s *ScreenTile) Triangles() []int {
	return s.trianglesBuffer
}
