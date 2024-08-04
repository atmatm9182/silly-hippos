package common

import "github.com/atmatm9182/silly-hippos/common/types"

type Vector2[T any] struct {
	X T
	Y T
}

func Vector2ToProto(v Vector2[float32]) types.Vector2 {
	return types.Vector2{
		X: v.X,
		Y: v.Y,
	}
}
