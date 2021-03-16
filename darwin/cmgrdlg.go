// cmgrdlg.go: Implements the CentralManagerDelegate interface.  CoreBluetooth
// communicates events asynchronously via callbacks.  This file implements a
// synchronous interface by translating these callbacks into channel
// operations.

package darwin

import (
	"fmt"
	"github.com/thomascriley/ble"
	"github.com/JuulLabs-OSS/cbgo"
)

func (d *Device) CentralManagerDidUpdateState(cmgr cbgo.CentralManager) {
	d.evl.stateChanged.RxSignal(struct{}{})
}

func (d *Device) DidDiscoverPeripheral(cmgr cbgo.CentralManager, prph cbgo.Peripheral,
	advFields cbgo.AdvFields, rssi int) {

	if d.advHandler == nil {
		return
	}

	a := &adv{
		localName: advFields.LocalName,
		rssi:      int(rssi),
		mfgData:   advFields.ManufacturerData,
	}
	if advFields.Connectable != nil {
		a.connectable = *advFields.Connectable
	}
	if advFields.TxPowerLevel != nil {
		a.powerLevel = *advFields.TxPowerLevel
	}
	for _, u := range advFields.ServiceUUIDs {
		a.svcUUIDs = append(a.svcUUIDs, ble.UUID(u))
	}
	for _, sd := range advFields.ServiceData {
		a.svcData = append(a.svcData, ble.ServiceData{
			UUID: ble.UUID(sd.UUID),
			Data: sd.Data,
		})
	}
	a.peerUUID = ble.UUID(prph.Identifier())

	d.advHandler(a)
}

func (d *Device) DidConnectPeripheral(_ cbgo.CentralManager, prph cbgo.Peripheral) {
	c, err := newCentralConn(d, prph)
	if err != nil {
		d.evl.connected.RxSignal(&eventConnected{err: err})
	}
	d.evl.connected.RxSignal(&eventConnected{conn: c})
}

func (d *Device) DidDisconnectPeripheral(cmgr cbgo.CentralManager, prph cbgo.Peripheral, err error) {
	fmt.Printf("device disconnected: %s", prph.Identifier())
	if c, ok := d.findConn(ble.NewAddr(prph.Identifier().String())); ok {
		select {
		case <-c.done:
		default:
			close(c.done)
		}
	} else {
		fmt.Printf("failed to find disconnected peripheral: %s", prph.Identifier())
	}
}
