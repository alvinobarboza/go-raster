package lighting

import (
	"image/color"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type LightType uint

const (
	Directional LightType = iota
)

type Light struct {
	LightType LightType
	Color     transforms.Vec3

	Direction      transforms.Vec3
	DirectionWorld transforms.Vec3
}

func NewLight(t LightType, d transforms.Vec3, color color.RGBA) Light {
	return Light{
		LightType: t,
		Direction: d.Normalized(),
		Color: transforms.NewVec3(
			float32(color.R)/255,
			float32(color.G)/255,
			float32(color.B)/255,
		),
	}
}
