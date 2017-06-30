package l2cap

// MTUOption [Vol 3, Part A, 5.1 ]
const MTUOptionType uint8 = 0x01

type MTUOption struct {
	MTU uint16
}

func (o *MTUOption) Type() uint8 { return 0x01 }
func (o *MTUOption) Len() uint8  { return 0x02 }
func (o *MTUOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *MTUOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// FlushTimeoutOption [Vol 3, Part A, 5.2 ]
const FlushTimeoutOptionType uint8 = 0x02

type FlushTimeoutOption struct {
	FlushTimeout uint16
}

func (o *FlushTimeoutOption) Type() uint8 { return 0x02 }
func (o *FlushTimeoutOption) Len() uint8  { return 0x02 }
func (o *FlushTimeoutOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *FlushTimeoutOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// QoSOption [Vol 3, Part A, 5.3 ]
const QoSOptionType uint8 = 0x03

type QoSOption struct {
	Flags           uint8
	ServiceType     uint8
	TokenBucketSize uint32
	PeakBandwidth   uint32
	Latency         uint32
	DelayVariation  uint32
}

func (o *QoSOption) Type() uint8 { return 0x03 }
func (o *QoSOption) Len() uint8  { return 0x02 }
func (o *QoSOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *QoSOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// RetransmissionAndFlowControlOption [Vol 3, Part A, 5.4 ]
const RetransmissionAndFlowControlOptionType uint8 = 0x04

type RetransmissionAndFlowControlOption struct {
	Mode                  uint8
	TxWindowSize          uint8
	MaxTransmit           uint8
	RetransmissionTimeout uint16
	MonitorTimeout        uint16
	MaximumPDUSize        uint16
}

func (o *RetransmissionAndFlowControlOption) Type() uint8 { return 0x04 }
func (o *RetransmissionAndFlowControlOption) Len() uint8  { return 0x09 }
func (o *RetransmissionAndFlowControlOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *RetransmissionAndFlowControlOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// FrameCheckSequenceOption [Vol 3, Part A, 5.5 ]
const FrameCheckSequenceOptionType uint8 = 0x05

type FrameCheckSequenceOption struct {
	FCSType uint8
}

func (o *FrameCheckSequenceOption) Type() uint8 { return 0x05 }
func (o *FrameCheckSequenceOption) Len() uint8  { return 0x01 }
func (o *FrameCheckSequenceOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *FrameCheckSequenceOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// ExtendedFlowSpecificationOption [Vol 3, Part A, 5.6 ]
const ExtendedFlowSpecificationOptionType uint8 = 0x06

type ExtendedFlowSpecificationOption struct {
	Identifier          uint8
	ServiceType         uint8
	MaximumSDUSize      uint16
	SDUInterarrivalTime uint32
	AccessLatency       uint32
	FlushTimeout        uint32
}

func (o *ExtendedFlowSpecificationOption) Type() uint8 { return 0x06 }
func (o *ExtendedFlowSpecificationOption) Len() uint8  { return 0x16 }
func (o *ExtendedFlowSpecificationOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *ExtendedFlowSpecificationOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}

// ExtendedWindowSizeOption [Vol 3, Part A, 5.7 ]
const ExtendedWindowSizeOptionType uint8 = 0x07

type ExtendedWindowSizeOption struct {
	MaxWindowSize uint16
}

func (o *ExtendedWindowSizeOption) Type() uint8 { return 0x07 }
func (o *ExtendedWindowSizeOption) Len() uint8  { return 0x02 }
func (o *ExtendedWindowSizeOption) Marshal() ([]byte, error) {
	b := make([]byte, 0, o.Len())
	if err := marshal(o, b); err != nil {
		return nil, err
	}
	return b, nil
}
func (o *ExtendedWindowSizeOption) Unmarshal(b []byte) error {
	if err := unmarshal(o, b); err != nil {
		return err
	}
	return nil
}
