package ble

import (
	"fmt"
	"strings"
)

type AddressType uint8

const (
	AddressTypePublic AddressType = iota
	AddressTypeRandom
)

// Addr represents a network end point address.
// It's MAC address on Linux or Device UUID on OS X.
type Addr interface {
	fmt.Stringer
}

// NewAddr creates an Addr from string
func NewAddr(s string) Addr {
	return addr(strings.ToLower(s))
}

type addr string

func (a addr) String() string {
	return string(a)
}
