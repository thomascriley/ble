package hci

import (
	"net"

	"github.com/thomascriley/ble"
)

// InquiryEvent interface that wraps methods common to all inquiry events
type InquiryEvent interface {
	NumResponses() uint8
	BDADDR(i int) [6]byte
	PageScanRepetitionMode(i int) uint8
	ClassOfDevice(i int) [3]byte
	ClockOffset(i int) uint16
	RSSI(i int) uint8
	ExtendedInquiryResponse() [240]byte
}

func newInquiry(e InquiryEvent, i int) *Inquiry {
	return &Inquiry{e: e, i: i}
}

// Inquiry implements ble.Inquiry and other functions that are only
// available on Linux.
type Inquiry struct {
	e  InquiryEvent
	i  int
	sr *Inquiry
}

// PageScanRepetitionMode returns the mode used for page scans
func (i *Inquiry) PageScanRepetitionMode() uint8 {
	return i.e.PageScanRepetitionMode(i.i)
}

// ClockOffset returns the difference in time between the host and client
func (i *Inquiry) ClockOffset() uint16 {
	return i.e.ClockOffset(i.i)
}

// ClassOfDevice returns a bit mask as defined here: https://www.bluetooth.com/specifications/assigned-numbers/baseband
func (i *Inquiry) ClassOfDevice() [3]byte {
	return i.e.ClassOfDevice(i.i)
}

// RSSI returns RSSI signal strength.
func (i *Inquiry) RSSI() int {
	return int(i.e.RSSI(i.i))
}

// Address returns the address of the remote peripheral.
func (i *Inquiry) Address() ble.Addr {
	b := i.e.BDADDR(i.i)
	return net.HardwareAddr([]byte{b[5], b[4], b[3], b[2], b[1], b[0]})
}

func (i *Inquiry) ExtendedInquiryResponse() [240]byte {
	return i.e.ExtendedInquiryResponse()
}
