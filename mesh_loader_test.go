package main

import "testing"

func TestMeshLoader(t *testing.T) {
	t.Run("loadFile", func(t *testing.T) {
		LoadMeshFromFile("./assets/cube.obj", "")
	})
}
