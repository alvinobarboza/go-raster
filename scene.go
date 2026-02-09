package main

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
		for _, v := range o.vertices {
			p := s.activeCam.ProjectVertex(v)
			p.color = o.color
			s.activeCam.PutPixel(p)
		}
	}
}
