package rfcomm

import (
	"context"
	"errors"
	"fmt"
	"github.com/thomascriley/ble/log"
	"io"
	"sync"
	"time"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/multiplexer"
)

var (
	serverChannelNumbers = newServerChannels()
)

// Client implements an Attribute Protocol Client.
type Client struct {
	sync.RWMutex
	sync.WaitGroup

	l2c   ble.Conn
	p2p   chan []byte
	chErr chan error

	rxBuf   []byte
	chTxBuf chan []byte

	dlci         uint8
	priority     uint8
	maxFrameSize uint16

	serverChannel uint8

	flowControl bool
	credits     uint16
	waitCredits chan struct{}
}

// NewClient returns an RFCOMM Client that has been initialized according to the
// RFCOMM specifications.
func NewClient(l2c ble.Conn, channel uint8) *Client {
	c := &Client{
		l2c:           l2c,
		p2p:           make(chan []byte, 1),
		chErr:         make(chan error),
		chTxBuf:       make(chan []byte, 1),
		waitCredits:   make(chan struct{}),
		rxBuf:         make([]byte, ble.MaxACLMTU),
		serverChannel: channel}
	c.chTxBuf <- make([]byte, l2c.TxMTU(), l2c.TxMTU())
	return c
}

func (c *Client) DialContext(ctx context.Context) (err error) {
	// start a separate thread to close the connection if the context closes before the connection can be established,
	// this will stop the connection process and exist. Use a separate wait group since the close call will call a wait
	// on the client's wait group

	group := sync.WaitGroup{}
	group.Add(1)
	defer group.Wait()

	success, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(timeout context.Context, success context.Context){
		defer group.Done()
		select {
		case <-timeout.Done():
			_ = c.close(nil)
			err = errors.New("could not dial: timed out")
		case <-success.Done():
		}
	}(ctx, success)

	c.Add(1)
	go func(){
		defer c.Done()
		c.loop()
	}()

	if err2 := c.connect(ctx); err2 != nil {
		_ = c.close(nil)
		if err == nil {
			err = fmt.Errorf("could not dial: %s", err2)
		}
	}
	return
}

func (c *Client) close(ctx context.Context) error {
	defer c.Wait()
	_ = c.sendDISC()
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
	}
	return c.l2c.Close(ctx)
}

// Connect ...
func (c *Client) connect(ctx context.Context) error {
	// first send a SABM on DLCI (0) and expect an UA frame. If rejected
	// device with send a DM frame
	if err := c.sendSABM(0); err != nil {
		return err
	}

	// add to list of connected RFCOMM devices
	// TODO: use sdp to figure out which channel supports RFCOMM
	/*serverChannel, err := serverChannelNumbers.Add(c)
	if err != nil {
		return err
	}
	c.serverChannel = serverChannel*/

	// send parameter negotiation [optional]
	if err := c.sendParameterNegotiation(Priority, MaxFrameSize); err != nil {
		log.Printf("Error negotiating the rfcomm parameters: %s\n", err)
	}

	// send SABM on (DLCI X)
	if err := c.sendSABM(c.serverChannel); err != nil {
		return err
	}

	// MSC FRAME
	if err := c.sendModemStatus(); err != nil {
		return err
	}

	// Exchange Credits
	if err := c.exchangeCredits(0x21); err != nil {
		return err
	}

	// Send connection test
	/*
		if err := c.sendTest(ctx, 0x08); err != nil {
			return err
		}*/

	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Address returns the address of the client.
func (c *Client) Address() ble.Addr {
	return c.l2c.RemoteAddr()
}

// Read ...
func (c *Client) Read(b []byte) (int, error) {
	for {
		select {
		case <-c.Disconnected():
			return 0, errors.New("disconnected")
		case err := <-c.chErr:
			return 0, err
		case p, ok := <-c.p2p:
			if !ok {
				return 0, fmt.Errorf("%w: input channel closed", io.ErrClosedPipe)
			}
			if len(p) == 0 {
				return 0, fmt.Errorf("%w: recieved empty packet", io.ErrUnexpectedEOF)
			}
			if len(p) > len(b) {
				return 0, fmt.Errorf("%w: payload recieved exceeds sdu buffer", io.ErrShortBuffer)
			}
			frm := &frame{}
			if err := frm.Unmarshal(p); err != nil {
				return 0, err
			}

			if frm.ServerChannel != c.serverChannel {
				return 0, errors.New("mismatching server channel")
			}
			if frm.ControlNumber != ControlNumberUIH {
				return 0, errors.New("invalid control number")
			}
			if c.flowControl && frm.PollFinal == 0x01 {
				c.credits = c.credits + uint16(frm.Credits)
				select {
				case <-c.waitCredits:
				case <-c.Disconnected():
					return 0, errors.New("disconnected")
				default:
					if c.credits > 0 {
						close(c.waitCredits)
					}
				}
			}
			copy(b[:], frm.Payload)
			return len(frm.Payload), nil
		}
	}
}

// Write RFCOMM Bluetooth specification Version 1.0 B
func (c *Client) Write(v []byte) (int, error) {
	select {
	case <-c.waitCredits:
	case <-c.Disconnected():
		return 0, errors.New("disconnected")
	default:
		if _, err := c.Read(make([]byte, c.l2c.RxMTU())); err != nil {
			return 0, err
		}
	}
	if c.flowControl {
		c.credits = c.credits - 1
		if c.credits == 0 {
			c.waitCredits = make(chan struct{})
		}
		log.Printf("Credits: %d", c.credits)
	}

	return len(v), c.sendFrame( &frame{
		ControlNumber:      ControlNumberUIH,
		CommmandResponse:   0x01,
		Direction:          0x00,
		ServerChannel:      c.serverChannel,
		PollFinal:          0x00,
		Payload:            v,
		FrameCheckSequence: 0x9a})
}

// CancelConnection disconnects the connection.
func (c *Client) CancelConnection(ctx context.Context) error {
	return c.close(ctx)
}

// Disconnected returns a receiving channel, which is closed when the client disconnects.
func (c *Client) Disconnected() <-chan struct{} {
	return c.l2c.Disconnected()
}

func (c *Client) sendDISC() error {
	defer serverChannelNumbers.Remove(c.serverChannel)
	if err := c.sendDISCFrame(c.serverChannel); err != nil {
		return err
	}

	for {
		frm, err := c.getFrame()
		if err != nil {
			return err
		}
		if frm.ControlNumber == ControlNumberUA {
			break
		}
	}

	if err := c.sendDISCFrame(0x00); err != nil {
		return err
	}

	for {
		frm, err := c.getFrame()
		if err != nil {
			return err
		}
		if frm.ControlNumber == ControlNumberUA {
			break
		}
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
		return nil
	default:
		return errors.New("received unexpected control number")
	}
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
		return errors.New("received unexpected multiplexer")
	} else if c.serverChannel != prm.ServerChannel {
		return errors.New("the device responded on a different server channel")
	} else if priority != prm.Priority {
		return errors.New("the device attempted to change the priority")
	} else if maxFrameSize < prm.MaxSize {
		return errors.New("the device is trying to use a max frame size greater than proposed")
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
		return errors.New("received unexpected multiplexer")
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
			return errors.New("responded with flow control on when sent an off type")
		}
	case *multiplexer.FlowControlOff:
		if !on {
			return errors.New("responded with flow control off when sent an on type")
		}
	default:
		return errors.New("received unexpected multiplexer")
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
		return errors.New("received unexpected multiplexer")
	}

	// acknowledge the remote's modem status
	ms.CommandResponse = 0x00
	ms.ServerChannel = c.serverChannel
	if err := c.sendMultiplexerFrame(ms); err != nil {
		return err
	}

	// get the acknowledgement from the remote
	if m, err = c.getMultiplexerFrame(); err != nil {
		return err
	}
	if _, ok = m.(*multiplexer.ModemStatus); !ok {
		return errors.New("received unexpected multiplexer")
	}
	return nil
}

// exangeCredits
func (c *Client) exchangeCredits(credits uint8) error {
	err := c.sendFrame(&frame{
		Direction:        0x00,
		ServerChannel:    c.serverChannel,
		ControlNumber:    ControlNumberUIH,
		CommmandResponse: 0x01,
		PollFinal:        0x01,
		Credits:          credits,
		Payload:          []byte{},
	})
	if err != nil {
		return err
	}
	frm, err := c.getFrame()
	if err != nil {
		return err
	}
	if frm.ControlNumber != ControlNumberUIH {
		return fmt.Errorf("received unexpected control number: %v", frm.ControlNumber)
	}
	if frm.Credits == 0 && len(frm.Payload) > 0 {
		c.credits = uint16(frm.Payload[0]) + 1
	} else {
		c.credits = uint16(frm.Credits) + 1
	}
	close(c.waitCredits)
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
		return c.sendUIHFrame( 0x00, c.serverChannel, 0x01, b)
	}
	return c.sendUIHFrame(0x00, 0x00, 0x00, b)
}

func (c *Client) getMultiplexerFrame() (multiplexer.Multiplexer, error) {
	frm, err := c.getFrame()
	if err != nil {
		return nil, err
	}
	if frm.ControlNumber != ControlNumberUIH {
		return nil, fmt.Errorf("received unexpected control number: %v", frm.ControlNumber)
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
	return c.sendFrame(&frame{
		Direction:        0x00, // 0x01
		ServerChannel:    serverChannel,
		ControlNumber:    ControlNumberSABM,
		CommmandResponse: 0x01,
		PollFinal:        0x01,
		Payload:          []byte{}})
}

func (c *Client) sendUIHFrame(direction uint8, serverChannel uint8, pollFinal uint8, data []byte) error {
	return c.sendFrame(&frame{
		Direction:        direction,
		ServerChannel:    serverChannel,
		ControlNumber:    ControlNumberUIH,
		CommmandResponse: 0x01,
		PollFinal:        pollFinal,
		Payload:          data})
}

func (c *Client) sendDISCFrame(channel uint8) error {
	return c.sendFrame(&frame{
		Direction:        0x00,
		ServerChannel:    channel,
		ControlNumber:    ControlNumberDISC,
		CommmandResponse: 0x01,
		PollFinal:        0x01,
		Payload:          []byte{}})
}

func (c *Client) sendFrame(frm *frame) error {

	var txBuf []byte
	select {
	case <-c.Disconnected():
		return errors.New("disconnected")
	case txBuf = <-c.chTxBuf:
		if len(txBuf) < int(c.maxFrameSize) {
			txBuf = make([]byte, c.maxFrameSize, c.maxFrameSize)
		}
	}
	defer func() {
		select {
		case c.chTxBuf <- txBuf:
		case <-c.Disconnected():
		}
	}()

	n, err := frm.Marshal(txBuf)
	if err != nil {
		return err
	}

	if n > c.l2c.TxMTU() {
		return fmt.Errorf("the frame is larger than the mtu %d > %d", n, c.l2c.TxMTU())
	}
	if n > int(c.maxFrameSize) && c.maxFrameSize > 0 {
		return fmt.Errorf("the frame is larger than the max frame size %d > %d", n, c.maxFrameSize)
	}

	_, err = c.l2c.Write(txBuf[:n])
	return err
}

func (c *Client) getFrame() (*frame, error) {
	// l2cap has been setup, now we need to setup the RFCOMM connection.
	// command timeouts are 60s, if times out then send a DISC frame
	// on the original SAMB channel
	for {
		select {
		case <-c.Disconnected():
			return nil, errors.New("disconnected")
		case p, ok := <-c.p2p:
			if !ok {
				return nil, errors.New("channel closed")
			}
			frm := &frame{}
			if err := frm.Unmarshal(p); err != nil {
				return nil, err
			}
			switch frm.ControlNumber {
			case ControlNumberDISC:
				serverChannelNumbers.Remove(c.serverChannel)
				ctx, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
				_ = c.l2c.Close(ctx)
				cancel()
				return nil, errors.New("received disconnect")
			case ControlNumberDM:
				serverChannelNumbers.Remove(c.serverChannel)
				ctx, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
				_ = c.l2c.Close(ctx)
				cancel()
				return nil, errors.New("received disconnect mode")
			case ControlNumberUIH:
				if c.flowControl && frm.PollFinal == 0x01 {
					c.credits = c.credits + uint16(frm.Credits)
					select {
					case <-c.waitCredits:
					case <-c.Disconnected():
						return nil, errors.New("disconnected")
					default:
						if c.credits > 0 {
							close(c.waitCredits)
						}
					}
				}
			}
			return frm, nil
		}
	}
}

// Loop ...
func (c *Client) loop() {
	for {
		n, err := c.l2c.Read(c.rxBuf)
		if err != nil {
			// We don't expect any error from the bearer (L2CAP ACL-U)
			// Pass it along to the pending read, if any, and escape.
			select {
			case c.chErr <- err:
			case <-c.Disconnected():
			}
			return
		}

		b := make([]byte, n)
		copy(b, c.rxBuf)
		select {
		case c.p2p <- b:
		case <-c.Disconnected():
			return
		}
	}
}
