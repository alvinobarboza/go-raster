package scene

import (
	"github.com/alvinobarboza/go-raster/internal/camera"
	"github.com/alvinobarboza/go-raster/internal/lighting"
	"github.com/alvinobarboza/go-raster/internal/mesh"
)

type Scene struct {
	ActiveCam            *camera.Camera
	Objects              []*mesh.Model
	SkyBox               *mesh.Model
	Lights               []lighting.Light
	AmbientLightStrength float32 // maybe later I"ll add an envioriment struct wiht skybox, this and others
}

func NewScene(c *camera.Camera) *Scene {
	return &Scene{
		ActiveCam:            c,
		Objects:              make([]*mesh.Model, 0),
		Lights:               make([]lighting.Light, 0),
		AmbientLightStrength: 0.2, // fixed for now
	}
}

func (s *Scene) AddMesh(o *mesh.Model) {
	s.Objects = append(s.Objects, o)
}

func (s *Scene) AddLight(l lighting.Light) {
	s.Lights = append(s.Lights, l)
}

func (s *Scene) UpdateLights() {
	for i := range s.Lights {
		// rotate and invert, as the dot expects the normal to be aligned wiht the light, in other worlds,
		// light normal is opposit of its direction
		s.Lights[i].DirectionWorld = s.ActiveCam.Transforms.RotationMat.MultiplyByVec3(s.Lights[i].Direction).Normalized().Scale(-1)
	}
}
