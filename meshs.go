package main

import "image/color"

type Mesh struct {
	vertices []Vec3
	color    color.RGBA
}

func NewMesh(v []Vec3, color color.RGBA) Mesh {
	return Mesh{
		vertices: v,
		color:    color,
	}
}
