package rfcomm

type frame struct {
	CommmandResponse   bool
	DLCI               bool
	ServerNum          int
	PollFinal          bool
	ControlNumber      int
	Payload            []byte
	FrameCheckSequence int
}

func (f *frame) Marshal(b []byte) (int, error) {
	if len(b) < len(f.Payload)+4 {
		return 0, ErrInvalidArgument
	}

	// Address [5.4]
	ea := true
	if ea {
		b[0] = 0x01
	} else {
		b[0] = 0x00
	}
	if f.CommmandResponse {
		b[0] = b[0] | 0x02
	}
	if f.DLCI {
		b[0] = b[0] | 0x04
	}
	b[0] = b[0] | uint8(f.ServerNum)<<3

	// control field
	b[1] = uint8(f.ControlNumber)
	if f.PollFinal {
		b[1] = b[1] | 0x10
	}

	// length
	ea = true
	if ea {
		b[2] = 0x00
	} else {
		b[2] = 0x01
	}
	b[2] = b[2] | uint8(len(f.Payload))<<1

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
	i := 0

	// Address
	if len(b) <= i {
		return ErrInvalidArgument
	}
	f.CommmandResponse = b[i]>>1&0x01 == 0x01
	f.DLCI = b[i]>>2&0x01 == 0x01
	f.ServerNum = int(b[i] >> 3)
	i++

	// Control Field
	if len(b) <= i {
		return ErrInvalidArgument
	}
	f.ControlNumber = int(b[i] & 0xF7)
	f.PollFinal = b[i]&0x08 == 0x08
	i++

	// Length
	var length int
	if len(b) <= i {
		return ErrInvalidArgument
	}
	if b[i]&0x01 == 0x01 {
		length = int(b[i] >> 1)
		i++
	} else {
		if len(b) < i {
			return ErrInvalidArgument
		}
		length = int(b[i])>>1 | int(b[i+2])<<7
		i = i + 2
	}

	// skip the credit if there is one
	if f.PollFinal {
		i++
	}

	// Payload
	if len(b) <= i+length {
		return ErrInvalidArgument
	}
	f.Payload = make([]byte, length)
	copy(f.Payload[:], b[i:i+length])
	i = i + length

	// FCS
	// TODO: Check and through error if doesn't match
	// fcs := b[i]
	return nil
}
