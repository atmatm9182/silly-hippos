package hwg

import (
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/atmatm9182/silly-hippos/common"
)

type GenerationParams struct {
	NoiseFunc func(x, y float64) float64
	Width     int
	Height    int
}

func GenerateHippoWorld(seed [32]byte, params GenerationParams) []common.Tile {
	x1 := binary.NativeEndian.Uint64(seed[:8])
	x2 := binary.NativeEndian.Uint64(seed[8:])

	y1 := binary.NativeEndian.Uint64(seed[8:16])
	y2 := binary.NativeEndian.Uint64(seed[16:])
	_, _ = x2, y2

	s := (x1 >> 9) ^ x2 ^ y1 ^ (y2 << 13)
	r := rand.New(rand.NewSource(int64(s)))

	g := r.Float64()
	for g == 0 {
		g = r.Float64()
	}

	tiles := make([]common.Tile, 0, params.Width*params.Height)
	for x := 0; x < params.Width; x++ {
		for y := 0; y < params.Height; y++ {
			n := params.NoiseFunc(float64(x)/g, float64(y)/g)
			t := getTileBasedOnNoise(n)
			tiles = append(tiles, t)
		}
	}

	return tiles
}

func getTileBasedOnNoise(noise float64) common.Tile {
	noise = noise/2 + 0.5
	h := float64(common.TileCount - 1)
	v := math.Round(lerp(noise, 0, h))
	return common.Tile(v)
}
