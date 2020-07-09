package ble

import (
	"io"
)

type ClientRFCOMM interface {
	Client

	io.ReadWriter
}
