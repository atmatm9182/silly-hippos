package common

import (
	"github.com/atmatm9182/silly-hippos/common/types"
	"unsafe"
)

type WorldState struct {
	Hippos []Hippo
	Tiles  []Tile
}

func (ws *WorldState) GetTileAt(x, y int) Tile {
	idx := y * WorldWidth + x
	return ws.Tiles[idx]
}

func (ws *WorldState) ToProto() *types.WorldState {
	hippos := make([]*types.Hippo, len(ws.Hippos))

	for _, h := range ws.Hippos {
		hippos = append(hippos, h.ToProto())
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

func WorldStateFromProto(state *types.WorldState) WorldState {
	hippos := make([]Hippo, 0, len(state.Hippos))
	for _, h := range state.Hippos {
		hippos = append(hippos, HippoFromProto(h))
	}

	tilePtr := (*Tile)(unsafe.Pointer(unsafe.SliceData(state.Tiles)))
	tiles := unsafe.Slice(tilePtr, len(state.Tiles))

	return WorldState{
		Hippos: hippos,
		Tiles:  tiles,
	}
}
