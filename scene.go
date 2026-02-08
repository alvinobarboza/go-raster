package main

type Scene struct {
	activeCam *Camera
	objects   []Mesh
}

func NewScene(c *Camera) Scene {
	return Scene{
		activeCam: c,
		objects:   make([]Mesh, 0),
	}
}

func (s *Scene) AddMesh(o Mesh) {
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
