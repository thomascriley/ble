package darwin

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/JuulLabs-OSS/cbgo"
	"github.com/thomascriley/ble"
	"time"

	"sync"
)

var (
	ErrClosed = errors.New("closed")
	ErrNotSupported = errors.New("not supported")
)

// Device is either a Peripheral or Central device.
type Device struct {
	// Embed these two bases so we don't have to override all the esoteric
	// functions defined by CoreBluetooth delegate interfaces.
	cbgo.CentralManagerDelegateBase
	cbgo.PeripheralManagerDelegateBase

	scannerOpts *cbgo.CentralManagerScanOpts

	cm  cbgo.CentralManager
	pm  cbgo.PeripheralManager
	evl deviceEventListener
	pc  profCache

	conns    map[string]*conn
	connLock sync.Mutex

	advHandler ble.AdvHandler

	closed chan struct{}
}

func (d *Device) Inquire(ctx context.Context, interval time.Duration, numResponses int, h ble.InqHandler) error {
	return ErrNotSupported
}

func (d *Device) RequestRemoteName(ctx context.Context, a ble.Addr) (string, error) {
	return "", ErrNotSupported
}

func (d *Device) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (ble.ClientRFCOMM, error) {
	return nil, ErrNotSupported
}

// NewDevice returns a BLE device.
func NewDevice() *Device {
	d := &Device{
		cm:    cbgo.NewCentralManager(nil),
		pm:    cbgo.NewPeripheralManager(nil),
		pc:    newProfCache(),
		conns: make(map[string]*conn),
		closed: make(chan struct{}),
	}
	return d
}

func (d *Device) Initialize(ctx context.Context) error {
	blockUntilStateChange := func(getState func() cbgo.ManagerState) {
		if getState() != cbgo.ManagerStateUnknown {
			return
		}

		// Wait until state changes or until one second passes (whichever
		// happens first).
		for {
			select {
			case <-d.evl.stateChanged.Listen():
				if getState() != cbgo.ManagerStateUnknown {
					return
				}

			case <-time.NewTimer(time.Second).C:
				return
			}
		}
	}

	// Ensure central manager is ready.
	d.cm.SetDelegate(d)
	blockUntilStateChange(d.cm.State)
	if d.cm.State() != cbgo.ManagerStatePoweredOn {
		return fmt.Errorf("central manager has invalid state: have=%d want=%d: is Bluetooth turned on",
			d.cm.State(), cbgo.ManagerStatePoweredOn)
	}

	// Ensure peripheral manager is ready.
	d.pm.SetDelegate(d)
	blockUntilStateChange(d.pm.State)
	if d.pm.State() != cbgo.ManagerStatePoweredOn {
		return fmt.Errorf("peripheral manager has invalid state: have=%d want=%d: is Bluetooth turned on",
			d.pm.State(), cbgo.ManagerStatePoweredOn)
	}

	return nil
}

// Scan ...
func (d *Device) Scan(ctx context.Context, allowDup bool, h ble.AdvHandler) (err error) {
	d.advHandler = h

	d.scannerOpts = &cbgo.CentralManagerScanOpts{
		AllowDuplicates: allowDup,
	}
	d.cm.Scan(nil, d.scannerOpts)

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-d.closed:
		err = ErrClosed
	}

	d.scannerOpts = nil
	d.cm.StopScan()

	return err
}

// Dial ...
func (d *Device) DialBLE(ctx context.Context, addr ble.Addr, addrType ble.AddressType) (ble.ClientBLE, error) {
	d.cm.StopScan()
	defer func(){
		if d.scannerOpts != nil {
			d.cm.Scan(nil, d.scannerOpts)
		}
	}()

	uuid, err := cbgo.ParseUUID(uuidStrWithDashes(addr.String()))
	if err != nil {
		return nil, fmt.Errorf("dial failed: invalid peer address: %s", addr)
	}

	prphs := d.cm.RetrievePeripheralsWithIdentifiers([]cbgo.UUID{uuid})
	if len(prphs) == 0 {
		return nil, fmt.Errorf("dial failed: no peer with address: %s", addr)
	}

	ch := d.evl.connected.Listen()
	defer d.evl.connected.Close()

	d.cm.Connect(prphs[0], nil)
	select {
	case <-ctx.Done():
		d.cm.CancelConnect(prphs[0])
		return nil, ctx.Err()
	case <-d.closed:
		d.cm.CancelConnect(prphs[0])
		return nil, ErrClosed
	case itf := <-ch:
		if itf == nil {
			return nil, fmt.Errorf("connect failed: aborted")
		}

		ev := itf.(*eventConnected)
		if ev.err != nil {
			return nil, ev.err
		} else {
			ev.conn.SetContext(ctx)
			return NewClient(d.cm, ev.conn)
		}
	}
}

// Stop ...
func (d *Device) Stop() error {
	return nil
}

func (d *Device) Closed() <-chan struct{} {
	return d.closed
}

func (d *Device) Close() error {
	d.connLock.Lock()
	defer d.connLock.Unlock()

	select {
	case <-d.closed:
	default:
		close(d.closed)
	}

	for _, c := range d.conns {
		_ = c.Close(context.Background())
	}
	return nil
}

func (d *Device) findConn(a ble.Addr) (cn *conn, ok bool) {
	d.connLock.Lock()
	defer d.connLock.Unlock()

	if cn, ok = d.conns[a.String()]; cn == nil {
		delete(d.conns,a.String())
		return nil, false
	}
	return cn, ok
}

func (d *Device) addConn(c *conn) error {
	d.connLock.Lock()
	defer d.connLock.Unlock()

	if conn, ok := d.conns[c.addr.String()]; ok && conn != nil {
		_ = conn.Close(context.Background())
	}
	d.conns[c.addr.String()] = c
	return nil
}

func (d *Device) delConn(a ble.Addr) {
	d.connLock.Lock()
	defer d.connLock.Unlock()
	delete(d.conns, a.String())
}

func (d *Device) connectFail(err error) {
	d.evl.connected.RxSignal(&eventConnected{
		err: err,
	})
}

func chrPropPerm(c *ble.Characteristic) (cbgo.CharacteristicProperties, cbgo.AttributePermissions) {
	var prop cbgo.CharacteristicProperties
	var perm cbgo.AttributePermissions

	if c.Property&ble.CharRead != 0 {
		prop |= cbgo.CharacteristicPropertyRead
		if ble.CharRead&c.Secure != 0 {
			perm |= cbgo.AttributePermissionsReadEncryptionRequired
		} else {
			perm |= cbgo.AttributePermissionsReadable
		}
	}
	if c.Property&ble.CharWriteNR != 0 {
		prop |= cbgo.CharacteristicPropertyWriteWithoutResponse
		if c.Secure&ble.CharWriteNR != 0 {
			perm |= cbgo.AttributePermissionsWriteEncryptionRequired
		} else {
			perm |= cbgo.AttributePermissionsWriteable
		}
	}
	if c.Property&ble.CharWrite != 0 {
		prop |= cbgo.CharacteristicPropertyWrite
		if c.Secure&ble.CharWrite != 0 {
			perm |= cbgo.AttributePermissionsWriteEncryptionRequired
		} else {
			perm |= cbgo.AttributePermissionsWriteable
		}
	}
	if c.Property&ble.CharNotify != 0 {
		if c.Secure&ble.CharNotify != 0 {
			prop |= cbgo.CharacteristicPropertyNotifyEncryptionRequired
		} else {
			prop |= cbgo.CharacteristicPropertyNotify
		}
	}
	if c.Property&ble.CharIndicate != 0 {
		if c.Secure&ble.CharIndicate != 0 {
			prop |= cbgo.CharacteristicPropertyIndicateEncryptionRequired
		} else {
			prop |= cbgo.CharacteristicPropertyIndicate
		}
	}

	return prop, perm
}

func (d *Device) AddService(svc *ble.Service) error {
	chrMap := make(map[*ble.Characteristic]cbgo.Characteristic)
	dscMap := make(map[*ble.Descriptor]cbgo.Descriptor)

	msvc := cbgo.NewMutableService(cbgo.UUID(svc.UUID), true)

	var mchrs []cbgo.MutableCharacteristic
	for _, c := range svc.Characteristics {
		prop, perm := chrPropPerm(c)
		mchr := cbgo.NewMutableCharacteristic(cbgo.UUID(c.UUID), prop, c.Value, perm)

		var mdscs []cbgo.MutableDescriptor
		for _, d := range c.Descriptors {
			mdsc := cbgo.NewMutableDescriptor(cbgo.UUID(d.UUID), d.Value)
			mdscs = append(mdscs, mdsc)
			dscMap[d] = mdsc.Descriptor()
		}
		mchr.SetDescriptors(mdscs)

		mchrs = append(mchrs, mchr)
		chrMap[c] = mchr.Characteristic()
	}
	msvc.SetCharacteristics(mchrs)

	ch := d.evl.svcAdded.Listen()
	d.pm.AddService(msvc)

	itf := <-ch
	if itf != nil {
		return itf.(error)
	}

	d.pc.addSvc(svc, msvc.Service())
	for chr, cbc := range chrMap {
		d.pc.addChr(chr, cbc)
	}
	for dsc, cbd := range dscMap {
		d.pc.addDsc(dsc, cbd)
	}

	return nil
}

func (d *Device) RemoveAllServices() error {
	d.pm.RemoveAllServices()
	return nil
}

func (d *Device) SetServices(svcs []*ble.Service) error {
	if err := d.RemoveAllServices(); err != nil {
		return fmt.Errorf("failed to remove existing services: %w", err)
	}
	for _, s := range svcs {
		if err := d.AddService(s); err != nil {
			return fmt.Errorf("failed to add service `%s`: %w", s.UUID, err)
		}
	}

	return nil
}

func (d *Device) stopAdvertising() error {
	d.pm.StopAdvertising()
	return nil
}

func (d *Device) advData(ctx context.Context, ad cbgo.AdvData) error {
	ch := d.evl.advStarted.Listen()
	d.pm.StartAdvertising(ad)

	itf := <-ch
	if itf != nil {
		return itf.(error)
	}

	<-ctx.Done()
	_ = d.stopAdvertising()
	return ctx.Err()
}

func (d *Device) Advertise(ctx context.Context, adv ble.Advertisement) error {
	ad := cbgo.AdvData{}

	ad.LocalName = adv.LocalName()
	for _, u := range adv.Services() {
		ad.ServiceUUIDs = append(ad.ServiceUUIDs, cbgo.UUID(u))
	}

	return d.advData(ctx, ad)
}

func (d *Device) AdvertiseNameAndServices(ctx context.Context, name string, uuids ...ble.UUID) error {
	a := &adv{
		localName: name,
		svcUUIDs:  uuids,
	}

	return d.Advertise(ctx, a)
}

func (d *Device) AdvertiseMfgData(ctx context.Context, id uint16, b []byte) error {
	// CoreBluetooth doesn't let you specify manufacturer data :(
	return ErrNotSupported
}

func (d *Device) AdvertiseServiceData16(ctx context.Context, id uint16, b []byte) error {
	// CoreBluetooth doesn't let you specify service data :(
	return ErrNotSupported
}

func (d *Device) AdvertiseIBeaconData(ctx context.Context, b []byte) error {
	ad := cbgo.AdvData{
		IBeaconData: b,
	}
	return d.advData(ctx, ad)
}

func (d *Device) AdvertiseIBeacon(ctx context.Context, u ble.UUID, major, minor uint16, pwr int8) error {
	b := make([]byte, 21)
	copy(b, ble.Reverse(u))                   // Big endian
	binary.BigEndian.PutUint16(b[16:], major) // Big endian
	binary.BigEndian.PutUint16(b[18:], minor) // Big endian
	b[20] = uint8(pwr)                        // Measured Tx Power
	return d.AdvertiseIBeaconData(ctx, b)
}
