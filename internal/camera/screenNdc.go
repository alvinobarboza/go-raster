package camera

type NDCPoint struct {
	X, Y float32
}

type ScreenPoint struct {
	X, Y float32
}

func (sp *ScreenPoint) IsTopOrLeft(sp2 ScreenPoint) bool {
	edge := ScreenPoint{X: sp2.X - sp.X, Y: sp2.Y - sp.Y}
	isTopEdge := edge.Y == 0 && edge.X > 0
	isLeftEdge := edge.Y < 0
	return isTopEdge || isLeftEdge
}

func EdgeCross(a, b, p ScreenPoint) float32 {
	abX := b.X - a.X
	abY := b.Y - a.Y

	apX := p.X - a.X
	apY := p.Y - a.Y

	return (abX * apY) - (abY * apX)
}
