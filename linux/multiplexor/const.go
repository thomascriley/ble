package multiplexor

const EA uint8 = 0x01

const HeaderSize int = 2

const (
	// TypeParameterNegotiation
	TypeParameterNegotiation uint8 = 0x20

	// TypeTest used to check the RFCOMM connection
	TypeTest uint8 = 0x08

	// TypeFlowControlOn flow control mechanism which applies to all channels between two RFCOMM entities
	TypeFlowControlOn uint8 = 0x28

	// TypeFlowControlOff no flow control mechanism which applies to all channels between two RFCOMM entities
	TypeFlowControlOff uint8 = 0x18

	// TypeModemStatus a flow control mechanism which applies to a single channel
	TypeModemStatus uint8 = 0x38

	// TypeRemotePortNegotiation The Remote Port Negotiation (RPN) command is used to set communication settings at the remote end of a data link connection.
	TypeRemotePortNegotiation uint8 = 0x24

	// TypeRemoteLineStatus Sent when device needs to tell the other end of the link about an error
	TypeRemoteLineStatus uint8 = 0x14

	// TypeNotSupported sent whenever a device receives a command it does not support.
	TypeNotSupported uint8 = 0x04
)
