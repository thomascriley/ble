// +build !linux

package socket

import (
	"fmt"
)

// NewSocket is a dummy function for non-Linux platform.
func NewSocket(id int) (Closer, error) {
	return nil, fmt.Errorf("only available on linux")
}
