package multiplexer

const TestLength uint8 = 1

// Test The test command is used to test the connection between MS and the TE. The length byte describes the number of
// values bytes, which are used as a verification pattern. The opposite entity shall respond with exactly the same value
// bytes. [ETSI TS 101 369 V7.1.0, 5.4.6.3.4 ]
type Test struct {
	CommandResponse uint8

	Data uint8
}

func (m *Test) MarshalBinary() ([]byte, error) {
	b, err := marshal(m)
	if err != nil {
		return nil, err
	}
	i := HeaderSize
	b[i] = m.Data
	return b, nil
}

func (m *Test) UnmarshalBinary(b []byte) error {
	err := unmarshal(m, b)
	if err != nil {
		return err
	}
	i := HeaderSize

	m.Data = b[i]
	return nil
}

func (m *Test) Len() uint8  { return TestLength }
func (m *Test) Type() uint8 { return TypeTest }

func (m *Test) GetCommandResponse() uint8  { return m.CommandResponse }
func (m *Test) SetCommandResponse(l uint8) { m.CommandResponse = l }
