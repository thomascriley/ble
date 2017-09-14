package dev

import (
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/darwin"
)

// DefaultDevice ...
func DefaultDevice() (d ble.Device, err error) {
	return darwin.NewDevice()
}
