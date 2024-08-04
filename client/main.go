package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/atmatm9182/silly-hippos/common"
	"github.com/atmatm9182/silly-hippos/common/message"
	rl "github.com/gen2brain/raylib-go/raylib"
	"google.golang.org/protobuf/proto"
)

const windowWidth = 800
const windowHeight = 600

var world = common.World{
	TilePositions: []common.Vector2[float32]{
		{0.3, 0.5},
		{0.5, 0.5},
	},
	TileTypes: []common.TileType{
		common.GrassTile,
		common.GrassTile,
	},
}

var tileTextures = [common.TileCount]rl.Texture2D{}

func Vec2ToScreen(v common.Vector2[float32]) rl.Vector2 {
	return rl.Vector2{X: v.X * float32(windowWidth), Y: v.Y * float32(windowHeight)}
}

func DrawTile(pos common.Vector2[float32], ty common.TileType) {
	tex := tileTextures[ty]
	rl.DrawTextureEx(tex, Vec2ToScreen(pos), 0, 0.1, rl.Red)
}

func DrawWorld(world common.World) {
	for i, pos := range world.TilePositions {
		ty := world.TileTypes[i]
		DrawTile(pos, ty)
	}
}

// this needs to be called AFTER the window has been initialized
func LoadTextures() {
	tileTextures[common.GrassTile] = rl.LoadTexture("./assets/grass-tile.png")
	tileTextures[common.DirtTile] = rl.Texture2D{}
	tileTextures[common.PlainTile] = rl.Texture2D{}
}

type PlayerEvent interface {
	ToMessage() proto.Message
}

type PlayerMovedEvent uint8

const (
	PlayerMovedUp PlayerMovedEvent = iota
	PlayerMovedDown
	PlayerMovedLeft
	PlayerMovedRight
)

func (p PlayerMovedEvent) ToMessage() proto.Message {
	var where message.MoveDirection
	switch p {
	case PlayerMovedUp:
		where = message.MoveDirection_UP
	case PlayerMovedDown:
		where = message.MoveDirection_DOWN
	case PlayerMovedLeft:
		where = message.MoveDirection_LEFT
	case PlayerMovedRight:
		where = message.MoveDirection_RIGHT
	}

	return &message.Moved{
		Where: where,
	}
}

var (
	conn net.Conn

	playerEvents   = make(chan PlayerEvent, 64)
	serverMessages = make(chan *message.Message, 64)
)

func ListenToServerMessages() {
	var (
		msg message.Message
		err error
	)

	for {
		err = common.ReadMessage(conn, &msg)
		if err != nil {
			log.Fatalln(err)
		}

		serverMessages <- &msg
	}
}

func ConnectToTheServer(addr string) {
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to the server successfully")

	lemmeIn := message.LemmeIn{
		Name: "hipik",
	}

	err = common.WriteMessage(conn, &lemmeIn)
	if err != nil {
		log.Fatalln(err)
	}

	var discover message.Discover
	err = common.ReadMessage(conn, &discover)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Got the discover message: %+v\n", discover)

	go ListenToServerMessages()
}

func CollectEvents() {
	if rl.IsKeyPressed(rl.KeyW) {
		playerEvents <- PlayerMovedUp
	}
	if rl.IsKeyPressed(rl.KeyS) {
		playerEvents <- PlayerMovedDown
	}
	if rl.IsKeyPressed(rl.KeyA) {
		playerEvents <- PlayerMovedLeft
	}
	if rl.IsKeyPressed(rl.KeyD) {
		playerEvents <- PlayerMovedRight
	}
}

func DispatchPlayerEvents() {
	for {
		select {
		case e := <-playerEvents:
			fmt.Printf("Dispatching player event %+v\n", e)
		default:
			return
		}
	}
}

func DispatchServerMessages() {
	for {
		select {
		case m := <-serverMessages:
			fmt.Printf("Dispatching server message %+v\n", *m)
		default:
			return
		}
	}
}

func UpdateWorld() {
	DispatchPlayerEvents()
	DispatchServerMessages()
}

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
		DrawWorld(world)

		rl.EndDrawing()
	}
}
