package common_test

import (
	"testing"
	"unsafe"

	"github.com/atmatm9182/silly-hippos/common"
	"github.com/atmatm9182/silly-hippos/common/types"
)

func sizeof[T any]() uintptr {
	var ptr *T
	return unsafe.Sizeof(*ptr)
}

func AssertEq[T comparable](t *testing.T, lhs, rhs T) {
	if lhs != rhs {
		t.Errorf("%v != %v", lhs, rhs)
	}
}

func TestTileSize(t *testing.T) {
	AssertEq(t, sizeof[common.Tile](), sizeof[types.Tile]())
}

func TestTileCount(t *testing.T) {
	AssertEq(t, int(common.TileCount), len(types.Tile_value))
}

func TestTileValues(t *testing.T) {
	AssertEq(t, int32(common.GrassTile), int32(types.Tile_GRASS))
	AssertEq(t, int32(common.PlainTile), int32(types.Tile_PLAIN))
	AssertEq(t, int32(common.DirtTile), int32(types.Tile_DIRT))
}
