package l2cap

import (
	"bytes"
	"encoding"
	"encoding/binary"
)

type Option interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Type() uint8
	Len() uint8
	Hint() uint8
	SetHint(hint uint8)
}

func optionTypeFromTypeHint(b uint8) uint8 {
	return b & 0x7F
}

func optionHintFromTypeHint(b uint8) uint8 {
	return (b >> 7 & 0x01)
}

func marshalBinary(o Option) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, o.Len()+2))
	if err := binary.Write(buf, binary.LittleEndian, o); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unmarshalBinary(o Option, b []byte) error {
	buf := bytes.NewBuffer(b)
	if err := binary.Read(buf, binary.LittleEndian, o); err != nil {
		return err
	}
	copy(b, buf.Bytes())
	return nil
}
