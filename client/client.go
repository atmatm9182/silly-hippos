package main

import (
	"log"
	"net"

	"github.com/atmatm9182/silly-hippos/common"
	"github.com/atmatm9182/silly-hippos/common/message"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth = 800
const windowHeight = 600

var world = common.World{
	TileTypes: []common.TileType{
		common.GrassTile,
		common.GrassTile,
	},
}

var tileTextures = [common.TileCount]rl.Texture2D{}

func Vec2ToScreen(v common.Vector2[float32]) rl.Vector2 {
	return rl.Vector2{X: v.X * float32(windowWidth), Y: v.Y * float32(windowHeight)}
}

func TilePosFromIdx(idx int) rl.Vector2 {
	return rl.Vector2{
		X: float32(idx % common.WorldWidth),
		Y: float32(idx / common.WorldHeight),
	}
}

func DrawTile(idx int, ty common.TileType) {
	tex := tileTextures[ty]
	rl.DrawTextureEx(tex, TilePosFromIdx(idx), 0, 0.1, rl.Red)
}

func DrawWorld(world common.World) {
	for i, ty := range world.TileTypes {
		DrawTile(i, ty)
	}
}

// this needs to be called AFTER the window has been initialized
func LoadTextures() {
	tileTextures[common.GrassTile] = rl.LoadTexture("./assets/grass-tile.png")
	tileTextures[common.DirtTile] = rl.Texture2D{}
	tileTextures[common.PlainTile] = rl.Texture2D{}
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
			log.Printf("Dispatching player event %+v\n", e)
		default:
			return
		}
	}
}

func DispatchServerMessages() {
	for {
		select {
		case m := <-serverMessages:
			log.Printf("Dispatching server message %+v\n", *m)
		default:
			return
		}
	}
}

func UpdateWorld() {
	DispatchPlayerEvents()
	DispatchServerMessages()
}
