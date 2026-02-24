package transforms

import (
	"fmt"
)

const (
	Pi       = 3.14159265358979323846
	DegToRad = Pi / 180
)

type Transforms struct {
	Scale            Vec3
	Rotation         Vec3
	Position         Vec3
	ForwardDirection Vec3

	ScaleMat         Matrix
	RotationMat      Matrix
	TranslationMat   Matrix
	MatrixTransforms Matrix
}

func NewTransforms(pos, scale, rot Vec3) Transforms {
	return Transforms{
		Position: pos,
		Scale:    scale,
		Rotation: rot,
	}
}

func (t *Transforms) UpdateModelTransforms() {
	t.RotationMat = NewRotationMatrix(t.Rotation)
	t.ScaleMat = NewScaleMatrix(t.Scale)
	t.TranslationMat = NewTranslationMatrix(t.Position)

	t.MatrixTransforms = t.RotationMat.MultiplyByMatrix(t.ScaleMat)
	t.MatrixTransforms = t.TranslationMat.MultiplyByMatrix(t.MatrixTransforms)
}

func (t *Transforms) UpdateCameraTransforms() {
	t.RotationMat = NewRotationMatrix(t.Rotation).Transposed()
	t.ScaleMat = NewScaleMatrix(t.Scale)
	t.TranslationMat = NewTranslationMatrix(t.Position.Scale(-1))

	t.MatrixTransforms = t.ScaleMat.MultiplyByMatrix(t.RotationMat)
	t.MatrixTransforms = t.MatrixTransforms.MultiplyByMatrix(t.TranslationMat)
}

func (m Matrix) Print(name string) {
	fmt.Println(name, ":")
	for row := range MatLength {
		for col := range MatLength {
			fmt.Printf(" %6.2f", m[col*MatLength+row])
		}
		fmt.Print("\n")
	}
}
