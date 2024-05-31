package smp

import (
	"crypto/aes"

	"github.com/thomascriley/ble/cmac"
)

var (
	salt  = [16]byte{0x6C, 0x88, 0x83, 0x91, 0xAA, 0xF5, 0xA5, 0x38, 0x60, 0x37, 0x0B, 0xDB, 0x5A, 0x60, 0x83, 0xBE}
	keyID = [4]byte{0x62, 0x74, 0x6c, 0x65}
)

// c1 During the LE legacy pairing process confirm values are exchanged. This
// confirm value generation function c1 is used to generate the confirm values.
// [ Vol 3, Part H 2.2.3 ]
func C1(k, r [16]byte, pres *PairingResponse, preq *PairingRequest, iat byte, ia [6]byte, rat byte, ra [6]byte) ([16]byte, error) {

	// pres, preq, rat’ and iat’ are concatenated to generate p1 which is
	// XORed with r and used as 128-bit input parameter plaintextData to
	// security function e:
	var p1 [16]byte
	p1[0] = PairingRequestCode
	copy(p1[1:], pres.Marshal())
	p1[7] = PairingResponseCode
	copy(p1[8:], preq.Marshal())
	p1[14] = rat
	p1[15] = iat

	for i := 0; i < 16; i++ {
		p1[i] = r[i] ^ p1[i]
	}

	e1, err := e(k, p1)
	if err != nil {
		return [16]byte{}, err
	}

	// ra is concatenated with ia and padding to generate p2 which is XORed
	// with the result of the security function e using p1 as the input parameter
	// plaintextData and is then used as the 128-bit input parameter plaintextData
	// to security function e:
	var p2 [16]byte
	copy(p2[4:], ia[:])
	copy(p2[10:], ra[:])

	for i := 0; i < 16; i++ {
		p2[i] = e1[i] ^ p2[i]
	}

	// c1 (k, r, preq, pres, iat, rat, ia, ra) = e(k, e(k, r XOR p1) XOR p2)
	return e(k, p2)
}

func e(k, p [16]byte) (out [16]byte, err error) {
	blk, err := aes.NewCipher(k[:])
	if err != nil {
		return
	}
	blk.Encrypt(out[:], p[:])
	return
}

func AES_CMAC(k [16]byte, m []byte) (code [16]byte, err error) {
	blk, err := aes.NewCipher(k[:])
	if err != nil {
		return
	}

	hash, err := cmac.New(blk)
	if err != nil {
		return
	}

	sum := hash.Sum(m)
	copy(code[:], sum)
	return
}

// S1 is the key generation function s1 for LE Legacy Pairing.
// The key generation function s1 is used to generate the STK during the LE
// legacy pairing process. [ Vol 3, Part H 2.2.4]
func S1(k, r1, r2 [16]byte) ([16]byte, error) {
	var r [16]byte
	copy(r[:], r1[8:])
	copy(r[8:], r2[8:])
	return e(k, r)
}

// F4 during the LE Secure Connections pairing process, confirm values are exchanged.
// These confirm values are computed using the confirm value generation function f4.
//
// Z is zero (i.e. 8 bits of zeros) for Numeric Comparison and OOB protocol. In
// the Passkey Entry protocol, the most significant bit of Z is set equal to one
// and the least significant bit is made up from one bit of the passkey e.g. if
// the passkey bit is 1, then Z = 0x81 and if the passkey bit is 0, then Z = 0x80.
//
// [ Vol 3, Part H 2.2.6]
func F4(U, V [32]byte, X [16]byte, Z byte) ([16]byte, error) {
	m := append(U[:], V[:]...)
	m = append(m, Z)
	return AES_CMAC(X, m)
}

// F5 The LE Secure Connections key generation function f5 is used to generate
// derived keying material in order to create the LTK and keys for the commitment
// function f6 during the LE Secure Connections pairing process.
//
// [ Vol 3, Part H 2.2.7]
func F5(W [32]byte, N1, N2 [16]byte, A1, A2 [7]byte) (out [16]byte, err error) {
	t, err := AES_CMAC(salt, W[:])
	if err != nil {
		return t, err
	}
	m1 := [48]byte{0x00}
	m2 := [48]byte{0x01}
	copy(m1[1:], N1[:])
	copy(m1[17:], N2[:])
	copy(m1[33:], A1[:])
	copy(m1[40:], A2[:])
	m1[47] = 0xFF
	copy(m2[:], m1[1:])
	cmac1, err := AES_CMAC(t, m1[:])
	if err != nil {
		return cmac1, err
	}
	cmac2, err := AES_CMAC(t, m2[:])
	if err != nil {
		return cmac2, err
	}
	copy(out[:], append(cmac1[:], cmac2[:]...))
	return
}

// F6 The LE Secure Connections check value generation function f6 is used to
// generate check values during authentication stage 2 of the LE Secure Connections
// pairing process.
//
// [ Vol 3, Part H 2.2.8]
func F6(W, N1, N2, R [16]byte, IOcap [3]byte, A1, A2 [7]byte) (out [16]byte, err error) {
	var m [65]byte
	copy(m[:], N1[:])
	copy(m[16:], N2[:])
	copy(m[32:], R[:])
	copy(m[48:], IOcap[:])
	copy(m[51:], A1[:])
	copy(m[58:], A2[:])
	return AES_CMAC(W, m[:])
}
