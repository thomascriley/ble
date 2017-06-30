package ble

// InqHandler handles BR/EDR inquiries
type InqHandler func(i Inquiry)

// Inquiry BR/EDR inquiry
type Inquiry interface {
	Address() Addr

	PageScanRepetitionMode() int

	ClassOfDevice() [3]byte

	ClockOffset() int

	RSSI() int

	ExtendedInquiryResponse() [240]byte
}
