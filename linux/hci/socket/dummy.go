// +build !linux

package socket

// NewSocket is a dummy function for non-Linux platform.
func NewSocket(id int) (Closer, error) {
	return nil, nil
}
