package ble

import "io"

// A Client is a GATT client.
type RFCOMMClient interface {
	io.ReadWriter

	// Address returns platform specific unique ID of the remote peripheral, e.g. MAC on Linux, Client UUID on OS X.
	Address() Addr

	// InformationRequest Information requests are used to request implementation specific information from a remote L2CAP entity [Vol 3, Part A, 4.10]
	InformationRequest(infoType uint16) error

	// CancelConnection disconnects the connection.
	CancelConnection() error

	// Disconnected returns a receiving channel, which is closed when the client disconnects.
	Disconnected() <-chan struct{}
}
