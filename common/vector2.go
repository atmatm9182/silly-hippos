package common

import "github.com/atmatm9182/silly-hippos/common/types"

type Vector2 struct {
	X float32
	Y float32
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2) MulScalar(x float32) Vector2 {
	return Vector2{
		X: v.X * x,
		Y: v.Y * x,
	}
}

func (v Vector2) AddScalar(x float32) Vector2 {
	return Vector2{
		X: v.X + x,
		Y: v.Y + x,
	}
}

func (v Vector2) ToProto() *types.Vector2 {
	return &types.Vector2{
		X: v.X,
		Y: v.Y,
	}
}

func Vector2FromProto(v *types.Vector2) Vector2 {
	return Vector2{
		X: v.X,
		Y: v.Y,
	}
}
