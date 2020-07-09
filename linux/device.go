package linux

import (
	"fmt"
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"time"
	"errors"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/att"
	"github.com/thomascriley/ble/linux/gatt"
	"github.com/thomascriley/ble/linux/hci"
)

// Device ...
type Device struct {
	sync.WaitGroup
	HCI           *hci.HCI
	Server        *gatt.Server
	scanCtx       context.Context
	scanCancel    context.CancelFunc
	inquireCtx    context.Context
	inquireCancel context.CancelFunc
	numResponses  int
	allowDup      bool
	error         error
}

// NewDevice returns the default HCI device.
func NewDevice() *Device {
	return &Device{HCI: hci.NewHCI()}
}

func (d *Device) Initialize(ctx context.Context) error {
	err := d.HCI.Init(ctx)
	switch {
	case errors.Is(err,ble.ErrAlreadyInitialized):
		return err
	case err == nil :
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
				log.Printf("can't accept: %s", err)
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
			fmt.Printf("can't create ATT server: %s", err)
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

// Scan starts scanning. Duplicated advertisements will be filtered out if allowDup is set to false.
func (d *Device) Scan(ctx context.Context, allowDup bool, h ble.AdvHandler) error {
	if err := d.HCI.SetAdvHandler(h); err != nil {
		return fmt.Errorf("unable to set advertisement handler: %s", err)
	}
	if err := d.HCI.Scan(ctx, allowDup); err != nil {
		return fmt.Errorf("unable to start scan: %w", err)
	}
	d.allowDup = allowDup

	// scan until the context or socket close
	d.scanCtx, d.scanCancel = context.WithCancel(ctx)
	defer d.scanCancel()
	select {
	case <-d.scanCtx.Done():
		return d.HCI.StopScanning(context.Background())
	case <-d.HCI.Closed():
		return errors.New("hci device closed")
	}
}

// Inquire starts inquiring for bluetooth devices broadcasting using br/edr
func (d *Device) Inquire(ctx context.Context, interval time.Duration, numResponses int, h ble.InqHandler) error {
	d.numResponses = numResponses

	if err := d.HCI.SetInqHandler(h); err != nil {
		return fmt.Errorf("unable to set inquiry handler: %w", err)
	}

	// inquiries do not run indefinitely but for specific intervals of time. To mimic indefinite scanning: periodically
	// restart the scan with the timeout ends
	for {
		d.inquireCtx, d.inquireCancel = context.WithTimeout(ctx,interval)
		defer d.inquireCancel()

		if err := d.HCI.Inquire(ctx, int(float64(interval) / 1.28), 255); err != nil {
			return fmt.Errorf("unable to start inquiry: %w", err)
		}

		// inquire until the context closes or the socket closes
		select {
		case <-d.inquireCtx.Done():
			// stop the current inquiry
			if err := d.HCI.StopInquiry(context.Background()); err != nil {
				return fmt.Errorf("unable to stop inquiry: %w", err)
			}
			// check if it's time to exit, otherwise restart the inquiry
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		case <-d.HCI.Closed():
			return errors.New("hci device closed")
		}
	}
}

// RequestRemoteName ...
func (d *Device) RequestRemoteName(ctx context.Context, a ble.Addr) (string, error) {
	return d.HCI.RequestRemoteName(ctx, a)
}

// tempStop temporarily stops scanning or inquiring (this is done when connecting)
func (d *Device) tempStop() error {
	if d.scanCtx != nil {
		select {
		case <-d.scanCtx.Done():
		default:
			if err := d.HCI.StopScanning(context.Background()); err != nil {
				return err
			}
		}
	}
	if d.inquireCtx != nil {
		select {
		case <-d.inquireCtx.Done():
		default:
			if err := d.HCI.StopInquiry(context.Background()); err != nil {
				return err
			}
		}
	}
	return nil
}

// tempStop restarts a temporary stoppage
func (d *Device) tempStart() error {
	errored := make([]string, 0)
	if d.scanCtx != nil {
		select {
		case <-d.scanCtx.Done():
		default:
			// restart scan, if unable to cancel the context to stop the original request
			if err := d.HCI.Scan(d.scanCtx, d.allowDup); err != nil {
				errored = append(errored, fmt.Sprintf("could not restart scan: %s", err.Error()))
			}
		}
	}
	if d.inquireCtx != nil {
		select {
		case <-d.inquireCtx.Done():
		default:
			deadline, _ := d.inquireCtx.Deadline()
			length := int(deadline.Sub(time.Now()).Seconds() / 1.28)
			if err := d.HCI.Inquire(d.inquireCtx, length, d.numResponses); err != nil {
				errored = append(errored, fmt.Sprintf("could not restart inquiry: %s", err.Error()))
			}
		}
	}
	if len(errored) > 0 {
		if d.scanCancel != nil {
			d.scanCancel()
		}
		if d.inquireCancel != nil {
			d.inquireCancel()
		}
		return fmt.Errorf("%w: %s", ble.ErrRestartScan, strings.Join(errored, ","))
	}
	return nil
}

// Dial ...
func (d *Device) DialBLE(ctx context.Context, address ble.Addr, addressType ble.AddressType) (ble.ClientBLE, error) {
	// d.HCI.Dial is a blocking call, although most of time it should return immediately.
	// But in case passing wrong device address or the device went non-connectable, it blocks.
	// stopping the scan will improve ability to connect
	select {
	case <-d.HCI.Closed():
		return nil, errors.New("hci device is down")
	default:
	}
	if err := d.tempStop(); err != nil {
		return nil, fmt.Errorf("unable to temporary stop scan: %w", err)
	}
	cln, err := d.HCI.Dial(ctx, address, addressType)
	if err := d.tempStart(); err != nil {
		return cln, err
	}
	return cln, fmt.Errorf("can't dial ble: %w", err)
}

// DialRFCOMM ...
// TODO: implement SDP to determine RFCOMM channel number
func (d *Device) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode uint8, channel uint8) (ble.ClientRFCOMM, error) {
	// d.HCI.DialRFCOMM is a blocking call, although most of time it should return immediately.
	// But in case passing wrong device address or the device went non-connectable, it blocks.
	select {
	case <-d.HCI.Closed():
		return nil, errors.New("hci device is down")
	default:
	}
	if err := d.tempStop(); err != nil {
		return nil, fmt.Errorf("unable to temporary stop scan: %w", err)
	}
	cln, err := d.HCI.DialRFCOMM(ctx, a, clockOffset, pageScanRepetitionMode, channel)
	if err := d.tempStart(); err != nil {
		return cln, err
	}
	return cln, fmt.Errorf( "can't dial rfcomm: %w", err)
}


