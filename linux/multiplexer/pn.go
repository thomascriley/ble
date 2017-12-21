package multiplexer

import (
	"encoding/binary"
)

// ParamsLength The length field in a PN message is always set to 8, and the value field contains 8 bytes
const ParameterNegotiationLength uint8 = 8

// ParameterNegotiation Before a DLC is set up using the mechanism in subclause 5.4.1, the TE and MS must agree on the parameters to be used
// for that DLC. These parameters are determined by parameter negotiation. [ETSI TS 101 369 V7.1.0, 5.4.6.3.1 ]
type ParameterNegotiation struct {

	// CommandResponse
	CommandResponse uint8

	// DLCI Data link connection for which parameters are being negotiated
	ServerChannel uint8

	// FrameTypes The type of frames used to carry information on the channel
	// UIH frames indicated by the value 0x1000
	FrameType uint8

	// ConvergenceLayer RFCOMM uses Type 1 (unstructured octet stream) 0x0000
	// in versions after 1.0b this may also be set to 0x0F to enable credit
	// based flow control
	ConvergenceLayer uint8

	// CreditBasedFlowControl overwrites values for ConvergenceLayer and FrameType
	CreditBasedFlowControl bool

	// Priority Assign a priority to the data link connection: 0 (lowest) - 63 (highest)
	Priority uint8

	// Timer In RFCOMM is the timer elapses, the connection is closed down. The
	// timers value is not negotiable but is fixed at 60s. This field is set to
	// 0 to indicate that the timer is not negotiable
	Timer uint8

	// MaxSize The maximum size of the frame
	MaxSize uint16

	// MaxRetransmissions The maximum number of retransmissions. Because the
	// Bluetooth baseband gives RFCOMM a reliable transport layer, RFCOMM will
	// not retransmit, so this value is set to zero
	MaxRetransmissions uint8

	// WindowSize The window size for error recovery mode. RFCOMM uses basic
	// mode, so these bits are not interpreted by RFCOMM
	WindowSize uint8
}

// Marshal ...
func (p *ParameterNegotiation) MarshalBinary() ([]byte, error) {

	b, err := marshal(p)
	if err != nil {
		return nil, err
	}
	i := HeaderSize

	b[i] = 0x00 | p.ServerChannel&0x1F<<1

	if p.CreditBasedFlowControl {
		b[i+1] = 0xF0
	} else {
		// first 4 bits are the FrameType, last 4 are the ConvergenceLayer
		b[i+1] = p.ConvergenceLayer<<4 | p.FrameType
	}

	// first 6 bits are Priority, last two are padding
	b[i+2] = p.Priority & 0x3F

	// Acknowledgement Timer is 8 bits
	b[i+3] = p.Timer

	// Maximum window size is 16 bits
	binary.LittleEndian.PutUint16(b[i+4:], p.MaxSize)

	// Maximum number of retransmisssions is 8 bits
	b[i+6] = p.MaxRetransmissions

	// K Error Recovery Window is the first 4 bytes, zero padded
	b[i+7] = p.WindowSize & 0x0F

	return b, nil
}

// Unmarshal ...
func (p *ParameterNegotiation) UnmarshalBinary(b []byte) error {

	err := unmarshal(p, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	p.ServerChannel = b[i] >> 1 & 0x1F

	// first 4 bits are the FrameType, last 4 are the ConvergenceLayer
	if b[i+1] == 0xF0 {
		p.CreditBasedFlowControl = true
	}
	p.ConvergenceLayer = b[i+1] >> 4 & 0x0F
	p.FrameType = b[i+1] & 0x0F

	// first 6 bits are Priority, last two are padding
	p.Priority = b[i+2] & 0x3F

	// Acknowledgement Timer is 8 bits
	p.Timer = b[i+3]

	// Maximum window size is 16 bits
	p.MaxSize = binary.LittleEndian.Uint16(b[i+4:])

	// Maximum number of retransmisssions is 8 bits
	p.MaxRetransmissions = b[i+6]

	// K Error Recovery Window is the first 4 bytes, zero padded
	p.WindowSize = b[i+7] & 0x0F

	return nil
}

func (p *ParameterNegotiation) Len() uint8  { return ParameterNegotiationLength }
func (p *ParameterNegotiation) Type() uint8 { return TypeParameterNegotiation }

func (p *ParameterNegotiation) GetCommandResponse() uint8  { return p.CommandResponse }
func (p *ParameterNegotiation) SetCommandResponse(l uint8) { p.CommandResponse = l }
