package ble

import "io"

// A Client is a GATT client.
type RFCOMMClient interface {
	io.ReadWriter

	// Address returns platform specific unique ID of the remote peripheral, e.g. MAC on Linux, Client UUID on OS X.
	Address() Addr

	// CancelConnection disconnects the connection.
	CancelConnection() error

	// Disconnected returns a receiving channel, which is closed when the client disconnects.
	Disconnected() <-chan struct{}
}
