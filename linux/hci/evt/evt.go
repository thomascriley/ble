package evt

import "encoding/binary"

func (e CommandComplete) NumHCICommandPackets() uint8 { return e[0] }
func (e CommandComplete) CommandOpcode() uint16       { return binary.LittleEndian.Uint16(e[1:]) }
func (e CommandComplete) ReturnParameters() []byte    { return e[3:] }

// Per-spec [Vol 2, Part E, 7.7.19], the packet structure should be:
//
//     NumOfHandle, HandleA, HandleB, CompPktNumA, CompPktNumB
//
// But we got the actual packet from BCM20702A1 with the following structure instead.
//
//     NumOfHandle, HandleA, CompPktNumA, HandleB, CompPktNumB
//              02,   40 00,       01 00,   41 00,       01 00

func (e NumberOfCompletedPackets) NumberOfHandles() uint8 { return e[0] }
func (e NumberOfCompletedPackets) ConnectionHandle(i int) uint16 {
	// return binary.LittleEndian.Uint16(e[1+i*2:])
	return binary.LittleEndian.Uint16(e[1+i*4:])
}
func (e NumberOfCompletedPackets) HCNumOfCompletedPackets(i int) uint16 {
	// return binary.LittleEndian.Uint16(e[1+int(e.NumberOfHandles())*2:])
	return binary.LittleEndian.Uint16(e[1+i*4+2:])
}
func (e LEAdvertisingReport) SubeventCode() uint8     { return e[0] }
func (e LEAdvertisingReport) NumReports() uint8       { return e[1] }
func (e LEAdvertisingReport) EventType(i int) uint8   { return e[2+i] }
func (e LEAdvertisingReport) AddressType(i int) uint8 { return e[2+int(e.NumReports())*1+i] }
func (e LEAdvertisingReport) Address(i int) [6]byte {
	e = e[2+int(e.NumReports())*2:]
	b := [6]byte{}
	copy(b[:], e[6*i:])
	return b
}

func (e LEAdvertisingReport) LengthData(i int) uint8 { return e[2+int(e.NumReports())*8+i] }

func (e LEAdvertisingReport) Data(i int) []byte {
	l := 0
	for j := 0; j < i; j++ {
		l += int(e.LengthData(j))
	}
	b := e[2+int(e.NumReports())*9+l:]
	return b[:e.LengthData(i)]
}

func (e LEAdvertisingReport) RSSI(i int) int8 {
	l := 0
	for j := 0; j < int(e.NumReports()); j++ {
		l += int(e.LengthData(j))
	}
	return int8(e[2+int(e.NumReports())*9+l+i])
}

func (e InquiryResult) NumResponses() uint8 { return e[0] }
func (e InquiryResult) BDADDR(i int) [6]byte {
	b := [6]byte{}
	copy(b[:], e[1+i*14:])
	return b
}
func (e InquiryResult) PageScanRepetitionMode(i int) uint8 {
	return uint8(e[1+6+i*14])
}
func (e InquiryResult) ClassOfDevice(i int) [3]byte {
	b := [3]byte{}
	copy(b[:], e[1+9+i*14:])
	return b
}
func (e InquiryResult) ClockOffset(i int) uint16 {
	return binary.LittleEndian.Uint16(e[1+12+i*14:])
}
func (e InquiryResult) RSSI(i int) uint8 {
	return 0
}
func (e InquiryResult) ExtendedInquiryResponse() [240]byte {
	return [240]byte{}
}

func (e InquiryResultwithRSSI) NumResponses() uint8 { return e[0] }
func (e InquiryResultwithRSSI) BDADDR(i int) [6]byte {
	b := [6]byte{}
	copy(b[:], e[1+i*15:])
	return b
}
func (e InquiryResultwithRSSI) PageScanRepetitionMode(i int) uint8 {
	return uint8(e[1+6+i*15])
}
func (e InquiryResultwithRSSI) ClassOfDevice(i int) [3]byte {
	b := [3]byte{}
	copy(b[:], e[1+9+i*15:])
	return b
}
func (e InquiryResultwithRSSI) ClockOffset(i int) uint16 {
	return binary.LittleEndian.Uint16(e[1+12+i*15:])
}
func (e InquiryResultwithRSSI) RSSI(i int) uint8 {
	return uint8(e[1+14+i*15])
}
func (e InquiryResultwithRSSI) ExtendedInquiryResponse() [240]byte {
	return [240]byte{}
}

func (e ExtendedInquiry) NumResponses() uint8 { return e[0] }
func (e ExtendedInquiry) BDADDR(i int) [6]byte {
	b := [6]byte{}
	copy(b[:], e[1:])
	return b
}
func (e ExtendedInquiry) PageScanRepetitionMode(i int) uint8 {
	return uint8(e[1+6])
}
func (e ExtendedInquiry) ClassOfDevice(i int) [3]byte {
	b := [3]byte{}
	copy(b[:], e[1+9:])
	return b
}
func (e ExtendedInquiry) ClockOffset(i int) uint16 {
	return binary.LittleEndian.Uint16(e[1+12:])
}
func (e ExtendedInquiry) RSSI(i int) uint8 {
	return uint8(e[1+14])
}
func (e ExtendedInquiry) ExtendedInquiryResponse() [240]byte {
	b := [240]byte{}
	copy(b[:], e[1+15:])
	return b
}

func (e ConnectionComplete) PeerAddress() [6]byte {
	return e.BDADDR()
}

func (e ConnectionComplete) Role() uint8 {
	return 0x00
}
