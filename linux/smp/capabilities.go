package smp

type KeyGenMethod int

const (
	KeyGenMethodOOB KeyGenMethod = iota + 1
	KeyGenMethodJustWorks
	KeyGenMethodNumCompare
	KeyGenMethodPKEntryInitiatorDisplay
	KeyGenMethodPKEntryResponderDisplay
	KeyGenMethodPKEntryBothInput
)

// Device can only display a 6 digit decimal number
type Display func(code string)

// Device has two buttons that easily map to yes or no and can display a 6 digit
// decimal number
type DisplayYesNo func(code string) bool

// Can only input numberic values 0-9 and a confirmation
type Keyboard func() string

// sends and returns external information for the pairing process
type OOBAuth func() ([16]byte, error)

// sends, receives and stores a long term key for later use
type ExchangeLongTermKeys func([]byte) ([]byte, error)

type Capabilities struct {
	Display              Display
	Keyboard             Keyboard
	DisplayYesNo         DisplayYesNo
	OOBAuth              OOBAuth
	ExchangeLongTermKeys ExchangeLongTermKeys
	LESecureConnection   bool
	ManInTheMiddle       bool
}

func (i *Capabilities) IOCapability() uint8 {
	if i.Keyboard != nil && i.Display != nil {
		return IOCapKeyboardDisplay
	} else if i.DisplayYesNo != nil {
		return IOCapDisplayYesNo
	} else if i.Keyboard != nil {
		return IOCapKeyboardOnly
	} else if i.Display != nil {
		return IOCapDisplayOnly
	}
	return IOCapNoInputNoOutput
}

func (i *Capabilities) OOBDataFlag() uint8 {
	if i.OOBAuth != nil {
		return OOBAuthPresent
	}
	return OOBAuthNotPresent
}

func (i *Capabilities) BondingFlags() uint8 {
	if i.ExchangeLongTermKeys != nil {
		return OOBAuthPresent
	}
	return OOBAuthNotPresent
}

func (i *Capabilities) MITM() uint8 {
	if i.ManInTheMiddle {
		return 0x01
	}
	return 0x00
}

func (i *Capabilities) SecureConnection() uint8 {
	if i.LESecureConnection {
		return 0x01
	}
	return 0x00
}

func (i *Capabilities) PairingRequest() *PairingRequest {
	return &PairingRequest{
		IOCapability: i.IOCapability(),
		AuthReq: AuthReq{
			MITM:               i.MITM(),
			BondingFlag:        i.BondingFlags(),
			Keypress:           0x00,
			LESecureConnection: i.SecureConnection(),
		},
		OOBDataFlag:          i.OOBDataFlag(),
		MaxEncryptionKeySize: 0x10,
		InitiatorKeyDist: KeyDist{
			EncryptionKey: 0x00,
			IDKey:         0x00,
			Sign:          0x00,
			LinkKey:       0x00,
		},
		ResponderKeyDist: KeyDist{
			EncryptionKey: 0x00,
			IDKey:         0x00,
			Sign:          0x00,
			LinkKey:       0x00,
		},
	}
}

func (i *Capabilities) PairingResponse(req *PairingRequest) *PairingResponse {
	return &PairingResponse{
		IOCapability: i.IOCapability(),
		AuthReq: AuthReq{
			MITM:               i.MITM(),
			BondingFlag:        i.BondingFlags(),
			Keypress:           0x00,
			LESecureConnection: i.SecureConnection(),
		},
		OOBDataFlag:          i.OOBDataFlag(),
		MaxEncryptionKeySize: 0x10,
		InitiatorKeyDist: KeyDist{
			EncryptionKey: 0x00,
			IDKey:         0x00,
			Sign:          0x00,
			LinkKey:       0x00,
		},
		ResponderKeyDist: KeyDist{
			EncryptionKey: 0x00,
			IDKey:         0x00,
			Sign:          0x00,
			LinkKey:       0x00,
		},
	}
}

func (i *Capabilities) GetTemporaryKey(method KeyGenMethod, initiator bool) (tk [16]byte, err error) {

	switch method {
	case KeyGenMethodPKEntryBothInput, KeyGenMethodPKEntryInitiatorDisplay, KeyGenMethodPKEntryResponderDisplay:
		pinCode := GeneratePINCode()
		switch method {
		case KeyGenMethodPKEntryInitiatorDisplay:
			if initiator {
				i.Display(pinCode)
			} else {
				pinCode = i.Keyboard()
			}
		case KeyGenMethodPKEntryResponderDisplay:
			if initiator {
				i.Display(pinCode)
			} else {
				pinCode = i.Keyboard()
			}
		case KeyGenMethodPKEntryBothInput:
			pinCode = i.Keyboard()
		}
		if tk, err = PinCodeToTempKey(pinCode); err != nil {
			return
		}
	case KeyGenMethodOOB:
		if tk, err = i.OOBAuth(); err != nil {
			return
		}
	case KeyGenMethodJustWorks:
		// empty tk
	case KeyGenMethodNumCompare:
		// TODO:
	}
	return
}
