package ble

// InqHandler handles BR/EDR inquiries
type InqHandler func(i Inquiry)

// Inquiry BR/EDR inquiry
type Inquiry interface {
	Address() Addr

	PageScanRepetitionMode() uint8

	ClassOfDevice() [3]byte

	ClockOffset() uint16

	RSSI() int

	ExtendedInquiryResponse() [240]byte
}
