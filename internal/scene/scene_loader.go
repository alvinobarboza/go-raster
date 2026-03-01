package scene

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alvinobarboza/go-raster/internal/mesh"
	"github.com/alvinobarboza/go-raster/internal/transforms"
)

type ModelData struct {
	MeshPath       string   `json:"meshpath"`
	TexturePath    string   `json:"meshtexturepath"`
	Position       Vec3Data `json:"position"`
	Scale          Vec3Data `json:"scale"`
	Rotation       Vec3Data `json:"rotation"`
	FlipNormals    bool     `json:"flipNormals"`
	WindingReorder bool     `json:"windingReorder"`
	ZNegative      bool     `json:"zNegative"`
}

type Vec3Data struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func LoadSceneFromJSON(filePath string) ([]mesh.Model, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scene file: %w", err)
	}

	var modelsData []ModelData
	if err := json.Unmarshal(fileBytes, &modelsData); err != nil {
		return nil, fmt.Errorf("failed to parse scene json: %w", err)
	}

	var models []mesh.Model

	for _, m := range modelsData {
		meshData, err := mesh.LoadMeshFromFile(
			m.MeshPath, m.TexturePath, m.ZNegative,
			m.WindingReorder, m.FlipNormals)
		if err != nil {
			return nil, fmt.Errorf("failed to load asset %s: %w", m.MeshPath, err)
		}

		fmt.Printf("tris: %d verts: %d norm: %d\n",
			len(meshData.Tris), len(meshData.Verts), len(meshData.Normals))

		transforms := transforms.NewTransforms(
			transforms.NewVec3(m.Position.X, m.Position.Y, m.Position.Z),
			transforms.NewVec3(m.Scale.X, m.Scale.Y, m.Scale.Z),
			transforms.NewVec3(m.Rotation.X, m.Rotation.Y, m.Rotation.Z),
		)

		model := mesh.NewModel(&meshData, transforms)
		models = append(models, model)
	}

	return models, nil
}
