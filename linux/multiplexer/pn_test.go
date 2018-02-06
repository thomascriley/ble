package multiplexer

import (
	"bytes"
	"testing"
)

var (
	pnTestBytes     []byte = []byte{0x83, 0x11, 0x02, 0xf0, 0x07, 0x00, 0xf0, 0x03, 0x00, 0x07}
	pnRespTestBytes []byte = []byte{0x81, 0x11, 0x02, 0xe0, 0x07, 0x00, 0x7f, 0x00, 0x00, 0x00}
)

const (
	testCommandTypePN        = 0x20
	testCRPN                 = 0x01
	testEAPN                 = 0x01
	textLenPN                = 0x08
	testDirectionPN          = 0x00
	testChannelPN            = 0x01
	testConvergencePN        = 0x0f
	testTypePN               = 0x00
	testPriorityPN           = 0x07
	testTimerPN              = 0x00
	testWindowSizePN         = 0x07
	testMaxRetransmissionsPN = 0x00
	testMaxSizePN            = 0x03F0
)

func TestPNUnmarshal(t *testing.T) {
	frm := &ParameterNegotiation{}
	if err := frm.UnmarshalBinary(pnTestBytes); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	if frm.CommandResponse != testCRPN {
		t.Fatalf("Command Response Expected: %X, Received: %X", testCRMSC, frm.CommandResponse)
	}
	if frm.ConvergenceLayer != testConvergencePN {
		t.Fatalf("ConvergenceLayer Expected: %X, Received: %X", testConvergencePN, frm.ConvergenceLayer)
	}
	if frm.FrameType != testTypePN {
		t.Fatalf("FrameType Expected: %X, Received: %X", testTypePN, frm.FrameType)
	}
	if frm.MaxRetransmissions != testMaxRetransmissionsPN {
		t.Fatalf("MaxRetransmissions Expected: %X, Received: %X", testMaxRetransmissionsPN, frm.MaxRetransmissions)
	}
	if frm.MaxSize != testMaxSizePN {
		t.Fatalf("MaxSize Expected: %X, Received: %X", testMaxSizePN, frm.MaxSize)
	}
	if frm.Priority != testPriorityPN {
		t.Fatalf("Priority Expected: %X, Received: %X", testPriorityPN, frm.Priority)
	}
	if frm.ServerChannel != testChannelPN {
		t.Fatalf("ServerChannel Expected: %X, Received: %X", testChannelPN, frm.ServerChannel)
	}
	if frm.Timer != testTimerPN {
		t.Fatalf("Timer Expected: %X, Received: %X", testTimerPN, frm.Timer)
	}
	if frm.WindowSize != testWindowSizePN {
		t.Fatalf("WindowSize Expected: %X, Received: %X", testWindowSizePN, frm.WindowSize)
	}
	t.Logf("Unmarshalled ParameterNegotiation: %X", pnTestBytes)
}

func TestPNRespUnmarshal(t *testing.T) {
	frm := &ParameterNegotiation{}
	if err := frm.UnmarshalBinary(pnRespTestBytes); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	t.Logf("Unmarshalled ParameterNegotiation: %X", pnRespTestBytes)
}

func TestPNMarshal(t *testing.T) {
	frm := &ParameterNegotiation{
		CommandResponse:    testCRPN,
		ConvergenceLayer:   testConvergencePN,
		FrameType:          testTypePN,
		MaxRetransmissions: testMaxRetransmissionsPN,
		MaxSize:            testMaxSizePN,
		Priority:           testPriorityPN,
		ServerChannel:      testChannelPN,
		Timer:              testTimerPN,
		WindowSize:         testWindowSizePN,
	}

	bs, err := frm.MarshalBinary()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(bs, pnTestBytes) {
		t.Fatalf("Exepected: %X, Received: %X", pnTestBytes, bs)
	}
}
