package smp

const (
	// IO Capability Table 3.3. See also section 2.3.2
	IOCapDisplayOnly     = 0x00
	IOCapDisplayYesNo    = 0x01
	IOCapKeyboardOnly    = 0x02
	IOCapNoInputNoOutput = 0x03
	IOCapKeyboardDisplay = 0x04

	// OOBAuthPresent OOB Authentication data not present. [ Section 3.5 Table 3.4 ]
	OOBAuthNotPresent = 0x00

	// OOBAuthPresent OOB Authentication data from remote device present. [ Section 3.5 Table 3.4 ]
	OOBAuthPresent = 0x01

	// BondingFlagBonding [ Section 3.5 Table 3.5 ]
	BondingFlagBonding = 0x01

	// BondingFlagNoBonding [ Section 3.5 Table 3.5 ]
	BondingFlagNoBonding = 0x00

	// ReasonReserved reserved for future use. [ Section 3.5 Table 3.6 ]
	ReasonReserved = 0x00

	// ReasonPassKeyEntyFailed The user input of passkey failed, for example,
	// the user cancelled the operation. [ Section 3.5 Table 3.6 ]
	ReasonPassKeyEntyFailed = 0x01

	// ReasonOOBNotAvailable The OOB data is not available. [ Section 3.5 Table 3.6 ]
	ReasonOOBNotAvailable = 0x02

	// ReasonAuthRequirements The pairing procedure cannot be performed as
	// authentication requirements cannot be met due to IO capabilities of
	// one or both devices. [ Section 3.5 Table 3.6 ]
	ReasonAuthRequirements = 0x03

	// ReasonConfirmValueFailed The confirm value does not match the calculated
	// compare value. [ Section 3.5 Table 3.6 ]
	ReasonConfirmValueFailed = 0x04

	// ReasonPairingNotSupported Pairing is not supported by the device. [ Section 3.5 Table 3.6 ]
	ReasonPairingNotSupported = 0x05

	// ReasonEncryptionKeySize The resultant encryption key size is insufficient
	// for the security requirements of this device. [ Section 3.5 Table 3.6 ]
	ReasonEncryptionKeySize = 0x06

	// ReasonCommandNotSupported The SMP command received is not sup- ported
	// on this device. [ Section 3.5 Table 3.6 ]
	ReasonCommandNotSupported = 0x07

	// ReasonUnspecifiedReason failed due to an unspecified reason. [ Section 3.5 Table 3.6 ]
	ReasonUnspecifiedReason = 0x08

	// ReasonRepeatedAttempts Pairing or authentication procedure is disallowed
	// because too little time has elapsed since last pairing request or security
	// request. [ Section 3.5 Table 3.6 ]
	ReasonRepeatedAttempts = 0x09

	// ReasonInvalidParams indicates the command length is invalid a parameter
	// is outside of the specified rane. [ Section 3.5 Table 3.6 ]
	ReasonInvalidParams = 0x0A

	// AddrTypePublic if the address type if public or BDADDR is set to all zeros [ Section 3.6.5 ]
	AddrTypePublic = 0x00

	// AddTypeRandom if BDADDR is a static random device address [ Section 3.6.5 ]
	AddTypeRandom = 0x01

	// NotificationTypePasskeyEntryStarted Passkey entry started
	NotificationTypePasskeyEntryStarted = 0x00

	// NotificationTypePasskeyDigitEntered Passkey digit entered
	NotificationTypePasskeyDigitEntered = 0x00

	// NotificationTypePasskeyDigitErased Passkey digit erased
	NotificationTypePasskeyDigitErased = 0x00

	// NotificationTypePasskeyCleared Passkey cleared
	NotificationTypePasskeyCleared = 0x00

	// NotificationTypePasskeyEntryCompleted Passkey completed
	NotificationTypePasskeyEntryCompleted = 0x00
)

type AuthReq struct {
	BondingFlag        uint8
	MITM               uint8
	LESecureConnection uint8
	Keypress           uint8
}

func (a *AuthReq) Unmarshal(b byte) error {
	a.BondingFlag = (b >> 6) & 0x03
	a.MITM = (b >> 5) & 0x01
	a.LESecureConnection = (b >> 4) & 0x01
	a.Keypress = (b >> 3) & 0x01
	return nil
}

func (a *AuthReq) Marshal() byte {
	return 0xFF&
		(a.BondingFlag&0x03<<6) |
		(a.MITM & 0x01 << 5) |
		(a.LESecureConnection & 0x01 << 4) |
		(a.Keypress & 0x01 << 3)
}

type KeyDist struct {
	EncryptionKey uint8
	IDKey         uint8
	Sign          uint8
	LinkKey       uint8
}

func (k *KeyDist) Unmarshal(b byte) error {
	k.EncryptionKey = (b >> 7) & 0x01
	k.IDKey = (b >> 6) & 0x01
	k.Sign = (b >> 5) & 0x01
	k.LinkKey = (b >> 4) & 0x01
	return nil
}

func (k *KeyDist) Marshal() byte {
	return 0xFF&
		(k.EncryptionKey&0x01>>7) |
		(k.IDKey & 0x01 >> 6) |
		(k.Sign & 0x01 >> 5) |
		(k.LinkKey & 0x01 >> 4)
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *PairingRequest) Marshal() []byte {
	return []byte{
		s.IOCapability,
		s.OOBDataFlag,
		s.AuthReq.Marshal(),
		s.MaxEncryptionKeySize,
		s.InitiatorKeyDist.Marshal(),
		s.ResponderKeyDist.Marshal(),
	}
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingRequest) Unmarshal(b []byte) error {
	s.IOCapability = b[0]
	s.OOBDataFlag = b[1]
	s.AuthReq.Unmarshal(b[2])
	s.MaxEncryptionKeySize = b[3]
	s.InitiatorKeyDist.Unmarshal(b[4])
	s.ResponderKeyDist.Unmarshal(b[5])
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *PairingResponse) Marshal() []byte {
	return []byte{
		s.IOCapability,
		s.OOBDataFlag,
		s.AuthReq.Marshal(),
		s.MaxEncryptionKeySize,
		s.InitiatorKeyDist.Marshal(),
		s.ResponderKeyDist.Marshal(),
	}
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingResponse) Unmarshal(b []byte) error {
	s.IOCapability = b[0]
	s.OOBDataFlag = b[1]
	s.AuthReq.Unmarshal(b[2])
	s.MaxEncryptionKeySize = b[3]
	s.InitiatorKeyDist.Unmarshal(b[4])
	s.ResponderKeyDist.Unmarshal(b[5])
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *PairingConfirm) Marshal() []byte {
	return s.ConfirmValue
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingConfirm) Unmarshal(b []byte) error {
	s.ConfirmValue = b
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *PairingRandom) Marshal() []byte {
	return s.RandomValue
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingRandom) Unmarshal(b []byte) error {
	s.RandomValue = b
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *SigningInformation) Marshal() []byte {
	return s.SignatureKey
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *SigningInformation) Unmarshal(b []byte) error {
	s.SignatureKey = b
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *IdentityIdentification) Marshal() []byte {
	return s.IdentityResolvingKey
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *IdentityIdentification) Unmarshal(b []byte) error {
	s.IdentityResolvingKey = b
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *EncryptionInformation) Marshal() []byte {
	return s.LongTermKey
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *EncryptionInformation) Unmarshal(b []byte) error {
	s.LongTermKey = b
	return nil
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *SecurityRequest) Marshal() []byte {
	return []byte{s.AuthReq.Marshal()}
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *SecurityRequest) Unmarshal(b []byte) error {
	return s.AuthReq.Unmarshal(b[0])
}

// Marshal Serializes the struct into binary data int LittleEndian order
func (s *IdentityAddressIdentification) Marshal() []byte {
	b := make([]byte, 7)
	b[0] = s.AddrType
	b[1] = s.BDADDR[5]
	b[2] = s.BDADDR[4]
	b[3] = s.BDADDR[3]
	b[4] = s.BDADDR[2]
	b[5] = s.BDADDR[1]
	b[6] = s.BDADDR[0]
	return b
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *IdentityAddressIdentification) Unmarshal(b []byte) error {
	s.AddrType = b[0]
	s.BDADDR = [6]byte{b[6], b[5], b[4], b[3], b[2], b[1]}
	return nil
}
