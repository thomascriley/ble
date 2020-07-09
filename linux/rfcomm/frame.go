package rfcomm

import (
	"fmt"
)

type frame struct {
	CommmandResponse   uint8
	Direction          uint8
	ServerChannel      uint8
	PollFinal          uint8
	ControlNumber      uint8
	Payload            []byte
	FrameCheckSequence uint8
	Credits            uint8
}

func (f *frame) Marshal(b []byte) (int, error) {
	if len(b) < len(f.Payload)+4 {
		return 0, fmt.Errorf("the byte array is longer than the frame size %d > %d", len(b), len(f.Payload)+4)
	}

	// Address [5.4]
	var ea uint8 = 0x01
	b[0] = ea&0x01 | f.CommmandResponse&0x01<<1 | f.Direction&0x01<<2 | f.ServerChannel&0x1F<<3

	// control field
	b[1] = f.ControlNumber | f.PollFinal&0x01<<4

	// length
	ea = 0x01
	b[2] = ea | uint8(len(f.Payload))<<1

	i := 3

	// add credits for UIH frames when PF bit is set
	if f.ControlNumber == ControlNumberUIH && f.PollFinal == 0x01 {
		b[3] = f.Credits
		i++
	}

	copy(b[i:], f.Payload)
	i = i + len(f.Payload)

	// fcs [5.1.1]
	// UIH frames: Address and Control Field
	// Other frames: Address, Control Field and length
	var fcsBytes []byte
	switch f.ControlNumber {
	case ControlNumberUIH:
		fcsBytes = b[0:2]
	default:
		fcsBytes = b[0:3]
	}
	b[i] = generateFCS(fcsBytes)
	i++

	return i, nil
}

func (f *frame) Unmarshal(b []byte) error {
	if len(b) < 3 {
		return fmt.Errorf("the frame must be at least 3 bytes long (%X)", b)
	}

	// Address
	f.CommmandResponse = b[0] >> 1 & 0x01
	f.Direction = b[0] >> 2 & 0x01
	f.ServerChannel = b[0] >> 3 & 0x1F

	// Control Field
	f.ControlNumber = b[1] & 0xEF
	f.PollFinal = b[1] >> 4 & 0x01

	// Length
	var length int
	var ea = b[2] & 0x01
	if ea == 0x01 {
		length = int(b[2]) >> 1
	} else if len(b) < 4 {
		return fmt.Errorf("the frame must be at least 4 bytes long when ea==0 (%X)", b)
	} else { // LittleEndian
		length = int(b[2])>>1 | int(b[3])<<7
	}
	var i = 3 + (int(ea)+1)%2

	if f.ControlNumber == ControlNumberUIH && f.PollFinal == 0x01 {
		f.Credits = b[i]
		i = i + 1
	}

	// Payload
	if len(b) <= i+length {
		return fmt.Errorf("the frame must be > %d+%d bytes long (%X)", i, length, b)
	}
	if length > 0 {
		f.Payload = make([]byte, length)
		copy(f.Payload[:], b[i:i+length])
	}

	// Frame check sequence
	var fcsBytes []byte
	switch f.ControlNumber {
	case ControlNumberUIH:
		fcsBytes = b[0:2]
	default:
		fcsBytes = b[0:3]
	}

	f.FrameCheckSequence = b[len(b)-1]
	if fcs := generateFCS(fcsBytes); fcs != f.FrameCheckSequence {
		return fmt.Errorf("the frame check sequence does not match. Expected: %X, Received: %X", fcs, f.FrameCheckSequence)
	}
	return nil
}
