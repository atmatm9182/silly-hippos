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

type HippoServer struct {
	WorldState   common.World
	Players      [MaxPlayers]PlayerState
	PlayersCount int32
}

func (hs *HippoServer) Start(port int) {
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

		if hs.PlayersCount >= MaxPlayers {
			// TODO: send an error message
			log.Fatalln("Player limit exceeded")
		}

		var lemmeIn message.LemmeIn
		err = common.ReadMessage(conn, &lemmeIn)
		if err != nil {
			log.Fatalln(err)
		}

		hs.Players[hs.PlayersCount] = PlayerState{
			Conn: conn,
			Hippo: common.Hippo{
				Position: common.Vector2[float32]{0.0, 0.0},
				Name: lemmeIn.Name,
			},
		}
		hs.PlayersCount++

		hs.SendDiscover(hs.PlayersCount - 1)

		go hs.HandleConnection(hs.PlayersCount - 1)
	}
}

// TODO: lock the world's state mutex
func (hs *HippoServer) CreateDiscoverMessage(id PlayerId) message.Discover {
	hippos := make([]*types.Hippo, 0, hs.PlayersCount)

	for i := int32(0); i < hs.PlayersCount; i++ {
		if i == id {
			continue
		}

		p := &hs.Players[i]
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

func (hs *HippoServer) SendDiscover(id PlayerId) {
	p := &hs.Players[id]
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	discover := hs.CreateDiscoverMessage(id)
	fmt.Printf("Created a discover message for player %d: %+v\n", id, discover)

	err := common.WriteMessage(p.Conn, &discover)
	if err != nil {
		log.Fatalln(err)
	}
}

func (hs *HippoServer) HandleConnection(id PlayerId) {
	p := &hs.Players[id]

	for {
		var msg message.Message
		err := common.ReadMessage(p.Conn, &msg)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Got a new message from player %d: %+v\n", id, &msg)
		os.Exit(1)

		hs.HandleMessage(&msg)
	}
}

func (hs *HippoServer) HandleMessage(msg *message.Message) {
	if moved := msg.GetMoved(); moved != nil {
		hs.UpdateHippoMoved(msg.Id, moved.Where)

		hs.NotifyMoved(msg.Id, moved)
		return
	}
}

func (hs *HippoServer) NotifyMoved(sender PlayerId, moved *message.Moved) {
	for i := range hs.Players {
		p := &hs.Players[i]

		if int32(i) == sender {
			continue
		}

		err := common.WriteMessage(p.Conn, moved)

		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (hs *HippoServer) UpdateHippoMoved(id PlayerId, where message.MoveDirection) {
	p := &hs.Players[id]
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

func main() {
	server := HippoServer{}
	server.Start(6969)
}
