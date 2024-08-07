package hci

import (
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/adv"
	"github.com/thomascriley/ble/linux/hci/evt"
	"net"
)

// [Vol 6, Part B, 4.4.2] [Vol 3, Part C, 11]
const (
	evtTypAdvInd        = 0x00 // Connectable undirected advertising (ADV_IND).
	evtTypAdvDirectInd  = 0x01 // Connectable directed advertising (ADV_DIRECT_IND).
	evtTypAdvScanInd    = 0x02 // Scannable undirected advertising (ADV_SCAN_IND).
	evtTypAdvNonconnInd = 0x03 // Non connectable undirected advertising (ADV_NONCONN_IND).
	evtTypScanRsp       = 0x04 // Scan Response (SCAN_RSP).
)

func newAdvertisement(e evt.LEAdvertisingReport, i int) *Advertisement {
	return &Advertisement{e: e, i: i}
}

// Advertisement implements ble.Advertisement and other functions that are only
// available on Linux.
type Advertisement struct {
	e  evt.LEAdvertisingReport
	i  int
	sr *Advertisement

	// cached packets.
	p *adv.Packet

	addr       ble.Addr
	addrString string
}

// setScanResponse associate the response to the existing advertisement.
func (a *Advertisement) setScanResponse(sr *Advertisement) {
	a.sr = sr
	a.p = nil // clear the cached.
}

// packets returns the combined advertising packet and scan response (if presents)
func (a *Advertisement) packets() *adv.Packet {
	// fmt.Println("packets")
	// defer fmt.Println("packets done")
	if a.p != nil {
		return a.p
	}
	a.p = adv.NewRawPacket(a.Data(), a.ScanResponse())
	return a.p
}

// LocalName returns the LocalName of the remote peripheral.
func (a *Advertisement) LocalName() string {
	// fmt.Println("local name")
	//defer fmt.Println("local name - done")
	return a.packets().LocalName()
}

// ManufacturerData returns the ManufacturerData of the advertisement.
func (a *Advertisement) ManufacturerData() []byte {
	//fmt.Println("mfg")
	// defer fmt.Println("mfg - done")
	return a.packets().ManufacturerData()
}

// ServiceData returns the service data of the advertisement.
func (a *Advertisement) ServiceData() []ble.ServiceData {
	// fmt.Println("svc")
	// defer fmt.Println("svc - done")
	return a.packets().ServiceData()
}

// Services returns the service UUIDs of the advertisement.
func (a *Advertisement) Services() []ble.UUID {
	// fmt.Println("svcs")
	// defer fmt.Println("svcs - done")
	return a.packets().UUIDs()
}

// OverflowService returns the UUIDs of overflowed service.
func (a *Advertisement) OverflowService() []ble.UUID {
	// fmt.Println("over")
	// defer fmt.Println("over - done")
	return a.packets().UUIDs()
}

// TxPowerLevel returns the tx power level of the remote peripheral.
func (a *Advertisement) TxPowerLevel() int {
	// fmt.Println("pwr")
	// defer fmt.Println("pwr - done")
	pwr, _ := a.packets().TxPower()
	return pwr
}

// SolicitedService returns UUIDs of solicited services.
func (a *Advertisement) SolicitedService() []ble.UUID {
	return a.packets().ServiceSol()
}

// Connectable indicates weather the remote peripheral is connectable.
func (a *Advertisement) Connectable() bool {
	// fmt.Println("conn")
	// defer fmt.Println("conn - done")
	return a.EventType() == evtTypAdvDirectInd || a.EventType() == evtTypAdvInd
}

// RSSI returns RSSI signal strength.
func (a *Advertisement) RSSI() int {
	// fmt.Println("rssi")
	// defer fmt.Println("rssi - done")
	if a.e == nil {
		return 0
	}
	return int(a.e.RSSI(a.i))
}

// Addr returns the address of the remote peripheral.
func (a *Advertisement) Address() ble.Addr {
	// fmt.Println("addr")
	// defer fmt.Println("addr - done")
	if a.addr == nil {
		if a.e == nil {
			return nil
		}
		b := a.e.Address(a.i)
		a.addr = net.HardwareAddr([]byte{b[5], b[4], b[3], b[2], b[1], b[0]})
	}
	return a.addr
}

func (a *Advertisement) AddressString() string {
	// fmt.Println("addrstr")
	// defer fmt.Println("addrstr - done")
	if a.addrString == "" {
		addr := a.Address()
		if addr != nil {
			a.addrString = addr.String()
		}
	}
	return a.addrString
}

// EventType returns the event type of Advertisement.
// This is linux sepcific.
func (a *Advertisement) EventType() uint8 {
	// fmt.Println("evt type")
	// defer fmt.Println("evt type - done")
	if a.e == nil {
		return 0
	}
	return a.e.EventType(a.i)
}

// AddressType returns the address type of the Advertisement.
// This is linux sepcific.
func (a *Advertisement) AddressType() ble.AddressType {
	// fmt.Println("addr type")
	// defer fmt.Println("addr type - done")
	if a.e == nil {
		return ble.AddressTypeRandom
	}
	return ble.AddressType(a.e.AddressType(a.i))
}

// Data returns the advertising data of the packet.
// This is linux sepcific.
func (a *Advertisement) Data() []byte {
	// fmt.Println("data")
	// defer fmt.Println("data - done")
	if a.e == nil {
		return nil
	}
	return a.e.Data(a.i)
}

// ScanResponse returns the scan response of the packet, if it presents.
// This is linux sepcific.
func (a *Advertisement) ScanResponse() []byte {
	// fmt.Println("sr")
	// defer fmt.Println("sr - done")
	if a.sr == nil {
		return nil
	}
	return a.sr.Data()
}
