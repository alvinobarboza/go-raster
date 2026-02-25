package mesh

import (
	"image/color"

	"github.com/alvinobarboza/go-raster/internal/maths"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type Texture struct {
	width, height int
	pixels        []color.RGBA
}

func (t *Texture) TexelColor(uv transforms.Vec2) color.RGBA {
	u := uv.X - maths.Floor32(uv.X)
	v := uv.Y - maths.Floor32(uv.Y)

	w := int(u * float32(t.width))
	h := int(v * float32(t.height))

	i := h*t.width + w

	if uint(i) < uint(len(t.pixels)) {
		return t.pixels[i]
	}

	return t.pixels[0]
}
