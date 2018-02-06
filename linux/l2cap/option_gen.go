package l2cap

// MTUType is the option type of MTU configuration option.
const MTUOptionType = 0x01

// MTUOption implements MTU (0x01) [Vol 3, Part A, 5.1].
type MTUOption struct {
	TypeHint uint8
	Length   uint8
	MTU      uint16
}

// Type ...
func (o *MTUOption) Type() uint8 { return 0x01 }

// Len returns the length of the object payload in bytes
func (o *MTUOption) Len() uint8 { return 0x02 }

// Hint returns if a bad value should cause the connection to fail
func (o *MTUOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *MTUOption) SetHint(hint uint8) { o.TypeHint = 0x01&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *MTUOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x01 & 0x7F)
	o.Length = 0x02
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *MTUOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// FlushTimeoutType is the option type of Flush Timeout configuration option.
const FlushTimeoutOptionType = 0x02

// FlushTimeoutOption implements Flush Timeout (0x02) [Vol 3, Part A, 5.2].
type FlushTimeoutOption struct {
	TypeHint     uint8
	Length       uint8
	FlushTimeout uint16
}

// Type ...
func (o *FlushTimeoutOption) Type() uint8 { return 0x02 }

// Len returns the length of the object payload in bytes
func (o *FlushTimeoutOption) Len() uint8 { return 0x02 }

// Hint returns if a bad value should cause the connection to fail
func (o *FlushTimeoutOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *FlushTimeoutOption) SetHint(hint uint8) { o.TypeHint = 0x02&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *FlushTimeoutOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x02 & 0x7F)
	o.Length = 0x02
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *FlushTimeoutOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// QoSType is the option type of QoS configuration option.
const QoSOptionType = 0x03

// QoSOption implements QoS (0x03) [Vol 3, Part A, 5.3].
type QoSOption struct {
	TypeHint        uint8
	Length          uint8
	Flags           uint8
	ServiceType     uint8
	TokenBucketSize uint32
	PeakBandwidth   uint32
	Latency         uint32
	DelayVariation  uint32
}

// Type ...
func (o *QoSOption) Type() uint8 { return 0x03 }

// Len returns the length of the object payload in bytes
func (o *QoSOption) Len() uint8 { return 0x13 }

// Hint returns if a bad value should cause the connection to fail
func (o *QoSOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *QoSOption) SetHint(hint uint8) { o.TypeHint = 0x03&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *QoSOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x03 & 0x7F)
	o.Length = 0x13
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *QoSOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// RetransmissionAndFlowControlType is the option type of Retransmission And Flow Control configuration option.
const RetransmissionAndFlowControlOptionType = 0x04

// RetransmissionAndFlowControlOption implements Retransmission And Flow Control (0x04) [Vol 3, Part A, 5.4].
type RetransmissionAndFlowControlOption struct {
	TypeHint              uint8
	Length                uint8
	Mode                  uint8
	TxWindowSize          uint8
	MaxTransmit           uint8
	RetransmissionTimeout uint16
	MonitorTimeout        uint16
	MaximumPDUSize        uint16
}

// Type ...
func (o *RetransmissionAndFlowControlOption) Type() uint8 { return 0x04 }

// Len returns the length of the object payload in bytes
func (o *RetransmissionAndFlowControlOption) Len() uint8 { return 0x09 }

// Hint returns if a bad value should cause the connection to fail
func (o *RetransmissionAndFlowControlOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *RetransmissionAndFlowControlOption) SetHint(hint uint8) {
	o.TypeHint = 0x04&0x7F | hint&0x01<<7
}

// Marshal serializes the command parameters into binary form.
func (o *RetransmissionAndFlowControlOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x04 & 0x7F)
	o.Length = 0x09
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *RetransmissionAndFlowControlOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// FrameCheckSequenceType is the option type of Frame Check Sequence configuration option.
const FrameCheckSequenceOptionType = 0x05

// FrameCheckSequenceOption implements Frame Check Sequence (0x05) [Vol 3, Part A, 5.5].
type FrameCheckSequenceOption struct {
	TypeHint uint8
	Length   uint8
	FCSType  uint8
}

// Type ...
func (o *FrameCheckSequenceOption) Type() uint8 { return 0x05 }

// Len returns the length of the object payload in bytes
func (o *FrameCheckSequenceOption) Len() uint8 { return 0x01 }

// Hint returns if a bad value should cause the connection to fail
func (o *FrameCheckSequenceOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *FrameCheckSequenceOption) SetHint(hint uint8) { o.TypeHint = 0x05&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *FrameCheckSequenceOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x05 & 0x7F)
	o.Length = 0x01
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *FrameCheckSequenceOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// ExtendedFlowSpecificationType is the option type of Extended Flow Specification configuration option.
const ExtendedFlowSpecificationOptionType = 0x06

// ExtendedFlowSpecificationOption implements Extended Flow Specification (0x06) [Vol 3, Part A, 5.6].
type ExtendedFlowSpecificationOption struct {
	TypeHint            uint8
	Length              uint8
	Identifier          uint8
	ServiceType         uint8
	MaximumSDUSize      uint8
	SDUInterarrivalTime uint16
	AccessLatency       uint16
	FlushTimeout        uint16
}

// Type ...
func (o *ExtendedFlowSpecificationOption) Type() uint8 { return 0x06 }

// Len returns the length of the object payload in bytes
func (o *ExtendedFlowSpecificationOption) Len() uint8 { return 0x10 }

// Hint returns if a bad value should cause the connection to fail
func (o *ExtendedFlowSpecificationOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *ExtendedFlowSpecificationOption) SetHint(hint uint8) { o.TypeHint = 0x06&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *ExtendedFlowSpecificationOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x06 & 0x7F)
	o.Length = 0x10
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *ExtendedFlowSpecificationOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}

// ExtendedWindowSizeType is the option type of Extended Window Size configuration option.
const ExtendedWindowSizeOptionType = 0x07

// ExtendedWindowSizeOption implements Extended Window Size (0x07) [Vol 3, Part A, 5.6].
type ExtendedWindowSizeOption struct {
	TypeHint      uint8
	Length        uint8
	MaxWindowSize uint8
}

// Type ...
func (o *ExtendedWindowSizeOption) Type() uint8 { return 0x07 }

// Len returns the length of the object payload in bytes
func (o *ExtendedWindowSizeOption) Len() uint8 { return 0x02 }

// Hint returns if a bad value should cause the connection to fail
func (o *ExtendedWindowSizeOption) Hint() uint8 { return optionHintFromTypeHint(o.TypeHint) }

// SetHint sets the Hint value based off of the MSB of the Type
func (o *ExtendedWindowSizeOption) SetHint(hint uint8) { o.TypeHint = 0x07&0x7F | hint&0x01<<7 }

// Marshal serializes the command parameters into binary form.
func (o *ExtendedWindowSizeOption) MarshalBinary() ([]byte, error) {
	o.TypeHint = o.TypeHint | (0x07 & 0x7F)
	o.Length = 0x02
	return marshalBinary(o)
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (o *ExtendedWindowSizeOption) UnmarshalBinary(b []byte) error {
	return unmarshalBinary(o, b)
}
