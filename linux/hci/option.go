package hci

import (
	"time"

	"github.com/thomascriley/ble/linux/hci/cmd"
	"github.com/thomascriley/ble/linux/smp"
)

// An Option is a configuration function, which configures the device.
type Option func(*HCI) error

// OptDeviceID sets HCI device ID.
func OptDeviceID(id int) Option {
	return func(h *HCI) error {
		h.id = id
		return nil
	}
}

// OptDialerTimeout sets dialing timeout for Dialer.
func OptDialerTimeout(d time.Duration) Option {
	return func(h *HCI) error {
		h.dialerTmo = d
		return nil
	}
}

// OptListenerTimeout sets dialing timeout for Listener.
func OptListenerTimeout(d time.Duration) Option {
	return func(h *HCI) error {
		h.listenerTmo = d
		return nil
	}
}

// OptConnParams overrides default connection parameters.
func OptConnParams(param cmd.LECreateConnection) Option {
	return func(h *HCI) error {
		h.params.connParams = param
		return nil
	}
}

// OptPairingIO allows input and output for the pairing process.
func OptPairingCapabilities(param smp.Capabilities) Option {
	return func(h *HCI) error {
		h.smpCapabilites = param
		return nil
	}
}
