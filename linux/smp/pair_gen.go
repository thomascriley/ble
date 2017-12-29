package smp

import (
	"bytes"
	"encoding/binary"
)

// PairingRequestCode is the code of Pairing Request signaling packet.
const PairingRequestCode = 0x01

// PairingRequest implements Pairing Request (0x01) [Vol 3, Part H, 3.5.1].
type PairingRequest struct {
	IOCapability         uint8
	OOBDataFlag          uint8
	MaxEncryptionKeySize uint8

	AuthReq AuthReq

	InitiatorKeyDist KeyDist
	ResponderKeyDist KeyDist
}

// Code returns the event code of the command.
func (s PairingRequest) Code() int { return 0x01 }

// PairingResponseCode is the code of Pairing Response signaling packet.
const PairingResponseCode = 0x02

// PairingResponse implements Pairing Response (0x02) [Vol 3, Part H, 3.5.2].
type PairingResponse struct {
	IOCapability         uint8
	OOBDataFlag          uint8
	MaxEncryptionKeySize uint8

	AuthReq AuthReq

	InitiatorKeyDist KeyDist
	ResponderKeyDist KeyDist
}

// Code returns the event code of the command.
func (s PairingResponse) Code() int { return 0x02 }

// PairingConfirmCode is the code of Pairing Confirm signaling packet.
const PairingConfirmCode = 0x03

// PairingConfirm implements Pairing Confirm (0x03) [Vol 3, Part H, 3.5.3].
type PairingConfirm struct {
	ConfirmValue []byte
}

// Code returns the event code of the command.
func (s PairingConfirm) Code() int { return 0x03 }

// PairingRandomCode is the code of Pairing Random signaling packet.
const PairingRandomCode = 0x04

// PairingRandom implements Pairing Random (0x04) [Vol 3, Part H, 3.5.4].
type PairingRandom struct {
	RandomValue []byte
}

// Code returns the event code of the command.
func (s PairingRandom) Code() int { return 0x04 }

// PairingFailedCode is the code of Pairing Failed signaling packet.
const PairingFailedCode = 0x05

// PairingFailed implements Pairing Failed (0x05) [Vol 3, Part H, 3.5.5].
type PairingFailed struct {
	Reason uint8
}

// Code returns the event code of the command.
func (s PairingFailed) Code() int { return 0x05 }

// Marshal serializes the command parameters into binary form.
func (s *PairingFailed) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingFailed) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// EncryptionInformationCode is the code of Encryption Information signaling packet.
const EncryptionInformationCode = 0x06

// EncryptionInformation implements Encryption Information (0x06) [Vol 3, Part H, 3.5.2].
type EncryptionInformation struct {
	LongTermKey []byte
}

// Code returns the event code of the command.
func (s EncryptionInformation) Code() int { return 0x06 }

// MasterIdentificationCode is the code of Master Identification signaling packet.
const MasterIdentificationCode = 0x07

// MasterIdentification implements Master Identification (0x07) [Vol 3, Part H, 3.5.3].
type MasterIdentification struct {
	EDIV uint16
	Rand uint64
}

// Code returns the event code of the command.
func (s MasterIdentification) Code() int { return 0x07 }

// Marshal serializes the command parameters into binary form.
func (s *MasterIdentification) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *MasterIdentification) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// IdentityIdentificationCode is the code of Identity Identification signaling packet.
const IdentityIdentificationCode = 0x08

// IdentityIdentification implements Identity Identification (0x08) [Vol 3, Part H, 3.5.4].
type IdentityIdentification struct {
	IdentityResolvingKey []byte
}

// Code returns the event code of the command.
func (s IdentityIdentification) Code() int { return 0x08 }

// IdentityAddressIdentificationCode is the code of Identity Address Identification signaling packet.
const IdentityAddressIdentificationCode = 0x09

// IdentityAddressIdentification implements Identity Address Identification (0x09) [Vol 3, Part H, 3.5.5].
type IdentityAddressIdentification struct {
	AddrType uint8
	BDADDR   [6]byte
}

// Code returns the event code of the command.
func (s IdentityAddressIdentification) Code() int { return 0x09 }

// SigningInformationCode is the code of Signing Information signaling packet.
const SigningInformationCode = 0x0A

// SigningInformation implements Signing Information (0x0A) [Vol 3, Part H, 3.5.6].
type SigningInformation struct {
	SignatureKey []byte
}

// Code returns the event code of the command.
func (s SigningInformation) Code() int { return 0x0A }

// SecurityRequestCode is the code of Security Request signaling packet.
const SecurityRequestCode = 0x0B

// SecurityRequest implements Security Request (0x0B) [Vol 3, Part H, 3.5.6].
type SecurityRequest struct {
	AuthReq AuthReq
}

// Code returns the event code of the command.
func (s SecurityRequest) Code() int { return 0x0B }

// PairingPublicKeyCode is the code of Pairing Public Key signaling packet.
const PairingPublicKeyCode = 0x0C

// PairingPublicKey implements Pairing Public Key (0x0C) [Vol 3, Part H, 3.5.6].
type PairingPublicKey struct {
	KeyX [32]byte
	KeyY [32]byte
}

// Code returns the event code of the command.
func (s PairingPublicKey) Code() int { return 0x0C }

// Marshal serializes the command parameters into binary form.
func (s *PairingPublicKey) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingPublicKey) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// PairingDHKeyCheckCode is the code of Pairing DHKey Check signaling packet.
const PairingDHKeyCheckCode = 0x0D

// PairingDHKeyCheck implements Pairing DHKey Check (0x0D) [Vol 3, Part H, 3.5.7].
type PairingDHKeyCheck struct {
	DHKeyCheck [16]byte
}

// Code returns the event code of the command.
func (s PairingDHKeyCheck) Code() int { return 0x0D }

// Marshal serializes the command parameters into binary form.
func (s *PairingDHKeyCheck) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *PairingDHKeyCheck) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// KeypressNotificationCode is the code of Keypress Notification signaling packet.
const KeypressNotificationCode = 0x0E

// KeypressNotification implements Keypress Notification (0x0E) [Vol 3, Part H, 3.5.8].
type KeypressNotification struct {
	NotificationType uint8
}

// Code returns the event code of the command.
func (s KeypressNotification) Code() int { return 0x0E }

// Marshal serializes the command parameters into binary form.
func (s *KeypressNotification) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *KeypressNotification) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}
