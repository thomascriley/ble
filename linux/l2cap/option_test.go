package l2cap

import (
	"bytes"
	"testing"
)

const testMTU = 0x04

var testMTUBytes []byte = []byte{0x01, 0x02, 0x04, 0x00}

func TestMarshalMTUOption(t *testing.T) {
	m := &MTUOption{MTU: testMTU}
	b, err := m.MarshalBinary()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if !bytes.Equal(testMTUBytes, b) {
		t.Fatalf("Expected: %X, Received: %X", testMTUBytes, b)
	}
}

func TestUnMarshalMTUOption(t *testing.T) {
	m := &MTUOption{}
	if err := m.UnmarshalBinary(testMTUBytes); err != nil {
		t.Fatalf("Error: %s", err)
	}
	if m.MTU != testMTU {
		t.Fatalf("Expected: %d, Received: %d", testMTU, m.MTU)
	}
}
