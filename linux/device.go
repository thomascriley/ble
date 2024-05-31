package linux

import (
	"context"
	"errors"
	"fmt"
	"github.com/thomascriley/ble/log"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/att"
	"github.com/thomascriley/ble/linux/gatt"
	"github.com/thomascriley/ble/linux/hci"
)

// Device ...
type Device struct {
	sync.WaitGroup
	HCI          *hci.HCI
	Server       *gatt.Server
	numResponses int
	allowDup     bool
	interval     time.Duration

	scanMutex       sync.Mutex
	scanErr         chan error
	scanning        bool
	scanTempStopped chan bool
	scanRequested   bool

	inquireMutex       sync.Mutex
	inquireErr         chan error
	inquiring          bool
	inquireTempStopped chan bool
	inquireRequested   bool

	log *slog.Logger
}

// NewDevice returns the default HCI device.
func NewDevice(log *slog.Logger) *Device {
	log = log.With("package", "github.com/thomascriley/ble")
	d := &Device{
		HCI:                hci.NewHCI(log),
		scanErr:            make(chan error, 1),
		inquireErr:         make(chan error, 1),
		scanTempStopped:    make(chan bool),
		inquireTempStopped: make(chan bool),
		log:                log,
	}
	close(d.scanTempStopped)
	close(d.inquireTempStopped)
	return d
}

func (d *Device) Initialize(ctx context.Context) error {
	err := d.HCI.Init(ctx)
	switch {
	case errors.Is(err, ble.ErrAlreadyInitialized):
		return err
	case err == nil:
		return nil
	default:
		_ = d.HCI.Close()
		return fmt.Errorf("can't init hci: %w", err)
	}
}

// Address returns the listener's device address.
func (d *Device) Address() ble.Addr {
	return d.HCI.Addr()
}

// Closes the gatt server.
func (d *Device) Close() error {
	defer d.Wait()
	return d.HCI.Close()
}

// Closed ...
func (d *Device) Closed() <-chan struct{} {
	return d.HCI.Closed()
}

// blocking call
func (d *Device) Serve(name string, handler ble.NotifyHandler) (err error) {
	if d.Server, err = gatt.NewServerWithNameAndHandler(name, handler); err != nil {
		return fmt.Errorf("can't create server: %w", err)
	}

	// mtu := ble.DefaultMTU
	mtu := ble.MaxMTU // TODO: get this from user using Option.
	//if mtu > ble.MaxMTU {
	//	return  fmt.Errorf( "maximum ATT_MTU is %d", ble.MaxMTU)
	//}

	for {
		l2c, err := d.HCI.Accept()
		if err != nil {
			// An EOF error indicates that the HCI socket was closed during
			// the read.  Don't report this as an error.
			if errors.Is(err, io.EOF) {
				d.log.Debug("can't accept", log.Error(err))
				continue
			}
			if err2 := d.HCI.Close(); err2 != nil {
				return fmt.Errorf("could not accept connections: %w and could not close: %s", err, err2)
			}
			return fmt.Errorf("could not accept connections: %w", err)
		}

		// Initialize the per-connection cccd values.
		//l2c.SetContext(context.WithValue(l2c.Context(), ble.ContextKeyCCC, make(map[uint16]uint16)))
		l2c.SetRxMTU(mtu)

		d.Server.Lock()
		as, err := att.NewServer(d.Server.DB(), l2c)
		d.Server.Unlock()
		if err != nil {
			d.log.Debug("can't create ATT server", log.Error(err))
			continue
		}

		d.Add(1)
		go func() {
			defer d.Done()
			as.Loop()
		}()
	}
}

// AddService adds a service to database.
func (d *Device) AddService(svc *ble.Service) error {
	return d.Server.AddService(svc)
}

// RemoveAllServices removes all services that are currently in the database.
func (d *Device) RemoveAllServices() error {
	return d.Server.RemoveAllServices()
}

// SetServices set the specified service to the database.
// It removes all currently added services, if any.
func (d *Device) SetServices(svcs []*ble.Service) error {
	return d.Server.SetServices(svcs)
}

func (d *Device) Advertise(ctx context.Context, adv ble.Advertisement) error {
	if err := d.HCI.AdvertiseAdv(ctx, adv); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)

}

// AdvertiseNameAndServices advertises device name, and specified service UUIDs.
// It tres to fit the UUIDs in the advertising packet as much as possible.
// If name doesn't fit in the advertising packet, it will be put in scan response.
func (d *Device) AdvertiseNameAndServices(ctx context.Context, name string, uuids ...ble.UUID) error {
	if err := d.HCI.AdvertiseNameAndServices(ctx, name, uuids...); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)
}

// AdvertiseMfgData avertises the given manufacturer data.
func (d *Device) AdvertiseMfgData(ctx context.Context, id uint16, b []byte) error {
	if err := d.HCI.AdvertiseMfgData(ctx, id, b); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)
}

// AdvertiseServiceData16 advertises data associated with a 16bit service uuid
func (d *Device) AdvertiseServiceData16(ctx context.Context, id uint16, b []byte) error {
	if err := d.HCI.AdvertiseServiceData16(ctx, id, b); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)
}

// AdvertiseIBeaconData advertise iBeacon with given manufacturer data.
func (d *Device) AdvertiseIBeaconData(ctx context.Context, b []byte) error {
	if err := d.HCI.AdvertiseIBeaconData(ctx, b); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)
}

// AdvertiseIBeacon advertises iBeacon with specified parameters.
func (d *Device) AdvertiseIBeacon(ctx context.Context, u ble.UUID, major, minor uint16, pwr int8) error {
	if err := d.HCI.AdvertiseIBeacon(ctx, u, major, minor, pwr); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case <-d.HCI.Closed():
		return errors.New("hci device is down")
	}
	return d.HCI.StopAdvertising(ctx)
}

// RequestRemoteName ...
func (d *Device) RequestRemoteName(ctx context.Context, a ble.Addr) (string, error) {
	return d.HCI.RequestRemoteName(ctx, a)
}

// Scan starts scanning. Duplicated advertisements will be filtered out if allowDup is set to false.
func (d *Device) Scan(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
	select {
	case <-d.scanTempStopped:
	case <-ctx.Done():
		return ctx.Err()
	}

	d.scanRequested = true
	defer func() { d.scanRequested = false }()

	if err := d.HCI.SetAdvHandler(h); err != nil {
		return fmt.Errorf("unable to set advertisement handler: %s", err)
	}
	if err := d.startScan(ctx, allowDup); err != nil {
		return err
	}

	// scan until the context or socket close
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	select {
	case err := <-d.scanErr:
		return err
	case <-ctx.Done():
		return d.stopScan()
	case <-d.HCI.Closed():
		return errors.New("hci device closed")
	}
}

// Inquire starts inquiring for bluetooth devices broadcasting using br/edr
func (d *Device) Inquire(ctx context.Context, interval time.Duration, numResponses int, h ble.InqHandler) error {
	select {
	case <-d.inquireTempStopped:
	case <-ctx.Done():
		return ctx.Err()
	}

	d.inquireRequested = true
	defer func() { d.inquireRequested = false }()

	d.numResponses = numResponses

	if err := d.HCI.SetInqHandler(h); err != nil {
		return fmt.Errorf("unable to set inquiry handler: %w", err)
	}

	// inquiries do not run indefinitely but for specific intervals of time. To mimic indefinite scanning: periodically
	// restart the scan with the timeout ends
	for {
		if err := d.inquire(ctx, interval); err != nil {
			return err
		}
		// check if it's time to exit, otherwise restart the inquiry
		select {
		case <-ctx.Done():
			return nil
		default:
			continue
		}
	}
}

// Dial ...
func (d *Device) DialBLE(ctx context.Context, address ble.Addr, addressType ble.AddressType) (cli ble.ClientBLE, err error) {
	// d.HCI.Dial is a blocking call, although most of time it should return immediately.
	// But in case passing wrong device address or the device went non-connectable, it blocks.
	// stopping the scan will improve ability to connect
	select {
	case <-d.HCI.Closed():
		return nil, errors.New("hci device is down")
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	fmt.Println("temporarily stopping scan")
	if err = d.tempStop(); err != nil {
		return nil, fmt.Errorf("failed to temporary stop scan: %w", err)
	}
	cli, err = d.HCI.Dial(ctx, address, addressType)
	fmt.Println("restarting scan")
	d.tempStart()
	return cli, err
}

// DialRFCOMM ...
// TODO: implement SDP to determine RFCOMM channel number
func (d *Device) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode uint8, channel uint8) (cli ble.ClientRFCOMM, err error) {
	// d.HCI.DialRFCOMM is a blocking call, although most of time it should return immediately.
	// But in case passing wrong device address or the device went non-connectable, it blocks.
	select {
	case <-d.HCI.Closed():
		return nil, errors.New("hci device is down")
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err = d.tempStop(); err != nil {
		return nil, fmt.Errorf("unable to temporary stop scan: %w", err)
	}
	cli, err = d.HCI.DialRFCOMM(ctx, a, clockOffset, pageScanRepetitionMode, channel)
	d.tempStart()
	return cli, err
}

func (d *Device) inquire(ctx context.Context, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, interval)
	defer cancel()

	if err := d.startInquiry(ctx, interval); err != nil {
		return err
	}

	// inquire until the context closes or the socket closes
	select {
	case err := <-d.inquireErr:
		return err
	case <-ctx.Done():
		return d.stopInquiry()
	case <-d.HCI.Closed():
		return errors.New("hci device closed")
	}
}

// tempStop temporarily stops scanning or inquiring (this is done when connecting)
func (d *Device) tempStop() error {
	select {
	case <-d.scanTempStopped:
		d.scanTempStopped = make(chan bool)
	default:
	}
	select {
	case <-d.inquireTempStopped:
		d.inquireTempStopped = make(chan bool)
	default:
	}
	if err := d.tempStopScan(); err != nil {
		return err
	}
	if err := d.tempStopInquiry(); err != nil {
		d.trigger(d.tempStartScan)
		return err
	}
	return nil
}

func (d *Device) tempStart() {
	select {
	case <-d.scanTempStopped:
	default:
		close(d.scanTempStopped)
	}
	select {
	case <-d.inquireTempStopped:
	default:
		close(d.inquireTempStopped)
	}
	d.trigger(d.tempStartScan, d.tempStartInquiry)
}

func (d *Device) tempStopScan() error {
	if !d.scanRequested {
		return nil
	}
	slog.Debug("BLE: temporarily stopping scan")
	if err := d.stopScan(); err != nil {
		return err
	}
	slog.Debug("BLE: temporarily stopped scan")
	return nil
}

func (d *Device) tempStopInquiry() error {
	if !d.inquireRequested {
		return nil
	}
	slog.Debug("BLE: temporarily stopping inquiry")
	if err := d.stopInquiry(); err != nil {
		return err
	}
	slog.Debug("BLE: temporarily stopped inquiry")
	return nil
}

func (d *Device) tempStartScan(ctx context.Context) {
	if !d.scanRequested {
		return
	}
	slog.Debug("BLE: temporarily starting scan")
	if err := d.startScan(ctx, d.allowDup); err != nil {
		select {
		case d.scanErr <- err:
		default:
		}
		slog.Debug("BLE: failed to temporarily start inquiry", slog.String("Error", err.Error()))
		return
	}
	slog.Debug("BLE: temporarily started scan")
}

func (d *Device) tempStartInquiry(ctx context.Context) {
	if !d.inquireRequested {
		return
	}
	slog.Debug("BLE: temporarily starting inquiry")
	if err := d.startInquiry(ctx, d.interval); err != nil {
		select {
		case d.scanErr <- err:
		default:
		}
		slog.Debug("BLE: failed to temporarily start inquiry", slog.String("Error", err.Error()))
		return
	}
	slog.Debug("BLE: temporarily started inquiry")
}

func (d *Device) trigger(fns ...func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(context.Background())

	// cancel context if device closes
	d.Add(1)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer d.Done()
		select {
		case <-ctx.Done():
		case <-d.Closed():
			cancel()
		}
	}(ctx, cancel)

	// trigger with context
	d.Add(1)
	go func(ctx context.Context) {
		defer d.Done()
		defer cancel()
		for _, fn := range fns {
			fn(ctx)
		}
	}(ctx)
}

func (d *Device) startScan(ctx context.Context, allowDup bool) error {
	d.scanMutex.Lock()
	defer d.scanMutex.Unlock()

	if d.scanning {
		return nil
	}
	if err := d.HCI.Scan(ctx, allowDup); err != nil {
		return fmt.Errorf("ble failed to start scan: %w", err)
	}
	d.allowDup = allowDup
	d.scanning = true
	return nil
}

func (d *Device) stopScan() error {
	d.scanMutex.Lock()
	defer d.scanMutex.Unlock()

	if !d.scanning {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.HCI.StopScanning(ctx); err != nil {
		return fmt.Errorf("ble failed to stop scanning: %w", err)
	}
	d.scanning = false
	return nil
}

func (d *Device) startInquiry(ctx context.Context, interval time.Duration) error {
	d.inquireMutex.Lock()
	defer d.inquireMutex.Unlock()

	if d.inquiring {
		return nil
	}
	if err := d.HCI.Inquire(ctx, int(float64(interval)/1.28), 255); err != nil {
		return fmt.Errorf("unable to start inquiry: %w", err)
	}
	d.interval = interval
	d.inquiring = true
	return nil
}

func (d *Device) stopInquiry() error {
	d.inquireMutex.Lock()
	defer d.inquireMutex.Unlock()

	if !d.inquiring {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := d.HCI.StopInquiry(ctx); err != nil {
		return fmt.Errorf("ble failed to stop inquiry: %w", err)
	}
	d.inquiring = false
	return nil
}
