package common

import (
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

func ReadMessage(r io.Reader, msg proto.Message) error {
	var msgLenBytes [4]byte
	n, err := io.ReadFull(r, msgLenBytes[:])
	if err != nil {
		fmt.Printf("ERROR AT READING FULL: %s\n", err)
		return err
	}

	if n != 4 {
		return fmt.Errorf("Expected to read %d bytes, read %d instead", 4, n)
	}

	msgLen := binary.BigEndian.Uint32(msgLenBytes[:])
	msgData := make([]byte, msgLen)
	n, err = io.ReadFull(r, msgData)
	if err != nil {
		return err
	}

	if uint32(n) != msgLen {
		return fmt.Errorf("Expected to read %d bytes, read %d instead", msgLen, n)
	}

	return proto.Unmarshal(msgData, msg)
}

func EncodeMessage(msg proto.Message) (data []byte, err error) {
	var msgData []byte
	msgData, err = proto.Marshal(msg)
	if err != nil {
		return
	}

	data = make([]byte, 4+len(msgData))
	binary.BigEndian.PutUint32(data, uint32(len(msgData)))
	for i := 4; i < len(data); i++ {
		data[i] = msgData[i-4]
	}

	return
}

func WriteMessage(w io.Writer, msg proto.Message) error {
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	n, err := w.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("Expected to read %d bytes, read %d instead", len(data), n)
	}

	return nil
}
