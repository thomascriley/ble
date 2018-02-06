package multiplexer

import (
	"encoding/binary"
)

const (
	RemotePortNegotiationSetupLength   uint8 = 8
	RemotePortNegotiationRequestLength uint8 = 1
)

const (
	MaskBaudRate uint8 = iota
	MaskDataBits
	MaskStopBits
	MaskParity
	MaskParityType
	MaskXONChar
	MaskXOFFChar
	MaskInputXONOFF
	MaskOutputXONOFF
	MaskInputRTR
	MaskOutputRTR
	MaskInputRTC
	MaskOutputRTC
)

// RemotePortNegotiation Remote Port Negotiation used to set communication settings
// at the remote end of the data link connection
type RemotePortNegotiation struct {

	// CommandResponse
	CommandResponse uint8

	// DLCI is composed of a direction bit and a 5 bit server number 1 - 30
	ServerChannel uint8

	baudRate uint8

	dataBits uint8

	stopBits uint8

	parity uint8

	parityType uint8

	flowControl uint8

	xON uint8

	xOFF uint8

	mask uint16

	Setup bool
}

// Marshal ...
func (p *RemotePortNegotiation) MarshalBinary() ([]byte, error) {

	b, err := marshal(p)
	if err != nil {
		return nil, err
	}
	i := HeaderSize

	b[i] = 0x01 | 0x01<<1 | p.ServerChannel&0x1F<<2

	if !p.Setup {
		return b, nil
	}

	p.baudRate = b[i+1]
	p.dataBits = b[i+2] & 0x03
	p.stopBits = b[i+2] >> 2 & 0x01
	p.parity = b[i+2] >> 3 & 0x01
	p.parityType = b[i+2] >> 4 & 0x03
	p.flowControl = b[i+3] & 0x3F
	p.xON = b[i+4]
	p.xOFF = b[i+5]
	p.mask = binary.LittleEndian.Uint16(b[i+6:])
	return b, nil
}

// Unmarshal ...
func (p *RemotePortNegotiation) UnmarshalBinary(b []byte) error {

	err := unmarshal(p, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	p.ServerChannel = b[i] >> 2 & 0x3F

	if b[1]>>1 == 0x01 {
		p.Setup = false
		return nil
	}
	p.Setup = true

	b[i+1] = p.baudRate
	b[i+2] = p.dataBits & 0x03
	b[i+2] = p.stopBits & 0x01 << 2
	b[i+2] = p.parity & 0x01 << 3
	b[i+2] = p.parityType & 0x03 << 4
	b[i+3] = p.flowControl & 0x3F
	b[i+4] = p.xON
	b[i+5] = p.xOFF
	binary.LittleEndian.PutUint16(b[i+6:], p.mask)
	return nil
}

func (p *RemotePortNegotiation) Change(mask uint8, value uint8) {
	switch mask {
	case MaskBaudRate:
		p.baudRate = value
	case MaskDataBits:
		p.dataBits = value
	case MaskStopBits:
		p.stopBits = value
	case MaskParity:
		p.parity = value
	case MaskParityType:
		p.parityType = value
	case MaskXONChar:
		p.xON = value
	case MaskXOFFChar:
		p.xOFF = value
	}
	p.mask = p.mask | 0x01<<mask
}

func (p *RemotePortNegotiation) IsChanged(mask uint8) bool {
	return p.mask>>mask&0x01 == 0x01
}

func (p *RemotePortNegotiation) Value(mask uint8) uint8 {
	switch mask {
	case MaskBaudRate:
		return p.baudRate
	case MaskDataBits:
		return p.dataBits
	case MaskStopBits:
		return p.stopBits
	case MaskParity:
		return p.parity
	case MaskParityType:
		return p.parityType
	case MaskXONChar:
		return p.xON
	case MaskXOFFChar:
		return p.xOFF
	}
	return uint8(p.mask >> mask & 0x01)
}

func (p *RemotePortNegotiation) Len() uint8 {
	if p.Setup {
		return RemotePortNegotiationSetupLength
	} else {
		return RemotePortNegotiationRequestLength
	}
}

func (p *RemotePortNegotiation) Type() uint8 { return TypeRemotePortNegotiation }

func (p *RemotePortNegotiation) GetCommandResponse() uint8  { return p.CommandResponse }
func (p *RemotePortNegotiation) SetCommandResponse(l uint8) { p.CommandResponse = l }
