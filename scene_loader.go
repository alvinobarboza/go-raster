package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ModelData struct {
	MeshPath    string   `json:"meshpath"`
	TexturePath string   `json:"meshtexturepath"`
	Position    Vec3Data `json:"position"`
	Scale       Vec3Data `json:"scale"`
	Rotation    Vec3Data `json:"rotation"`
}

type Vec3Data struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func LoadSceneFromJSON(filePath string) ([]Model, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scene file: %w", err)
	}

	var modelsData []ModelData
	if err := json.Unmarshal(fileBytes, &modelsData); err != nil {
		return nil, fmt.Errorf("failed to parse scene json: %w", err)
	}

	var models []Model

	for _, m := range modelsData {
		mesh, err := LoadMeshFromFile(m.MeshPath, m.TexturePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load asset %s: %w", m.MeshPath, err)
		}

		fmt.Printf("tris: %d verts: %d\n", len(mesh.tris), len(mesh.verts))

		transforms := NewTransforms(
			NewVec3(m.Position.X, m.Position.Y, m.Position.Z),
			NewVec3(m.Scale.X, m.Scale.Y, m.Scale.Z),
			NewVec3(m.Rotation.X, m.Rotation.Y, m.Rotation.Z),
		)

		model := NewModel(&mesh, transforms)
		models = append(models, model)
	}

	return models, nil
}
