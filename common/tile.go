package common

type Tile int32

const (
	GrassTile Tile = iota
	PlainTile
	DirtTile
	WaterTile
	MountainTile
	TileCount
)
