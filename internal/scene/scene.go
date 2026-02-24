package scene

import (
	"github.com/alvinobarboza/go-raster/internal/camera"
	"github.com/alvinobarboza/go-raster/internal/mesh"
)

type Scene struct {
	ActiveCam *camera.Camera
	Objects   []*mesh.Model
}

func NewScene(c *camera.Camera) *Scene {
	return &Scene{
		ActiveCam: c,
		Objects:   make([]*mesh.Model, 0),
	}
}

func (s *Scene) AddMesh(o *mesh.Model) {
	s.Objects = append(s.Objects, o)
}
