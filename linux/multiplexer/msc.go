package multiplexer

const ModemStatusLength uint8 = 1

// ModemStatus It is desired to convey virtual V.24 control signals to a data stream, this is done by sending the MSC command. The
// MSC command has one mandatory control signal byte and an optional break signal byte. This command is only relevant
// when the basic option is chosen. [ETSI TS 101 369 V7.1.0, 5.4.6.3.7 ]
type ModemStatus struct {
	CommandResponse uint8

	ServerChannel uint8

	FlowControl uint8

	ReadyToCommunicate uint8

	ReadyToReceive uint8

	IncomingCall uint8

	DataValid uint8
}

func (m *ModemStatus) MarshalBinary() ([]byte, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}
	i := HeaderSize

	b[i] = 0x01 | 0x01<<1 | 0x01<<2 | m.ServerChannel&0x1F<<3
	b[i+1] = 0x01
	b[i+1] = b[i+1] | m.FlowControl&0x01<<1
	b[i+1] = b[i+1] | m.ReadyToCommunicate&0x01<<2
	b[i+1] = b[i+1] | m.ReadyToReceive&0x01<<3
	b[i+1] = b[i+1] | m.IncomingCall&0x01<<6
	b[i+1] = b[i+1] | m.DataValid&0x01<<7
	return b, nil
}

func (m *ModemStatus) UnmarshalBinary(b []byte) error {
	err := unmarshal(m, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	m.ServerChannel = b[i] >> 3 & 0x1F
	m.FlowControl = b[i+1] >> 1 & 0x01
	m.ReadyToCommunicate = b[i+1] >> 1 & 0x01
	m.ReadyToReceive = b[i+1] >> 1 & 0x01
	m.IncomingCall = b[i+1] >> 1 & 0x01
	m.DataValid = b[i+1] >> 1 & 0x01
	return nil
}

func (m *ModemStatus) Len() uint8                 { return ModemStatusLength }
func (m *ModemStatus) Type() uint8                { return TypeModemStatus }
func (m *ModemStatus) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *ModemStatus) SetCommandResponse(l uint8) { m.CommandResponse = l }
