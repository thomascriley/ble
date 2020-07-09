package ble

import (
	"context"
	"time"
)

// Device ...
type Device interface {
	// AddService adds a service to database.
	AddService(svc *Service) error

	// RemoveAllServices removes all services that are currently in the database.
	RemoveAllServices() error

	// SetServices set the specified service to the database.
	// It removes all currently added services, if any.
	SetServices(svcs []*Service) error

	// Closed the underlying hci socket and waits for all connections to close and goroutines to finish
	Close() error

	// Advertise advertises a given Advertisement
	Advertise(ctx context.Context, adv Advertisement) error

	// AdvertiseNameAndServices advertises device name, and specified service UUIDs.
	// It tres to fit the UUIDs in the advertising packet as much as possi
	// If name doesn't fit in the advertising packet, it will be put in scan response.
	AdvertiseNameAndServices(ctx context.Context, name string, uuids ...UUID) error

	// AdvertiseMfgData avertises the given manufacturer data.
	AdvertiseMfgData(ctx context.Context, id uint16, b []byte) error

	// AdvertiseServiceData16 advertises data associated with a 16bit service uuid
	AdvertiseServiceData16(ctx context.Context, id uint16, b []byte) error

	// AdvertiseIBeaconData advertise iBeacon with given manufacturer data.
	AdvertiseIBeaconData(ctx context.Context, b []byte) error

	// AdvertiseIBeacon advertises iBeacon with specified parameters.
	AdvertiseIBeacon(ctx context.Context, u UUID, major, minor uint16, pwr int8) error

	// Scan starts scanning. Duplicated advertisements will be filtered out if allowDup is set to false.
	Scan(ctx context.Context, allowDup bool, h AdvHandler) error

	// Inquire starts a BR/EDR scan
	Inquire(ctx context.Context, interval time.Duration, numResponses int, h InqHandler) error

	// RequestRemoteName queries the remote BR/EDR device for its name
	RequestRemoteName(ctx context.Context, a Addr) (string, error)

	// Dial ...
	DialBLE(context.Context, Addr, AddressType) (ClientBLE, error)

	// Dial
	DialRFCOMM(ctx context.Context, a Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (ClientRFCOMM, error)

	// Initializes the underlying hardware
	Initialize(context.Context) error

	// Closed returns a channel that is closed when the underlying socket closes
	Closed() <-chan struct{}
}
