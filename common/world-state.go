package common

import (
	"github.com/atmatm9182/silly-hippos/common/types"
	"unsafe"
)

type WorldState struct {
	Hippos []Hippo
	Tiles  []Tile
}

func (ws *WorldState) ToProto() *types.WorldState {
	hippos := make([]*types.Hippo, len(ws.Hippos))

	for _, h := range ws.Hippos {
		pos := new(types.Vector2)
		*pos = Vector2ToProto(h.Pos)

		hippos = append(hippos, &types.Hippo{
			Pos:  pos,
			Name: h.Name,
		})
	}

	// Do this to avoid allocating a new slice.
	// This is 100% safe if the tests of this package are passing
	tilePtr := (*types.Tile)(unsafe.Pointer(unsafe.SliceData(ws.Tiles)))
	tiles := unsafe.Slice(tilePtr, len(ws.Tiles))

	return &types.WorldState{
		Hippos: hippos,
		Tiles:  tiles,
	}
}
