package ble

import "context"

type Client interface {
	// Address returns platform specific unique ID of the remote peripheral, e.g. MAC on Linux, Client UUID on OS X.
	Address() Addr

	// CancelConnection disconnects the connection.
	CancelConnection(ctx context.Context) error

	// Disconnected returns a receiving channel, which is closed when the client disconnects.
	Disconnected() <-chan struct{}
}
