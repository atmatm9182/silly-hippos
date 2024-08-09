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

const tileSize float32 = 64.0

var (
	tileTexture rl.Texture2D
	camera      rl.Camera2D
)

var (
	world common.WorldState
	me    common.Hippo
)

func Vec2ToScreen(v common.Vector2) rl.Vector2 {
	return rl.Vector2{X: v.X * float32(windowWidth), Y: v.Y * float32(windowHeight)}
}

func TilePosFromIdx(idx int) rl.Vector2 {
	return rl.Vector2{
		X: float32(idx % common.WorldWidth),
		Y: float32(idx / common.WorldHeight),
	}
}

func GetTileRect(tile common.Tile) rl.Rectangle {
	return rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  tileSize,
		Height: tileSize,
	}
}

func DrawTile(at rl.Vector2, tile common.Tile) {
	rl.DrawTextureRec(tileTexture, GetTileRect(tile), at, rl.White)
}

func DrawWorld(world common.WorldState) {
	xc := windowWidth/tileSize / 2
	yc := windowHeight/tileSize / 2

	xcc := xc + 1
	ycc := yc + 1

	for x := -xcc; x < xcc; x++ {
		for y := -ycc; y < ycc; y++ {
			tileX := me.Pos.X + x
			if tileX < 0 {
				continue
			}

			tileY := me.Pos.Y + y
			if tileY < 0 {
				continue
			}

			tile := world.GetTileAt(int(tileX), int(tileY))

			ax := tileX - xc //  + float32(xc)
			ay := tileY - yc //  + float32(yc)

			at := rl.Vector2{
				X: ax * tileSize,
				Y: ay * tileSize,
			}

			at = rl.Vector2Add(at, camera.Offset)

			log.Println(ax, ay)

			DrawTile(at, tile)
		}
	}
}

// this needs to be called AFTER the window has been initialized
func LoadTextures() {
	tileTexture = rl.LoadTexture("./assets/tile-texture.png")
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

	worldSize := common.WorldHeight * common.WorldWidth
	if len(discover.WorldState.Tiles) != worldSize {
		log.Fatalf("Expected to get %d tiles, got %d\n", worldSize, len(discover.WorldState.Tiles))
	}

	world = common.WorldStateFromProto(discover.WorldState)
	me.Pos = common.Vector2FromProto(discover.YourPos)

	camera = rl.NewCamera2D(
		rl.Vector2{
			X: windowWidth / 2,
			Y: windowHeight / 2,
		},
		Vector2ToRl(me.Pos),
		0,
		1.0)

	go ListenToServerMessages()
}

func Vector2ToRl(v common.Vector2) rl.Vector2 {
	return rl.Vector2{
		X: v.X,
		Y: v.Y,
	}
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
			e.Apply(&me)
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
