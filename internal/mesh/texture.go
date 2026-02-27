package mesh

import (
	"image/color"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

// Must be powers of two 256, 512 ...
type Texture struct {
	width, height         int
	fWidth, fHeight       float32
	widthMask, heightMask int
	pixels                []color.RGBA
}

func (t *Texture) TexelColor(uv transforms.Vec2) color.RGBA {
	w := int(uv.X*t.fWidth) & t.widthMask
	h := int(uv.Y*t.fHeight) & t.heightMask

	return t.pixels[h*t.width+w]
}
