package main

import "testing"

func TestBoundary(t *testing.T) {
	verts := []Vec3{
		NewVec3(1, 0, 0),
		NewVec3(-1, 0, 0),
	}

	bounds := NewBoundingSphere()

	bounds.CalculateBoundaries(verts, NewIdentityMatrix())

	if bounds.radius != 1 {
		t.Errorf("Expercted r=1, got %+v", bounds)
	}
}
