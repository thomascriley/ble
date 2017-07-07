package l2cap

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"io"
)

type Option interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Type() uint8
	Len() uint8
	Hint() uint8
	SetHint(hint uint8)
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
