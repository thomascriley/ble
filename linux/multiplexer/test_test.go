package multiplexer

import (
	"bytes"
	"testing"
)

var testTestBytes []byte = []byte{0x23, 0x03, 0x08}

func TestTestMarshal(t *testing.T) {
	m := &Test{
		CommandResponse: 0x01,
		Data:            0x08,
	}

	bs, err := m.MarshalBinary()
	if err != nil {
		t.Fatal(err.Error())
	}

	if !bytes.Equal(bs, testTestBytes) {
		t.Fatalf("Exepected: %X, Received: %X", testTestBytes, bs)
	}
}
