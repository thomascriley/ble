package rfcomm

import (
	"io"
	"sync"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux/l2cap"
	"github.com/currantlabs/ble/linux/multiplexor"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidArgument means one or more of the arguments are invalid.
	ErrInvalidArgument = errors.New("invalid argument")
)

var (
	serverChannelNumbers = newServerChannels()
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

	dlci         uint8
	priority     uint8
	maxFrameSize uint16

	serverChannel int
}

// NewClient returns an RFCOMM Client that has been initialized according to the
// RFCOMM specifications.
func NewClient(l2c ble.Conn) *Client {
	c := &Client{
		l2c:     l2c,
		p2p:     make(chan []byte),
		chErr:   make(chan error),
		chTxBuf: make(chan []byte, 1),
		rxBuf:   make([]byte, ble.MaxACLMTU),
	}
	go c.loop()

	if err := c.l2c.InformationRequest(l2cap.InfoTypeConnectionlessMTU, timeout); err != nil {
		return nil, err
	}
	if err := c.l2c.InformationRequest(l2cap.InfoTypeExtendedFeatures, timeout); err != nil {
		return nil, err
	}
	// 1.2 - 2.1 + EDR will return not supported
	c.l2c.InformationRequest(l2cap.InfoTypeFixedChannels, timeout)

	if err := c.l2c.ConnectionRequest(psmRFCOMM, timeout); err != nil {
		return nil, err
	}

	// Even if all default values are acceptable, a Configuration Request
	// packet with no options shall be sent. [Vol 3, Part A, 4.4]
	// TODO: make this non-static
	options := []l2cap.Option{&l2cap.MTUOption{MTU: 0x03f5}}
	if err := c.l2c.ConfigurationRequest(options, timeout); err != nil {
		return nil, err
	}

	// first send a SABM frame and expect an UA frame. If rejected
	// device with send a DM frame
	if err := c.sendSABM(); err != nil {
		return nil, err
	}

	// add to list of connected RFCOMM devices
	serverChannel, err := serverChannelNumbers.Add(c)
	if err != nil {
		c.CancelConnection()
	}
	c.serverChannel = serverChannel

	// send parameter negotiation
	if err = c.sendPN(); err != nil {
		return nil, err
	}

	// send SABM on (DLCI X)
	if err = c.sendSABM(); err != nil {
		return nil, err
	}

	// LMP Authentication

	// receive UA Frame on (DCLI X)

	// MSC FRAME

	// Data on UIH Frame
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
		ServerNum:          c.serverChannel,
		PollFinal:          false,
		Payload:            v,
		FrameCheckSequence: 0x9a})
}

// CancelConnection disconnects the connection.
func (c *Client) CancelConnection() error {
	c.Lock()
	defer c.Unlock()
	c.sendDISC()
	return c.l2c.Close()
}

// Disconnected returns a receiving channel, which is closed when the client disconnects.
func (c *Client) Disconnected() <-chan struct{} {
	c.Lock()
	defer c.Unlock()
	return c.l2c.Disconnected()
}

func (c *Client) sendDISC() error {
	defer serverChannelNumbers.Remove(c.serverChannel)
	if err = c.sendDISCFrame(); err != nil {
		return err
	}
	select {
	case <-c.p2p:
	case <-time.After(60 * time.Second):
	}
	return nil
}

func (c *Client) sendSABM() error {

	if err := c.sendSABMFrame(); err != nil {
		return err
	}

	frm, err := c.getFrame()
	if err != nil {
		return err
	}

	switch frm.ControlNumber {
	case ControlNumberUA:
	default:
		c.sendDISC()
		c.l2c.Close()
		return errors.New("Received unexpected control number")
	}
	return nil
}

func (c *Client) sendParameterNegotiation() error {

	priority = 0x07
	maxFrameSize = 0x03f0
	commandResponse = 0x01
	direction = 0x01

	err := c.sendMultiplexorFrame(multiplexor.ParameterNegotiation{
		CommandResponse:        commandResponse,
		Direction:              direction,
		ServerNumber:           c.serverNum,
		FrameType:              FrameTypeUIH,
		ConvergenceLayer:       ConvergenceLayer,
		CreditBasedFrowControl: true,
		Priority:               priority,
		Timer:                  Timer,
		MaxSize:                maxFrameSize,
		MaxRetransmissions:     MaxRetransmissions,
		WindowSize:             WindowSize,
	})
	if err != nil {
		return err
	}

	// One device sends a PN message, and the other responds with another PN message.
	// The response may not change the DLCI, the priority, the convergence layer, or the timer
	// value. The response may send back a different timer value. In this case, the device which
	// sent the first PN messages will still use the timer it proposed, but the device at the other
	// end of the connection will use the value it sent in its message.
	// The response may have a smaller value for the maximum frame size, but it may not
	// propose a larger value for this parameter [ pg. 183 ]

	m, err := c.getMultiplexorFrame()
	if err != nil {
		return err
	}

	switch m.(type) {
	case multiplexor.ParameterNegotiation:
		prm := m.(*multiplexor.ParameterNegotiation)
		if c.serverNum != prm.ServerNumber {
			return errors.New("The device attempted to change the DLCI")
		} else if priority != prm.Priority {
			return errors.New("The device attempted to change the priority")
		} else if maxFrameSize < prm.MaxSize {
			return errors.New("The device is trying to use a max frame size greater than proposed")
		}
		c.priority = prm.Priority
		c.maxFrameSize = prm.MaxSize
	default:
		c.sendDISC()
		c.l2c.Close()
		return errors.New("Received unexpected multiplexor")
	}
	return nil
}

// Multiplexor commands and responses are sent on DLCI = 0. The multiplexor commands
// and responses are carried as messages inside an RFCOMM UIH frame as shown in
// Figure 10–10. It is possible to send several multiplexor command messages in
// one RFCOMM frame, or to split a multiplexor command message over more than one frame.
// [pg. 181]
func (c *Client) sendMultiplexorFrame(m multiplexor.Multiplexor) error {
	b, err := m.MarshalBinary()
	if err != nil {
		return nil
	}
	return c.sendUIHFrame(0x00, 0x00, b)
}

func (c *Client) getMultiplexorFrame() (multiplexor.Multiplexor, error) {
	frm, err := c.getFrame()
	if err != nil {
		return nil, err
	}
	if frm.ControlNumber == ControlNumberUIH {
		return multiplexor.UnmarshalBinary(frm.Payload)
	}
	c.sendDISC()
	c.l2c.Close()
	return errors.New("Received unexpected control number")
}

func (c *Client) sendSABMFrame(serverChannel uint8) error {
	return c.sendFrame(frame{
		Direction:          0x01,
		ServerChannel:      serverChannel,
		ControlNumber:      ControlNumberSABM,
		CommmandResponse:   0x01,
		Direction:          0x01,
		PollFinal:          0x01,
		Payload:            []byte{},
		FrameCheckSequence: 0x1c})
}

func (c *Client) sendUIHFrame(direction uint8, serverChannel uint8, data []byte) error {
	return c.sendFrame(frame{
		Direction:          direction,
		ServerChannel:      serverChannel,
		ControlNumber:      ControlNumberUIH,
		CommmandResponse:   0x01,
		PollFinal:          0x01,
		Payload:            data,
		FrameCheckSequence: 0x70})
}

func (c *Client) sendDISCFrame() error {
	return c.sendFrame(frame{
		Direction:          0x01,
		ServerChannel:      c.serverChannel,
		ControlNumber:      ControlNumberDISC,
		CommmandResponse:   0x01,
		Direction:          0x01,
		PollFinal:          0x01,
		Payload:            []byte{},
		FrameCheckSequence: 0x00})
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

func (c *Client) getFrame() (*frame, error) {
	// l2cap has been setup, now we need to setup the RFCOMM connection.
	// command timeouts are 60s, if times out then send a DISC frame
	// on the original SAMB channel
	select {
	case p, ok := <-c.p2p:
		if !ok {
			return nil, errors.New("Channel closed")
		}
		var frm *frame
		if err = frm.Unmarshal(p); err != nil {
			return nil, err
		}
		switch int(frm.ControlNumber) {
		case ControlNumberDISC:
			serverChannelNumbers.Remove(c.serverChannel)
			c.l2c.Close()
			return nil, errors.New("Received disconnect")
		case ControlNumberDM:
			serverChannelNumbers.Remove(c.serverChannel)
			c.l2c.Close()
			return nil, errors.New("Received disconnect mode")
		}
		return frm, nil
	case <-time.After(60 * time.Second):
		c.sendDISC()
		c.l2c.Close()
		return nil, errors.New("Timed out")
	}
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
