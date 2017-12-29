package smp

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var tkRunes = []rune("0123456789")
var randBytes = []byte("0123456789")

func GeneratePINCode() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = tkRunes[rand.Intn(len(tkRunes))]
	}
	return string(b)
}

func GenerateRand() (b [16]byte) {
	for i := range b {
		b[i] = randBytes[rand.Intn(len(tkRunes))]
	}
	return b
}

func PinCodeToTempKey(pinCode string) (tk [16]byte, err error) {
	pin, err := strconv.Atoi(pinCode)
	if err != nil {
		return tk, err
	}
	binary.BigEndian.PutUint32(tk[:], uint32(pin))
	return
}
func IsLESecureConnection(req *PairingRequest, rsp *PairingResponse) bool {
	return req.AuthReq.LESecureConnection == 0x01 && rsp.AuthReq.LESecureConnection == 0x001
}

func KeyGenMethodToUse(req *PairingRequest, rsp *PairingResponse, initiator bool) KeyGenMethod {
	if IsLESecureConnection(req, rsp) {
		if req.OOBDataFlag == OOBAuthPresent || rsp.OOBDataFlag == OOBAuthPresent {
			return KeyGenMethodOOB
		} else if req.AuthReq.MITM == 0x01 || rsp.AuthReq.MITM == 0x01 {
			return KeyGenMethodToUseForIOCap(req.IOCapability, rsp.IOCapability, initiator, true)
		}
		return KeyGenMethodJustWorks
	}

	// use LE Legacy pairing
	if req.OOBDataFlag == OOBAuthPresent && rsp.OOBDataFlag == OOBAuthPresent {
		return KeyGenMethodOOB
	} else if req.AuthReq.MITM == 0x01 || rsp.AuthReq.MITM == 0x01 {
		return KeyGenMethodToUseForIOCap(req.IOCapability, rsp.IOCapability, initiator, false)
	}
	return KeyGenMethodJustWorks
}

func KeyGenMethodToUseForIOCap(reqIOCap, rspIOCap uint8, initiator, secure bool) KeyGenMethod {
	switch reqIOCap {
	case IOCapDisplayOnly:
		switch rspIOCap {
		case IOCapDisplayOnly, IOCapDisplayYesNo, IOCapNoInputNoOutput:
			return KeyGenMethodJustWorks
		case IOCapKeyboardOnly, IOCapKeyboardDisplay:
			if initiator {
				return KeyGenMethodPKEntryInitiatorDisplay
			}
			return KeyGenMethodPKEntryResponderDisplay
		}
	case IOCapDisplayYesNo:
		switch rspIOCap {
		case IOCapDisplayOnly, IOCapNoInputNoOutput:
			return KeyGenMethodJustWorks
		case IOCapDisplayYesNo:
			if secure {
				return KeyGenMethodNumCompare
			}
			return KeyGenMethodJustWorks
		case IOCapKeyboardOnly:
			if initiator {
				return KeyGenMethodPKEntryInitiatorDisplay
			}
			return KeyGenMethodPKEntryResponderDisplay
		case IOCapKeyboardDisplay:
			if secure {
				return KeyGenMethodNumCompare
			}
			if initiator {
				return KeyGenMethodPKEntryInitiatorDisplay
			}
			return KeyGenMethodPKEntryResponderDisplay
		}
	case IOCapKeyboardOnly:
		switch rspIOCap {
		case IOCapDisplayOnly, IOCapDisplayYesNo:
			if initiator {
				return KeyGenMethodPKEntryResponderDisplay
			}
			return KeyGenMethodPKEntryInitiatorDisplay
		case IOCapKeyboardOnly:
			return KeyGenMethodPKEntryBothInput
		case IOCapKeyboardDisplay:
			if initiator {
				return KeyGenMethodPKEntryResponderDisplay
			}
			return KeyGenMethodPKEntryInitiatorDisplay
		case IOCapNoInputNoOutput:
			return KeyGenMethodJustWorks
		}
	case IOCapKeyboardDisplay:
		switch rspIOCap {
		case IOCapDisplayOnly:
			if initiator {
				return KeyGenMethodPKEntryResponderDisplay
			}
			return KeyGenMethodPKEntryInitiatorDisplay
		case IOCapDisplayYesNo:
			if secure {
				return KeyGenMethodNumCompare
			}
			if initiator {
				return KeyGenMethodPKEntryResponderDisplay
			}
			return KeyGenMethodPKEntryInitiatorDisplay
		case IOCapKeyboardOnly:
			if initiator {
				return KeyGenMethodPKEntryInitiatorDisplay
			}
			return KeyGenMethodPKEntryResponderDisplay
		case IOCapNoInputNoOutput:
			return KeyGenMethodJustWorks
		case IOCapKeyboardDisplay:
			if secure {
				return KeyGenMethodNumCompare
			}
			if initiator {
				return KeyGenMethodPKEntryInitiatorDisplay
			}
			return KeyGenMethodPKEntryResponderDisplay
		}
	case IOCapNoInputNoOutput:
		return KeyGenMethodJustWorks
	}
	// should never reach this point unless new IO capabilities are added to
	// bluetooth - i.e. Face Recognition Only
	return KeyGenMethodJustWorks
}
