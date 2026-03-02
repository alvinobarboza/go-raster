package lighting

import "github.com/alvinobarboza/go-raster/internal/transforms"

type LightType uint

const (
	Directional LightType = iota
)

type Light struct {
	LightType          LightType
	IntesityMultiplier float32

	Direction      transforms.Vec3
	DirectionWorld transforms.Vec3
}

func NewLight(t LightType, i float32, d transforms.Vec3) Light {
	return Light{
		LightType:          t,
		IntesityMultiplier: i,
		Direction:          d.Normalized(),
	}
}
