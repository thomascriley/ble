package multiplexor

const ModemStatusLength uint8 = 1

type ModemStatus struct {
	CommandResponse uint8

	Direction uint8

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
	i := HeaderSize + 1

	b[i] = EA & 0x01
	b[i] = b[i] | m.FlowControl&0x01<<1
	b[i] = b[i] | m.ReadyToCommunicate&0x01<<2
	b[i] = b[i] | m.ReadyToReceive&0x01<<3
	b[i] = b[i] | m.IncomingCall&0x01<<6
	b[i] = b[i] | m.DataValid&0x01<<7
	return b, nil
}

func (m *ModemStatus) UnmarshalBinary(b []byte) error {
	err := unmarshal(m, b)
	if err != nil {
		return err
	}
	i := HeaderSize + 1

	m.FlowControl = b[i] >> 1 & 0x01
	m.ReadyToCommunicate = b[i] >> 1 & 0x01
	m.ReadyToReceive = b[i] >> 1 & 0x01
	m.IncomingCall = b[i] >> 1 & 0x01
	m.DataValid = b[i] >> 1 & 0x01
	return nil
}

func (m *ModemStatus) Len() uint8  { return ModemStatusLength }
func (m *ModemStatus) Type() uint8 { return TypeModemStatus }

func (m *ModemStatus) GetCommandResponse() uint8 { return m.CommandResponse }
func (m *ModemStatus) GetDirection() uint8       { return m.Direction }
func (m *ModemStatus) GetServerChannel() uint8   { return m.ServerChannel }

func (m *ModemStatus) SetCommandResponse(l uint8) { m.CommandResponse = l }
func (m *ModemStatus) SetDirection(l uint8)       { m.Direction = l }
func (m *ModemStatus) SetServerChannel(l uint8)   { m.ServerChannel = l }
