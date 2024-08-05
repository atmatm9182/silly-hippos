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

var numOfTileTests = 0

func AssertTilesEq(t *testing.T, lhs common.Tile, rhs types.Tile) {
	AssertEq(t, int32(lhs), int32(rhs))
	numOfTileTests++
}

func TestTileSize(t *testing.T) {
	AssertEq(t, sizeof[common.Tile](), sizeof[types.Tile]())
}

func TestTileCount(t *testing.T) {
	AssertEq(t, int(common.TileCount), len(types.Tile_value))
}

func TestTileValues(t *testing.T) {
	AssertTilesEq(t, common.GrassTile, types.Tile_GRASS)
	AssertTilesEq(t, common.PlainTile, types.Tile_PLAIN)
	AssertTilesEq(t, common.DirtTile, types.Tile_DIRT)
	AssertTilesEq(t, common.WaterTile, types.Tile_WATER)
	AssertTilesEq(t, common.MountainTile, types.Tile_MOUNTAIN)

	AssertEq(t, numOfTileTests, int(common.TileCount))
}
