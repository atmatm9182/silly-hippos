package main

import (
	"github.com/atmatm9182/silly-hippos/common/message"
	"google.golang.org/protobuf/proto"
)

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

