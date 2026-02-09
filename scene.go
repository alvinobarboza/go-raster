package main

import "log"

type Scene struct {
	activeCam *Camera
	objects   []*Model
}

func NewScene(c *Camera) Scene {
	return Scene{
		activeCam: c,
		objects:   make([]*Model, 0),
	}
}

func (s *Scene) AddMesh(o *Model) {
	s.objects = append(s.objects, o)
}

func (s Scene) Render() {
	for _, o := range s.objects {
		for _, v := range o.mesh.verts {
			v = o.transforms.matrixTransforms.MultiplyByVec3(v)
			log.Println(o.boundingSphere)

			p := s.activeCam.ProjectVertex(v)
			p.color = o.mesh.tris[0].color
			s.activeCam.PutPixel(p)
		}
	}
}
