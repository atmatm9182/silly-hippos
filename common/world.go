package common

type TileType int

const (
	GrassTile TileType = iota
	PlainTile
	DirtTile
	TileCount
)

type World struct {
	Hippos    []Hippo
	TileTypes []TileType
}
