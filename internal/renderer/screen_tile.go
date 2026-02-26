package renderer

type ScreenTile struct {
	Width, Height float32
	OffW, OffH    float32

	FullWidth, FullHeight float32

	minX, minY float32
	maxX, maxY float32

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
	s.minX = 0 + offx
	s.minY = 0 + offy
	s.maxX = w + offx
	s.maxY = h + offy
}

func (s *ScreenTile) UpdateBufferSize(size int) {
	s.trianglesBuffer = make([]int, size)
}

func (s *ScreenTile) TileTriangleCollision(minX, minY, maxX, maxY float32) bool {
	return s.minX <= maxX && s.maxX >= minX && s.minY <= maxY && s.maxY >= minY
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

// minX, minY, maxX, maxY
func (s *ScreenTile) Bounduries() (float32, float32, float32, float32) {
	return s.minX, s.minY, s.maxX, s.maxY
}
