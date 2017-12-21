package multiplexer

const RemoteLineStatusLength uint8 = 2

type RemoteLineStatus struct {
	CommandResponse uint8

	ServerChannel uint8

	LineStatus uint8
}

func (m *RemoteLineStatus) MarshalBinary() ([]byte, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}
	i := HeaderSize

	b[i] = 0x01 | 0x01<<1 | m.ServerChannel&0x3F<<2
	b[i+1] = m.LineStatus & 0x0F
	return b, nil
}

func (m *RemoteLineStatus) UnmarshalBinary(b []byte) error {
	err := unmarshal(m, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	m.ServerChannel = b[i] >> 2 & 0x3F
	m.LineStatus = b[i+1] & 0x0F
	return nil
}

func (m *RemoteLineStatus) Len() uint8  { return RemoteLineStatusLength }
func (m *RemoteLineStatus) Type() uint8 { return TypeRemoteLineStatus }

func (m *RemoteLineStatus) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *RemoteLineStatus) SetCommandResponse(l uint8) { m.CommandResponse = l }
