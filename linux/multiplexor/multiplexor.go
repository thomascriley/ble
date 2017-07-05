package multiplexor

import (
	"encoding"
	"fmt"
)

type Multiplexor interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Type() uint8
	Len() uint8

	GetCommandResponse() uint8
	GetDirection() uint8
	GetServerChannel() uint8

	SetCommandResponse(l uint8)
	SetDirection(l uint8)
	SetServerChannel(l uint8)
}

func MarshalBinary(p Multiplexor) ([]byte, error) {
	return p.MarshalBinary()
}

func UnmarshalBinary(data []byte) (Multiplexor, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("The byte buffer must be at least %d bytes long", HeaderSize)
	}
	var p Multiplexor
	switch data[0] >> 3 {
	case TypeFlowControlOff:
		//p = &FlowControlOff{}
	case TypeFlowControlOn:
		//p = &FlowControlOn{}
	case TypeModemStatus:
		p = &ModemStatus{}
	case TypeNotSupported:
		//p = &NotSupported{}
	case TypeParameterNegotiation:
		p = &ParameterNegotiation{}
	case TypeRemoteLineStatus:
		//p = &RemoteLineStatus{}
	case TypeRemotePortNegotiation:
		//p = &RemotePortNegotiation{}
	case TypeTest:
		//p = &Test{}
	default:
		return nil, fmt.Errorf("Unknown multiplexor type %X", data[0]>>3)
	}
	return p, p.UnmarshalBinary(data)
}

func marshal(p Multiplexor) ([]byte, error) {
	b := make([]byte, HeaderSize+int(p.Len()))

	// Parameter Negotiation Type
	b[0] = EA&0x01 | p.GetCommandResponse()&0x01<<1 | p.Type()<<3

	// Length of the parameter values
	b[1] = EA&0x01 | p.Len()<<1

	// first 6 bits are DLCI, last two are padding
	// DLCI is composed of a direction bit and a 5 bit server number 1 - 30
	b[2] = p.GetDirection()&0x01 | p.GetServerChannel()&0x1F<<1

	return b, nil
}

func unmarshal(p Multiplexor, b []byte) error {
	if len(b) < HeaderSize+1 {
		return fmt.Errorf("The byte buffer must be at least %d bytes long", 3)
	}

	if p.Type() != b[0]>>3 {
		return fmt.Errorf("The multiplexor types do ")
	}
	p.SetCommandResponse(b[0] >> 1 & 0x01)
	p.SetDirection(b[2] & 0x01)
	p.SetServerChannel(b[2] >> 1 & 0x1F)

	length := b[1] >> 1
	if int(length)+HeaderSize > len(b) {
		return fmt.Errorf("The byte buffer must be at least %d bytes long", int(length)+HeaderSize)
	}
	return nil
}
