package spp

import (
	"fmt"
	"io"
	"sync"
	"time"

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

	l2c ble.Conn
	p2p chan []byte

	rxBuf   []byte
	chTxBuf chan []byte
}

// NewClient returns an Attribute Protocol Client.
func NewClient(l2c ble.Conn) *Client {
	c := &Client{
		l2c:     l2c,
		p2p:     make(chan []byte),
		chTxBuf: make(chan []byte, 1),
		rxBuf:   make([]byte, ble.MaxMTU),
	}
	go c.loop()
	return c
}

// Address returns the address of the client.
func (p *Client) Address() ble.Addr {
	p.RLock()
	defer p.RUnlock()
	return p.conn.RemoteAddr()
}

// Read ...
func (c *Client) Read(b []byte) (int, error) {
	p, ok := <-c.p2p
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
}

// Write ...
func (c *Client) Write(v []byte) (int, error) {
	if len(v) > c.l2c.TxMTU() {
		return 0, ErrInvalidArgument
	}

	txBuf := <-c.chTxBuf
	defer func() { c.chTxBuf <- txBuf }()

	copy(txBuf, len(v)+4)
	copy(txBuf[2:], dcid)
	copy(txBuf[4:], scid)

	return len(v), c.sendCmd(txBuf)
}

// CancelConnection disconnects the connection.
func (p *Client) CancelConnection() error {
	p.Lock()
	defer p.Unlock()
	return p.conn.Close()
}

// Disconnected returns a receiving channel, which is closed when the client disconnects.
func (p *Client) Disconnected() <-chan struct{} {
	p.Lock()
	defer p.Unlock()
	return p.conn.Disconnected()
}

func (c *Client) sendCmd(b []byte) error {
	_, err := c.l2c.Write(b)
	return err
}

func (c *Client) sendReq(b []byte) (rsp []byte, err error) {
	if _, err := c.l2c.Write(b); err != nil {
		return nil, errors.Wrap(err, "send ATT request failed")
	}
	for {
		select {
		case rsp := <-c.rspc:
			if rsp[0] == ErrorResponseCode || rsp[0] == rspOfReq[b[0]] {
				return rsp, nil
			}
			// Sometimes when we connect to an Apple device, it sends
			// ATT requests asynchronously to us. // In this case, we
			// return an ErrReqNotSupp response, and continue to wait
			// the response to our request.
			errRsp := newErrorResponse(rsp[0], 0x0000, ble.ErrReqNotSupp)
			logger.Debug("client", "req", fmt.Sprintf("% X", b))
			_, err := c.l2c.Write(errRsp)
			if err != nil {
				return nil, errors.Wrap(err, "unexpected ATT response recieved")
			}
		case err := <-c.chErr:
			return nil, errors.Wrap(err, "ATT request failed")
		case <-time.After(30 * time.Second):
			return nil, errors.Wrap(ErrSeqProtoTimeout, "ATT request timeout")
		}
	}
}

// Loop ...
func (c *Client) loop() {

	for {
		n, err := c.l2c.Read(c.rxBuf)
		logger.Debug("client", "rsp", fmt.Sprintf("% X", c.rxBuf[:n]))
		if err != nil {
			// We don't expect any error from the bearer (L2CAP ACL-U)
			// Pass it along to the pending request, if any, and escape.
			c.chErr <- err
			return
		}

		b := make([]byte, n)
		copy(b, c.rxBuf)
		c.p2p <- b
	}
}
