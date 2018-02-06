package rfcomm

import "testing"

type fcsTest struct {
	bs  []byte
	fcs uint8
}

var fcsTestBytes []fcsTest = []fcsTest{
	fcsTest{bs: []byte{0x03, 0x3f, 0x01}, fcs: 0x1c},
	fcsTest{bs: []byte{0x03, 0x73, 0x01}, fcs: 0xd7},
}

func TestFCS(t *testing.T) {

	for _, b := range fcsTestBytes {
		fcs := generateFCS(b.bs)
		if fcs != b.fcs {
			t.Fatalf("Expected: %X, Received: %X", b.fcs, fcs)
		}
	}
}
