package renderer

import "github.com/alvinobarboza/go-raster/internal/mesh"

type ScreenTile struct {
	Width, Height float32
	OffW, OffH    float32

	FullWidth, FullHeight float32

	minX, minY float32
	maxX, maxY float32

	trianglesBuffer []mesh.FullTriangle
}

func NewScreenTile(w, h, fw, fh, offx, offy float32, buffSize int) *ScreenTile {
	s := &ScreenTile{}
	s.UpdateTileSize(w, h, fw, fh, offx, offy)
	s.UpdateBufferSize(buffSize)
	s.ResetBuff()
	return s
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
	s.trianglesBuffer = make([]mesh.FullTriangle, size)
}

func (s *ScreenTile) TileTriangleCollision(minX, minY, maxX, maxY float32) bool {
	return s.minX <= maxX && s.maxX >= minX && s.minY <= maxY && s.maxY >= minY
}

func (s *ScreenTile) AddTriangle(t mesh.FullTriangle) {
	s.trianglesBuffer = append(s.trianglesBuffer, t)
}

func (s *ScreenTile) ResetBuff() {
	s.trianglesBuffer = s.trianglesBuffer[:0]
}

func (s *ScreenTile) Triangles() []mesh.FullTriangle {
	return s.trianglesBuffer
}

// minX, minY, maxX, maxY
func (s *ScreenTile) Bounduries() (float32, float32, float32, float32) {
	return s.minX, s.minY, s.maxX, s.maxY
}
