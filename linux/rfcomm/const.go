package rfcomm

const (
	// ControlNumberSABM Set Asynchronous Balanced Mode
	ControlNumberSABM uint8 = 0x2F

	// ControlNumberUA Unnumbered Acknowledgement
	ControlNumberUA uint8 = 0x63

	// ControlNumberDM Disconnect Mode
	ControlNumberDM uint8 = 0x0F

	// ControlNumberDISC Disconnect
	ControlNumberDISC uint8 = 0x43

	// ControlNumberUIH Unnumbered Information with Header Check
	ControlNumberUIH uint8 = 0xEF
)

const (
	// Priority ...
	Priority uint8 = 0x07

	// MaxFrameSize ...
	MaxFrameSize uint16 = 0x03f0

	// FlowControl set to 1 when a device is unable to accept any RFCOMM frames
	FlowControl uint8 = 0x00

	// ReadyToCommunicate set to 1 when the device is ready to communicate
	ReadyToCommunicate uint8 = 0x01

	// ReadyToReceive set to 0 when the device cannot receive data and 1 when it can
	ReadyToReceive uint8 = 0x01

	// IncomingCall 1 indicates an incoming call
	IncomingCall uint8 = 0x00

	// ValidData 1 indicates that valid data is being sent
	ValidData uint8 = 0x01
)

const (
	// FrameTypeUIH In RFCOMM UIH frames indicated by the value 0b1000 are used.
	FrameTypeUIH uint8 = 0x01

	// ConvergenceLayer RFCOMM uses Type 1 (unstructured octet stream) 0x0000
	ConvergenceLayer uint8 = 0x00

	// Timer in RFCOMM, if the timer elapses, the connection is closed down. The timer’s value is not negotiable, but is fixed at 60s. This field is set to 0 to indicate that the timer is not negotiable.
	Timer uint8 = 0x00

	// MaxRetransmissions Because the Bluetooth baseband gives RFCOMM a reliable transport layer, RFCOMM will not retransmit, so this value is set to zero
	MaxRetransmissions uint8 = 0x00

	// WindowSize RFCOMM uses basic mode, so these bits are not interpreted by RFCOMM.
	WindowSize uint8 = 0x00
)
