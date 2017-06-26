package cmd

// Inquiry implements Inquiry (0x01|0x0001) [Vol 2, Part E, 7.1.1]
type Inquiry struct {
	LAP           [3]byte
	InquiryLength uint8
	NumResponses  uint8
}

func (c *Inquiry) String() string {
	return "Inquiry (0x01|0x0001)"
}

// OpCode returns the opcode of the command.
func (c *Inquiry) OpCode() int { return 0x01<<10 | 0x0001 }

// Len returns the length of the command.
func (c *Inquiry) Len() int { return 5 }

// Marshal serializes the command parameters into binary form.
func (c *Inquiry) Marshal(b []byte) error {
	return marshal(c, b)
}

// InquiryCancel implements Inquiry Cancel (0x01|0x0002) [Vol 2, Part E, 7.1.2]
type InquiryCancel struct {
}

func (c *InquiryCancel) String() string {
	return "Inquiry Cancel (0x01|0x0002)"
}

// OpCode returns the opcode of the command.
func (c *InquiryCancel) OpCode() int { return 0x01<<10 | 0x0002 }

// Len returns the length of the command.
func (c *InquiryCancel) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *InquiryCancel) Marshal(b []byte) error {
	return marshal(c, b)
}

// InquiryCancelRP returns the return parameter of Inquiry Cancel
type InquiryCancelRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *InquiryCancelRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// PeriodicInquiryMode implements Periodic Inquiry Mode (0x01|0x0003) [Vol 2, Part E, 7.1.3]
type PeriodicInquiryMode struct {
	MaxPeriodLength uint16
	MinPeriodLength uint16
	LAP             [3]byte
	InquiryLength   uint8
	NumResponses    uint8
}

func (c *PeriodicInquiryMode) String() string {
	return "Periodic Inquiry Mode (0x01|0x0003)"
}

// OpCode returns the opcode of the command.
func (c *PeriodicInquiryMode) OpCode() int { return 0x01<<10 | 0x0003 }

// Len returns the length of the command.
func (c *PeriodicInquiryMode) Len() int { return 9 }

// Marshal serializes the command parameters into binary form.
func (c *PeriodicInquiryMode) Marshal(b []byte) error {
	return marshal(c, b)
}

// PeriodicInquiryModeRP returns the return parameter of Periodic Inquiry Mode
type PeriodicInquiryModeRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *PeriodicInquiryModeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ExitPeriodicInquiryMode implements Exit Periodic Inquiry Mode (0x01|0x0004) [Vol 2, Part E, 7.1.4]
type ExitPeriodicInquiryMode struct {
}

func (c *ExitPeriodicInquiryMode) String() string {
	return "Exit Periodic Inquiry Mode (0x01|0x0004)"
}

// OpCode returns the opcode of the command.
func (c *ExitPeriodicInquiryMode) OpCode() int { return 0x01<<10 | 0x0004 }

// Len returns the length of the command.
func (c *ExitPeriodicInquiryMode) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ExitPeriodicInquiryMode) Marshal(b []byte) error {
	return marshal(c, b)
}

// ExitPeriodicInquiryModeRP returns the return parameter of Exit Periodic Inquiry Mode
type ExitPeriodicInquiryModeRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ExitPeriodicInquiryModeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// CreateConnection implements Create Connection (0x01|0x0005) [Vol 2, Part E, 7.1.5]
type CreateConnection struct {
	BDADDR                 [6]byte
	PacketType             uint16
	PageScanRepetitionMode uint8
	Reserved               uint8
	ClockOffset            uint16
	AllowRoleSwitch        uint8
}

func (c *CreateConnection) String() string {
	return "Create Connection (0x01|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *CreateConnection) OpCode() int { return 0x01<<10 | 0x0005 }

// Len returns the length of the command.
func (c *CreateConnection) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *CreateConnection) Marshal(b []byte) error {
	return marshal(c, b)
}

// Disconnect implements Disconnect (0x01|0x0006) [Vol 2, Part E, 7.1.6]
type Disconnect struct {
	ConnectionHandle uint16
	Reason           uint8
}

func (c *Disconnect) String() string {
	return "Disconnect (0x01|0x0006)"
}

// OpCode returns the opcode of the command.
func (c *Disconnect) OpCode() int { return 0x01<<10 | 0x0006 }

// Len returns the length of the command.
func (c *Disconnect) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *Disconnect) Marshal(b []byte) error {
	return marshal(c, b)
}

// CreateConnectionCancel implements Create Connection Cancel (0x01|0x0008) [Vol 2, Part E, 7.1.7]
type CreateConnectionCancel struct {
	BDADDR [6]byte
}

func (c *CreateConnectionCancel) String() string {
	return "Create Connection Cancel (0x01|0x0008)"
}

// OpCode returns the opcode of the command.
func (c *CreateConnectionCancel) OpCode() int { return 0x01<<10 | 0x0008 }

// Len returns the length of the command.
func (c *CreateConnectionCancel) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *CreateConnectionCancel) Marshal(b []byte) error {
	return marshal(c, b)
}

// CreateConnectionCancelRP returns the return parameter of Create Connection Cancel
type CreateConnectionCancelRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *CreateConnectionCancelRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// AcceptConnectionRequest implements Accept Connection Request (0x01|0x0009) [Vol 2, Part E, 7.1.8]
type AcceptConnectionRequest struct {
	BDADDR [6]byte
	Role   uint8
}

func (c *AcceptConnectionRequest) String() string {
	return "Accept Connection Request (0x01|0x0009)"
}

// OpCode returns the opcode of the command.
func (c *AcceptConnectionRequest) OpCode() int { return 0x01<<10 | 0x0009 }

// Len returns the length of the command.
func (c *AcceptConnectionRequest) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *AcceptConnectionRequest) Marshal(b []byte) error {
	return marshal(c, b)
}

// RejectConnectionRequest implements Reject Connection Request (0x01|0x000A) [Vol 2, Part E, 7.1.9]
type RejectConnectionRequest struct {
	BDADDR [6]byte
	Reason uint8
}

func (c *RejectConnectionRequest) String() string {
	return "Reject Connection Request (0x01|0x000A)"
}

// OpCode returns the opcode of the command.
func (c *RejectConnectionRequest) OpCode() int { return 0x01<<10 | 0x000A }

// Len returns the length of the command.
func (c *RejectConnectionRequest) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *RejectConnectionRequest) Marshal(b []byte) error {
	return marshal(c, b)
}

// LinkKeyRequestReply implements Link Key Request Reply (0x01|0x000B) [Vol 2, Part E, 7.1.10]
type LinkKeyRequestReply struct {
	BDADDR  [6]byte
	LinkKey [16]byte
}

func (c *LinkKeyRequestReply) String() string {
	return "Link Key Request Reply (0x01|0x000B)"
}

// OpCode returns the opcode of the command.
func (c *LinkKeyRequestReply) OpCode() int { return 0x01<<10 | 0x000B }

// Len returns the length of the command.
func (c *LinkKeyRequestReply) Len() int { return 22 }

// Marshal serializes the command parameters into binary form.
func (c *LinkKeyRequestReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LinkKeyRequestReplyRP returns the return parameter of Link Key Request Reply
type LinkKeyRequestReplyRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LinkKeyRequestReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LinkKeyRequestNegativeReply implements Link Key Request Negative Reply (0x01|0x000C) [Vol 2, Part E, 7.1.11]
type LinkKeyRequestNegativeReply struct {
	BDADDR [6]byte
}

func (c *LinkKeyRequestNegativeReply) String() string {
	return "Link Key Request Negative Reply (0x01|0x000C)"
}

// OpCode returns the opcode of the command.
func (c *LinkKeyRequestNegativeReply) OpCode() int { return 0x01<<10 | 0x000C }

// Len returns the length of the command.
func (c *LinkKeyRequestNegativeReply) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *LinkKeyRequestNegativeReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LinkKeyRequestNegativeReplyRP returns the return parameter of Link Key Request Negative Reply
type LinkKeyRequestNegativeReplyRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LinkKeyRequestNegativeReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// PinCodeRequestReply implements Pin Code Request Reply (0x01|0x000D) [Vol 2, Part E, 7.1.12]
type PinCodeRequestReply struct {
	BDADDR        [6]byte
	PINCodeLength uint8
	PINCode       [16]byte
}

func (c *PinCodeRequestReply) String() string {
	return "Pin Code Request Reply (0x01|0x000D)"
}

// OpCode returns the opcode of the command.
func (c *PinCodeRequestReply) OpCode() int { return 0x01<<10 | 0x000D }

// Len returns the length of the command.
func (c *PinCodeRequestReply) Len() int { return 23 }

// Marshal serializes the command parameters into binary form.
func (c *PinCodeRequestReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// PinCodeRequestReplyRP returns the return parameter of Pin Code Request Reply
type PinCodeRequestReplyRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *PinCodeRequestReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// PinCodeRequestNegativeReply implements Pin Code Request Negative Reply (0x01|0x000E) [Vol 2, Part E, 7.1.13]
type PinCodeRequestNegativeReply struct {
	BDADDR [6]byte
}

func (c *PinCodeRequestNegativeReply) String() string {
	return "Pin Code Request Negative Reply (0x01|0x000E)"
}

// OpCode returns the opcode of the command.
func (c *PinCodeRequestNegativeReply) OpCode() int { return 0x01<<10 | 0x000E }

// Len returns the length of the command.
func (c *PinCodeRequestNegativeReply) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *PinCodeRequestNegativeReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// PinCodeRequestNegativeReplyRP returns the return parameter of Pin Code Request Negative Reply
type PinCodeRequestNegativeReplyRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *PinCodeRequestNegativeReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ChangeConnectionPacketType implements Change Connection Packet Type (0x01|0x000F) [Vol 2, Part E, 7.1.14]
type ChangeConnectionPacketType struct {
	ConnectionHandle uint16
	PacketType       uint16
}

func (c *ChangeConnectionPacketType) String() string {
	return "Change Connection Packet Type (0x01|0x000F)"
}

// OpCode returns the opcode of the command.
func (c *ChangeConnectionPacketType) OpCode() int { return 0x01<<10 | 0x000F }

// Len returns the length of the command.
func (c *ChangeConnectionPacketType) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *ChangeConnectionPacketType) Marshal(b []byte) error {
	return marshal(c, b)
}

// AuthenticationRequested implements Authentication Requested (0x01|0x0011) [Vol 2, Part E, 7.1.15]
type AuthenticationRequested struct {
	ConnectionHandle uint16
}

func (c *AuthenticationRequested) String() string {
	return "Authentication Requested (0x01|0x0011)"
}

// OpCode returns the opcode of the command.
func (c *AuthenticationRequested) OpCode() int { return 0x01<<10 | 0x0011 }

// Len returns the length of the command.
func (c *AuthenticationRequested) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *AuthenticationRequested) Marshal(b []byte) error {
	return marshal(c, b)
}

// SetConnectionEncryption implements Set Connection Encryption (0x01|0x0013) [Vol 2, Part E, 7.1.16]
type SetConnectionEncryption struct {
	ConnectionHandle uint16
	EncryptionEnable uint8
}

func (c *SetConnectionEncryption) String() string {
	return "Set Connection Encryption (0x01|0x0013)"
}

// OpCode returns the opcode of the command.
func (c *SetConnectionEncryption) OpCode() int { return 0x01<<10 | 0x0013 }

// Len returns the length of the command.
func (c *SetConnectionEncryption) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *SetConnectionEncryption) Marshal(b []byte) error {
	return marshal(c, b)
}

// ChangeConnectionLinkKey implements Change Connection Link Key (0x01|0x0015) [Vol 2, Part E, 7.1.17]
type ChangeConnectionLinkKey struct {
	ConnectionHandle uint16
}

func (c *ChangeConnectionLinkKey) String() string {
	return "Change Connection Link Key (0x01|0x0015)"
}

// OpCode returns the opcode of the command.
func (c *ChangeConnectionLinkKey) OpCode() int { return 0x01<<10 | 0x0015 }

// Len returns the length of the command.
func (c *ChangeConnectionLinkKey) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ChangeConnectionLinkKey) Marshal(b []byte) error {
	return marshal(c, b)
}

// MasterLinkKey implements Master Link Key (0x01|0x0017) [Vol 2, Part E, 7.1.18]
type MasterLinkKey struct {
	KeyFlag uint8
}

func (c *MasterLinkKey) String() string {
	return "Master Link Key (0x01|0x0017)"
}

// OpCode returns the opcode of the command.
func (c *MasterLinkKey) OpCode() int { return 0x01<<10 | 0x0017 }

// Len returns the length of the command.
func (c *MasterLinkKey) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *MasterLinkKey) Marshal(b []byte) error {
	return marshal(c, b)
}

// RemoteNameRequest implements Remote Name Request (0x01|0x0019) [Vol 2, Part E, 7.1.19]
type RemoteNameRequest struct {
	BDADDR               [6]byte
	PageScanRepitionMode uint8
	Reserved             uint8
	ClockOffset          uint16
}

func (c *RemoteNameRequest) String() string {
	return "Remote Name Request (0x01|0x0019)"
}

// OpCode returns the opcode of the command.
func (c *RemoteNameRequest) OpCode() int { return 0x01<<10 | 0x0019 }

// Len returns the length of the command.
func (c *RemoteNameRequest) Len() int { return 10 }

// Marshal serializes the command parameters into binary form.
func (c *RemoteNameRequest) Marshal(b []byte) error {
	return marshal(c, b)
}

// RemoteNameRequestCancel implements Remote Name Request Cancel (0x01|0x001A) [Vol 2, Part E, 7.1.20]
type RemoteNameRequestCancel struct {
	BDADDR [6]byte
}

func (c *RemoteNameRequestCancel) String() string {
	return "Remote Name Request Cancel (0x01|0x001A)"
}

// OpCode returns the opcode of the command.
func (c *RemoteNameRequestCancel) OpCode() int { return 0x01<<10 | 0x001A }

// Len returns the length of the command.
func (c *RemoteNameRequestCancel) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *RemoteNameRequestCancel) Marshal(b []byte) error {
	return marshal(c, b)
}

// RemoteNameRequestCancelRP returns the return parameter of Remote Name Request Cancel
type RemoteNameRequestCancelRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *RemoteNameRequestCancelRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadRemoteSupportedFeatures implements Read Remote Supported Features (0x01|0x001B) [Vol 2, Part E, 7.1.21]
type ReadRemoteSupportedFeatures struct {
	ConnectionHandle uint16
}

func (c *ReadRemoteSupportedFeatures) String() string {
	return "Read Remote Supported Features (0x01|0x001B)"
}

// OpCode returns the opcode of the command.
func (c *ReadRemoteSupportedFeatures) OpCode() int { return 0x01<<10 | 0x001B }

// Len returns the length of the command.
func (c *ReadRemoteSupportedFeatures) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadRemoteSupportedFeatures) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadRemoteExtendedFeatures implements Read Remote Extended Features (0x01|0x001C) [Vol 2, Part E, 7.1.22]
type ReadRemoteExtendedFeatures struct {
	ConnectionHandle uint16
	PageNumber       uint8
}

func (c *ReadRemoteExtendedFeatures) String() string {
	return "Read Remote Extended Features (0x01|0x001C)"
}

// OpCode returns the opcode of the command.
func (c *ReadRemoteExtendedFeatures) OpCode() int { return 0x01<<10 | 0x001C }

// Len returns the length of the command.
func (c *ReadRemoteExtendedFeatures) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *ReadRemoteExtendedFeatures) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadRemoteVersionInformation implements Read Remote Version Information (0x01|0x001D) [Vol 2, Part E, 7.1.23]
type ReadRemoteVersionInformation struct {
	ConnectionHandle uint16
}

func (c *ReadRemoteVersionInformation) String() string {
	return "Read Remote Version Information (0x01|0x001D)"
}

// OpCode returns the opcode of the command.
func (c *ReadRemoteVersionInformation) OpCode() int { return 0x01<<10 | 0x001D }

// Len returns the length of the command.
func (c *ReadRemoteVersionInformation) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadRemoteVersionInformation) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadClockOffset implements Read Clock Offset (0x01|0x001F) [Vol 2, Part E, 7.1.24]
type ReadClockOffset struct {
	ConnectionHandle uint16
}

func (c *ReadClockOffset) String() string {
	return "Read Clock Offset (0x01|0x001F)"
}

// OpCode returns the opcode of the command.
func (c *ReadClockOffset) OpCode() int { return 0x01<<10 | 0x001F }

// Len returns the length of the command.
func (c *ReadClockOffset) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadClockOffset) Marshal(b []byte) error {
	return marshal(c, b)
}

// HoldMode implements Hold Mode (0x02|0x0001) [Vol 2, Part E, 7.2.1]
type HoldMode struct {
	ConnectionHandle    uint16
	HoldModeMaxInterval uint16
	HoldModeMinInterval uint16
}

func (c *HoldMode) String() string {
	return "Hold Mode (0x02|0x0001)"
}

// OpCode returns the opcode of the command.
func (c *HoldMode) OpCode() int { return 0x02<<10 | 0x0001 }

// Len returns the length of the command.
func (c *HoldMode) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *HoldMode) Marshal(b []byte) error {
	return marshal(c, b)
}

// SniffMode implements Sniff Mode (0x02|0x0003) [Vol 2, Part E, 7.2.2]
type SniffMode struct {
	ConnectionHandle     uint16
	SniffModeMaxInterval uint16
	SniffModeMinInterval uint16
	SniffAttempt         uint16
	SniffTimeout         uint16
}

func (c *SniffMode) String() string {
	return "Sniff Mode (0x02|0x0003)"
}

// OpCode returns the opcode of the command.
func (c *SniffMode) OpCode() int { return 0x02<<10 | 0x0003 }

// Len returns the length of the command.
func (c *SniffMode) Len() int { return 10 }

// Marshal serializes the command parameters into binary form.
func (c *SniffMode) Marshal(b []byte) error {
	return marshal(c, b)
}

// ExitSniffMode implements Exit Sniff Mode (0x02|0x0004) [Vol 2, Part E, 7.2.3]
type ExitSniffMode struct {
	ConnectionHandle uint16
}

func (c *ExitSniffMode) String() string {
	return "Exit Sniff Mode (0x02|0x0004)"
}

// OpCode returns the opcode of the command.
func (c *ExitSniffMode) OpCode() int { return 0x02<<10 | 0x0004 }

// Len returns the length of the command.
func (c *ExitSniffMode) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ExitSniffMode) Marshal(b []byte) error {
	return marshal(c, b)
}

// ParkState implements Park State (0x02|0x0005) [Vol 2, Part E, 7.2.4]
type ParkState struct {
	ConnectionHandle  uint16
	BeaconMaxInterval uint16
	BeaconMinInterval uint16
}

func (c *ParkState) String() string {
	return "Park State (0x02|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *ParkState) OpCode() int { return 0x02<<10 | 0x0005 }

// Len returns the length of the command.
func (c *ParkState) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *ParkState) Marshal(b []byte) error {
	return marshal(c, b)
}

// ExitParkState implements Exit Park State (0x02|0x0006) [Vol 2, Part E, 7.2.5]
type ExitParkState struct {
	ConnectionHandle uint16
}

func (c *ExitParkState) String() string {
	return "Exit Park State (0x02|0x0006)"
}

// OpCode returns the opcode of the command.
func (c *ExitParkState) OpCode() int { return 0x02<<10 | 0x0006 }

// Len returns the length of the command.
func (c *ExitParkState) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ExitParkState) Marshal(b []byte) error {
	return marshal(c, b)
}

// QoSSetup implements QoS Setup (0x02|0x0007) [Vol 2, Part E, 7.2.6]
type QoSSetup struct {
	ConnectionHandle uint16
	Flags            uint8
	ServiceType      uint8
	TokenRate        uint32
	BeakBandwidth    uint32
	Latency          uint32
	DelayVariation   uint32
}

func (c *QoSSetup) String() string {
	return "QoS Setup (0x02|0x0007)"
}

// OpCode returns the opcode of the command.
func (c *QoSSetup) OpCode() int { return 0x02<<10 | 0x0007 }

// Len returns the length of the command.
func (c *QoSSetup) Len() int { return 20 }

// Marshal serializes the command parameters into binary form.
func (c *QoSSetup) Marshal(b []byte) error {
	return marshal(c, b)
}

// RoleDiscovery implements Role Discovery (0x02|0x0009) [Vol 2, Part E, 7.2.7]
type RoleDiscovery struct {
	ConnectionHandle uint16
}

func (c *RoleDiscovery) String() string {
	return "Role Discovery (0x02|0x0009)"
}

// OpCode returns the opcode of the command.
func (c *RoleDiscovery) OpCode() int { return 0x02<<10 | 0x0009 }

// Len returns the length of the command.
func (c *RoleDiscovery) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *RoleDiscovery) Marshal(b []byte) error {
	return marshal(c, b)
}

// RoleDiscoveryRP returns the return parameter of Role Discovery
type RoleDiscoveryRP struct {
	Status           uint8
	ConnectionHandle uint16
	CurrentRole      uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *RoleDiscoveryRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// SwitchRole implements Switch Role (0x02|0x000B) [Vol 2, Part E, 7.2.8]
type SwitchRole struct {
	BDADDR [6]byte
	Role   uint8
}

func (c *SwitchRole) String() string {
	return "Switch Role (0x02|0x000B)"
}

// OpCode returns the opcode of the command.
func (c *SwitchRole) OpCode() int { return 0x02<<10 | 0x000B }

// Len returns the length of the command.
func (c *SwitchRole) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *SwitchRole) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadLinkPolicySettings implements Read Link Policy Settings (0x02|0x000C) [Vol 2, Part E, 7.2.9]
type ReadLinkPolicySettings struct {
	ConnectionHandle uint16
}

func (c *ReadLinkPolicySettings) String() string {
	return "Read Link Policy Settings (0x02|0x000C)"
}

// OpCode returns the opcode of the command.
func (c *ReadLinkPolicySettings) OpCode() int { return 0x02<<10 | 0x000C }

// Len returns the length of the command.
func (c *ReadLinkPolicySettings) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadLinkPolicySettings) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadLinkPolicySettingsRP returns the return parameter of Read Link Policy Settings
type ReadLinkPolicySettingsRP struct {
	Status             uint8
	ConnectionHandle   uint16
	LinkPolicySettings uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadLinkPolicySettingsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteLinkPolicySettings implements Write Link Policy Settings (0x02|0x000D) [Vol 2, Part E, 7.2.10]
type WriteLinkPolicySettings struct {
	ConnectionHandle   uint16
	LinkPolicySettings uint16
}

func (c *WriteLinkPolicySettings) String() string {
	return "Write Link Policy Settings (0x02|0x000D)"
}

// OpCode returns the opcode of the command.
func (c *WriteLinkPolicySettings) OpCode() int { return 0x02<<10 | 0x000D }

// Len returns the length of the command.
func (c *WriteLinkPolicySettings) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *WriteLinkPolicySettings) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteLinkPolicySettingsRP returns the return parameter of Write Link Policy Settings
type WriteLinkPolicySettingsRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteLinkPolicySettingsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadDefaultLinkPolicySettings implements Read Default Link Policy Settings (0x02|0x000E) [Vol 2, Part E, 7.2.11]
type ReadDefaultLinkPolicySettings struct {
}

func (c *ReadDefaultLinkPolicySettings) String() string {
	return "Read Default Link Policy Settings (0x02|0x000E)"
}

// OpCode returns the opcode of the command.
func (c *ReadDefaultLinkPolicySettings) OpCode() int { return 0x02<<10 | 0x000E }

// Len returns the length of the command.
func (c *ReadDefaultLinkPolicySettings) Len() int { return -1 }

// ReadDefaultLinkPolicySettingsRP returns the return parameter of Read Default Link Policy Settings
type ReadDefaultLinkPolicySettingsRP struct {
	Status                    uint8
	DefaultLinkPolicySettings uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadDefaultLinkPolicySettingsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteDefaultLinkPolicySettings implements Write Default Link Policy Settings (0x02|0x000D) [Vol 2, Part E, 7.2.12]
type WriteDefaultLinkPolicySettings struct {
	DefaultLinkPolicySettings uint16
}

func (c *WriteDefaultLinkPolicySettings) String() string {
	return "Write Default Link Policy Settings (0x02|0x000D)"
}

// OpCode returns the opcode of the command.
func (c *WriteDefaultLinkPolicySettings) OpCode() int { return 0x02<<10 | 0x000D }

// Len returns the length of the command.
func (c *WriteDefaultLinkPolicySettings) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *WriteDefaultLinkPolicySettings) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteDefaultLinkPolicySettingsRP returns the return parameter of Write Default Link Policy Settings
type WriteDefaultLinkPolicySettingsRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteDefaultLinkPolicySettingsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// FlowSpecification implements Flow Specification (0x02|0x0010) [Vol 2, Part E, 7.2.13]
type FlowSpecification struct {
	ConnectionHandle uint16
	Flags            uint8
	FlowDirection    uint8
	ServiceType      uint8
	TokenRate        uint32
	TokenBucketSize  uint32
	PeakBandwidth    uint32
	AccessLatency    uint32
}

func (c *FlowSpecification) String() string {
	return "Flow Specification (0x02|0x0010)"
}

// OpCode returns the opcode of the command.
func (c *FlowSpecification) OpCode() int { return 0x02<<10 | 0x0010 }

// Len returns the length of the command.
func (c *FlowSpecification) Len() int { return 21 }

// Marshal serializes the command parameters into binary form.
func (c *FlowSpecification) Marshal(b []byte) error {
	return marshal(c, b)
}

// SniffSubrating implements Sniff Subrating (0x02|0x0011) [Vol 2, Part E, 7.2.14]
type SniffSubrating struct {
	ConnectionHandle     uint16
	MaximumLatency       uint16
	MinimumRemoteLatency uint16
	MinimumLocalTimeout  uint16
}

func (c *SniffSubrating) String() string {
	return "Sniff Subrating (0x02|0x0011)"
}

// OpCode returns the opcode of the command.
func (c *SniffSubrating) OpCode() int { return 0x02<<10 | 0x0011 }

// Len returns the length of the command.
func (c *SniffSubrating) Len() int { return 8 }

// Marshal serializes the command parameters into binary form.
func (c *SniffSubrating) Marshal(b []byte) error {
	return marshal(c, b)
}

// SniffSubratingRP returns the return parameter of Sniff Subrating
type SniffSubratingRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *SniffSubratingRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// SetEventMask implements Set Event Mask (0x03|0x0001) [Vol 2, Part E, 7.3.1]
type SetEventMask struct {
	EventMask uint64
}

func (c *SetEventMask) String() string {
	return "Set Event Mask (0x03|0x0001)"
}

// OpCode returns the opcode of the command.
func (c *SetEventMask) OpCode() int { return 0x03<<10 | 0x0001 }

// Len returns the length of the command.
func (c *SetEventMask) Len() int { return 8 }

// Marshal serializes the command parameters into binary form.
func (c *SetEventMask) Marshal(b []byte) error {
	return marshal(c, b)
}

// SetEventMaskRP returns the return parameter of Set Event Mask
type SetEventMaskRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *SetEventMaskRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// Reset implements Reset (0x03|0x0003) [Vol 2, Part E, 7.3.2]
type Reset struct {
}

func (c *Reset) String() string {
	return "Reset (0x03|0x0003)"
}

// OpCode returns the opcode of the command.
func (c *Reset) OpCode() int { return 0x03<<10 | 0x0003 }

// Len returns the length of the command.
func (c *Reset) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *Reset) Marshal(b []byte) error {
	return marshal(c, b)
}

// ResetRP returns the return parameter of Reset
type ResetRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ResetRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// SetEventFilter implements Set Event Filter (0x03|0x0005) [Vol 2, Part E, 7.3.3]
type SetEventFilter struct {
	FilterType          uint8
	FilterConditionType uint8
	Condition           [7]byte
}

func (c *SetEventFilter) String() string {
	return "Set Event Filter (0x03|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *SetEventFilter) OpCode() int { return 0x03<<10 | 0x0005 }

// Len returns the length of the command.
func (c *SetEventFilter) Len() int { return 9 }

// Marshal serializes the command parameters into binary form.
func (c *SetEventFilter) Marshal(b []byte) error {
	return marshal(c, b)
}

// SetEventFilterRP returns the return parameter of Set Event Filter
type SetEventFilterRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *SetEventFilterRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// Flush implements Flush (0x03|0x0008) [Vol 2, Part E, 7.3.4]
type Flush struct {
	ConnectionHandle uint16
}

func (c *Flush) String() string {
	return "Flush (0x03|0x0008)"
}

// OpCode returns the opcode of the command.
func (c *Flush) OpCode() int { return 0x03<<10 | 0x0008 }

// Len returns the length of the command.
func (c *Flush) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *Flush) Marshal(b []byte) error {
	return marshal(c, b)
}

// FlushRP returns the return parameter of Flush
type FlushRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *FlushRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadPINType implements Read PIN Type (0x03|0x0009) [Vol 2, Part E, 7.3.5]
type ReadPINType struct {
}

func (c *ReadPINType) String() string {
	return "Read PIN Type (0x03|0x0009)"
}

// OpCode returns the opcode of the command.
func (c *ReadPINType) OpCode() int { return 0x03<<10 | 0x0009 }

// Len returns the length of the command.
func (c *ReadPINType) Len() int { return -1 }

// ReadPINTypeRP returns the return parameter of Read PIN Type
type ReadPINTypeRP struct {
	Status  uint8
	PINType uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadPINTypeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WritePINType implements Write PIN Type (0x03|0x000A) [Vol 2, Part E, 7.3.6]
type WritePINType struct {
	PINType uint8
}

func (c *WritePINType) String() string {
	return "Write PIN Type (0x03|0x000A)"
}

// OpCode returns the opcode of the command.
func (c *WritePINType) OpCode() int { return 0x03<<10 | 0x000A }

// Len returns the length of the command.
func (c *WritePINType) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *WritePINType) Marshal(b []byte) error {
	return marshal(c, b)
}

// WritePINTypeRP returns the return parameter of Write PIN Type
type WritePINTypeRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WritePINTypeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// CreateNewUnitKey implements Create New Unit Key (0x03|0x000B) [Vol 2, Part E, 7.3.7]
type CreateNewUnitKey struct {
}

func (c *CreateNewUnitKey) String() string {
	return "Create New Unit Key (0x03|0x000B)"
}

// OpCode returns the opcode of the command.
func (c *CreateNewUnitKey) OpCode() int { return 0x03<<10 | 0x000B }

// Len returns the length of the command.
func (c *CreateNewUnitKey) Len() int { return -1 }

// CreateNewUnitKeyRP returns the return parameter of Create New Unit Key
type CreateNewUnitKeyRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *CreateNewUnitKeyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadStoredLinkKey implements Read Stored Link Key (0x03|0x000D) [Vol 2, Part E, 7.3.8]
type ReadStoredLinkKey struct {
	BDADDR      [6]byte
	ReadAllFlag uint8
}

func (c *ReadStoredLinkKey) String() string {
	return "Read Stored Link Key (0x03|0x000D)"
}

// OpCode returns the opcode of the command.
func (c *ReadStoredLinkKey) OpCode() int { return 0x03<<10 | 0x000D }

// Len returns the length of the command.
func (c *ReadStoredLinkKey) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *ReadStoredLinkKey) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadStoredLinkKeyRP returns the return parameter of Read Stored Link Key
type ReadStoredLinkKeyRP struct {
	Status      uint8
	MaxNumKeys  uint16
	NumKeysRead uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadStoredLinkKeyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteStoredLinkKey implements Write Stored Link Key (0x03|0x0011) [Vol 2, Part E, 7.3.9]
type WriteStoredLinkKey struct {
	NumKeysToWrite uint8
	BDADDR         [6]byte
	LinkKey        [16]byte
}

func (c *WriteStoredLinkKey) String() string {
	return "Write Stored Link Key (0x03|0x0011)"
}

// OpCode returns the opcode of the command.
func (c *WriteStoredLinkKey) OpCode() int { return 0x03<<10 | 0x0011 }

// Len returns the length of the command.
func (c *WriteStoredLinkKey) Len() int { return 23 }

// Marshal serializes the command parameters into binary form.
func (c *WriteStoredLinkKey) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteStoredLinkKeyRP returns the return parameter of Write Stored Link Key
type WriteStoredLinkKeyRP struct {
	Status         uint8
	NumKeysWritten uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteStoredLinkKeyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// DeleteStoredLinkKey implements Delete Stored Link Key (0x03|0x0012) [Vol 2, Part E, 7.3.10]
type DeleteStoredLinkKey struct {
	BDADDR        [6]byte
	DeleteAllFlag uint8
}

func (c *DeleteStoredLinkKey) String() string {
	return "Delete Stored Link Key (0x03|0x0012)"
}

// OpCode returns the opcode of the command.
func (c *DeleteStoredLinkKey) OpCode() int { return 0x03<<10 | 0x0012 }

// Len returns the length of the command.
func (c *DeleteStoredLinkKey) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *DeleteStoredLinkKey) Marshal(b []byte) error {
	return marshal(c, b)
}

// DeleteStoredLinkKeyRP returns the return parameter of Delete Stored Link Key
type DeleteStoredLinkKeyRP struct {
	Status         uint8
	NumKeysDeleted uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *DeleteStoredLinkKeyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteLocalName implements Write Local Name (0x03|0x0013) [Vol 2, Part E, 7.3.11]
type WriteLocalName struct {
	LocalName [248]byte
}

func (c *WriteLocalName) String() string {
	return "Write Local Name (0x03|0x0013)"
}

// OpCode returns the opcode of the command.
func (c *WriteLocalName) OpCode() int { return 0x03<<10 | 0x0013 }

// Len returns the length of the command.
func (c *WriteLocalName) Len() int { return 248 }

// Marshal serializes the command parameters into binary form.
func (c *WriteLocalName) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteLocalNameRP returns the return parameter of Write Local Name
type WriteLocalNameRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteLocalNameRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadLocalName implements Read Local Name (0x03|0x0014) [Vol 2, Part E, 7.3.12]
type ReadLocalName struct {
}

func (c *ReadLocalName) String() string {
	return "Read Local Name (0x03|0x0014)"
}

// OpCode returns the opcode of the command.
func (c *ReadLocalName) OpCode() int { return 0x03<<10 | 0x0014 }

// Len returns the length of the command.
func (c *ReadLocalName) Len() int { return -1 }

// ReadLocalNameRP returns the return parameter of Read Local Name
type ReadLocalNameRP struct {
	Status    uint8
	LocalName [248]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadLocalNameRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadConnectionAcceptTimeout implements Read Connection Accept Timeout (0x03|0x0015) [Vol 2, Part E, 7.3.13]
type ReadConnectionAcceptTimeout struct {
}

func (c *ReadConnectionAcceptTimeout) String() string {
	return "Read Connection Accept Timeout (0x03|0x0015)"
}

// OpCode returns the opcode of the command.
func (c *ReadConnectionAcceptTimeout) OpCode() int { return 0x03<<10 | 0x0015 }

// Len returns the length of the command.
func (c *ReadConnectionAcceptTimeout) Len() int { return -1 }

// ReadConnectionAcceptTimeoutRP returns the return parameter of Read Connection Accept Timeout
type ReadConnectionAcceptTimeoutRP struct {
	Status            uint8
	ConnAcceptTimeout uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadConnectionAcceptTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteConnectionAcceptTimeout implements Write Connection Accept Timeout (0x03|0x0016) [Vol 2, Part E, 7.3.14]
type WriteConnectionAcceptTimeout struct {
	ConnAcceptTimeout uint16
}

func (c *WriteConnectionAcceptTimeout) String() string {
	return "Write Connection Accept Timeout (0x03|0x0016)"
}

// OpCode returns the opcode of the command.
func (c *WriteConnectionAcceptTimeout) OpCode() int { return 0x03<<10 | 0x0016 }

// Len returns the length of the command.
func (c *WriteConnectionAcceptTimeout) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *WriteConnectionAcceptTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteConnectionAcceptTimeoutRP returns the return parameter of Write Connection Accept Timeout
type WriteConnectionAcceptTimeoutRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteConnectionAcceptTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadPageTimeout implements Read Page Timeout (0x03|0x0017) [Vol 2, Part E, 7.3.15]
type ReadPageTimeout struct {
}

func (c *ReadPageTimeout) String() string {
	return "Read Page Timeout (0x03|0x0017)"
}

// OpCode returns the opcode of the command.
func (c *ReadPageTimeout) OpCode() int { return 0x03<<10 | 0x0017 }

// Len returns the length of the command.
func (c *ReadPageTimeout) Len() int { return -1 }

// ReadPageTimeoutRP returns the return parameter of Read Page Timeout
type ReadPageTimeoutRP struct {
	Status      uint8
	PageTimeout uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadPageTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WritePageTimeout implements Write Page Timeout (0x03|0x0018) [Vol 2, Part E, 7.3.16]
type WritePageTimeout struct {
	PageTimeout uint16
}

func (c *WritePageTimeout) String() string {
	return "Write Page Timeout (0x03|0x0018)"
}

// OpCode returns the opcode of the command.
func (c *WritePageTimeout) OpCode() int { return 0x03<<10 | 0x0018 }

// Len returns the length of the command.
func (c *WritePageTimeout) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *WritePageTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// WritePageTimeoutRP returns the return parameter of Write Page Timeout
type WritePageTimeoutRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WritePageTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadScanEnable implements Read Scan Enable (0x03|0x0019) [Vol 2, Part E, 7.3.17]
type ReadScanEnable struct {
}

func (c *ReadScanEnable) String() string {
	return "Read Scan Enable (0x03|0x0019)"
}

// OpCode returns the opcode of the command.
func (c *ReadScanEnable) OpCode() int { return 0x03<<10 | 0x0019 }

// Len returns the length of the command.
func (c *ReadScanEnable) Len() int { return -1 }

// ReadScanEnableRP returns the return parameter of Read Scan Enable
type ReadScanEnableRP struct {
	Status     uint8
	ScanEnable uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadScanEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteScanEnable implements Write Scan Enable (0x03|0x001A) [Vol 2, Part E, 7.3.18]
type WriteScanEnable struct {
	ScanEnable uint8
}

func (c *WriteScanEnable) String() string {
	return "Write Scan Enable (0x03|0x001A)"
}

// OpCode returns the opcode of the command.
func (c *WriteScanEnable) OpCode() int { return 0x03<<10 | 0x001A }

// Len returns the length of the command.
func (c *WriteScanEnable) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *WriteScanEnable) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteScanEnableRP returns the return parameter of Write Scan Enable
type WriteScanEnableRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteScanEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadPageScanActivity implements Read Page Scan Activity (0x03|0x001B) [Vol 2, Part E, 7.3.19]
type ReadPageScanActivity struct {
}

func (c *ReadPageScanActivity) String() string {
	return "Read Page Scan Activity (0x03|0x001B)"
}

// OpCode returns the opcode of the command.
func (c *ReadPageScanActivity) OpCode() int { return 0x03<<10 | 0x001B }

// Len returns the length of the command.
func (c *ReadPageScanActivity) Len() int { return -1 }

// ReadPageScanActivityRP returns the return parameter of Read Page Scan Activity
type ReadPageScanActivityRP struct {
	Status           uint8
	PageScanInterval uint16
	PageScanWindow   uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadPageScanActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WritePageScanActivity implements Write Page Scan Activity (0x03|0x001C) [Vol 2, Part E, 7.3.20]
type WritePageScanActivity struct {
	PageScanInterval uint16
	PageScanWindow   uint16
}

func (c *WritePageScanActivity) String() string {
	return "Write Page Scan Activity (0x03|0x001C)"
}

// OpCode returns the opcode of the command.
func (c *WritePageScanActivity) OpCode() int { return 0x03<<10 | 0x001C }

// Len returns the length of the command.
func (c *WritePageScanActivity) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *WritePageScanActivity) Marshal(b []byte) error {
	return marshal(c, b)
}

// WritePageScanActivityRP returns the return parameter of Write Page Scan Activity
type WritePageScanActivityRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WritePageScanActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadInquiryScanActivity implements Read Inquiry Scan Activity (0x03|0x001D) [Vol 2, Part E, 7.3.21]
type ReadInquiryScanActivity struct {
}

func (c *ReadInquiryScanActivity) String() string {
	return "Read Inquiry Scan Activity (0x03|0x001D)"
}

// OpCode returns the opcode of the command.
func (c *ReadInquiryScanActivity) OpCode() int { return 0x03<<10 | 0x001D }

// Len returns the length of the command.
func (c *ReadInquiryScanActivity) Len() int { return -1 }

// ReadInquiryScanActivityRP returns the return parameter of Read Inquiry Scan Activity
type ReadInquiryScanActivityRP struct {
	Status              uint8
	InquiryScanInterval uint16
	InquiryScanWindow   uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadInquiryScanActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteInquiryScanActivity implements Write Inquiry Scan Activity (0x03|0x001E) [Vol 2, Part E, 7.3.22]
type WriteInquiryScanActivity struct {
	InquiryScanInterval uint16
	InquiryScanWindow   uint16
}

func (c *WriteInquiryScanActivity) String() string {
	return "Write Inquiry Scan Activity (0x03|0x001E)"
}

// OpCode returns the opcode of the command.
func (c *WriteInquiryScanActivity) OpCode() int { return 0x03<<10 | 0x001E }

// Len returns the length of the command.
func (c *WriteInquiryScanActivity) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *WriteInquiryScanActivity) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteInquiryScanActivityRP returns the return parameter of Write Inquiry Scan Activity
type WriteInquiryScanActivityRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteInquiryScanActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadAuthenticationEnable implements Read Authentication Enable (0x03|0x001F) [Vol 2, Part E, 7.3.23]
type ReadAuthenticationEnable struct {
}

func (c *ReadAuthenticationEnable) String() string {
	return "Read Authentication Enable (0x03|0x001F)"
}

// OpCode returns the opcode of the command.
func (c *ReadAuthenticationEnable) OpCode() int { return 0x03<<10 | 0x001F }

// Len returns the length of the command.
func (c *ReadAuthenticationEnable) Len() int { return -1 }

// ReadAuthenticationEnableRP returns the return parameter of Read Authentication Enable
type ReadAuthenticationEnableRP struct {
	Status               uint8
	AuthenticationEnable uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadAuthenticationEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteAuthenticationEnable implements Write Authentication Enable (0x03|0x0020) [Vol 2, Part E, 7.3.24]
type WriteAuthenticationEnable struct {
	AuthenticationEnable uint8
}

func (c *WriteAuthenticationEnable) String() string {
	return "Write Authentication Enable (0x03|0x0020)"
}

// OpCode returns the opcode of the command.
func (c *WriteAuthenticationEnable) OpCode() int { return 0x03<<10 | 0x0020 }

// Len returns the length of the command.
func (c *WriteAuthenticationEnable) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *WriteAuthenticationEnable) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteAuthenticationEnableRP returns the return parameter of Write Authentication Enable
type WriteAuthenticationEnableRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteAuthenticationEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadClassOfDevice implements Read Class Of Device (0x03|0x0023) [Vol 2, Part E, 7.3.25]
type ReadClassOfDevice struct {
}

func (c *ReadClassOfDevice) String() string {
	return "Read Class Of Device (0x03|0x0023)"
}

// OpCode returns the opcode of the command.
func (c *ReadClassOfDevice) OpCode() int { return 0x03<<10 | 0x0023 }

// Len returns the length of the command.
func (c *ReadClassOfDevice) Len() int { return -1 }

// ReadClassOfDeviceRP returns the return parameter of Read Class Of Device
type ReadClassOfDeviceRP struct {
	Status        uint8
	ClassOfDevice [3]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadClassOfDeviceRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteClassOfDevice implements Write Class Of Device (0x03|0x0024) [Vol 2, Part E, 7.3.26]
type WriteClassOfDevice struct {
	ClassOfDevice [3]byte
}

func (c *WriteClassOfDevice) String() string {
	return "Write Class Of Device (0x03|0x0024)"
}

// OpCode returns the opcode of the command.
func (c *WriteClassOfDevice) OpCode() int { return 0x03<<10 | 0x0024 }

// Len returns the length of the command.
func (c *WriteClassOfDevice) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *WriteClassOfDevice) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteClassOfDeviceRP returns the return parameter of Write Class Of Device
type WriteClassOfDeviceRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteClassOfDeviceRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadVoiceSetting implements Read Voice Setting (0x03|0x0025) [Vol 2, Part E, 7.3.27]
type ReadVoiceSetting struct {
}

func (c *ReadVoiceSetting) String() string {
	return "Read Voice Setting (0x03|0x0025)"
}

// OpCode returns the opcode of the command.
func (c *ReadVoiceSetting) OpCode() int { return 0x03<<10 | 0x0025 }

// Len returns the length of the command.
func (c *ReadVoiceSetting) Len() int { return -1 }

// ReadVoiceSettingRP returns the return parameter of Read Voice Setting
type ReadVoiceSettingRP struct {
	Status       uint8
	VoiceSetting uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadVoiceSettingRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteVoiceSetting implements Write Voice Setting (0x03|0x0026) [Vol 2, Part E, 7.3.28]
type WriteVoiceSetting struct {
	VoiceSetting uint16
}

func (c *WriteVoiceSetting) String() string {
	return "Write Voice Setting (0x03|0x0026)"
}

// OpCode returns the opcode of the command.
func (c *WriteVoiceSetting) OpCode() int { return 0x03<<10 | 0x0026 }

// Len returns the length of the command.
func (c *WriteVoiceSetting) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *WriteVoiceSetting) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteVoiceSettingRP returns the return parameter of Write Voice Setting
type WriteVoiceSettingRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteVoiceSettingRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadAutomaticFlushTimeout implements Read Automatic Flush Timeout (0x03|0x0027) [Vol 2, Part E, 7.3.29]
type ReadAutomaticFlushTimeout struct {
	ConnectionHandle uint16
}

func (c *ReadAutomaticFlushTimeout) String() string {
	return "Read Automatic Flush Timeout (0x03|0x0027)"
}

// OpCode returns the opcode of the command.
func (c *ReadAutomaticFlushTimeout) OpCode() int { return 0x03<<10 | 0x0027 }

// Len returns the length of the command.
func (c *ReadAutomaticFlushTimeout) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadAutomaticFlushTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadAutomaticFlushTimeoutRP returns the return parameter of Read Automatic Flush Timeout
type ReadAutomaticFlushTimeoutRP struct {
	Status           uint8
	ConnectionHandle uint16
	FlushTimeout     uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadAutomaticFlushTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteAutomaticFlushTimeout implements Write Automatic Flush Timeout (0x03|0x0028) [Vol 2, Part E, 7.3.30]
type WriteAutomaticFlushTimeout struct {
	ConnectionHandle uint16
	FlushTimeout     uint16
}

func (c *WriteAutomaticFlushTimeout) String() string {
	return "Write Automatic Flush Timeout (0x03|0x0028)"
}

// OpCode returns the opcode of the command.
func (c *WriteAutomaticFlushTimeout) OpCode() int { return 0x03<<10 | 0x0028 }

// Len returns the length of the command.
func (c *WriteAutomaticFlushTimeout) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *WriteAutomaticFlushTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteAutomaticFlushTimeoutRP returns the return parameter of Write Automatic Flush Timeout
type WriteAutomaticFlushTimeoutRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteAutomaticFlushTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadNumBroadcastRetransmissions implements Read Num Broadcast Retransmissions (0x03|0x0029) [Vol 2, Part E, 7.3.31]
type ReadNumBroadcastRetransmissions struct {
}

func (c *ReadNumBroadcastRetransmissions) String() string {
	return "Read Num Broadcast Retransmissions (0x03|0x0029)"
}

// OpCode returns the opcode of the command.
func (c *ReadNumBroadcastRetransmissions) OpCode() int { return 0x03<<10 | 0x0029 }

// Len returns the length of the command.
func (c *ReadNumBroadcastRetransmissions) Len() int { return -1 }

// ReadNumBroadcastRetransmissionsRP returns the return parameter of Read Num Broadcast Retransmissions
type ReadNumBroadcastRetransmissionsRP struct {
	Status                      uint8
	NumBroadcastRetransmissions uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadNumBroadcastRetransmissionsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteNumBroadcastRetransmissions implements Write Num Broadcast Retransmissions (0x03|0x002A) [Vol 2, Part E, 7.3.32]
type WriteNumBroadcastRetransmissions struct {
	NumBroadcastRetransmissions uint8
}

func (c *WriteNumBroadcastRetransmissions) String() string {
	return "Write Num Broadcast Retransmissions (0x03|0x002A)"
}

// OpCode returns the opcode of the command.
func (c *WriteNumBroadcastRetransmissions) OpCode() int { return 0x03<<10 | 0x002A }

// Len returns the length of the command.
func (c *WriteNumBroadcastRetransmissions) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *WriteNumBroadcastRetransmissions) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteNumBroadcastRetransmissionsRP returns the return parameter of Write Num Broadcast Retransmissions
type WriteNumBroadcastRetransmissionsRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteNumBroadcastRetransmissionsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadHoldModeActivity implements Read Hold Mode Activity (0x03|0x002B) [Vol 2, Part E, 7.3.33]
type ReadHoldModeActivity struct {
}

func (c *ReadHoldModeActivity) String() string {
	return "Read Hold Mode Activity (0x03|0x002B)"
}

// OpCode returns the opcode of the command.
func (c *ReadHoldModeActivity) OpCode() int { return 0x03<<10 | 0x002B }

// Len returns the length of the command.
func (c *ReadHoldModeActivity) Len() int { return -1 }

// ReadHoldModeActivityRP returns the return parameter of Read Hold Mode Activity
type ReadHoldModeActivityRP struct {
	Status           uint8
	HoldModeActivity uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadHoldModeActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteHoldModeActivity implements Write Hold Mode Activity (0x03|0x002C) [Vol 2, Part E, 7.3.34]
type WriteHoldModeActivity struct {
	HoldModeActivity uint8
}

func (c *WriteHoldModeActivity) String() string {
	return "Write Hold Mode Activity (0x03|0x002C)"
}

// OpCode returns the opcode of the command.
func (c *WriteHoldModeActivity) OpCode() int { return 0x03<<10 | 0x002C }

// Len returns the length of the command.
func (c *WriteHoldModeActivity) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *WriteHoldModeActivity) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteHoldModeActivityRP returns the return parameter of Write Hold Mode Activity
type WriteHoldModeActivityRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteHoldModeActivityRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadTransmitPowerLevel implements Read Transmit Power Level (0x03|0x002D) [Vol 2, Part E, 7.3.35]
type ReadTransmitPowerLevel struct {
	ConnectionHandle uint16
	Type             uint8
}

func (c *ReadTransmitPowerLevel) String() string {
	return "Read Transmit Power Level (0x03|0x002D)"
}

// OpCode returns the opcode of the command.
func (c *ReadTransmitPowerLevel) OpCode() int { return 0x03<<10 | 0x002D }

// Len returns the length of the command.
func (c *ReadTransmitPowerLevel) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *ReadTransmitPowerLevel) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadTransmitPowerLevelRP returns the return parameter of Read Transmit Power Level
type ReadTransmitPowerLevelRP struct {
	Status             uint8
	ConnectionHandle   uint16
	TransmitPowerLevel uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadTransmitPowerLevelRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadSynchronousFlowControlEnable implements Read Synchronous Flow Control Enable (0x03|0x002E) [Vol 2, Part E, 7.3.36]
type ReadSynchronousFlowControlEnable struct {
}

func (c *ReadSynchronousFlowControlEnable) String() string {
	return "Read Synchronous Flow Control Enable (0x03|0x002E)"
}

// OpCode returns the opcode of the command.
func (c *ReadSynchronousFlowControlEnable) OpCode() int { return 0x03<<10 | 0x002E }

// Len returns the length of the command.
func (c *ReadSynchronousFlowControlEnable) Len() int { return -1 }

// ReadSynchronousFlowControlEnableRP returns the return parameter of Read Synchronous Flow Control Enable
type ReadSynchronousFlowControlEnableRP struct {
	Status                       uint8
	SynchronousFlowControlEnable uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadSynchronousFlowControlEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteSynchronousFlowControlEnable implements Write Synchronous Flow Control Enable (0x03|0x002F) [Vol 2, Part E, 7.3.37]
type WriteSynchronousFlowControlEnable struct {
	SynchronousFlowControlEnable uint8
}

func (c *WriteSynchronousFlowControlEnable) String() string {
	return "Write Synchronous Flow Control Enable (0x03|0x002F)"
}

// OpCode returns the opcode of the command.
func (c *WriteSynchronousFlowControlEnable) OpCode() int { return 0x03<<10 | 0x002F }

// Len returns the length of the command.
func (c *WriteSynchronousFlowControlEnable) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *WriteSynchronousFlowControlEnable) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteSynchronousFlowControlEnableRP returns the return parameter of Write Synchronous Flow Control Enable
type WriteSynchronousFlowControlEnableRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteSynchronousFlowControlEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// SetControllertoHostFlowControl implements Set Controller to Host Flow Control (0x03|0x0031) [Vol 2, Part E, 7.3.38]
type SetControllertoHostFlowControl struct {
	FlowControlEnable uint8
}

func (c *SetControllertoHostFlowControl) String() string {
	return "Set Controller to Host Flow Control (0x03|0x0031)"
}

// OpCode returns the opcode of the command.
func (c *SetControllertoHostFlowControl) OpCode() int { return 0x03<<10 | 0x0031 }

// Len returns the length of the command.
func (c *SetControllertoHostFlowControl) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *SetControllertoHostFlowControl) Marshal(b []byte) error {
	return marshal(c, b)
}

// SetControllertoHostFlowControlRP returns the return parameter of Set Controller to Host Flow Control
type SetControllertoHostFlowControlRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *SetControllertoHostFlowControlRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// HostBufferSize implements Host Buffer Size (0x03|0x0033) [Vol 2, Part E, 7.3.39]
type HostBufferSize struct {
	HostACLDataPacketLength            uint16
	HostSynchronousDataPacketLength    uint8
	HostTotalNumACLDataPackets         uint16
	HostTotalNumSynchronousDataPackets uint16
}

func (c *HostBufferSize) String() string {
	return "Host Buffer Size (0x03|0x0033)"
}

// OpCode returns the opcode of the command.
func (c *HostBufferSize) OpCode() int { return 0x03<<10 | 0x0033 }

// Len returns the length of the command.
func (c *HostBufferSize) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *HostBufferSize) Marshal(b []byte) error {
	return marshal(c, b)
}

// HostBufferSizeRP returns the return parameter of Host Buffer Size
type HostBufferSizeRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *HostBufferSizeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// HostNumberOfCompletedPackets implements Host Number Of Completed Packets (0x03|0x0035) [Vol 2, Part E, 7.3.40]
type HostNumberOfCompletedPackets struct {
	NumberOfHandles           uint8
	ConnectionHandle          []uint16
	HostNumOfCompletedPackets []uint16
}

func (c *HostNumberOfCompletedPackets) String() string {
	return "Host Number Of Completed Packets (0x03|0x0035)"
}

// OpCode returns the opcode of the command.
func (c *HostNumberOfCompletedPackets) OpCode() int { return 0x03<<10 | 0x0035 }

// Len returns the length of the command.
func (c *HostNumberOfCompletedPackets) Len() int { return -1 }

// SetEventMaskPage2 implements Set Event Mask Page 2 (0x03|0x0063) [Vol 2, Part E, 7.3.69]
type SetEventMaskPage2 struct {
	EventMaskPage2 uint64
}

func (c *SetEventMaskPage2) String() string {
	return "Set Event Mask Page 2 (0x03|0x0063)"
}

// OpCode returns the opcode of the command.
func (c *SetEventMaskPage2) OpCode() int { return 0x03<<10 | 0x0063 }

// Len returns the length of the command.
func (c *SetEventMaskPage2) Len() int { return 8 }

// Marshal serializes the command parameters into binary form.
func (c *SetEventMaskPage2) Marshal(b []byte) error {
	return marshal(c, b)
}

// SetEventMaskPage2RP returns the return parameter of Set Event Mask Page 2
type SetEventMaskPage2RP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *SetEventMaskPage2RP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteLEHostSupport implements Write LE Host Support (0x03|0x006D) [Vol 2, Part E, 7.3.79]
type WriteLEHostSupport struct {
	LESupportedHost    uint8
	SimultaneousLEHost uint8
}

func (c *WriteLEHostSupport) String() string {
	return "Write LE Host Support (0x03|0x006D)"
}

// OpCode returns the opcode of the command.
func (c *WriteLEHostSupport) OpCode() int { return 0x03<<10 | 0x006D }

// Len returns the length of the command.
func (c *WriteLEHostSupport) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *WriteLEHostSupport) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteLEHostSupportRP returns the return parameter of Write LE Host Support
type WriteLEHostSupportRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteLEHostSupportRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadAuthenticatedPayloadTimeout implements Read Authenticated Payload Timeout (0x03|0x007B) [Vol 2, Part E, 7.3.93]
type ReadAuthenticatedPayloadTimeout struct {
	ConnectionHandle uint16
}

func (c *ReadAuthenticatedPayloadTimeout) String() string {
	return "Read Authenticated Payload Timeout (0x03|0x007B)"
}

// OpCode returns the opcode of the command.
func (c *ReadAuthenticatedPayloadTimeout) OpCode() int { return 0x03<<10 | 0x007B }

// Len returns the length of the command.
func (c *ReadAuthenticatedPayloadTimeout) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadAuthenticatedPayloadTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadAuthenticatedPayloadTimeoutRP returns the return parameter of Read Authenticated Payload Timeout
type ReadAuthenticatedPayloadTimeoutRP struct {
	Status                      uint8
	ConnectionHandle            uint16
	AuthenticatedPayloadTimeout uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadAuthenticatedPayloadTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// WriteAuthenticatedPayloadTimeout implements Write Authenticated Payload Timeout (0x01|0x007C) [Vol 2, Part E, 7.3.94]
type WriteAuthenticatedPayloadTimeout struct {
	ConnectionHandle            uint16
	AuthenticatedPayloadTimeout uint16
}

func (c *WriteAuthenticatedPayloadTimeout) String() string {
	return "Write Authenticated Payload Timeout (0x01|0x007C)"
}

// OpCode returns the opcode of the command.
func (c *WriteAuthenticatedPayloadTimeout) OpCode() int { return 0x01<<10 | 0x007C }

// Len returns the length of the command.
func (c *WriteAuthenticatedPayloadTimeout) Len() int { return 4 }

// Marshal serializes the command parameters into binary form.
func (c *WriteAuthenticatedPayloadTimeout) Marshal(b []byte) error {
	return marshal(c, b)
}

// WriteAuthenticatedPayloadTimeoutRP returns the return parameter of Write Authenticated Payload Timeout
type WriteAuthenticatedPayloadTimeoutRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *WriteAuthenticatedPayloadTimeoutRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadLocalVersionInformation implements Read Local Version Information (0x04|0x0001) [Vol 2, Part E, 7.4.1]
type ReadLocalVersionInformation struct {
}

func (c *ReadLocalVersionInformation) String() string {
	return "Read Local Version Information (0x04|0x0001)"
}

// OpCode returns the opcode of the command.
func (c *ReadLocalVersionInformation) OpCode() int { return 0x04<<10 | 0x0001 }

// Len returns the length of the command.
func (c *ReadLocalVersionInformation) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ReadLocalVersionInformation) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadLocalVersionInformationRP returns the return parameter of Read Local Version Information
type ReadLocalVersionInformationRP struct {
	Status           uint8
	HCIVersion       uint8
	HCIRevision      uint16
	LMPPAMVersion    uint8
	ManufacturerName uint16
	LMPPAMSubversion uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadLocalVersionInformationRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadLocalSupportedCommands implements Read Local Supported Commands (0x04|0x0002) [Vol 2, Part E, 7.4.2]
type ReadLocalSupportedCommands struct {
}

func (c *ReadLocalSupportedCommands) String() string {
	return "Read Local Supported Commands (0x04|0x0002)"
}

// OpCode returns the opcode of the command.
func (c *ReadLocalSupportedCommands) OpCode() int { return 0x04<<10 | 0x0002 }

// Len returns the length of the command.
func (c *ReadLocalSupportedCommands) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ReadLocalSupportedCommands) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadLocalSupportedCommandsRP returns the return parameter of Read Local Supported Commands
type ReadLocalSupportedCommandsRP struct {
	Status     uint8
	Supporteds uint64
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadLocalSupportedCommandsRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadLocalSupportedFeatures implements Read Local Supported Features (0x04|0x0003) [Vol 2, Part E, 7.4.3]
type ReadLocalSupportedFeatures struct {
}

func (c *ReadLocalSupportedFeatures) String() string {
	return "Read Local Supported Features (0x04|0x0003)"
}

// OpCode returns the opcode of the command.
func (c *ReadLocalSupportedFeatures) OpCode() int { return 0x04<<10 | 0x0003 }

// Len returns the length of the command.
func (c *ReadLocalSupportedFeatures) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ReadLocalSupportedFeatures) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadLocalSupportedFeaturesRP returns the return parameter of Read Local Supported Features
type ReadLocalSupportedFeaturesRP struct {
	Status      uint8
	LMPFeatures uint64
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadLocalSupportedFeaturesRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadBufferSize implements Read Buffer Size (0x04|0x0005) [Vol 2, Part E, 7.4.5]
type ReadBufferSize struct {
}

func (c *ReadBufferSize) String() string {
	return "Read Buffer Size (0x04|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *ReadBufferSize) OpCode() int { return 0x04<<10 | 0x0005 }

// Len returns the length of the command.
func (c *ReadBufferSize) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ReadBufferSize) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadBufferSizeRP returns the return parameter of Read Buffer Size
type ReadBufferSizeRP struct {
	Status                           uint8
	HCACLDataPacketLength            uint16
	HCSynchronousDataPacketLength    uint8
	HCTotalNumACLDataPackets         uint16
	HCTotalNumSynchronousDataPackets uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadBufferSizeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadBDADDR implements Read BD_ADDR (0x04|0x0009) [Vol 2, Part E, 7.4.6]
type ReadBDADDR struct {
}

func (c *ReadBDADDR) String() string {
	return "Read BD_ADDR (0x04|0x0009)"
}

// OpCode returns the opcode of the command.
func (c *ReadBDADDR) OpCode() int { return 0x04<<10 | 0x0009 }

// Len returns the length of the command.
func (c *ReadBDADDR) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *ReadBDADDR) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadBDADDRRP returns the return parameter of Read BD_ADDR
type ReadBDADDRRP struct {
	Status uint8
	BDADDR [6]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadBDADDRRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// ReadRSSI implements Read RSSI (0x05|0x0005) [Vol 2, Part E, 7.5.4]
type ReadRSSI struct {
	Handle uint16
}

func (c *ReadRSSI) String() string {
	return "Read RSSI (0x05|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *ReadRSSI) OpCode() int { return 0x05<<10 | 0x0005 }

// Len returns the length of the command.
func (c *ReadRSSI) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *ReadRSSI) Marshal(b []byte) error {
	return marshal(c, b)
}

// ReadRSSIRP returns the return parameter of Read RSSI
type ReadRSSIRP struct {
	Status           uint8
	ConnectionHandle uint16
	RSSI             int8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *ReadRSSIRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetEventMask implements LE Set Event Mask (0x08|0x0001) [Vol 2, Part E, 7.8.1]
type LESetEventMask struct {
	LEEventMask uint64
}

func (c *LESetEventMask) String() string {
	return "LE Set Event Mask (0x08|0x0001)"
}

// OpCode returns the opcode of the command.
func (c *LESetEventMask) OpCode() int { return 0x08<<10 | 0x0001 }

// Len returns the length of the command.
func (c *LESetEventMask) Len() int { return 8 }

// Marshal serializes the command parameters into binary form.
func (c *LESetEventMask) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetEventMaskRP returns the return parameter of LE Set Event Mask
type LESetEventMaskRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetEventMaskRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadBufferSize implements LE Read Buffer Size (0x08|0x0002) [Vol 2, Part E, 7.8.2]
type LEReadBufferSize struct {
}

func (c *LEReadBufferSize) String() string {
	return "LE Read Buffer Size (0x08|0x0002)"
}

// OpCode returns the opcode of the command.
func (c *LEReadBufferSize) OpCode() int { return 0x08<<10 | 0x0002 }

// Len returns the length of the command.
func (c *LEReadBufferSize) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadBufferSize) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadBufferSizeRP returns the return parameter of LE Read Buffer Size
type LEReadBufferSizeRP struct {
	Status                  uint8
	HCLEDataPacketLength    uint16
	HCTotalNumLEDataPackets uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadBufferSizeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadLocalSupportedFeatures implements LE Read Local Supported Features (0x08|0x0003) [Vol 2, Part E, 7.8.3]
type LEReadLocalSupportedFeatures struct {
}

func (c *LEReadLocalSupportedFeatures) String() string {
	return "LE Read Local Supported Features (0x08|0x0003)"
}

// OpCode returns the opcode of the command.
func (c *LEReadLocalSupportedFeatures) OpCode() int { return 0x08<<10 | 0x0003 }

// Len returns the length of the command.
func (c *LEReadLocalSupportedFeatures) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadLocalSupportedFeatures) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadLocalSupportedFeaturesRP returns the return parameter of LE Read Local Supported Features
type LEReadLocalSupportedFeaturesRP struct {
	Status     uint8
	LEFeatures uint64
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadLocalSupportedFeaturesRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetRandomAddress implements LE Set Random Address (0x08|0x0005) [Vol 2, Part E, 7.8.4]
type LESetRandomAddress struct {
	RandomAddress [6]byte
}

func (c *LESetRandomAddress) String() string {
	return "LE Set Random Address (0x08|0x0005)"
}

// OpCode returns the opcode of the command.
func (c *LESetRandomAddress) OpCode() int { return 0x08<<10 | 0x0005 }

// Len returns the length of the command.
func (c *LESetRandomAddress) Len() int { return 6 }

// Marshal serializes the command parameters into binary form.
func (c *LESetRandomAddress) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetRandomAddressRP returns the return parameter of LE Set Random Address
type LESetRandomAddressRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetRandomAddressRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetAdvertisingParameters implements LE Set Advertising Parameters (0x08|0x0006) [Vol 2, Part E, 7.8.5]
type LESetAdvertisingParameters struct {
	AdvertisingIntervalMin  uint16
	AdvertisingIntervalMax  uint16
	AdvertisingType         uint8
	OwnAddressType          uint8
	DirectAddressType       uint8
	DirectAddress           [6]byte
	AdvertisingChannelMap   uint8
	AdvertisingFilterPolicy uint8
}

func (c *LESetAdvertisingParameters) String() string {
	return "LE Set Advertising Parameters (0x08|0x0006)"
}

// OpCode returns the opcode of the command.
func (c *LESetAdvertisingParameters) OpCode() int { return 0x08<<10 | 0x0006 }

// Len returns the length of the command.
func (c *LESetAdvertisingParameters) Len() int { return 15 }

// Marshal serializes the command parameters into binary form.
func (c *LESetAdvertisingParameters) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetAdvertisingParametersRP returns the return parameter of LE Set Advertising Parameters
type LESetAdvertisingParametersRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetAdvertisingParametersRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadAdvertisingChannelTxPower implements LE Read Advertising Channel Tx Power (0x08|0x0007) [Vol 2, Part E, 7.8.6]
type LEReadAdvertisingChannelTxPower struct {
}

func (c *LEReadAdvertisingChannelTxPower) String() string {
	return "LE Read Advertising Channel Tx Power (0x08|0x0007)"
}

// OpCode returns the opcode of the command.
func (c *LEReadAdvertisingChannelTxPower) OpCode() int { return 0x08<<10 | 0x0007 }

// Len returns the length of the command.
func (c *LEReadAdvertisingChannelTxPower) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadAdvertisingChannelTxPower) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadAdvertisingChannelTxPowerRP returns the return parameter of LE Read Advertising Channel Tx Power
type LEReadAdvertisingChannelTxPowerRP struct {
	Status             uint8
	TransmitPowerLevel uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadAdvertisingChannelTxPowerRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetAdvertisingData implements LE Set Advertising Data (0x08|0x0008) [Vol 2, Part E, 7.8.7]
type LESetAdvertisingData struct {
	AdvertisingDataLength uint8
	AdvertisingData       [31]byte
}

func (c *LESetAdvertisingData) String() string {
	return "LE Set Advertising Data (0x08|0x0008)"
}

// OpCode returns the opcode of the command.
func (c *LESetAdvertisingData) OpCode() int { return 0x08<<10 | 0x0008 }

// Len returns the length of the command.
func (c *LESetAdvertisingData) Len() int { return 32 }

// Marshal serializes the command parameters into binary form.
func (c *LESetAdvertisingData) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetAdvertisingDataRP returns the return parameter of LE Set Advertising Data
type LESetAdvertisingDataRP struct {
	Status                  uint8
	HCLEDataPacketLength    uint16
	HCTotalNumLEDataPackets uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetAdvertisingDataRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetScanResponseData implements LE Set Scan Response Data (0x08|0x0009) [Vol 2, Part E, 7.8.8]
type LESetScanResponseData struct {
	ScanResponseDataLength uint8
	ScanResponseData       [31]byte
}

func (c *LESetScanResponseData) String() string {
	return "LE Set Scan Response Data (0x08|0x0009)"
}

// OpCode returns the opcode of the command.
func (c *LESetScanResponseData) OpCode() int { return 0x08<<10 | 0x0009 }

// Len returns the length of the command.
func (c *LESetScanResponseData) Len() int { return 32 }

// Marshal serializes the command parameters into binary form.
func (c *LESetScanResponseData) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetScanResponseDataRP returns the return parameter of LE Set Scan Response Data
type LESetScanResponseDataRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetScanResponseDataRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetAdvertiseEnable implements LE Set Advertise Enable (0x08|0x000A) [Vol 2, Part E, 7.8.9]
type LESetAdvertiseEnable struct {
	AdvertisingEnable uint8
}

func (c *LESetAdvertiseEnable) String() string {
	return "LE Set Advertise Enable (0x08|0x000A)"
}

// OpCode returns the opcode of the command.
func (c *LESetAdvertiseEnable) OpCode() int { return 0x08<<10 | 0x000A }

// Len returns the length of the command.
func (c *LESetAdvertiseEnable) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *LESetAdvertiseEnable) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetAdvertiseEnableRP returns the return parameter of LE Set Advertise Enable
type LESetAdvertiseEnableRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetAdvertiseEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetScanParameters implements LE Set Scan Parameters (0x08|0x000B) [Vol 2, Part E, 7.8.10]
type LESetScanParameters struct {
	LEScanType           uint8
	LEScanInterval       uint16
	LEScanWindow         uint16
	OwnAddressType       uint8
	ScanningFilterPolicy uint8
}

func (c *LESetScanParameters) String() string {
	return "LE Set Scan Parameters (0x08|0x000B)"
}

// OpCode returns the opcode of the command.
func (c *LESetScanParameters) OpCode() int { return 0x08<<10 | 0x000B }

// Len returns the length of the command.
func (c *LESetScanParameters) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *LESetScanParameters) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetScanParametersRP returns the return parameter of LE Set Scan Parameters
type LESetScanParametersRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetScanParametersRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LESetScanEnable implements LE Set Scan Enable (0x08|0x000C) [Vol 2, Part E, 7.8.11]
type LESetScanEnable struct {
	LEScanEnable     uint8
	FilterDuplicates uint8
}

func (c *LESetScanEnable) String() string {
	return "LE Set Scan Enable (0x08|0x000C)"
}

// OpCode returns the opcode of the command.
func (c *LESetScanEnable) OpCode() int { return 0x08<<10 | 0x000C }

// Len returns the length of the command.
func (c *LESetScanEnable) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *LESetScanEnable) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetScanEnableRP returns the return parameter of LE Set Scan Enable
type LESetScanEnableRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetScanEnableRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LECreateConnection implements LE Create Connection (0x08|0x000D) [Vol 2, Part E, 7.8.12]
type LECreateConnection struct {
	LEScanInterval        uint16
	LEScanWindow          uint16
	InitiatorFilterPolicy uint8
	PeerAddressType       uint8
	PeerAddress           [6]byte
	OwnAddressType        uint8
	ConnIntervalMin       uint16
	ConnIntervalMax       uint16
	ConnLatency           uint16
	SupervisionTimeout    uint16
	MinimumCELength       uint16
	MaximumCELength       uint16
}

func (c *LECreateConnection) String() string {
	return "LE Create Connection (0x08|0x000D)"
}

// OpCode returns the opcode of the command.
func (c *LECreateConnection) OpCode() int { return 0x08<<10 | 0x000D }

// Len returns the length of the command.
func (c *LECreateConnection) Len() int { return 25 }

// Marshal serializes the command parameters into binary form.
func (c *LECreateConnection) Marshal(b []byte) error {
	return marshal(c, b)
}

// LECreateConnectionCancel implements LE Create Connection Cancel (0x08|0x000E) [Vol 2, Part E, 7.8.13]
type LECreateConnectionCancel struct {
}

func (c *LECreateConnectionCancel) String() string {
	return "LE Create Connection Cancel (0x08|0x000E)"
}

// OpCode returns the opcode of the command.
func (c *LECreateConnectionCancel) OpCode() int { return 0x08<<10 | 0x000E }

// Len returns the length of the command.
func (c *LECreateConnectionCancel) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LECreateConnectionCancel) Marshal(b []byte) error {
	return marshal(c, b)
}

// LECreateConnectionCancelRP returns the return parameter of LE Create Connection Cancel
type LECreateConnectionCancelRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LECreateConnectionCancelRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadWhiteListSize implements LE Read White List Size (0x08|0x000F) [Vol 2, Part E, 7.8.14]
type LEReadWhiteListSize struct {
}

func (c *LEReadWhiteListSize) String() string {
	return "LE Read White List Size (0x08|0x000F)"
}

// OpCode returns the opcode of the command.
func (c *LEReadWhiteListSize) OpCode() int { return 0x08<<10 | 0x000F }

// Len returns the length of the command.
func (c *LEReadWhiteListSize) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadWhiteListSize) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadWhiteListSizeRP returns the return parameter of LE Read White List Size
type LEReadWhiteListSizeRP struct {
	Status        uint8
	WhiteListSize uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadWhiteListSizeRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEClearWhiteList implements LE Clear White List (0x08|0x0010) [Vol 2, Part E, 7.8.15]
type LEClearWhiteList struct {
}

func (c *LEClearWhiteList) String() string {
	return "LE Clear White List (0x08|0x0010)"
}

// OpCode returns the opcode of the command.
func (c *LEClearWhiteList) OpCode() int { return 0x08<<10 | 0x0010 }

// Len returns the length of the command.
func (c *LEClearWhiteList) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEClearWhiteList) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEClearWhiteListRP returns the return parameter of LE Clear White List
type LEClearWhiteListRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEClearWhiteListRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEAddDeviceToWhiteList implements LE Add Device To White List (0x08|0x0011) [Vol 2, Part E, 7.8.16]
type LEAddDeviceToWhiteList struct {
	AddressType uint8
	Address     [6]byte
}

func (c *LEAddDeviceToWhiteList) String() string {
	return "LE Add Device To White List (0x08|0x0011)"
}

// OpCode returns the opcode of the command.
func (c *LEAddDeviceToWhiteList) OpCode() int { return 0x08<<10 | 0x0011 }

// Len returns the length of the command.
func (c *LEAddDeviceToWhiteList) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *LEAddDeviceToWhiteList) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEAddDeviceToWhiteListRP returns the return parameter of LE Add Device To White List
type LEAddDeviceToWhiteListRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEAddDeviceToWhiteListRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LERemoveDeviceFromWhiteList implements LE Remove Device From White List (0x08|0x0012) [Vol 2, Part E, 7.8.17]
type LERemoveDeviceFromWhiteList struct {
	AddressType uint8
	Address     [6]byte
}

func (c *LERemoveDeviceFromWhiteList) String() string {
	return "LE Remove Device From White List (0x08|0x0012)"
}

// OpCode returns the opcode of the command.
func (c *LERemoveDeviceFromWhiteList) OpCode() int { return 0x08<<10 | 0x0012 }

// Len returns the length of the command.
func (c *LERemoveDeviceFromWhiteList) Len() int { return 7 }

// Marshal serializes the command parameters into binary form.
func (c *LERemoveDeviceFromWhiteList) Marshal(b []byte) error {
	return marshal(c, b)
}

// LERemoveDeviceFromWhiteListRP returns the return parameter of LE Remove Device From White List
type LERemoveDeviceFromWhiteListRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LERemoveDeviceFromWhiteListRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEConnectionUpdate implements LE Connection Update (0x08|0x0013) [Vol 2, Part E, 7.8.18]
type LEConnectionUpdate struct {
	ConnectionHandle   uint16
	ConnIntervalMin    uint16
	ConnIntervalMax    uint16
	ConnLatency        uint16
	SupervisionTimeout uint16
	MinimumCELength    uint16
	MaximumCELength    uint16
}

func (c *LEConnectionUpdate) String() string {
	return "LE Connection Update (0x08|0x0013)"
}

// OpCode returns the opcode of the command.
func (c *LEConnectionUpdate) OpCode() int { return 0x08<<10 | 0x0013 }

// Len returns the length of the command.
func (c *LEConnectionUpdate) Len() int { return 14 }

// Marshal serializes the command parameters into binary form.
func (c *LEConnectionUpdate) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetHostChannelClassification implements LE Set Host Channel Classification (0x08|0x0014) [Vol 2, Part E, 7.8.19]
type LESetHostChannelClassification struct {
	ChannelMap [5]byte
}

func (c *LESetHostChannelClassification) String() string {
	return "LE Set Host Channel Classification (0x08|0x0014)"
}

// OpCode returns the opcode of the command.
func (c *LESetHostChannelClassification) OpCode() int { return 0x08<<10 | 0x0014 }

// Len returns the length of the command.
func (c *LESetHostChannelClassification) Len() int { return 5 }

// Marshal serializes the command parameters into binary form.
func (c *LESetHostChannelClassification) Marshal(b []byte) error {
	return marshal(c, b)
}

// LESetHostChannelClassificationRP returns the return parameter of LE Set Host Channel Classification
type LESetHostChannelClassificationRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LESetHostChannelClassificationRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadChannelMap implements LE Read Channel Map (0x08|0x0015) [Vol 2, Part E, 7.8.20]
type LEReadChannelMap struct {
	ConnectionHandle uint16
}

func (c *LEReadChannelMap) String() string {
	return "LE Read Channel Map (0x08|0x0015)"
}

// OpCode returns the opcode of the command.
func (c *LEReadChannelMap) OpCode() int { return 0x08<<10 | 0x0015 }

// Len returns the length of the command.
func (c *LEReadChannelMap) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadChannelMap) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadChannelMapRP returns the return parameter of LE Read Channel Map
type LEReadChannelMapRP struct {
	Status           uint8
	ConnectionHandle uint16
	ChannelMap       [5]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadChannelMapRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadRemoteUsedFeatures implements LE Read Remote Used Features (0x08|0x0016) [Vol 2, Part E, 7.8.21]
type LEReadRemoteUsedFeatures struct {
	ConnectionHandle uint16
}

func (c *LEReadRemoteUsedFeatures) String() string {
	return "LE Read Remote Used Features (0x08|0x0016)"
}

// OpCode returns the opcode of the command.
func (c *LEReadRemoteUsedFeatures) OpCode() int { return 0x08<<10 | 0x0016 }

// Len returns the length of the command.
func (c *LEReadRemoteUsedFeatures) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadRemoteUsedFeatures) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEEncrypt implements LE Encrypt (0x08|0x0017) [Vol 2, Part E, 7.8.22]
type LEEncrypt struct {
	Key           [16]byte
	PlaintextData [16]byte
}

func (c *LEEncrypt) String() string {
	return "LE Encrypt (0x08|0x0017)"
}

// OpCode returns the opcode of the command.
func (c *LEEncrypt) OpCode() int { return 0x08<<10 | 0x0017 }

// Len returns the length of the command.
func (c *LEEncrypt) Len() int { return 32 }

// Marshal serializes the command parameters into binary form.
func (c *LEEncrypt) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEEncryptRP returns the return parameter of LE Encrypt
type LEEncryptRP struct {
	Status        uint8
	EncryptedData [16]byte
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEEncryptRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LERand implements LE Rand (0x08|0x0018) [Vol 2, Part E, 7.8.23]
type LERand struct {
}

func (c *LERand) String() string {
	return "LE Rand (0x08|0x0018)"
}

// OpCode returns the opcode of the command.
func (c *LERand) OpCode() int { return 0x08<<10 | 0x0018 }

// Len returns the length of the command.
func (c *LERand) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LERand) Marshal(b []byte) error {
	return marshal(c, b)
}

// LERandRP returns the return parameter of LE Rand
type LERandRP struct {
	Status       uint8
	RandomNumber uint64
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LERandRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEStartEncryption implements LE Start Encryption (0x08|0x0019) [Vol 2, Part E, 7.8.24]
type LEStartEncryption struct {
	ConnectionHandle     uint16
	RandomNumber         uint64
	EncryptedDiversifier uint16
	LongTermKey          [16]byte
}

func (c *LEStartEncryption) String() string {
	return "LE Start Encryption (0x08|0x0019)"
}

// OpCode returns the opcode of the command.
func (c *LEStartEncryption) OpCode() int { return 0x08<<10 | 0x0019 }

// Len returns the length of the command.
func (c *LEStartEncryption) Len() int { return 28 }

// Marshal serializes the command parameters into binary form.
func (c *LEStartEncryption) Marshal(b []byte) error {
	return marshal(c, b)
}

// LELongTermKeyRequestReply implements LE Long Term Key Request Reply (0x08|0x001A) [Vol 2, Part E, 7.8.25]
type LELongTermKeyRequestReply struct {
	ConnectionHandle uint16
	LongTermKey      [16]byte
}

func (c *LELongTermKeyRequestReply) String() string {
	return "LE Long Term Key Request Reply (0x08|0x001A)"
}

// OpCode returns the opcode of the command.
func (c *LELongTermKeyRequestReply) OpCode() int { return 0x08<<10 | 0x001A }

// Len returns the length of the command.
func (c *LELongTermKeyRequestReply) Len() int { return 18 }

// Marshal serializes the command parameters into binary form.
func (c *LELongTermKeyRequestReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LELongTermKeyRequestReplyRP returns the return parameter of LE Long Term Key Request Reply
type LELongTermKeyRequestReplyRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LELongTermKeyRequestReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LELongTermKeyRequestNegativeReply implements LE Long Term Key Request Negative Reply (0x08|0x001B) [Vol 2, Part E, 7.8.26]
type LELongTermKeyRequestNegativeReply struct {
	ConnectionHandle uint16
}

func (c *LELongTermKeyRequestNegativeReply) String() string {
	return "LE Long Term Key Request Negative Reply (0x08|0x001B)"
}

// OpCode returns the opcode of the command.
func (c *LELongTermKeyRequestNegativeReply) OpCode() int { return 0x08<<10 | 0x001B }

// Len returns the length of the command.
func (c *LELongTermKeyRequestNegativeReply) Len() int { return 2 }

// Marshal serializes the command parameters into binary form.
func (c *LELongTermKeyRequestNegativeReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LELongTermKeyRequestNegativeReplyRP returns the return parameter of LE Long Term Key Request Negative Reply
type LELongTermKeyRequestNegativeReplyRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LELongTermKeyRequestNegativeReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReadSupportedStates implements LE Read Supported States (0x08|0x001C) [Vol 2, Part E, 7.8.27]
type LEReadSupportedStates struct {
}

func (c *LEReadSupportedStates) String() string {
	return "LE Read Supported States (0x08|0x001C)"
}

// OpCode returns the opcode of the command.
func (c *LEReadSupportedStates) OpCode() int { return 0x08<<10 | 0x001C }

// Len returns the length of the command.
func (c *LEReadSupportedStates) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LEReadSupportedStates) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReadSupportedStatesRP returns the return parameter of LE Read Supported States
type LEReadSupportedStatesRP struct {
	Status   uint8
	LEStates uint64
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReadSupportedStatesRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LEReceiverTest implements LE Receiver Test (0x08|0x001D) [Vol 2, Part E, 7.8.28]
type LEReceiverTest struct {
	RXChannel uint8
}

func (c *LEReceiverTest) String() string {
	return "LE Receiver Test (0x08|0x001D)"
}

// OpCode returns the opcode of the command.
func (c *LEReceiverTest) OpCode() int { return 0x08<<10 | 0x001D }

// Len returns the length of the command.
func (c *LEReceiverTest) Len() int { return 1 }

// Marshal serializes the command parameters into binary form.
func (c *LEReceiverTest) Marshal(b []byte) error {
	return marshal(c, b)
}

// LEReceiverTestRP returns the return parameter of LE Receiver Test
type LEReceiverTestRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LEReceiverTestRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LETransmitterTest implements LE Transmitter Test (0x08|0x001E) [Vol 2, Part E, 7.8.29]
type LETransmitterTest struct {
	TXChannel        uint8
	LengthOfTestData uint8
	PacketPayload    uint8
}

func (c *LETransmitterTest) String() string {
	return "LE Transmitter Test (0x08|0x001E)"
}

// OpCode returns the opcode of the command.
func (c *LETransmitterTest) OpCode() int { return 0x08<<10 | 0x001E }

// Len returns the length of the command.
func (c *LETransmitterTest) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *LETransmitterTest) Marshal(b []byte) error {
	return marshal(c, b)
}

// LETransmitterTestRP returns the return parameter of LE Transmitter Test
type LETransmitterTestRP struct {
	Status uint8
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LETransmitterTestRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LETestEnd implements LE Test End (0x08|0x001F) [Vol 2, Part E, 7.8.30]
type LETestEnd struct {
}

func (c *LETestEnd) String() string {
	return "LE Test End (0x08|0x001F)"
}

// OpCode returns the opcode of the command.
func (c *LETestEnd) OpCode() int { return 0x08<<10 | 0x001F }

// Len returns the length of the command.
func (c *LETestEnd) Len() int { return 0 }

// Marshal serializes the command parameters into binary form.
func (c *LETestEnd) Marshal(b []byte) error {
	return marshal(c, b)
}

// LETestEndRP returns the return parameter of LE Test End
type LETestEndRP struct {
	Status          uint8
	NumberOfPackats uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LETestEndRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LERemoteConnectionParameterRequestReply implements LE Remote Connection Parameter Request Reply (0x08|0x0020) [Vol 2, Part E, 7.8.31]
type LERemoteConnectionParameterRequestReply struct {
	ConnectionHandle uint16
	IntervalMin      uint16
	IntervalMax      uint16
	Latency          uint16
	Timeout          uint16
	MinimumCELength  uint16
	MaximumCELength  uint16
}

func (c *LERemoteConnectionParameterRequestReply) String() string {
	return "LE Remote Connection Parameter Request Reply (0x08|0x0020)"
}

// OpCode returns the opcode of the command.
func (c *LERemoteConnectionParameterRequestReply) OpCode() int { return 0x08<<10 | 0x0020 }

// Len returns the length of the command.
func (c *LERemoteConnectionParameterRequestReply) Len() int { return 14 }

// Marshal serializes the command parameters into binary form.
func (c *LERemoteConnectionParameterRequestReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LERemoteConnectionParameterRequestReplyRP returns the return parameter of LE Remote Connection Parameter Request Reply
type LERemoteConnectionParameterRequestReplyRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LERemoteConnectionParameterRequestReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}

// LERemoteConnectionParameterRequestNegativeReply implements LE Remote Connection Parameter Request Negative Reply (0x08|0x0021) [Vol 2, Part E, 7.8.32]
type LERemoteConnectionParameterRequestNegativeReply struct {
	ConnectionHandle uint16
	Reason           uint8
}

func (c *LERemoteConnectionParameterRequestNegativeReply) String() string {
	return "LE Remote Connection Parameter Request Negative Reply (0x08|0x0021)"
}

// OpCode returns the opcode of the command.
func (c *LERemoteConnectionParameterRequestNegativeReply) OpCode() int { return 0x08<<10 | 0x0021 }

// Len returns the length of the command.
func (c *LERemoteConnectionParameterRequestNegativeReply) Len() int { return 3 }

// Marshal serializes the command parameters into binary form.
func (c *LERemoteConnectionParameterRequestNegativeReply) Marshal(b []byte) error {
	return marshal(c, b)
}

// LERemoteConnectionParameterRequestNegativeReplyRP returns the return parameter of LE Remote Connection Parameter Request Negative Reply
type LERemoteConnectionParameterRequestNegativeReplyRP struct {
	Status           uint8
	ConnectionHandle uint16
}

// Unmarshal de-serializes the binary data and stores the result in the receiver.
func (c *LERemoteConnectionParameterRequestNegativeReplyRP) Unmarshal(b []byte) error {
	return unmarshal(c, b)
}
