package rfcomm

type frame struct {
	CommmandResponse   uint8
	Direction          uint8
	ServerChannel      uint8
	PollFinal          uint8
	ControlNumber      uint8
	Payload            []byte
	FrameCheckSequence uint8
}

func (f *frame) Marshal(b []byte) (int, error) {
	if len(b) < len(f.Payload)+4 {
		return 0, ErrInvalidArgument
	}

	// Address [5.4]
	ea := 0x01
	b[0] = ea&0x01 | f.CommmandResponse&0x01<<1 | f.Direction&0x01<<2 | f.ServerChannel&0x1F<<3

	// control field
	b[1] = f.ControlNumber | f.PollFinal&0x01<<4

	// length
	ea = 0x01
	b[2] = ea | uint8(len(f.Payload))<<1

	// payload
	copy(b[3:], f.Payload)

	// fcs [5.1.1]
	// TODO: actual calculation based on control number  as outlined:
	// 	https://books.google.com/books?id=mUyiREePWHwC&pg=PT239&lpg=PT239&dq=SABM+command+bluetooth&source=bl&ots=0AwIW9roPo&sig=Ze-LT5IK5-qLoETQLlRv8apLJR8&hl=en&sa=X&ved=0ahUKEwiw_sHHpObUAhUSyGMKHZ62AnUQ6AEIMDAC#v=onepage&q=SABM%20command%20bluetooth&f=false
	// UIH frames: Address and Control Field
	// Other frames: Address, Control Field and length
	copy(b[len(f.Payload)+3:], []byte{uint8(f.FrameCheckSequence)})

	return len(f.Payload) + 4, nil
}

func (f *frame) Unmarshal(b []byte) error {
	if len(b) < 3 {
		return ErrInvalidArgument
	}

	// Address
	f.CommmandResponse = b[0] >> 1 & 0x01
	f.Direction = b[0] >> 2 & 0x01
	f.ServerChannel = b[0] >> 3 & 0x1F

	// Control Field
	f.ControlNumber = b[1] & 0xF7
	f.PollFinal = b[1] >> 4 & 0x01

	// Length
	var length int
	ea := b[2] & 0x01
	if ea == 0x01 {
		length = int(b[2] >> 1)
	} else if len(b) < 4 {
		return ErrInvalidArgument
	} else { // LittleEndian
		length = int(b[2])>>1 | int(b[3])<<7
	}

	// TODO: Process credit if PollFile = 0x01

	// Payload
	i := 3 + ea + f.PollFinal
	if len(b) <= i+length {
		return ErrInvalidArgument
	}
	f.Payload = make([]byte, length)
	copy(f.Payload[:], b[i:i+length])

	// FCS
	// TODO: Check and through error if doesn't match
	// fcs := b[i+length]
	return nil
}
