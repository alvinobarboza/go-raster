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
	Intensity float32

	Direction      transforms.Vec3
	DirectionWorld transforms.Vec3
}

func NewLight(t LightType, d transforms.Vec3, i float32, color color.RGBA) Light {
	return Light{
		LightType: t,
		Direction: d.Normalized(),
		Intensity: i,
		Color: transforms.NewVec3(
			float32(color.R)/255,
			float32(color.G)/255,
			float32(color.B)/255,
		),
	}
}
