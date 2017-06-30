package l2cap

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Option interface {
	Type() uint8
	Len() uint8
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
}

func marshal(o Option, b []byte) error {
	buf := bytes.NewBuffer(b)
	buf.Reset()
	if buf.Cap() < int(o.Len()) {
		return io.ErrShortBuffer
	}
	return binary.Write(buf, binary.LittleEndian, o)
}

func unmarshal(o Option, b []byte) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, binary.LittleEndian, o)
}
