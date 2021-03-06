package hci

// HCI Packet types
const (
	pktTypeCommand uint8 = 0x01
	pktTypeACLData uint8 = 0x02
	pktTypeSCOData uint8 = 0x03
	pktTypeEvent   uint8 = 0x04
	pktTypeVendor  uint8 = 0xFF
)

// Packet boundary flags of HCI ACL Data Packet [Vol 2, Part E, 5.4.2].
const (
	pbfHostToControllerStart = 0x00 // Start of a non-automatically-flushable from host to controller.
	pbfContinuing            = 0x01 // Continuing fragment.
	pbfControllerToHostStart = 0x02 // Start of a non-automatically-flushable from controller to host.
	pbfCompleteL2CAPPDU      = 0x03 // A automatically flushable complete PDU. (Not used in LE-U).
)

// L2CAP Channel Identifier namespace for LE-U logical link [Vol 3, Part A, 2.1].
const (
	cidSignal       uint16 = 0x01 // L2CAP Signaling Channel [Vol 3, Part A, 4].
	cidLEAtt        uint16 = 0x04 // Attribute Protocol [Vol 3, Part F].
	cidLESignal     uint16 = 0x05 // Low Energy L2CAP Signaling channel [Vol 3, Part A, 4].
	cidSMP          uint16 = 0x06 // SecurityManager Protocol [Vol 3, Part H].
	cidDynamicStart uint16 = 0x40 // Dyncamically Allocated [Vol 3 Section 7.1]
)

const (
	roleMaster = 0x00
	roleSlave  = 0x01
)

// Assigned Numbers
// https://www.bluetooth.com/specifications/assigned-numbers/logical-link-control
const (
	psmSDP    = 0x0001 // Service Discovery Protocol https://www.bluetooth.com/specifications/adopted-specifications
	psmRFCOMM = 0x0003 // RFCOMM with TS 07.10 https://www.bluetooth.com/specifications/adopted-specifications
)
