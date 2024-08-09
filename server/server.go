package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/atmatm9182/silly-hippos/common"
	"github.com/atmatm9182/silly-hippos/common/message"
	"github.com/atmatm9182/silly-hippos/common/types"
)

const MaxPlayers = 8

// TODO: add a way to track if the player is still connected
type PlayerState struct {
	Conn  net.Conn
	Hippo common.Hippo
	Mutex sync.RWMutex
}

type PlayerId = int32

var (
	worldState   common.WorldState
	players      [MaxPlayers]PlayerState
	playersCount int32
)

func StartHippoServer(port int, worldTiles []common.Tile) {
	worldState.Tiles = worldTiles

	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)

	log.Printf("Server started on address %s\n", addr)

	if err != nil {
		log.Fatalln(err)
	}

	for {
		// TODO: check if player has disconnected previously
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Got a new connection: %s", conn)

		if playersCount >= MaxPlayers {
			// TODO: send an error message
			log.Fatalln("Player limit exceeded")
		}

		var lemmeIn message.LemmeIn
		err = common.ReadMessage(conn, &lemmeIn)
		if err != nil {
			log.Fatalln(err)
		}

		players[playersCount] = PlayerState{
			Conn: conn,
			Hippo: common.Hippo{
				Pos:  common.Vector2{0.0, 0.0},
				Name: lemmeIn.Name,
			},
		}
		playersCount++

		SendDiscover(playersCount - 1)

		go HandleConnection(playersCount - 1)
	}
}

// TODO: lock the world's state mutex
func CreateDiscoverMessage(id PlayerId) message.Discover {
	hippos := make([]*types.Hippo, 0, playersCount)

	for i := int32(0); i < playersCount; i++ {
		if i == id {
			continue
		}

		p := &players[i]
		p.Mutex.RLock()

		pos := new(types.Vector2)
		pos = p.Hippo.Pos.ToProto()

		hippos = append(hippos, &types.Hippo{
			Name: p.Hippo.Name,
			Pos:  pos,
		})

		p.Mutex.RUnlock()
	}

	return message.Discover{
		YourId:     id,
		YourPos:    &types.Vector2{X: 0.0, Y: 0.0},
		WorldState: worldState.ToProto(),
	}
}

func SendDiscover(id PlayerId) {
	p := &players[id]

	discover := CreateDiscoverMessage(id)

	err := common.WriteMessage(p.Conn, &discover)
	if err != nil {
		log.Fatalln(err)
	}
}

func HandleConnection(id PlayerId) {
	p := &players[id]

	for {
		var msg message.Message
		err := common.ReadMessage(p.Conn, &msg)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Got a new message from player %d: %+v\n", id, &msg)
		os.Exit(1)

		HandleMessage(&msg)
	}
}

func HandleMessage(msg *message.Message) {
	if moved := msg.GetMoved(); moved != nil {
		UpdateHippoMoved(msg.Id, moved.Where)

		NotifyMoved(msg.Id, moved)
		return
	}
}

func NotifyMoved(sender PlayerId, moved *message.Moved) {
	for i := range players {
		p := &players[i]

		if int32(i) == sender {
			continue
		}

		err := common.WriteMessage(p.Conn, moved)

		if err != nil {
			log.Fatalln(err)
		}
	}
}

func UpdateHippoMoved(id PlayerId, where message.MoveDirection) {
	p := &players[id]
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	h := &p.Hippo

	switch where {
	case message.MoveDirection_UP:
		h.Pos.Y += 1.0
	case message.MoveDirection_DOWN:
		h.Pos.Y -= 1.0
	case message.MoveDirection_LEFT:
		h.Pos.X -= 1.0
	case message.MoveDirection_RIGHT:
		h.Pos.X += 1.0
	}
}
