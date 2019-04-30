package socket

import "io"

type Closer interface {
	io.ReadWriteCloser
	Closed() chan struct{}
}
