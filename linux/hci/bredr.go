package hci

import (
	"github.com/currantlabs/ble"
)

// SetAdvHandler ...
func (h *HCI) SetInqHandler(ah ble.InqHandler) error {
	h.inqHandler = ah
	return nil
}

// Inquire starts inquiring for BR/EDR devices.
func (h *HCI) Inquire(length int, numResponses int) error {
	if numResponses > 255 || numResponses < 0 {
		numResponses = 0x00 // unlimited
	}
	if length > 0x30 || length <= 0 {
		length = 0x30 // max length (61.44 sec)
	}

	// General/Unlimited Inquiry Access Code (GIAC)
	// see: https://www.bluetooth.com/specifications/assigned-numbers/baseband
	h.params.inquiry.LAP = [3]byte{0x9e, 0x8b, 0x33}
	h.params.inquiry.InquiryLength = uint8(length)
	h.params.inquiry.NumResponses = uint8(numResponses)
	return h.Send(&h.params.inquiry, nil)
}

// StopInquiry stops inquiring for BR/EDR devices by sending an InquiryCancel command
func (h *HCI) StopInquiry() error {
	return h.Send(&h.params.inquiryCancel, nil)
}
