package main

import (
	"math/rand"

	"github.com/atmatm9182/silly-hippos/common"
	"github.com/atmatm9182/silly-hippos/hwg"
)

func main() {
	seed := getRandomSeed()
	params := hwg.GenerationParams{
		Width:     common.WorldWidth,
		Height:    common.WorldHeight,
		NoiseFunc: hwg.PerlinNoise,
	}
	tiles := hwg.GenerateHippoWorld(seed, params)
	StartHippoServer(6969, tiles)
}

func getRandomSeed() (seed [32]byte) {
	for i := 0; i < 32; i++ {
		seed[i] = byte(rand.Intn(255))
	}

	return
}
