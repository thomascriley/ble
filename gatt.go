package ble

import (
	"context"
	"fmt"
	"errors"
	"time"
)

// ErrDefaultDevice ...
var ErrDefaultDevice = errors.New("default device is not set")

var defaultDevice Device

// SetDefaultDevice returns the default HCI device.
func SetDefaultDevice(d Device) {
	defaultDevice = d
}

// AddService adds a service to database.
func AddService(svc *Service) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.AddService(svc)
}

// RemoveAllServices removes all services that are currently in the database.
func RemoveAllServices() error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.RemoveAllServices()
}

// SetServices set the specified service to the database.
// It removes all currently added services, if any.
func SetServices(svcs []*Service) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.SetServices(svcs)
}

// Close stop the GATT server
func Close() error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.Close()
}

// AdvertiseNameAndServices advertises device name, and specified service UUIDs.
// It tres to fit the UUIDs in the advertising packet as much as possi
// If name doesn't fit in the advertising packet, it will be put in scan response.
func AdvertiseNameAndServices(ctx context.Context, name string, uuids ...UUID) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.AdvertiseNameAndServices(ctx, name, uuids...)
}

// AdvertiseIBeaconData advertise iBeacon with given manufacturer data.
func AdvertiseIBeaconData(ctx context.Context, b []byte) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.AdvertiseIBeaconData(ctx, b)
}

// AdvertiseIBeacon advertises iBeacon with specified parameters.
func AdvertiseIBeacon(ctx context.Context, u UUID, major, minor uint16, pwr int8) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.AdvertiseIBeacon(ctx, u, major, minor, pwr)
}

// Scan starts scanning. Duplicated advertisements will be filtered out if allowDup is set to false.
func Scan(ctx context.Context, allowDup bool, h AdvHandler, f AdvFilter) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}

	if f == nil {
		return defaultDevice.Scan(ctx, allowDup, h)
	}

	h2 := func(a Advertisement) {
		if f(a) {
			h(a)
		}
	}
	return defaultDevice.Scan(ctx, allowDup, h2)
}

func Inquire(ctx context.Context, interval time.Duration, numResponses int, h InqHandler) error {
	if defaultDevice == nil {
		return ErrDefaultDevice
	}
	return defaultDevice.Inquire(ctx, interval, numResponses, h)
}

// Find ...
func Find(ctx context.Context, allowDup bool, f AdvFilter) ([]Advertisement, error) {
	if defaultDevice == nil {
		return nil, ErrDefaultDevice
	}
	var advs []Advertisement
	h := func(a Advertisement) {
		advs = append(advs, a)
	}
	return advs, Scan(ctx, allowDup, h, f)
}

// Dial ...
func DialBLE(ctx context.Context, address Addr, addressType AddressType) (Client, error) {
	if defaultDevice == nil {
		return nil, ErrDefaultDevice
	}
	return defaultDevice.DialBLE(ctx, address, addressType)
}

// DialRFCOMM ...
func DialRFCOMM(ctx context.Context, a Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (ClientRFCOMM, error) {
	if defaultDevice == nil {
		return nil, ErrDefaultDevice
	}
	return defaultDevice.DialRFCOMM(ctx, a, clockOffset, pageScanRepetitionMode, channel)
}

// Connect searches for and connects to a Peripheral which matches specified condition.
func Connect(ctx context.Context, f AdvFilter) (Client, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan Advertisement)
	fn := func(a Advertisement) {
		select {
		case ch <- a:
			cancel()
		case <-ctx.Done():
			return
		}
	}
	if err := Scan(ctx2, false, fn, f); err != nil {
		if err != context.Canceled {
			return nil, fmt.Errorf("can't scan: %w", err)
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case adv := <-ch:
		cln, err := DialBLE(ctx, adv.Address(), adv.AddressType())
		return cln, fmt.Errorf("can't dial: %w", err)
	}
}

// A NotificationHandler handles notification or indication from a server.
type NotificationHandler func(req []byte)

