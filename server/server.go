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
	Mutex sync.Mutex
}

type PlayerId = int32

var (
	worldState common.World
	players      [MaxPlayers]PlayerState
	playersCount int32
)

func StartHippoServer(port int) {
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
				Position: common.Vector2[float32]{0.0, 0.0},
				Name:     lemmeIn.Name,
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
		p.Mutex.Lock()

		pos := new(types.Vector2)
		*pos = common.Vector2ToProto(p.Hippo.Position)

		hippos = append(hippos, &types.Hippo{
			Id:   i,
			Name: p.Hippo.Name,
			Pos:  pos,
		})

		p.Mutex.Unlock()
	}

	return message.Discover{
		Hippos: hippos,
	}
}

func SendDiscover(id PlayerId) {
	p := &players[id]
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	discover := CreateDiscoverMessage(id)
	fmt.Printf("Created a discover message for player %d: %+v\n", id, discover)

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

func  UpdateHippoMoved(id PlayerId, where message.MoveDirection) {
	p := &players[id]
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	h := &p.Hippo

	switch where {
	case message.MoveDirection_UP:
		h.Position.Y += 1.0
	case message.MoveDirection_DOWN:
		h.Position.Y -= 1.0
	case message.MoveDirection_LEFT:
		h.Position.X -= 1.0
	case message.MoveDirection_RIGHT:
		h.Position.X += 1.0
	}
}
