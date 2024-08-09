package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Not enough arguments\n")
		os.Exit(1)
	}

	rl.InitWindow(800, 600, "Silly Hippos")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	LoadTextures()

	serverAddr := args[0]
	ConnectToTheServer(serverAddr)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Red)

		CollectEvents()
		UpdateWorld()

		camera.Target = Vector2ToRl(me.Pos.AddScalar(0.5).MulScalar(tileSize))

		rl.BeginMode2D(camera)
		DrawWorld(world)
		rl.EndMode2D()

		rl.EndDrawing()
	}
}
