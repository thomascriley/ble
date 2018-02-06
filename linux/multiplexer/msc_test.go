package multiplexer

import (
	"bytes"
	"testing"
)

var mscTestBytes []byte = []byte{0xe3, 0x05, 0x0b, 0x8d}

const (
	testServerChannelMSC = 0x01
	testCommandTypeMSC   = 0x38
	testCRMSC            = 0x01
	testEAMSC            = 0x01
	textLenMSC           = 0x02
	testDirectionMSC     = 0x00
	testChannelMSC       = 0x01
	testFCMSC            = 0x00
	testRTCMSC           = 0x01
	testRTRMSC           = 0x01
	testICMSC            = 0x00
	testDVMSC            = 0x01
)

func TestMSCUnmarshal(t *testing.T) {
	frm := &ModemStatus{}
	if err := frm.UnmarshalBinary(mscTestBytes); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}
	if frm.CommandResponse != testCRMSC {
		t.Fatalf("Command Response Expected: %X, Received: %X", testCRMSC, frm.CommandResponse)
	}
	if frm.DataValid != testDVMSC {
		t.Fatalf("Data Valid Expected: %X, Received: %X", testDVMSC, frm.DataValid)
	}
	if frm.FlowControl != testFCMSC {
		t.Fatalf("Flow Control Expected: %X, Received: %X", testFCMSC, frm.FlowControl)
	}
	if frm.ReadyToCommunicate != testRTCMSC {
		t.Fatalf("Ready To Communicate Expected: %X, Received: %X", testRTCMSC, frm.ReadyToCommunicate)
	}
	if frm.ReadyToReceive != testRTRMSC {
		t.Fatalf("Ready To Receive Expected: %X, Received: %X", testRTRMSC, frm.ReadyToReceive)
	}
	if frm.IncomingCall != testICMSC {
		t.Fatalf("IncomingCall Expected: %X, Received: %X", testICMSC, frm.IncomingCall)
	}
	t.Logf("Unmarshalled Modem Status: %X", mscTestBytes)
}

func TestMSCMarshal(t *testing.T) {
	frm := &ModemStatus{
		ServerChannel:      testServerChannelMSC,
		CommandResponse:    testCRMSC,
		DataValid:          testDVMSC,
		FlowControl:        testFCMSC,
		ReadyToCommunicate: testRTCMSC,
		ReadyToReceive:     testRTRMSC,
		IncomingCall:       testICMSC,
	}

	bs, err := frm.MarshalBinary()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(bs, mscTestBytes) {
		t.Fatalf("Exepected: %X, Received: %X", mscTestBytes, bs)
	}
}
