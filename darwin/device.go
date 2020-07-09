package darwin

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"errors"
	"github.com/raff/goble/xpc"
	"github.com/thomascriley/ble"
)

// Device is either a Peripheral or Central device.
type Device struct {
	pm xpc.XPC // peripheralManager
	cm xpc.XPC // centralManager

	role int // 1: peripheralManager (server), 0: centralManager (client)

	rspc chan msg

	conns    map[string]*conn
	connLock sync.Mutex

	// Only used in client/centralManager implementation
	advHandler ble.AdvHandler
	chConn     chan *conn

	// Only used in server/peripheralManager implementation
	chars map[int]*ble.Characteristic
	base  int

	closed chan struct{}
}

// NewDeviceWithNameAndHandler returns a BLE device.
func NewDeviceWithNameAndHandler(name string, handler ble.NotifyHandler, opts ...ble.Option) (*Device, error) {
	err := initXpcIDs()
	if err != nil {
		return nil, err
	}

	d := &Device{
		rspc:   make(chan msg),
		conns:  make(map[string]*conn),
		chConn: make(chan *conn),
		chars:  make(map[int]*ble.Characteristic),
		base:   1,
		closed: make(chan struct{}, 1),
	}
	if err := d.Option(opts...); err != nil {
		return nil, err
	}

	d.pm = xpc.XpcConnect(serviceID, d)
	d.cm = xpc.XpcConnect(serviceID, d)

	if err := d.Init(); err != nil {
		return nil, fmt.Errorf("can't init: %w", err)
	}
	return d, nil
}

// Option sets the options specified.
func (d *Device) Option(opts ...ble.Option) error {
	var err error
	for _, opt := range opts {
		err = opt(d)
	}
	return err
}

// Init ...
func (d *Device) Init() error {
	rsp, err := d.sendReq(d.cm, cmdInit, xpc.Dict{
		"kCBMsgArgName": fmt.Sprintf("gopher-%v", time.Now().Unix()),
		"kCBMsgArgOptions": xpc.Dict{
			"kCBInitOptionShowPowerAlert": 1,
		},
		"kCBMsgArgType": 0,
	})
	if err != nil {
		return err
	}
	s := State(rsp.state())
	if s != StatePoweredOn {
		return fmt.Errorf("state: %s", s)
	}

	rsp, err = d.sendReq(d.pm, cmdInit, xpc.Dict{
		"kCBMsgArgName": fmt.Sprintf("gopher-%v", time.Now().Unix()),
		"kCBMsgArgOptions": xpc.Dict{
			"kCBInitOptionShowPowerAlert": 1,
		},
		"kCBMsgArgType": 1,
	})
	if err != nil {
		return err
	}
	s = State(rsp.state())
	if s != StatePoweredOn {
		return fmt.Errorf("state: %s", s)
	}
	return nil
}

// Advertise advertises the given Advertisement
func (d *Device) Advertise(ctx context.Context, adv ble.Advertisement) error {
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStart, xpc.Dict{
		"kCBAdvDataLocalName":    adv.LocalName(),
		"kCBAdvDataServiceUUIDs": adv.Services(),
		"kCBAdvDataAppleMfgData": adv.ManufacturerData(),
	})
	if err != nil {
		return err
	}
	if err := rsp.err(); err != nil {
		return err
	}
	<-ctx.Done()
	_ = d.stopAdvertising()
	return ctx.Err()

}

// AdvertiseMfgData ...
func (d *Device) AdvertiseMfgData(ctx context.Context, id uint16, md []byte) error {
	l := len(md)
	b := []byte{byte(l + 3), 0xFF, uint8(id), uint8(id >> 8)}
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStart, xpc.Dict{
		"kCBAdvDataAppleMfgData": append(b, md...),
	})
	if err != nil {
		return err
	}
	if err := rsp.err(); err != nil {
		return fmt.Errorf("can't advertise: %w", err)
	}
	<-ctx.Done()
	return ctx.Err()
}

// AdvertiseServiceData16 advertises data associated with a 16bit service uuid
func (d *Device) AdvertiseServiceData16(ctx context.Context, id uint16, b []byte) error {
	l := len(b)
	prefix := []byte{
		0x03, 0x03, uint8(id), uint8(id >> 8),
		byte(l + 3), 0x16, uint8(id), uint8(id >> 8),
	}
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStart, xpc.Dict{
		"kCBAdvDataAppleMfgData": append(prefix, b...),
	})
	if err != nil {
		return err
	}
	if err := rsp.err(); err != nil {
		return fmt.Errorf( "can't advertise: %w", err)
	}
	<-ctx.Done()
	return ctx.Err()
}

// AdvertiseNameAndServices advertises name and specifid service UUIDs.
func (d *Device) AdvertiseNameAndServices(ctx context.Context, name string, ss ...ble.UUID) error {
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStart, xpc.Dict{
		"kCBAdvDataLocalName":    name,
		"kCBAdvDataServiceUUIDs": uuidSlice(ss)},
	)
	if err != nil {
		return err
	}
	if err := rsp.err(); err != nil {
		return err
	}
	<-ctx.Done()
	_ = d.stopAdvertising()
	return ctx.Err()
}

// AdvertiseIBeaconData advertises iBeacon packet with specified manufacturer data.
func (d *Device) AdvertiseIBeaconData(ctx context.Context, md []byte) error {
	var utsname xpc.Utsname
	err := xpc.Uname(&utsname)
	if err != nil {
		return err
	}

	if utsname.Release >= "14." {
		ibeaconCode := []byte{0x02, 0x15}
		return d.AdvertiseMfgData(ctx, 0x004C, append(ibeaconCode, md...))
	}
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStart, xpc.Dict{"kCBAdvDataAppleBeaconKey": md})
	if err != nil {
		return err
	}
	if err := rsp.err(); err != nil {
		return err
	}
	<-ctx.Done()
	return d.stopAdvertising()
}

// AdvertiseIBeacon advertises iBeacon packet.
func (d *Device) AdvertiseIBeacon(ctx context.Context, u ble.UUID, major, minor uint16, pwr int8) error {
	b := make([]byte, 21)
	copy(b, ble.Reverse(u))                   // Big endian
	binary.BigEndian.PutUint16(b[16:], major) // Big endian
	binary.BigEndian.PutUint16(b[18:], minor) // Big endian
	b[20] = uint8(pwr)                        // Measured Tx Power
	return d.AdvertiseIBeaconData(ctx, b)
}

// stopAdvertising stops advertising.
func (d *Device) stopAdvertising() error {
	rsp, err := d.sendReq(d.pm, cmdAdvertiseStop, nil)
	if err != nil {
		return fmt.Errorf( "can't send stop advertising: %w", err)
	}
	if err := rsp.err(); err != nil {
		return fmt.Errorf( "can't stop advertising: %w", err)
	}
	return nil
}

// Scans for BR/EDR devices
func (d *Device) Inquire(ctx context.Context, numResponses int, h ble.InqHandler) error {
	return errors.New("Not implemented")
}

// Scan ...
func (d *Device) Scan(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
	d.advHandler = h
	if err := d.sendCmd(d.cm, cmdScanningStart, xpc.Dict{
		// "kCBMsgArgUUIDs": uuidSlice(ss),
		"kCBMsgArgOptions": xpc.Dict{
			"kCBScanOptionAllowDuplicates": map[bool]int{true: 1, false: 0}[allowDup],
		},
	}); err != nil {
		return err
	}
	<-ctx.Done()
	if err := d.stopScanning(); err != nil {
		return fmt.Errorf("%s: %s",ctx.Err(), err)
	}
	return ctx.Err()
}

// stopAdvertising stops advertising.
func (d *Device) stopScanning() error {
	if err := d.sendCmd(d.cm, cmdScanningStop, nil); err != nil {
		return fmt.Errorf( "can't stop scanning: %w", err)
	}
	return nil
}

// RemoveAllServices removes all services of device's
func (d *Device) RemoveAllServices() error {
	return d.sendCmd(d.pm, cmdServicesRemove, nil)
}

// AddService adds a service to device's database.
// The following services are ignored as they are provided by OS X.
//
// 0x1800 (Generic Access)
// 0x1801 (Generic Attribute)
// 0x1805 (Current Time Service)
// 0x180A (Device Information)
// 0x180F (Battery Service)
// 0x1812 (Human Interface Device)
func (d *Device) AddService(s *ble.Service) error {
	if s.UUID.Equal(ble.GAPUUID) ||
		s.UUID.Equal(ble.GATTUUID) ||
		s.UUID.Equal(ble.CurrentTimeUUID) ||
		s.UUID.Equal(ble.DeviceInfoUUID) ||
		s.UUID.Equal(ble.BatteryUUID) ||
		s.UUID.Equal(ble.HIDUUID) {
		return nil
	}
	xs := xpc.Dict{
		"kCBMsgArgAttributeID":     d.base,
		"kCBMsgArgAttributeIDs":    []int{},
		"kCBMsgArgCharacteristics": nil,
		"kCBMsgArgType":            1, // 1 => primary, 0 => excluded
		"kCBMsgArgUUID":            ble.Reverse(s.UUID),
	}
	d.base++

	xcs := xpc.Array{}
	for _, c := range s.Characteristics {
		props := 0
		perm := 0
		if c.Property&ble.CharRead != 0 {
			props |= 0x02
			if ble.CharRead&c.Secure != 0 {
				perm |= 0x04
			} else {
				perm |= 0x01
			}
		}
		if c.Property&ble.CharWriteNR != 0 {
			props |= 0x04
			if c.Secure&ble.CharWriteNR != 0 {
				perm |= 0x08
			} else {
				perm |= 0x02
			}
		}
		if c.Property&ble.CharWrite != 0 {
			props |= 0x08
			if c.Secure&ble.CharWrite != 0 {
				perm |= 0x08
			} else {
				perm |= 0x02
			}
		}
		if c.Property&ble.CharNotify != 0 {
			if c.Secure&ble.CharNotify != 0 {
				props |= 0x100
			} else {
				props |= 0x10
			}
		}
		if c.Property&ble.CharIndicate != 0 {
			if c.Secure&ble.CharIndicate != 0 {
				props |= 0x200
			} else {
				props |= 0x20
			}
		}

		xc := xpc.Dict{
			"kCBMsgArgAttributeID":              d.base,
			"kCBMsgArgUUID":                     ble.Reverse(c.UUID),
			"kCBMsgArgAttributePermissions":     perm,
			"kCBMsgArgCharacteristicProperties": props,
			"kCBMsgArgData":                     c.Value,
		}
		c.Handle = uint16(d.base)
		d.chars[d.base] = c
		d.base++

		xds := xpc.Array{}
		for _, d := range c.Descriptors {
			if d.UUID.Equal(ble.ClientCharacteristicConfigUUID) {
				// skip CCCD
				continue
			}
			xd := xpc.Dict{
				"kCBMsgArgData": d.Value,
				"kCBMsgArgUUID": ble.Reverse(d.UUID),
			}
			xds = append(xds, xd)
		}
		xc["kCBMsgArgDescriptors"] = xds
		xcs = append(xcs, xc)
	}
	xs["kCBMsgArgCharacteristics"] = xcs

	rsp, err := d.sendReq(d.pm, cmdServicesAdd, xs)
	if err != nil {
		return err
	}
	return rsp.err()
}

// SetServices ...
func (d *Device) SetServices(ss []*ble.Service) error {
	if err := d.RemoveAllServices(); err != nil {
		return nil
	}
	for _, s := range ss {
		if err := d.AddService(s); err != nil {
			return err
		}
	}
	return nil
}

// Dial ...
func (d *Device) Dial(ctx context.Context, a ble.Addr) (ble.Client, error) {
	err := d.sendCmd(d.cm, cmdConnect, xpc.Dict{
		"kCBMsgArgDeviceUUID": xpc.MakeUUID(a.String()),
		"kCBMsgArgOptions": xpc.Dict{
			"kCBConnectOptionNotifyOnDisconnection": 1,
		},
	})
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case c := <-d.chConn:
		c.SetContext(ctx)
		return NewClient(c)
	}
}

// DialRFCOMM ...
func (d *Device) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (ble.RFCOMMClient, error) {
	return nil, errors.New("Not implemented")
}

// RequestRemoteName
func (h *Device) RequestRemoteName(a ble.Addr) (string, error) {
	return "", errors.New("Not implemented")
}

// Stop ...
func (d *Device) Stop() error {
	return nil
}

// HandleXpcEvent process Device events and asynchronous errors.
func (d *Device) HandleXpcEvent(event xpc.Dict, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	m := msg(event)
	args := msg(msg(event).args())
	//logger.Info("recv", "id", m.id(), "args", fmt.Sprintf("%v", m.args()))

	switch m.id() {
	case // Device event
		evtStateChanged,
		evtAdvertisingStarted,
		evtAdvertisingStopped,
		evtServiceAdded:
		d.rspc <- args

	case evtPeripheralDiscovered:
		if d.advHandler == nil {
			break
		}
		a := &adv{args: m.args(), ad: args.advertisementData()}
		go d.advHandler(a)

	case evtConfirmation:
		// log.Printf("confirmed: %d", args.attributeID())

	case evtATTMTU:
		d.conn(args).SetTxMTU(args.attMTU())

	case evtSlaveConnectionComplete:
		// remote peripheral is connected.
		fallthrough
	case evtMasterConnectionComplete:
		// remote central is connected.

		// Could be LEConnectionComplete or LEConnectionUpdateComplete.
		c := d.conn(args)
		c.connInterval = args.connectionInterval()
		c.connLatency = args.connectionLatency()
		c.supervisionTimeout = args.supervisionTimeout()

	case evtReadRequest:
		aid := args.attributeID()
		char := d.chars[aid]
		v := char.Value
		if v == nil {
			c := d.conn(args)
			req := ble.NewRequest(c, nil, args.offset())
			buf := bytes.NewBuffer(make([]byte, 0, c.txMTU-1))
			rsp := ble.NewResponseWriter(buf)
			char.ReadHandler.ServeRead(req, rsp)
			v = buf.Bytes()
		}

		err := d.sendCmd(d.pm, cmdSendData, xpc.Dict{
			"kCBMsgArgAttributeID":   aid,
			"kCBMsgArgData":          v,
			"kCBMsgArgTransactionID": args.transactionID(),
			"kCBMsgArgResult":        0,
		})
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
	case evtWriteRequest:
		for _, xxw := range args.attWrites() {
			xw := msg(xxw.(xpc.Dict))
			aid := xw.attributeID()
			char := d.chars[aid]
			req := ble.NewRequest(d.conn(args), xw.data(), xw.offset())
			char.WriteHandler.ServeWrite(req, nil)
			if xw.ignoreResponse() == 1 {
				continue
			}
			err := d.sendCmd(d.pm, cmdSendData, xpc.Dict{
				"kCBMsgArgAttributeID":   aid,
				"kCBMsgArgData":          nil,
				"kCBMsgArgTransactionID": args.transactionID(),
				"kCBMsgArgResult":        0,
			})
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
		}

	case evtSubscribe:
		// characteristic is subscribed by remote central.
		d.conn(args).subscribed(d.chars[args.attributeID()])

	case evtUnsubscribe:
		// characteristic is unsubscribed by remote central.
		d.conn(args).unsubscribed(d.chars[args.attributeID()])

	case evtPeripheralConnected:
		c := d.conn(args)
		if !c.isConnected {
			c.isConnected = true
			d.chConn <- c
		}

	case evtPeripheralDisconnected:
		c := d.conn(args)
		c.isConnected = false
		select {
		case c.rspc <- m:
			// Canceled by local central synchronously
		default:
			// Canceled by remote peripheral asynchronously.
		}
		d.connLock.Lock()
		delete(d.conns, c.RemoteAddr().String())
		d.connLock.Unlock()
		close(c.done)

	case evtCharacteristicRead:
		// Notification
		c := d.conn(args)

		sub := c.subs[uint16(args.characteristicHandle())]
		if sub == nil {
			fmt.Println("notified by unsubscribed handle")
			// FIXME: should terminate the connection?
		} else {
			sub.fn(args.data())
		}
		break

	case // Peripheral events
		evtRSSIRead,
		evtServiceDiscovered,
		evtIncludedServicesDiscovered,
		evtCharacteristicsDiscovered,
		evtCharacteristicWritten,
		evtNotificationValueSet,
		evtDescriptorsDiscovered,
		evtDescriptorRead,
		evtDescriptorWritten:

		d.conn(args).rspc <- m

	default:
		fmt.Printf("Unhandled event: %#v\n", event)
	}
}

// Closed ...
func (d *Device) Closed() <-chan struct{} {
	return d.closed
}

// Error ...
func (d *Device) Error() error {
	return nil
}

func (d *Device) conn(m msg) *conn {
	// Convert xpc.UUID to ble.UUID.
	a := ble.MustParse(m.deviceUUID().String())
	d.connLock.Lock()
	c, ok := d.conns[a.String()]
	if !ok {
		c = newConn(d, a, m.attMTU())
		d.conns[a.String()] = c
	}
	d.connLock.Unlock()
	return c
}

// sendReq sends a message and waits for its reply.
func (d *Device) sendReq(x xpc.XPC, id int, args xpc.Dict) (msg, error) {
	err := d.sendCmd(x, id, args)
	if err != nil {
		return msg{}, err
	}
	return <-d.rspc, nil
}

func (d *Device) sendCmd(x xpc.XPC, id int, args xpc.Dict) error {
	//logger.Info("send", "id", id, "args", fmt.Sprintf("%v", args))
	x.Send(xpc.Dict{"kCBMsgId": id, "kCBMsgArgs": args}, false)
	return nil
}
