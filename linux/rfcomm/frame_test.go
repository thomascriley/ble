package rfcomm

import (
	"bytes"
	"testing"
)

var frameBytes1 []byte = []byte{0x03, 0x73, 0x01, 0xd7}
var frameBytes2 []byte = []byte{0x09, 0xff, 0x01, 0x0f, 0x5c}
var frameBytes3 []byte = []byte{0x0b, 0xff, 0x01, 0x21, 0x86}
var frameBytes4 []byte = []byte{
	0x49, 0xFF, 0xB7, 0x01, 0x01, 0x01, 0x00, 0x57, 0x01, 0x01, 0x01, 0x00,
	0x52, 0x01, 0x01, 0x00, 0x4E, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x57, 0x69,
	0x74, 0x68, 0x69, 0x6E, 0x67, 0x73, 0x20, 0x57, 0x53, 0x35, 0x30, 0x11,
	0x30, 0x30, 0x3A, 0x32, 0x34, 0x3A, 0x65, 0x34, 0x3A, 0x31, 0x34, 0x3A,
	0x65, 0x39, 0x3A, 0x64, 0x61, 0x10, 0x39, 0x38, 0x37, 0x38, 0x31, 0x33,
	0x35, 0x61, 0x62, 0x33, 0x32, 0x39, 0x35, 0x63, 0x34, 0x34, 0x00, 0xFF,
	0xFF, 0xFF, 0x08, 0x30, 0x30, 0x30, 0x38, 0x30, 0x30, 0x31, 0x33, 0x00,
	0x00, 0x00, 0x09, 0x00, 0x00, 0x03, 0x0D, 0x00, 0xFF, 0xFF, 0xFF, 0x08,
}

var payload []byte = []byte{
	0x01, 0x01, 0x00, 0x57, 0x01, 0x01, 0x01, 0x00, 0x52, 0x01, 0x01,
	0x00, 0x4E, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x57, 0x69, 0x74, 0x68, 0x69,
	0x6E, 0x67, 0x73, 0x20, 0x57, 0x53, 0x35, 0x30, 0x11, 0x30, 0x30, 0x3A,
	0x32, 0x34, 0x3A, 0x65, 0x34, 0x3A, 0x31, 0x34, 0x3A, 0x65, 0x39, 0x3A,
	0x64, 0x61, 0x10, 0x39, 0x38, 0x37, 0x38, 0x31, 0x33, 0x35, 0x61, 0x62,
	0x33, 0x32, 0x39, 0x35, 0x63, 0x34, 0x34, 0x00, 0xFF, 0xFF, 0xFF, 0x08,
	0x30, 0x30, 0x30, 0x38, 0x30, 0x30, 0x31, 0x33, 0x00, 0x00, 0x00, 0x09,
	0x00, 0x00, 0x03, 0x0D, 0x00, 0xFF, 0xFF, 0xFF,
}

func TestFrameUnmarshal(t *testing.T) {
	frm := &frame{}
	if err := frm.Unmarshal(frameBytes1); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	if frm.CommmandResponse != 0x01 {
		t.Fatalf("Command Response Expected: %X, Received: %X", 0x01, frm.CommmandResponse)
	}
	if frm.ControlNumber != ControlNumberUA {
		t.Fatalf("ControlNumber Expected: %X, Received: %X", ControlNumberUA, frm.ControlNumber)
	}
	if len(frm.Payload) != 0 {
		t.Fatalf("ControlNumber Expected: %X, Received: %X", 0, len(frm.Payload))
	}
	if frm.Direction != 0x00 {
		t.Fatalf("Direction Expected: %X, Received: %X", 0x00, frm.Direction)
	}
	if frm.ServerChannel != 0x00 {
		t.Fatalf("Server Channel Expected: %X, Received: %X", 0x00, frm.ServerChannel)
	}
	if frm.PollFinal != 0x01 {
		t.Fatalf("Poll Final Expected: %X, Received: %X", 0x01, frm.PollFinal)
	}
	if frm.FrameCheckSequence != 0xD7 {
		t.Fatalf("Frame Check Sequence Expected: %X, Received: %X", 0xD7, frm.FrameCheckSequence)
	}
}

func TestFrameCreditUnmarshal(t *testing.T) {
	frm := &frame{}
	if err := frm.Unmarshal(frameBytes2); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	if frm.Credits != 15 {
		t.Fatalf("Credits Expected: 15, Received: %d", frm.Credits)
	}
}

func TestFramePayloadUnmarshal(t *testing.T) {
	frm := &frame{}
	if err := frm.Unmarshal(frameBytes4); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	if !bytes.Equal(frm.Payload, payload) {
		t.Fatalf("Payload\nExpected: %X\nReceived: %X", payload, frm.Payload)
	}
}

func TestFrameCreditMarshal(t *testing.T) {
	frm := &frame{
		Direction:        0x00,
		ServerChannel:    0x01,
		ControlNumber:    ControlNumberUIH,
		CommmandResponse: 0x01,
		PollFinal:        0x01,
		Credits:          33,
		Payload:          []byte{},
	}
	bs := make([]byte, 4096)
	n, err := frm.Marshal(bs)
	if err != nil {
		t.Fatalf("Error marshalling: %s", err)
	}
	if !bytes.Equal(bs[:n], frameBytes3) {
		t.Fatalf("Exepected: %X, Received: %X", frameBytes3, bs[:n])
	}
}
