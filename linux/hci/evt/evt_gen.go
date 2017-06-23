package evt

import "encoding/binary"

const InquiryCompleteCode = 0x01

// InquiryComplete implements Inquiry Complete (0x01) [Vol 2, Part E, 7.7.1].
type InquiryComplete []byte

func (r InquiryComplete) Status() uint8 { return r[0] }

const InquiryResultCode = 0x02

// InquiryResult implements Inquiry Result (0x02) [Vol 2, Part E, 7.7.2].
type InquiryResult []byte

const ConnectionCompleteCode = 0x03

// ConnectionComplete implements Connection Complete (0x03) [Vol 2, Part E, 7.7.3].
type ConnectionComplete []byte

func (r ConnectionComplete) Status() uint8 { return r[0] }

func (r ConnectionComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r ConnectionComplete) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[3:])
	return b
}

func (r ConnectionComplete) LinkType() uint8 { return r[9] }

func (r ConnectionComplete) EncryptionEnabled() uint8 { return r[10] }

const ConnectionRequestCode = 0x04

// ConnectionRequest implements Connection Request (0x04) [Vol 2, Part E, 7.7.6].
type ConnectionRequest []byte

func (r ConnectionRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r ConnectionRequest) ClassOfDevice() [3]byte {
	b := [3]byte{}
	copy(b[:], r[6:])
	return b
}

func (r ConnectionRequest) LinkType() uint8 { return r[9] }

const DisconnectionCompleteCode = 0x05

// DisconnectionComplete implements Disconnection Complete (0x05) [Vol 2, Part E, 7.7.5].
type DisconnectionComplete []byte

func (r DisconnectionComplete) Status() uint8 { return r[0] }

func (r DisconnectionComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r DisconnectionComplete) Reason() uint8 { return r[3] }

const AuthenticationCompleteCode = 0x06

// AuthenticationComplete implements Authentication Complete (0x06) [Vol 2, Part E, 7.7.6].
type AuthenticationComplete []byte

func (r AuthenticationComplete) Status() uint8 { return r[0] }

func (r AuthenticationComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

const RemoteNameRequestCompleteCode = 0x07

// RemoteNameRequestComplete implements Remote Name Request Complete (0x07) [Vol 2, Part E, 7.7.7].
type RemoteNameRequestComplete []byte

func (r RemoteNameRequestComplete) Status() uint8 { return r[0] }

func (r RemoteNameRequestComplete) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[1:])
	return b
}

func (r RemoteNameRequestComplete) RemoteName() [248]byte {
	b := [248]byte{}
	copy(b[:], r[7:])
	return b
}

const EncryptionChangeCode = 0x08

// EncryptionChange implements Encryption Change (0x08) [Vol 2, Part E, 7.7.8].
type EncryptionChange []byte

func (r EncryptionChange) Status() uint8 { return r[0] }

func (r EncryptionChange) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r EncryptionChange) EncryptionEnabled() uint8 { return r[3] }

const ChangeConnectionLinkKeyCompleteCode = 0x09

// ChangeConnectionLinkKeyComplete implements Change Connection Link Key Complete (0x09) [Vol 2, Part E, 7.7.9].
type ChangeConnectionLinkKeyComplete []byte

func (r ChangeConnectionLinkKeyComplete) Status() uint8 { return r[0] }

func (r ChangeConnectionLinkKeyComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

const MasterLinkKeyCompleteCode = 0x0A

// MasterLinkKeyComplete implements Master Link Key Complete (0x0A) [Vol 2, Part E, 7.7.10].
type MasterLinkKeyComplete []byte

func (r MasterLinkKeyComplete) Status() uint8 { return r[0] }

func (r MasterLinkKeyComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r MasterLinkKeyComplete) KeyFlag() uint8 { return r[3] }

const ReadRemoteSupportedFeaturesCompleteCode = 0x0B

// ReadRemoteSupportedFeaturesComplete implements Read Remote Supported Features Complete (0x0B) [Vol 2, Part E, 7.7.11].
type ReadRemoteSupportedFeaturesComplete []byte

func (r ReadRemoteSupportedFeaturesComplete) Status() uint8 { return r[0] }

func (r ReadRemoteSupportedFeaturesComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r ReadRemoteSupportedFeaturesComplete) LMPFeatures() uint64 {
	return binary.LittleEndian.Uint64(r[3:])
}

const ReadRemoteVersionInformationCompleteCode = 0x0C

// ReadRemoteVersionInformationComplete implements Read Remote Version Information Complete (0x0C) [Vol 2, Part E, 7.7.12].
type ReadRemoteVersionInformationComplete []byte

func (r ReadRemoteVersionInformationComplete) Status() uint8 { return r[0] }

func (r ReadRemoteVersionInformationComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r ReadRemoteVersionInformationComplete) Version() uint8 { return r[3] }

func (r ReadRemoteVersionInformationComplete) ManufacturerName() uint16 {
	return binary.LittleEndian.Uint16(r[4:])
}

func (r ReadRemoteVersionInformationComplete) Subversion() uint16 {
	return binary.LittleEndian.Uint16(r[6:])
}

const QoSSetupCompleteCode = 0x0D

// QoSSetupComplete implements QoS Setup Complete (0x0D) [Vol 2, Part E, 7.7.13].
type QoSSetupComplete []byte

func (r QoSSetupComplete) Status() uint8 { return r[0] }

func (r QoSSetupComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r QoSSetupComplete) Flags() uint8 { return r[3] }

func (r QoSSetupComplete) ServiceType() uint8 { return r[4] }

func (r QoSSetupComplete) TokenRate() uint32 { return binary.LittleEndian.Uint32(r[5:]) }

func (r QoSSetupComplete) PeakBandwidth() uint32 { return binary.LittleEndian.Uint32(r[9:]) }

func (r QoSSetupComplete) Latency() uint32 { return binary.LittleEndian.Uint32(r[13:]) }

func (r QoSSetupComplete) DelayVariation() uint32 { return binary.LittleEndian.Uint32(r[17:]) }

const CommandCompleteCode = 0x0E

// CommandComplete implements Command Complete (0x0E) [Vol 2, Part E, 7.7.14].
type CommandComplete []byte

const CommandStatusCode = 0x0F

// CommandStatus implements Command Status (0x0F) [Vol 2, Part E, 7.7.15].
type CommandStatus []byte

func (r CommandStatus) Status() uint8 { return r[0] }

func (r CommandStatus) NumHCICommandPackets() uint8 { return r[1] }

func (r CommandStatus) CommandOpcode() uint16 { return binary.LittleEndian.Uint16(r[2:]) }

const HardwareErrorCode = 0x10

// HardwareError implements Hardware Error (0x10) [Vol 2, Part E, 7.7.16].
type HardwareError []byte

func (r HardwareError) HardwareCode() uint8 { return r[0] }

const FlushOccurredCode = 0x11

// FlushOccurred implements Flush Occurred (0x11) [Vol 2, Part E, 7.7.17].
type FlushOccurred []byte

func (r FlushOccurred) Handle() uint16 { return binary.LittleEndian.Uint16(r[0:]) }

const RoleChangeCode = 0x12

// RoleChange implements Role Change (0x12) [Vol 2, Part E, 7.7.18].
type RoleChange []byte

func (r RoleChange) Status() uint8 { return r[0] }

func (r RoleChange) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[1:])
	return b
}

func (r RoleChange) NewRole() uint8 { return r[7] }

const NumberOfCompletedPacketsCode = 0x13

// NumberOfCompletedPackets implements Number Of Completed Packets (0x13) [Vol 2, Part E, 7.7.19].
type NumberOfCompletedPackets []byte

const ModeChangeCode = 0x14

// ModeChange implements Mode Change (0x14) [Vol 2, Part E, 7.7.20].
type ModeChange []byte

const ReturnLinkKeysCode = 0x15

// ReturnLinkKeys implements Return Link Keys (0x15) [Vol 2, Part E, 7.7.21].
type ReturnLinkKeys []byte

const PinCodeRequestCode = 0x16

// PinCodeRequest implements Pin Code Request (0x16) [Vol 2, Part E, 7.7.22].
type PinCodeRequest []byte

func (r PinCodeRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

const LinkKeyRequestCode = 0x17

// LinkKeyRequest implements Link Key Request (0x17) [Vol 2, Part E, 7.7.23].
type LinkKeyRequest []byte

func (r LinkKeyRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

const LinkKeyNotificationCode = 0x18

// LinkKeyNotification implements Link Key Notification (0x18) [Vol 2, Part E, 7.7.24].
type LinkKeyNotification []byte

func (r LinkKeyNotification) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r LinkKeyNotification) LinkKey() [16]byte {
	b := [16]byte{}
	copy(b[:], r[6:])
	return b
}

func (r LinkKeyNotification) LinkType() uint8 { return r[22] }

const LoopbackCommandCode = 0x19

// LoopbackCommand implements Loopback Command (0x19) [Vol 2, Part E, 7.7.25].
type LoopbackCommand []byte

const DataBufferOverflowCode = 0x1A

// DataBufferOverflow implements Data Buffer Overflow (0x1A) [Vol 2, Part E, 7.7.26].
type DataBufferOverflow []byte

func (r DataBufferOverflow) LinkType() uint8 { return r[0] }

const MaxSlotsChangeCode = 0x1B

// MaxSlotsChange implements Max Slots Change (0x1B) [Vol 2, Part E, 7.7.27].
type MaxSlotsChange []byte

func (r MaxSlotsChange) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[0:]) }

func (r MaxSlotsChange) LPMaxSlots() uint8 { return r[2] }

const ReadClockOffsetCompleteCode = 0x1C

// ReadClockOffsetComplete implements Read Clock Offset Complete (0x1C) [Vol 2, Part E, 7.7.28].
type ReadClockOffsetComplete []byte

func (r ReadClockOffsetComplete) Status() uint8 { return r[0] }

func (r ReadClockOffsetComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r ReadClockOffsetComplete) ClockOffset() uint16 { return binary.LittleEndian.Uint16(r[3:]) }

const ConnectionPacketTypeChangedCode = 0x1D

// ConnectionPacketTypeChanged implements Connection Packet Type Changed (0x1D) [Vol 2, Part E, 7.7.29].
type ConnectionPacketTypeChanged []byte

func (r ConnectionPacketTypeChanged) Status() uint8 { return r[0] }

func (r ConnectionPacketTypeChanged) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r ConnectionPacketTypeChanged) PacketType() uint16 { return binary.LittleEndian.Uint16(r[3:]) }

const QoSViolationCode = 0x1E

// QoSViolation implements QoS Violation (0x1E) [Vol 2, Part E, 7.7.30].
type QoSViolation []byte

func (r QoSViolation) Handle() uint16 { return binary.LittleEndian.Uint16(r[0:]) }

const PageScanRepetitionModeChangeCode = 0x20

// PageScanRepetitionModeChange implements Page Scan Repetition Mode Change (0x20) [Vol 2, Part E, 7.7.31].
type PageScanRepetitionModeChange []byte

func (r PageScanRepetitionModeChange) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r PageScanRepetitionModeChange) PageScanRepetitionMode() uint8 { return r[6] }

const FlowSpecificationCompleteCode = 0x21

// FlowSpecificationComplete implements Flow Specification Complete (0x21) [Vol 2, Part E, 7.7.32].
type FlowSpecificationComplete []byte

func (r FlowSpecificationComplete) Status() uint8 { return r[0] }

func (r FlowSpecificationComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r FlowSpecificationComplete) Flags() uint8 { return r[3] }

func (r FlowSpecificationComplete) FlowDirection() uint8 { return r[4] }

func (r FlowSpecificationComplete) ServiceType() uint8 { return r[5] }

func (r FlowSpecificationComplete) TokenRate() uint32 { return binary.LittleEndian.Uint32(r[6:]) }

func (r FlowSpecificationComplete) TokenBucketSize() uint32 {
	return binary.LittleEndian.Uint32(r[10:])
}

func (r FlowSpecificationComplete) PeakBandwidth() uint32 { return binary.LittleEndian.Uint32(r[14:]) }

func (r FlowSpecificationComplete) AccessLatency() uint32 { return binary.LittleEndian.Uint32(r[18:]) }

const InquiryResultwithRSSICode = 0x22

// InquiryResultwithRSSI implements Inquiry Result with RSSI (0x22) [Vol 2, Part E, 7.7.33].
type InquiryResultwithRSSI []byte

const ReadRemoteExtendedFeaturesCompleteCode = 0x23

// ReadRemoteExtendedFeaturesComplete implements Read Remote Extended Features Complete (0x23) [Vol 2, Part E, 7.7.34].
type ReadRemoteExtendedFeaturesComplete []byte

func (r ReadRemoteExtendedFeaturesComplete) Status() uint8 { return r[0] }

func (r ReadRemoteExtendedFeaturesComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r ReadRemoteExtendedFeaturesComplete) PageNumber() uint8 { return r[3] }

func (r ReadRemoteExtendedFeaturesComplete) MaximumPageNumber() uint8 { return r[4] }

func (r ReadRemoteExtendedFeaturesComplete) ExtendedLMPFeatures() uint64 {
	return binary.LittleEndian.Uint64(r[5:])
}

const SynchronousConnectionCompleteCode = 0x2C

// SynchronousConnectionComplete implements Synchronous Connection Complete (0x2C) [Vol 2, Part E, 7.7.35].
type SynchronousConnectionComplete []byte

func (r SynchronousConnectionComplete) Status() uint8 { return r[0] }

func (r SynchronousConnectionComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r SynchronousConnectionComplete) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[3:])
	return b
}

func (r SynchronousConnectionComplete) LinkType() uint8 { return r[9] }

func (r SynchronousConnectionComplete) TransmissionInterval() uint8 { return r[10] }

func (r SynchronousConnectionComplete) RetransmissionWindow() uint8 { return r[11] }

func (r SynchronousConnectionComplete) RxPacketLength() uint16 {
	return binary.LittleEndian.Uint16(r[12:])
}

func (r SynchronousConnectionComplete) TxPacketLength() uint16 {
	return binary.LittleEndian.Uint16(r[14:])
}

func (r SynchronousConnectionComplete) AirMode() uint8 { return r[16] }

const SynchronousConnectionChangedCode = 0x2D

// SynchronousConnectionChanged implements Synchronous Connection Changed (0x2D) [Vol 2, Part E, 7.7.36].
type SynchronousConnectionChanged []byte

func (r SynchronousConnectionChanged) Status() uint8 { return r[0] }

func (r SynchronousConnectionChanged) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r SynchronousConnectionChanged) TransmissionInterval() uint8 { return r[3] }

func (r SynchronousConnectionChanged) RetransmissionWindow() uint8 { return r[4] }

func (r SynchronousConnectionChanged) RxPacketLength() uint16 {
	return binary.LittleEndian.Uint16(r[5:])
}

func (r SynchronousConnectionChanged) TxPacketLength() uint16 {
	return binary.LittleEndian.Uint16(r[7:])
}

const SniffSubratingCode = 0x2E

// SniffSubrating implements Sniff Subrating (0x2E) [Vol 2, Part E, 7.7.37].
type SniffSubrating []byte

func (r SniffSubrating) Status() uint8 { return r[0] }

func (r SniffSubrating) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r SniffSubrating) MaximumTransmitLatency() uint16 { return binary.LittleEndian.Uint16(r[3:]) }

func (r SniffSubrating) MinimumTransmitLatency() uint16 { return binary.LittleEndian.Uint16(r[5:]) }

func (r SniffSubrating) MinimumRemoteTimeout() uint16 { return binary.LittleEndian.Uint16(r[7:]) }

func (r SniffSubrating) MinimumLocalTimeout() uint16 { return binary.LittleEndian.Uint16(r[9:]) }

const ExtendedInquiryCode = 0x2F

// ExtendedInquiry implements Extended Inquiry (0x2F) [Vol 2, Part E, 7.7.38].
type ExtendedInquiry []byte

const EncryptionKeyRefreshCompleteCode = 0x30

// EncryptionKeyRefreshComplete implements Encryption Key Refresh Complete (0x30) [Vol 2, Part E, 7.7.39].
type EncryptionKeyRefreshComplete []byte

func (r EncryptionKeyRefreshComplete) Status() uint8 { return r[0] }

func (r EncryptionKeyRefreshComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

const IOCapabilityRequestCode = 0x31

// IOCapabilityRequest implements IO Capability Request (0x31) [Vol 2, Part E, 7.7.40].
type IOCapabilityRequest []byte

func (r IOCapabilityRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

const IOCapabilityResponseCode = 0x32

// IOCapabilityResponse implements IO Capability Response (0x32) [Vol 2, Part E, 7.7.41].
type IOCapabilityResponse []byte

func (r IOCapabilityResponse) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r IOCapabilityResponse) IOCapability() uint8 { return r[6] }

func (r IOCapabilityResponse) OOBDataPresent() uint8 { return r[7] }

func (r IOCapabilityResponse) AuthenticationRequirements() uint8 { return r[8] }

const UserConfirmationRequestCode = 0x33

// UserConfirmationRequest implements User Confirmation Request (0x33) [Vol 2, Part E, 7.7.42].
type UserConfirmationRequest []byte

func (r UserConfirmationRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r UserConfirmationRequest) NumericValue() uint32 { return binary.LittleEndian.Uint32(r[6:]) }

const UserPasskeyRequestCode = 0x34

// UserPasskeyRequest implements User Passkey Request (0x34) [Vol 2, Part E, 7.7.43].
type UserPasskeyRequest []byte

func (r UserPasskeyRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

const RemoteOOBDataRequestCode = 0x35

// RemoteOOBDataRequest implements Remote OOB Data Request (0x35) [Vol 2, Part E, 7.7.44].
type RemoteOOBDataRequest []byte

func (r RemoteOOBDataRequest) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

const SimplePairingCompleteCode = 0x36

// SimplePairingComplete implements Simple Pairing Complete (0x36) [Vol 2, Part E, 7.7.45].
type SimplePairingComplete []byte

func (r SimplePairingComplete) Status() uint8 { return r[0] }

func (r SimplePairingComplete) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[1:])
	return b
}

const LinkSupervisionTimeoutChangedCode = 0x38

// LinkSupervisionTimeoutChanged implements Link Supervision Timeout Changed (0x38) [Vol 2, Part E, 7.7.46].
type LinkSupervisionTimeoutChanged []byte

func (r LinkSupervisionTimeoutChanged) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[0:])
}

func (r LinkSupervisionTimeoutChanged) LinkSupervisionTimeout() uint16 {
	return binary.LittleEndian.Uint16(r[2:])
}

const EnhancedFlushCompleteCode = 0x39

// EnhancedFlushComplete implements Enhanced Flush Complete (0x39) [Vol 2, Part E, 7.7.47].
type EnhancedFlushComplete []byte

func (r EnhancedFlushComplete) Handle() uint16 { return binary.LittleEndian.Uint16(r[0:]) }

const UserPasskeyNotificationCode = 0x3B

// UserPasskeyNotification implements User Passkey Notification (0x3B) [Vol 2, Part E, 7.7.48].
type UserPasskeyNotification []byte

func (r UserPasskeyNotification) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r UserPasskeyNotification) Passkey() uint32 { return binary.LittleEndian.Uint32(r[6:]) }

const KeypressNotificationCode = 0x3C

// KeypressNotification implements Keypress Notification (0x3C) [Vol 2, Part E, 7.7.49].
type KeypressNotification []byte

func (r KeypressNotification) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r KeypressNotification) NotificationType() uint8 { return r[6] }

const RemoteHostSupportedFeaturesNotificationCode = 0x3D

// RemoteHostSupportedFeaturesNotification implements Remote Host Supported Features Notification (0x3D) [Vol 2, Part E, 7.7.50].
type RemoteHostSupportedFeaturesNotification []byte

func (r RemoteHostSupportedFeaturesNotification) BDADDR() [6]byte {
	b := [6]byte{}
	copy(b[:], r[0:])
	return b
}

func (r RemoteHostSupportedFeaturesNotification) HostSupportedFeatures() uint64 {
	return binary.LittleEndian.Uint64(r[6:])
}

const PhysicalLinkCompleteCode = 0x40

// PhysicalLinkComplete implements Physical Link Complete (0x40) [Vol 2, Part E, 7.7.51].
type PhysicalLinkComplete []byte

func (r PhysicalLinkComplete) Status() uint8 { return r[0] }

func (r PhysicalLinkComplete) PhysicalLinkHandle() uint8 { return r[1] }

const PhysicalLinkHandleCode = 0x41

// PhysicalLinkHandle implements Physical Link Handle (0x41) [Vol 2, Part E, 7.7.52].
type PhysicalLinkHandle []byte

func (r PhysicalLinkHandle) PhysicalLinkHandle() uint8 { return r[0] }

const DisconnectionPhysicalLinkCompleteCode = 0x42

// DisconnectionPhysicalLinkComplete implements Disconnection Physical Link Complete (0x42) [Vol 2, Part E, 7.7.53].
type DisconnectionPhysicalLinkComplete []byte

func (r DisconnectionPhysicalLinkComplete) Status() uint8 { return r[0] }

func (r DisconnectionPhysicalLinkComplete) PhysicalLinkHandle() uint8 { return r[1] }

func (r DisconnectionPhysicalLinkComplete) Reason() uint8 { return r[2] }

const PhysicalLinkLossEarlyWarningCode = 0x43

// PhysicalLinkLossEarlyWarning implements Physical Link Loss Early Warning (0x43) [Vol 2, Part E, 7.7.54].
type PhysicalLinkLossEarlyWarning []byte

func (r PhysicalLinkLossEarlyWarning) PhysicalLinkHandle() uint8 { return r[0] }

func (r PhysicalLinkLossEarlyWarning) LinkLossReason() uint8 { return r[1] }

const PhysicalLinkRecoveryCode = 0x44

// PhysicalLinkRecovery implements Physical Link Recovery (0x44) [Vol 2, Part E, 7.7.55].
type PhysicalLinkRecovery []byte

func (r PhysicalLinkRecovery) PhysicalLinkHandle() uint8 { return r[0] }

const LogicalLinkCompleteCode = 0x45

// LogicalLinkComplete implements Logical Link Complete (0x45) [Vol 2, Part E, 7.7.56].
type LogicalLinkComplete []byte

func (r LogicalLinkComplete) Status() uint8 { return r[0] }

func (r LogicalLinkComplete) LogicalLinkHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r LogicalLinkComplete) PhysicalLinkHandle() uint8 { return r[3] }

func (r LogicalLinkComplete) TXFlowSpecID() uint8 { return r[4] }

const DisconnectionLogicalLinkCompleteCode = 0x46

// DisconnectionLogicalLinkComplete implements Disconnection Logical Link Complete (0x46) [Vol 2, Part E, 7.7.57].
type DisconnectionLogicalLinkComplete []byte

func (r DisconnectionLogicalLinkComplete) Status() uint8 { return r[0] }

func (r DisconnectionLogicalLinkComplete) LogicalLinkHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r DisconnectionLogicalLinkComplete) Reason() uint8 { return r[3] }

const FlowSpecModifyCompleteCode = 0x47

// FlowSpecModifyComplete implements Flow Spec Modify Complete (0x47) [Vol 2, Part E, 7.7.58].
type FlowSpecModifyComplete []byte

func (r FlowSpecModifyComplete) Status() uint8 { return r[0] }

func (r FlowSpecModifyComplete) Handle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

const NumberofCompletedDataBlocksCode = 0x48

// NumberofCompletedDataBlocks implements Number of Completed Data Blocks (0x48) [Vol 2, Part E, 7.7.59].
type NumberofCompletedDataBlocks []byte

const ShortRangeModeChangeCompleteCode = 0x4C

// ShortRangeModeChangeComplete implements Short Range Mode Change Complete (0x4C) [Vol 2, Part E, 7.7.60].
type ShortRangeModeChangeComplete []byte

func (r ShortRangeModeChangeComplete) Status() uint8 { return r[0] }

func (r ShortRangeModeChangeComplete) PhysicalLinkHandle() uint8 { return r[1] }

func (r ShortRangeModeChangeComplete) ShortRangeModeState() uint8 { return r[2] }

const AMPStatusChangeCode = 0x4D

// AMPStatusChange implements AMP Status Change (0x4D) [Vol 2, Part E, 7.7.61].
type AMPStatusChange []byte

func (r AMPStatusChange) Status() uint8 { return r[0] }

func (r AMPStatusChange) AMPStatus() uint8 { return r[1] }

const AMPStartTestCode = 0x49

// AMPStartTest implements AMP Start Test (0x49) [Vol 2, Part E, 7.7.62].
type AMPStartTest []byte

func (r AMPStartTest) Status() uint8 { return r[0] }

func (r AMPStartTest) TestScenario() uint8 { return r[1] }

const AMPTestEndCode = 0x4A

// AMPTestEnd implements AMP Test End (0x4A) [Vol 2, Part E, 7.7.63].
type AMPTestEnd []byte

func (r AMPTestEnd) Status() uint8 { return r[0] }

func (r AMPTestEnd) TestScenario() uint8 { return r[1] }

const AMPReceiverReportCode = 0x4B

// AMPReceiverReport implements AMP Receiver Report (0x4B) [Vol 2, Part E, 7.7.64].
type AMPReceiverReport []byte

func (r AMPReceiverReport) ControllerType() uint8 { return r[0] }

func (r AMPReceiverReport) Reason() uint8 { return r[1] }

func (r AMPReceiverReport) EventType() uint32 { return binary.LittleEndian.Uint32(r[2:]) }

func (r AMPReceiverReport) NumberOfFrames() uint16 { return binary.LittleEndian.Uint16(r[6:]) }

func (r AMPReceiverReport) NumberOfErrorFrames() uint16 { return binary.LittleEndian.Uint16(r[8:]) }

func (r AMPReceiverReport) NumberOfBits() uint32 { return binary.LittleEndian.Uint32(r[10:]) }

func (r AMPReceiverReport) NumberofErrorBits() uint32 { return binary.LittleEndian.Uint32(r[14:]) }

const LEConnectionCompleteCode = 0x3E

const LEConnectionCompleteSubCode = 0x01

// LEConnectionComplete implements LE Connection Complete (0x3E:0x01) [Vol 2, Part E, 7.7.65.1].
type LEConnectionComplete []byte

func (r LEConnectionComplete) SubeventCode() uint8 { return r[0] }

func (r LEConnectionComplete) Status() uint8 { return r[1] }

func (r LEConnectionComplete) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[2:]) }

func (r LEConnectionComplete) Role() uint8 { return r[4] }

func (r LEConnectionComplete) PeerAddressType() uint8 { return r[5] }

func (r LEConnectionComplete) PeerAddress() [6]byte {
	b := [6]byte{}
	copy(b[:], r[6:])
	return b
}

func (r LEConnectionComplete) ConnInterval() uint16 { return binary.LittleEndian.Uint16(r[12:]) }

func (r LEConnectionComplete) ConnLatency() uint16 { return binary.LittleEndian.Uint16(r[14:]) }

func (r LEConnectionComplete) SupervisionTimeout() uint16 { return binary.LittleEndian.Uint16(r[16:]) }

func (r LEConnectionComplete) MasterClockAccuracy() uint8 { return r[18] }

const LEAdvertisingReportCode = 0x3E

const LEAdvertisingReportSubCode = 0x02

// LEAdvertisingReport implements LE Advertising Report (0x3E:0x02) [Vol 2, Part E, 7.7.65.2].
type LEAdvertisingReport []byte

const LEConnectionUpdateCompleteCode = 0x0E

const LEConnectionUpdateCompleteSubCode = 0x03

// LEConnectionUpdateComplete implements LE Connection Update Complete (0x0E:0x03) [Vol 2, Part E, 7.7.65.3].
type LEConnectionUpdateComplete []byte

func (r LEConnectionUpdateComplete) SubeventCode() uint8 { return r[0] }

func (r LEConnectionUpdateComplete) Status() uint8 { return r[1] }

func (r LEConnectionUpdateComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[2:])
}

func (r LEConnectionUpdateComplete) ConnInterval() uint16 { return binary.LittleEndian.Uint16(r[4:]) }

func (r LEConnectionUpdateComplete) ConnLatency() uint16 { return binary.LittleEndian.Uint16(r[6:]) }

func (r LEConnectionUpdateComplete) SupervisionTimeout() uint16 {
	return binary.LittleEndian.Uint16(r[8:])
}

const LEReadRemoteUsedFeaturesCompleteCode = 0x3E

const LEReadRemoteUsedFeaturesCompleteSubCode = 0x04

// LEReadRemoteUsedFeaturesComplete implements LE Read Remote Used Features Complete (0x3E:0x04) [Vol 2, Part E, 7.7.65.4].
type LEReadRemoteUsedFeaturesComplete []byte

func (r LEReadRemoteUsedFeaturesComplete) SubeventCode() uint8 { return r[0] }

func (r LEReadRemoteUsedFeaturesComplete) Status() uint8 { return r[1] }

func (r LEReadRemoteUsedFeaturesComplete) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[2:])
}

func (r LEReadRemoteUsedFeaturesComplete) LEFeatures() uint64 {
	return binary.LittleEndian.Uint64(r[4:])
}

const LELongTermKeyRequestCode = 0x3E

const LELongTermKeyRequestSubCode = 0x05

// LELongTermKeyRequest implements LE Long Term Key Request (0x3E:0x05) [Vol 2, Part E, 7.7.65.5].
type LELongTermKeyRequest []byte

func (r LELongTermKeyRequest) SubeventCode() uint8 { return r[0] }

func (r LELongTermKeyRequest) ConnectionHandle() uint16 { return binary.LittleEndian.Uint16(r[1:]) }

func (r LELongTermKeyRequest) RandomNumber() uint64 { return binary.LittleEndian.Uint64(r[3:]) }

func (r LELongTermKeyRequest) EncryptionDiversifier() uint16 {
	return binary.LittleEndian.Uint16(r[11:])
}

const LERemoteConnectionParameterRequestCode = 0x3E

const LERemoteConnectionParameterRequestSubCode = 0x06

// LERemoteConnectionParameterRequest implements LE Remote Connection Parameter Request (0x3E:0x06) [Vol 2, Part E, 7.7.65.6].
type LERemoteConnectionParameterRequest []byte

func (r LERemoteConnectionParameterRequest) SubeventCode() uint8 { return r[0] }

func (r LERemoteConnectionParameterRequest) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[1:])
}

func (r LERemoteConnectionParameterRequest) IntervalMin() uint16 {
	return binary.LittleEndian.Uint16(r[3:])
}

func (r LERemoteConnectionParameterRequest) IntervalMax() uint16 {
	return binary.LittleEndian.Uint16(r[5:])
}

func (r LERemoteConnectionParameterRequest) Latency() uint16 {
	return binary.LittleEndian.Uint16(r[7:])
}

func (r LERemoteConnectionParameterRequest) Timeout() uint16 {
	return binary.LittleEndian.Uint16(r[9:])
}

const AuthenticatedPayloadTimeoutExpiredCode = 0x57

// AuthenticatedPayloadTimeoutExpired implements Authenticated Payload Timeout Expired (0x57) [Vol 2, Part E, 7.7.75].
type AuthenticatedPayloadTimeoutExpired []byte

func (r AuthenticatedPayloadTimeoutExpired) ConnectionHandle() uint16 {
	return binary.LittleEndian.Uint16(r[0:])
}
