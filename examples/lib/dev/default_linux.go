package dev

import (
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux"
)

// DefaultDevice ...
func DefaultDevice() (d ble.Device, err error) {
	return linux.NewDevice()
}
