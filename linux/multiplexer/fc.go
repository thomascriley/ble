package multiplexer

const FlowControlLength uint8 = 0

// FlowControlOff The flow control command is used to handle the aggregate flow. When either entity is not able to receive information it
// transmits the FCoff command. The opposite entity is not allowed to transmit frames except on the control channel
// (DLC=0). [ETSI TS 101 369 V7.1.0, 5.4.6.3.6 ]
type FlowControlOff struct {
	CommandResponse uint8
}

func (m *FlowControlOff) MarshalBinary() ([]byte, error) { return marshal(m) }
func (m *FlowControlOff) UnmarshalBinary(b []byte) error { return unmarshal(m, b) }

func (m *FlowControlOff) Len() uint8  { return FlowControlLength }
func (m *FlowControlOff) Type() uint8 { return TypeFlowControlOff }

func (m *FlowControlOff) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *FlowControlOff) SetCommandResponse(l uint8) { m.CommandResponse = l }

// FlowControlOn The flow control command is used to handle the aggregate flow. When either entity is able to receive new information it
// transmits this command. [ETSI TS 101 369 V7.1.0, 5.4.6.3.5 ]
type FlowControlOn struct {
	CommandResponse uint8
}

func (m *FlowControlOn) MarshalBinary() ([]byte, error) { return marshal(m) }
func (m *FlowControlOn) UnmarshalBinary(b []byte) error { return unmarshal(m, b) }

func (m *FlowControlOn) Len() uint8  { return FlowControlLength }
func (m *FlowControlOn) Type() uint8 { return TypeFlowControlOn }

func (m *FlowControlOn) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *FlowControlOn) SetCommandResponse(l uint8) { m.CommandResponse = l }
