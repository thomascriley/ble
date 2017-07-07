package l2cap

import (
	"bytes"
	"encoding/binary"
)

// SignalCommandReject is the code of Command Reject signaling packet.
const SignalCommandReject = 0x01

// CommandReject implements Command Reject (0x01) [Vol 3, Part A, 4.1].
type CommandReject struct {
	Reason         uint16
	ActualSigMTU   uint16
	SourceCID      uint16
	DestinationCID uint16
}

// Code returns the event code of the command.
func (s CommandReject) Code() int { return 0x01 }

// SignalConnectionRequest is the code of Connection Request signaling packet.
const SignalConnectionRequest = 0x02

// ConnectionRequest implements Connection Request (0x02) [Vol 3, Part A, 4.2].
type ConnectionRequest struct {
	PSM       uint16
	SourceCID uint16
}

// Code returns the event code of the command.
func (s ConnectionRequest) Code() int { return 0x02 }

// Marshal serializes the command parameters into binary form.
func (s *ConnectionRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConnectionRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalConnectionResponse is the code of Connection Response signaling packet.
const SignalConnectionResponse = 0x03

// ConnectionResponse implements Connection Response (0x03) [Vol 3, Part A, 4.3].
type ConnectionResponse struct {
	DestinationCID uint16
	SourceCID      uint16
	Result         uint16
	Status         uint16
}

// Code returns the event code of the command.
func (s ConnectionResponse) Code() int { return 0x03 }

// Marshal serializes the command parameters into binary form.
func (s *ConnectionResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConnectionResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalConfigurationRequest is the code of Configuration Request signaling packet.
const SignalConfigurationRequest = 0x04

// ConfigurationRequest implements Configuration Request (0x04) [Vol 3, Part A, 4.4].
type ConfigurationRequest struct {
	DestinationCID       uint16
	Flags                uint16
	ConfigurationOptions []Option
}

// Code returns the event code of the command.
func (s ConfigurationRequest) Code() int { return 0x04 }

// SignalConfigurationResponse is the code of Configuration Response signaling packet.
const SignalConfigurationResponse = 0x05

// ConfigurationResponse implements Configuration Response (0x05) [Vol 3, Part A, 4.5].
type ConfigurationResponse struct {
	SourceCID            uint16
	Flags                uint16
	Result               uint16
	ConfigurationOptions []Option
}

// Code returns the event code of the command.
func (s ConfigurationResponse) Code() int { return 0x05 }

// SignalDisconnectRequest is the code of Disconnect Request signaling packet.
const SignalDisconnectRequest = 0x06

// DisconnectRequest implements Disconnect Request (0x06) [Vol 3, Part A, 4.6].
type DisconnectRequest struct {
	DestinationCID uint16
	SourceCID      uint16
}

// Code returns the event code of the command.
func (s DisconnectRequest) Code() int { return 0x06 }

// Marshal serializes the command parameters into binary form.
func (s *DisconnectRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *DisconnectRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalDisconnectResponse is the code of Disconnect Response signaling packet.
const SignalDisconnectResponse = 0x07

// DisconnectResponse implements Disconnect Response (0x07) [Vol 3, Part A, 4.7].
type DisconnectResponse struct {
	DestinationCID uint16
	SourceCID      uint16
}

// Code returns the event code of the command.
func (s DisconnectResponse) Code() int { return 0x07 }

// Marshal serializes the command parameters into binary form.
func (s *DisconnectResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *DisconnectResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalEchoRequest is the code of Echo Request signaling packet.
const SignalEchoRequest = 0x08

// EchoRequest implements Echo Request (0x08) [Vol 3, Part A, 4.8].
type EchoRequest struct {
	Data uint16
}

// Code returns the event code of the command.
func (s EchoRequest) Code() int { return 0x08 }

// Marshal serializes the command parameters into binary form.
func (s *EchoRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *EchoRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalEchoResponse is the code of Echo Response signaling packet.
const SignalEchoResponse = 0x09

// EchoResponse implements Echo Response (0x09) [Vol 3, Part A, 4.9].
type EchoResponse struct {
	Data uint16
}

// Code returns the event code of the command.
func (s EchoResponse) Code() int { return 0x09 }

// Marshal serializes the command parameters into binary form.
func (s *EchoResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *EchoResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalInformationRequest is the code of Information Request signaling packet.
const SignalInformationRequest = 0x0A

// InformationRequest implements Information Request (0x0A) [Vol 3, Part A, 4.10].
type InformationRequest struct {
	InfoType uint16
}

// Code returns the event code of the command.
func (s InformationRequest) Code() int { return 0x0A }

// Marshal serializes the command parameters into binary form.
func (s *InformationRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *InformationRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalInformationResponse is the code of Information Response signaling packet.
const SignalInformationResponse = 0x0B

// InformationResponse implements Information Response (0x0B) [Vol 3, Part A, 4.11].
type InformationResponse struct {
	InfoType            uint16
	Result              uint16
	Data                []byte
	ConnectionlessMTU   uint16
	ExtendedFeatureMask uint32
	FixedChannels       uint64
}

// Code returns the event code of the command.
func (s InformationResponse) Code() int { return 0x0B }

// SignalCreateChannelRequest is the code of Create Channel Request signaling packet.
const SignalCreateChannelRequest = 0x0C

// CreateChannelRequest implements Create Channel Request (0x0C) [Vol 3, Part A, 4.14].
type CreateChannelRequest struct {
	PSM          uint16
	SourceCID    uint16
	ControllerID uint8
}

// Code returns the event code of the command.
func (s CreateChannelRequest) Code() int { return 0x0C }

// Marshal serializes the command parameters into binary form.
func (s *CreateChannelRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *CreateChannelRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalCreateChannelResponse is the code of Create Channel Response signaling packet.
const SignalCreateChannelResponse = 0x0D

// CreateChannelResponse implements Create Channel Response (0x0D) [Vol 3, Part A, 4.15].
type CreateChannelResponse struct {
	DestinationCID uint16
	SourceCID      uint16
	Result         uint16
	Status         uint16
}

// Code returns the event code of the command.
func (s CreateChannelResponse) Code() int { return 0x0D }

// Marshal serializes the command parameters into binary form.
func (s *CreateChannelResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *CreateChannelResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalMoveChannelRequest is the code of Move Channel Request signaling packet.
const SignalMoveChannelRequest = 0x0E

// MoveChannelRequest implements Move Channel Request (0x0E) [Vol 3, Part A, 4.16].
type MoveChannelRequest struct {
	InitiatorCID     uint16
	DestControllerID uint8
}

// Code returns the event code of the command.
func (s MoveChannelRequest) Code() int { return 0x0E }

// Marshal serializes the command parameters into binary form.
func (s *MoveChannelRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *MoveChannelRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalMoveChannelResponse is the code of Move Channel Response signaling packet.
const SignalMoveChannelResponse = 0x0F

// MoveChannelResponse implements Move Channel Response (0x0F) [Vol 3, Part A, 4.17].
type MoveChannelResponse struct {
	InitiatorCID uint16
	Result       uint16
}

// Code returns the event code of the command.
func (s MoveChannelResponse) Code() int { return 0x0F }

// Marshal serializes the command parameters into binary form.
func (s *MoveChannelResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *MoveChannelResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalMoveChannelConfirmation is the code of Move Channel Confirmation signaling packet.
const SignalMoveChannelConfirmation = 0x10

// MoveChannelConfirmation implements Move Channel Confirmation (0x10) [Vol 3, Part A, 4.18].
type MoveChannelConfirmation struct {
	InitiatorCID uint16
	Result       uint16
}

// Code returns the event code of the command.
func (s MoveChannelConfirmation) Code() int { return 0x10 }

// Marshal serializes the command parameters into binary form.
func (s *MoveChannelConfirmation) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *MoveChannelConfirmation) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalMoveChannelConfirmationResponse is the code of Move Channel Confirmation Response signaling packet.
const SignalMoveChannelConfirmationResponse = 0x11

// MoveChannelConfirmationResponse implements Move Channel Confirmation Response (0x11) [Vol 3, Part A, 4.19].
type MoveChannelConfirmationResponse struct {
	InitiatorCID uint16
}

// Code returns the event code of the command.
func (s MoveChannelConfirmationResponse) Code() int { return 0x11 }

// Marshal serializes the command parameters into binary form.
func (s *MoveChannelConfirmationResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *MoveChannelConfirmationResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalConnectionParameterUpdateRequest is the code of Connection Parameter Update Request signaling packet.
const SignalConnectionParameterUpdateRequest = 0x12

// ConnectionParameterUpdateRequest implements Connection Parameter Update Request (0x12) [Vol 3, Part A, 4.20].
type ConnectionParameterUpdateRequest struct {
	IntervalMin       uint16
	IntervalMax       uint16
	SlaveLatency      uint16
	TimeoutMultiplier uint16
}

// Code returns the event code of the command.
func (s ConnectionParameterUpdateRequest) Code() int { return 0x12 }

// Marshal serializes the command parameters into binary form.
func (s *ConnectionParameterUpdateRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConnectionParameterUpdateRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalConnectionParameterUpdateResponse is the code of Connection Parameter Update Response signaling packet.
const SignalConnectionParameterUpdateResponse = 0x13

// ConnectionParameterUpdateResponse implements Connection Parameter Update Response (0x13) [Vol 3, Part A, 4.21].
type ConnectionParameterUpdateResponse struct {
	Result uint16
}

// Code returns the event code of the command.
func (s ConnectionParameterUpdateResponse) Code() int { return 0x13 }

// Marshal serializes the command parameters into binary form.
func (s *ConnectionParameterUpdateResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *ConnectionParameterUpdateResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalLECreditBasedConnectionRequest is the code of LE Credit Based Connection Request signaling packet.
const SignalLECreditBasedConnectionRequest = 0x14

// LECreditBasedConnectionRequest implements LE Credit Based Connection Request (0x14) [Vol 3, Part A, 4.22].
type LECreditBasedConnectionRequest struct {
	LEPSM          uint16
	SourceCID      uint16
	MTU            uint16
	MPS            uint16
	InitialCredits uint16
}

// Code returns the event code of the command.
func (s LECreditBasedConnectionRequest) Code() int { return 0x14 }

// Marshal serializes the command parameters into binary form.
func (s *LECreditBasedConnectionRequest) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *LECreditBasedConnectionRequest) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalLECreditBasedConnectionResponse is the code of LE Credit Based Connection Response signaling packet.
const SignalLECreditBasedConnectionResponse = 0x15

// LECreditBasedConnectionResponse implements LE Credit Based Connection Response (0x15) [Vol 3, Part A, 4.23].
type LECreditBasedConnectionResponse struct {
	DestinationCID    uint16
	MTU               uint16
	MPS               uint16
	InitialCreditsCID uint16
	Result            uint16
}

// Code returns the event code of the command.
func (s LECreditBasedConnectionResponse) Code() int { return 0x15 }

// Marshal serializes the command parameters into binary form.
func (s *LECreditBasedConnectionResponse) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *LECreditBasedConnectionResponse) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}

// SignalLEFlowControlCredit is the code of LE Flow Control Credit signaling packet.
const SignalLEFlowControlCredit = 0x16

// LEFlowControlCredit implements LE Flow Control Credit (0x16) [Vol 3, Part A, 4.24].
type LEFlowControlCredit struct {
	CID     uint16
	Credits uint16
}

// Code returns the event code of the command.
func (s LEFlowControlCredit) Code() int { return 0x16 }

// Marshal serializes the command parameters into binary form.
func (s *LEFlowControlCredit) Marshal() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (s *LEFlowControlCredit) Unmarshal(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, s)
}
