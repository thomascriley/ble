package multiplexer

const NotSupportedLength uint8 = 0x01

type NotSupported struct {
	CommandResponse uint8

	NSCommandResponse uint8

	CommandType uint8
}

func (m *NotSupported) MarshalBinary() ([]byte, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}
	i := HeaderSize

	b[i] = 0x01 | m.NSCommandResponse&0x01<<1 | m.CommandType&0x3F<<2

	return b, nil
}

func (m *NotSupported) UnmarshalBinary(b []byte) error {
	err := unmarshal(m, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	m.NSCommandResponse = b[i] >> 1 & 0x01
	m.CommandType = b[i] >> 2 & 0x3F

	return nil
}

func (m *NotSupported) Len() uint8  { return NotSupportedLength }
func (m *NotSupported) Type() uint8 { return TypeNotSupported }

func (m *NotSupported) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *NotSupported) SetCommandResponse(l uint8) { m.CommandResponse = l }
