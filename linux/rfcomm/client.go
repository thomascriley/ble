package rfcomm

import (
	"fmt"
	"io"
	"sync"

	"github.com/currantlabs/ble"
	"github.com/pkg/errors"
)

// TODO: Added the following LMP Methods
// TODO: Authentication Requested [Vol 2, Part F, 4.1]
// TODO: Simple Pairing Message Sequence [Vol 2, Part F, 4.2]
// TODO: Link Supervision Timeout Event may be triggered if supported [Vol 2, Part F, 4.3]
// TODO: Set Connection Encryption [Vol 2, Part F, 4.4]
// TODO: Change Connection Link Key [Vol 2, Part F, 4.5]
// TODO: Change Connection Link Key With Encryption Pause and Resume [Vol 2, Part F, 4.6]
// TODO: Master Link Key [Vol 2, Part F, 4.7]
// TODO: Read Remote Supported Features [Vol 2, 4.8]
// TODO: Read Remote Extended Features  (if supported) [Vol 2, 4.9]
// TODO: Read Clock Offset [Vol 2, 4.10]
// TODO: Role Switch on an Encrypted Link Using Encryption Pause and Resume [Vol 2, 4.11]
// TODO: Refreshing Encryption Keys [Vol 2, 4.12]
// TODO: Read Remote Version Information [Vol 2, 4.13]
// TODO: QoS Setup [Vol 2, 4.14]
// TODO: Switch Role [Vol 2, 4.15]
// TODO: AMP Physical Link Creation and Disconnect [Vol 2, 4.16]

// Client implementa an Attribute Protocol Client.
type Client struct {
	sync.RWMutex

	l2c   ble.Conn
	p2p   chan []byte
	chErr chan error

	rxBuf   []byte
	chTxBuf chan []byte
}

// NewClient returns an Attribute Protocol Client.
func NewClient(l2c ble.Conn) *Client {
	c := &Client{
		l2c:     l2c,
		p2p:     make(chan []byte),
		chErr:   make(chan error),
		chTxBuf: make(chan []byte, 1),
		rxBuf:   make([]byte, ble.MaxACLMTU),
	}
	go c.loop()
	return c
}

// Address returns the address of the client.
func (c *Client) Address() ble.Addr {
	c.RLock()
	defer c.RUnlock()
	return c.l2c.RemoteAddr()
}

// Read ...
func (c *Client) Read(b []byte) (int, error) {
	select {
	case p, ok := <-c.p2p:
		if !ok {
			return 0, errors.Wrap(io.ErrClosedPipe, "input channel closed")
		}
		if len(p) == 0 {
			return 0, errors.Wrap(io.ErrUnexpectedEOF, "recieved empty packet")
		}
		if len(p) > len(b) {
			return 0, errors.Wrapf(io.ErrShortBuffer, "payload recieved exceeds sdu buffer")
		}
		copy(b[:], p)
		return len(p), nil
	case err := <-c.chErr:
		return 0, err
	}
}

// Write ...
func (c *Client) Write(v []byte) (int, error) {
	if len(v) > c.l2c.TxMTU() {
		return 0, ErrInvalidArgument
	}

	txBuf := <-c.chTxBuf
	defer func() { c.chTxBuf <- txBuf }()

	copy(txBuf, len(v)+4)
	copy(txBuf[2:], c.l2c.DestinationID)
	copy(txBuf[4:], c.l2c.SourceID)

	_, err := c.l2c.Write(b)
	return len(v), err
}

// CancelConnection disconnects the connection.
func (c *Client) CancelConnection() error {
	c.Lock()
	defer c.Unlock()
	return c.l2c.Close()
}

// Disconnected returns a receiving channel, which is closed when the client disconnects.
func (c *Client) Disconnected() <-chan struct{} {
	c.Lock()
	defer c.Unlock()
	return c.l2c.Disconnected()
}

// Loop ...
func (c *Client) loop() {

	for {
		n, err := c.l2c.Read(c.rxBuf)
		logger.Debug("client", "rsp", fmt.Sprintf("% X", c.rxBuf[:n]))
		if err != nil {
			// We don't expect any error from the bearer (L2CAP ACL-U)
			// Pass it along to the pending read, if any, and escape.
			c.chErr <- err
			return
		}

		b := make([]byte, n)
		copy(b, c.rxBuf)
		c.p2p <- b
	}
}
