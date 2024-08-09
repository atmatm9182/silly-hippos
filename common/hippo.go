package common

import "github.com/atmatm9182/silly-hippos/common/types"

type Hippo struct {
	Name string
	Pos  Vector2
}

func (h *Hippo) ToProto() *types.Hippo {
	return &types.Hippo{
		Name: h.Name,
		Pos:  h.Pos.ToProto(),
	}
}

func HippoFromProto(h *types.Hippo) Hippo {
	return Hippo{
		Name: h.Name,
		Pos:  Vector2FromProto(h.Pos),
	}
}
