package common

type TileType int

const (
	GrassTile TileType = iota
	PlainTile
	DirtTile
	TileCount
)

type World struct {
	TilePositions []Vector2[float32]
	TileTypes []TileType
}
