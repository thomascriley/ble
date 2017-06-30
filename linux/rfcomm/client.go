package rfcomm

import (
	"io"
	"sync"
	"time"

	"github.com/currantlabs/ble"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidArgument means one or more of the arguments are invalid.
	ErrInvalidArgument = errors.New("invalid argument")
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
		// TODO: check for matching server number
		// TODO: check ea, cr and dlci values
		if address := p[0]; address != 0x09 {
			return 0, errors.New("Invalid address byte")
		}
		controlField := p[1]
		if int(controlField&uint8(ControlNumberUIH)) != ControlNumberUIH {
			return 0, errors.New("Invalid control number")
		}
		// TODO: check that the length byte matches actual payload length
		// check for credit
		if controlField&0x10 == 0x10 {
			// TODO: handle credits
			copy(b[:], p[4:len(p)-1])
			return len(p) - 5, nil
		}
		copy(b[:], p[3:len(p)-1])
		return len(p) - 4, nil
	case err := <-c.chErr:
		return 0, err
	}
}

// Write RFCOMM Bluetooth specification Version 1.0 B
func (c *Client) Write(v []byte) (int, error) {
	return len(v), c.sendFrame(frame{
		ControlNumber:      ControlNumberUIH,
		CommmandResponse:   true,
		DLCI:               false,
		ServerNum:          1,
		PollFinal:          false,
		Payload:            v,
		FrameCheckSequence: 0x9a})
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

func (c *Client) sendSABMFrame() error {

	err := c.sendFrame(frame{
		ControlNumber:      ControlNumberSABM,
		CommmandResponse:   true,
		DLCI:               false,
		ServerNum:          0,
		PollFinal:          true,
		Payload:            []byte{},
		FrameCheckSequence: 0x1c})
	if err != nil {
		return err
	}

	select {
	case p, ok := <-c.p2p:
		if !ok {
			return errors.New("Channel closed")
		}
		var frm frame
		if err = frm.Unmarshal(p); err != nil {
			return err
		}
		switch int(frm.ControlNumber) {
		case ControlNumberDISC:
			return errors.New("Receveived disconnect")
		case ControlNumberDM:
			return errors.New("Received disconnect mode")
		case ControlNumberUIH:
			fallthrough
		case ControlNumberSABM:
			return errors.New("Received unexepcted control number")
		case ControlNumberUA:
			return nil
		}
	case <-time.After(60 * time.Second):
		// TODO: must send a ControlNumberDISC to the SABM channel
		return errors.New("Timed out")
	}
	return nil
}

func (c *Client) sendFrame(frm frame) error {
	txBuf := <-c.chTxBuf
	defer func() { c.chTxBuf <- txBuf }()

	n, err := frm.Marshal(txBuf)
	if err != nil {
		return err
	}

	if n > c.l2c.TxMTU() {
		return ErrInvalidArgument
	}

	_, err = c.l2c.Write(txBuf[:n])
	return err
}

// Loop ...
func (c *Client) loop() {

	for {
		n, err := c.l2c.Read(c.rxBuf)
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
