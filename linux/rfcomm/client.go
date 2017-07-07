package rfcomm

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux/multiplexer"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidArgument means one or more of the arguments are invalid.
	ErrInvalidArgument = errors.New("invalid argument")
)

var (
	serverChannelNumbers = newServerChannels()
)

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

	serverChannel uint8
}

// NewClient returns an RFCOMM Client that has been initialized according to the
// RFCOMM specifications.
func NewClient(l2c ble.Conn) (*Client, error) {
	c := &Client{
		l2c:     l2c,
		p2p:     make(chan []byte),
		chErr:   make(chan error),
		chTxBuf: make(chan []byte, 1),
		rxBuf:   make([]byte, ble.MaxACLMTU),
	}
	c.Lock()
	if err := c.connect(); err != nil {
		c.sendDISC()
		c.l2c.Close()
		return nil, err
	}
	return c, nil
}

// Connect ...
func (c *Client) connect() error {
	// first send a SABM on DLCI (0) and expect an UA frame. If rejected
	// device with send a DM frame
	if err := c.sendSABM(c.serverChannel); err != nil {
		return err
	}

	// add to list of connected RFCOMM devices
	serverChannel, err := serverChannelNumbers.Add(c)
	if err != nil {
		return err
	}
	c.serverChannel = serverChannel

	// send parameter negotiation [optional]
	if err = c.sendParameterNegotiation(Priority, MaxFrameSize); err != nil {
		fmt.Printf("Error negotiating the rfcomm parameters: %s\n", err)
	}

	// send SABM on (DLCI X)
	if err = c.sendSABM(c.serverChannel); err != nil {
		return err
	}

	// MSC FRAME
	if err = c.sendModemStatus(); err != nil {
		return err
	}

	// Send connection test
	if err := c.sendTest(0x08); err != nil {
		return err
	}
	return nil
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
		if controlField&ControlNumberUIH != ControlNumberUIH {
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
		CommmandResponse:   0x01,
		Direction:          0x00,
		ServerChannel:      c.serverChannel,
		PollFinal:          0x00,
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
	if err := c.sendDISCFrame(); err != nil {
		return err
	}
	select {
	case <-c.p2p:
	case <-time.After(60 * time.Second):
	}
	return nil
}

func (c *Client) sendSABM(serverChan uint8) error {

	if err := c.sendSABMFrame(serverChan); err != nil {
		return err
	}

	frm, err := c.getFrame()
	if err != nil {
		return err
	}

	// TODO: check for LMP Authentication
	switch frm.ControlNumber {
	case ControlNumberUA:
	default:
		return errors.New("Received unexpected control number")
	}
	return nil
}

func (c *Client) sendParameterNegotiation(priority uint8, maxFrameSize uint16) error {

	err := c.sendMultiplexerFrame(&multiplexer.ParameterNegotiation{
		CommandResponse:        0x01,
		ServerChannel:          c.serverChannel,
		FrameType:              FrameTypeUIH,
		ConvergenceLayer:       ConvergenceLayer,
		CreditBasedFlowControl: true,
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

	m, err := c.getMultiplexerFrame()
	if err != nil {
		return err
	}

	prm, ok := m.(*multiplexer.ParameterNegotiation)
	if !ok {
		return errors.New("Received unexpected multiplexer")
	} else if c.serverChannel != prm.ServerChannel {
		return errors.New("The device responded on a different server channel")
	} else if priority != prm.Priority {
		return errors.New("The device attempted to change the priority")
	} else if maxFrameSize < prm.MaxSize {
		return errors.New("The device is trying to use a max frame size greater than proposed")
	}
	c.priority = prm.Priority
	c.maxFrameSize = prm.MaxSize

	return nil
}

// sendTest ...
func (c *Client) sendTest(data uint8) error {

	err := c.sendMultiplexerFrame(&multiplexer.Test{
		CommandResponse: 0x01,
		Data:            data,
	})
	if err != nil {
		return err
	}

	// The test command is used to check the RFCOMM connection. As is normal, the length byte
	// gives the number of value bytes which follow. The number of value bytes is not
	// fixed, and is used to hold a test pattern. The remote end of the link echoes
	// the same value bytes back. [ pg. 184 ]

	m, err := c.getMultiplexerFrame()
	if err != nil {
		return err
	}

	_, ok := m.(*multiplexer.Test)
	if !ok {
		return errors.New("Received unexpected multiplexer")
	} // else if test.Data != data {
	//	return errors.New("The device echoed back a data stream that does not match")
	//}
	return nil
}

// sendFlowControl Applies the flow control mechanism to all connections
func (c *Client) sendFlowControl(on bool) error {
	var m multiplexer.Multiplexer
	if on {
		m = &multiplexer.FlowControlOn{}
	} else {
		m = &multiplexer.FlowControlOff{}
	}
	m.SetCommandResponse(0x01)

	if err := c.sendMultiplexerFrame(m); err != nil {
		return err
	}

	m, err := c.getMultiplexerFrame()
	if err != nil {
		return err
	}

	switch m.(type) {
	case *multiplexer.FlowControlOn:
		if !on {
			return errors.New("Responded with flow control on when sent an off type")
		}
	case *multiplexer.FlowControlOff:
		if !on {
			return errors.New("Responded with flow control off when sent an on type")
		}
	default:
		return errors.New("Received unexpected multiplexer")
	}
	return nil
}

// sendModemStatus a flow control mechanism which can be applied to just one channel at
// a time
func (c *Client) sendModemStatus() error {
	frm := &multiplexer.ModemStatus{
		CommandResponse:    0x01,
		ServerChannel:      c.serverChannel,
		FlowControl:        FlowControl,
		ReadyToCommunicate: ReadyToCommunicate,
		ReadyToReceive:     ReadyToReceive,
		IncomingCall:       IncomingCall,
		DataValid:          ValidData,
	}
	if err := c.sendMultiplexerFrame(frm); err != nil {
		return err
	}

	// Both the DTE and DCE uses this command to notify each other of the status of their own V.24 control signals.
	m, err := c.getMultiplexerFrame()
	if err != nil {
		return err
	}
	ms, ok := m.(*multiplexer.ModemStatus)
	if !ok {
		return errors.New("Received unexpected multiplexer")
	}

	// acknowledge the remote's modem status
	ms.CommandResponse = 0x00
	if c.sendMultiplexerFrame(ms); err != nil {
		return err
	}

	// get the acknowledgement from the remote
	if m, err = c.getMultiplexerFrame(); err != nil {
		return err
	}
	if _, ok = m.(*multiplexer.ModemStatus); !ok {
		return errors.New("Received unexpected multiplexer")
	}
	return nil
}

// multiplexer commands and responses are sent on DLCI = 0. The multiplexer commands
// and responses are carried as messages inside an RFCOMM UIH frame as shown in
// Figure 10–10. It is possible to send several multiplexer command messages in
// one RFCOMM frame, or to split a multiplexer command message over more than one frame.
// [pg. 181]
func (c *Client) sendMultiplexerFrame(m multiplexer.Multiplexer) error {
	b, err := m.MarshalBinary()
	if err != nil {
		return err
	}
	if _, ok := m.(*multiplexer.Test); ok {
		return c.sendUIHFrame(0x00, c.serverChannel, 0x01, b)
	}
	return c.sendUIHFrame(0x00, 0x00, 0x00, b)
}

func (c *Client) getMultiplexerFrame() (multiplexer.Multiplexer, error) {
	frm, err := c.getFrame()
	if err != nil {
		return nil, err
	}
	if frm.ControlNumber != ControlNumberUIH {
		return nil, errors.New("Received unexpected control number")
	}
	m, err := multiplexer.UnmarshalBinary(frm.Payload)
	if err != nil {
		return nil, err
	}
	switch m.(type) {
	case *multiplexer.NotSupported:
		return nil, errors.New("multiplexer is not supported by this device")
	}
	return m, nil
}

func (c *Client) sendSABMFrame(serverChannel uint8) error {
	return c.sendFrame(frame{
		Direction:          0x01,
		ServerChannel:      serverChannel,
		ControlNumber:      ControlNumberSABM,
		CommmandResponse:   0x01,
		PollFinal:          0x01,
		Payload:            []byte{},
		FrameCheckSequence: 0x1c})
}

func (c *Client) sendUIHFrame(direction uint8, serverChannel uint8, pollFinal uint8, data []byte) error {
	return c.sendFrame(frame{
		Direction:          direction,
		ServerChannel:      serverChannel,
		ControlNumber:      ControlNumberUIH,
		CommmandResponse:   0x01,
		PollFinal:          pollFinal,
		Payload:            data,
		FrameCheckSequence: 0x70})
}

func (c *Client) sendDISCFrame() error {
	return c.sendFrame(frame{
		Direction:          0x01,
		ServerChannel:      c.serverChannel,
		ControlNumber:      ControlNumberDISC,
		CommmandResponse:   0x01,
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
		if err := frm.Unmarshal(p); err != nil {
			return nil, err
		}
		switch frm.ControlNumber {
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