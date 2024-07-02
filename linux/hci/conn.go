package hci

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/thomascriley/ble/log"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/hci/cmd"
	"github.com/thomascriley/ble/linux/hci/evt"
	"github.com/thomascriley/ble/linux/smp"
)

type ConnectionCompleteEvent interface {
	PeerAddress() [6]byte
	ConnectionHandle() uint16
	Role() uint8
}

// Conn ...
type Conn struct {
	sync.WaitGroup

	hci *HCI

	param ConnectionCompleteEvent

	// LMP Supported Features as reported by the Read Remote Supported Features
	// Command [Vol 2, Part C, 3.3]
	lmpFeatures uint64

	// The channel identifiers. This may be a fixed ID for some protocols (LE)
	// but dynamic for others (BR/EDR) [Vol 3, Part A, 2.1]
	SourceID      uint16
	DestinationID uint16

	// While MTU is the maximum size of payload data that the upper layer (ATT)
	// can accept, the MPS is the maximum PDU payload size this L2CAP implementation
	// supports. When segmantation is not used, the MPS should be made to the same
	// values of MTUs [Vol 3, Part A, 1.4].
	//
	// For LE-U logical transport, the L2CAP implementations should support
	// a minimum of 23 bytes, which are also the default values before the
	// upper layer (ATT) optionally reconfigures them [Vol 3, Part A, 3.2.8].
	//
	// All L2CAP implementations shall support a minimum MTU of 48 octets over
	// the ACL-U logical link and 23 octets over the LE-U logical link [Vol 3, Part A, 5.1]
	rxMTU int
	txMTU int
	rxMPS int

	// Signaling MTUs are The maximum size of command information that the
	// L2CAP layer entity is capable of accepting.
	// A L2CAP implementations supporting LE-U should support at least 23 bytes.
	// Currently, we support 512 bytes, which should be more than sufficient.
	// The sigTxMTU is discovered via when we sent a signaling pkt that is
	// larger thean the remote device can handle, and get a response of "Command
	// Reject" indicating "Signaling MTU exceeded" along with the actual
	// signaling MTU [Vol 3, Part A, 4.1].
	sigRxMTU int
	sigTxMTU int

	// sigID is used to match responses with signaling requests.
	// The requesting device sets this field and the responding device uses the
	// same value in its response. Within each signalling channel a different
	// Identifier shall be used for each successive command. [Vol 3, Part A, 4]
	sigID uint8

	// sigCID The signaling channel for managing channels over ACL-U logical
	// links shall use CID 0x0001 and the signaling channel for managing channels
	// over LE-U logical links shall use CID 0x0005 [Vol 3, Part A, 4]
	sigCID uint16

	sigSent chan []byte
	// smpSent chan []byte

	// cfgRequest closes when the RFCOMM connection responds to a configuration
	// request
	cfgRequest chan struct{}

	// chDisconnect called with disconnect complete event from hci
	chDisconnect chan []byte

	chInPkt chan packet
	chInPDU chan pdu

	chRecombine chan error

	// Host to Controller Data Flow Control pkt-based Data flow control for LE-U [Vol 2, Part E, 4.1.1]
	// chSentBufs tracks the HCI buffer occupied by this connection.
	txBuffer *Client

	// extendedFeatures BR/EDR/LE supported extended features as determined by a signaling request
	extendedFeatures uint32

	// fixedChannels BR/EDR/LE supported fixed channels as determined by a signaling request
	fixedChannels uint64

	// set of agreed upon parameters to use for the pairing process
	smpPairingResp *smp.PairingResponse
	smpPairingReq  *smp.PairingRequest

	// smpInitiator if this device initiated the pairing process
	smpInitiator bool

	// Confirm and Random values for function c1 used during the LE Legacy Pairing Phase 2
	// store both the Slave and Master versions
	smpMConfirm [16]byte
	smpMRand    [16]byte
	smpSConfirm [16]byte
	smpSRand    [16]byte

	chDone chan struct{}

	// leFrame is set to be true when the LE Credit based flow control is used.
	leFrame bool

	log *slog.Logger

	disconnectHandler func(*Conn)
}

func newConn(h *HCI, param ConnectionCompleteEvent, disconnectHandler func(*Conn), log *slog.Logger) *Conn {
	var (
		sigCID     uint16
		defaultMTU int
	)

	if _, ok := param.(evt.LEConnectionComplete); ok {
		sigCID = cidLESignal
		defaultMTU = ble.DefaultMTU
	} else {
		sigCID = cidSignal
		defaultMTU = ble.DefaultACLMTU
	}

	c := &Conn{
		hci:   h,
		param: param,

		rxMTU: defaultMTU,
		txMTU: defaultMTU,

		rxMPS: defaultMTU,

		sigCID:   sigCID,
		sigRxMTU: defaultMTU,
		sigTxMTU: defaultMTU,

		cfgRequest:   make(chan struct{}),
		chDisconnect: make(chan []byte),

		chInPkt: make(chan packet, 16),
		chInPDU: make(chan pdu, 16),

		txBuffer: NewClient(h.pool),

		chDone:            make(chan struct{}),
		chRecombine:       make(chan error),
		log:               log,
		disconnectHandler: disconnectHandler,
	}

	c.Add(1)
	go func() {
		defer c.Done()
		c.log.Debug("Recombine loop exited")
		var err error
		for {
			if err = c.recombine(); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				_ = c.close(context.Background())
				return
			}
		}
	}()
	return c
}

// Read copies re-assembled L2CAP PDUs into sdu.
func (c *Conn) Read(sdu []byte) (n int, err error) {
	var p pdu
	var ok bool

	select {
	case p, ok = <-c.chInPDU:
		if !ok {
			return 0, fmt.Errorf("input channel closed: %w", io.ErrClosedPipe)
		}
		if len(p) == 0 {
			return 0, fmt.Errorf("received empty packet: %w", io.ErrUnexpectedEOF)
		}
	case <-c.Disconnected():
		return 0, io.ErrClosedPipe
	case <-c.hci.Closed():
		return 0, io.ErrClosedPipe
	}

	// Assume it's a B-Frame.
	sLen := p.dlen()
	data := p.payload()
	if c.leFrame {
		// LE-Frame.
		sLen = leFrameHdr(p).slen()
		data = leFrameHdr(p).payload()
	}
	if cap(sdu) < sLen {
		return 0, fmt.Errorf("payload received exceeds sdu buffer: %w", io.ErrShortBuffer)
	}
	buf := bytes.NewBuffer(sdu)
	buf.Reset()
	if _, err = buf.Write(data); err != nil {
		return 0, fmt.Errorf("unable to write payload to buffer: %w", err)
	}
	for buf.Len() < sLen {
		select {
		case p = <-c.chInPDU:
			if _, err = buf.Write(p.payload()); err != nil {
				return 0, fmt.Errorf("unable to write payload to buffer: %w", err)
			}
		case <-c.Disconnected():
			return 0, io.ErrClosedPipe
		case <-c.hci.Closed():
			return 0, io.ErrClosedPipe
		}
	}
	return sLen, nil
}

// Write breaks down a L2CAP SDU into segments [Vol 3, Part A, 7.3.1]
func (c *Conn) Write(sdu []byte) (int, error) {
	if len(sdu) > c.txMTU {
		return 0, fmt.Errorf("payload exceeds mtu: %w", io.ErrShortWrite)
	}

	pLen := len(sdu)
	if pLen > c.txMTU {
		pLen = c.txMTU
	}
	b := make([]byte, 4+pLen)
	binary.LittleEndian.PutUint16(b[0:2], uint16(len(sdu)))
	binary.LittleEndian.PutUint16(b[2:4], c.SourceID)
	if c.leFrame {
		binary.LittleEndian.PutUint16(b[4:6], uint16(len(sdu)))
		copy(b[6:], sdu)
	} else {
		copy(b[4:], sdu)
	}
	sent, err := c.writePDU(b)
	if err != nil {
		return sent, fmt.Errorf("unable to write pdu: %w", err)
	}
	sdu = sdu[pLen:]

	var n int
	for len(sdu) > 0 {
		pLen = len(sdu)
		if pLen > c.txMTU {
			pLen = c.txMTU
		}
		n, err = c.writePDU(sdu[:pLen])
		sent += n
		if err != nil {
			return sent, fmt.Errorf("unable to write pdu: %w", err)
		}
		sdu = sdu[pLen:]
	}
	return sent, nil
}

// writePDU breaks down a L2CAP PDU into fragments if it's larger than the HCI buffer size. [Vol 3, Part A, 7.2.1]
func (c *Conn) writePDU(pdu []byte) (sent int, err error) {
	flags := uint16(pbfHostToControllerStart << 4) // ACL boundary flags

	// All L2CAP fragments associated with an L2CAP PDU shall be processed for
	// transmission by the Controller before any other L2CAP PDU for the same
	// logical transport shall be processed.
	c.txBuffer.LockPool()
	defer c.txBuffer.UnlockPool()

	// Fail immediately if the connection is already closed
	// Check this with the pool locked to avoid race conditions
	// with handleDisconnectionComplete
	select {
	case <-c.Disconnected():
		return 0, io.ErrClosedPipe
	case <-c.hci.Closed():
		return 0, io.ErrClosedPipe
	default:
	}

	var buf *bytes.Buffer
	for len(pdu) > 0 {

		// Get a buffer from our pre-allocated and flow-controlled pool.
		buf = c.txBuffer.Get(c.chDone) // ACL pkt
		if buf == nil {
			return sent, errors.New("timed out")
		}

		flen := len(pdu) // fragment length
		if flen > buf.Cap()-1-4 {
			flen = buf.Cap() - 1 - 4
		}

		// Prepare the Headers

		// HCI Header: pkt Type
		if err = binary.Write(buf, binary.LittleEndian, pktTypeACLData); err != nil {
			return 0, fmt.Errorf("unable to write ACL data: %w", err)
		}
		// ACL Header: handle and flags
		if err = binary.Write(buf, binary.LittleEndian, c.param.ConnectionHandle()|(flags<<8)); err != nil {
			return 0, fmt.Errorf("unable to write ACL header with handle and flags: %w", err)
		}
		// ACL Header: data len
		if err = binary.Write(buf, binary.LittleEndian, uint16(flen)); err != nil {
			return 0, fmt.Errorf("unable to write ACL header with data len: %w", err)
		}
		// Append payload
		if err = binary.Write(buf, binary.LittleEndian, pdu[:flen]); err != nil {
			return 0, fmt.Errorf("unable to write pdu: %w", err)
		}

		// Flush the pkt to HCI
		select {
		case <-c.Disconnected():
			return 0, io.ErrClosedPipe
		case <-c.hci.Closed():
			return 0, io.ErrClosedPipe
		default:
		}

		if _, err = c.hci.skt.Write(buf.Bytes()); err != nil {
			return sent, fmt.Errorf("unable to write packet to socket: %w", err)
		}
		sent += flen

		flags = pbfContinuing << 4 // Set "continuing" in the boundary flags for the rest of fragments, if any.
		pdu = pdu[flen:]           // Advance the point
	}
	return sent, nil
}

// Recombines fragments into a L2CAP PDU. [Vol 3, Part A, 7.2.2]
func (c *Conn) recombine() error {
	var (
		pkt packet
		ok  bool
	)
	select {
	case <-c.hci.Closed():
		return io.EOF
	case <-c.Disconnected():
		return io.EOF
	case pkt, ok = <-c.chInPkt:
		if !ok {
			return io.EOF
		}
	}

	p := pdu(pkt.data())

	// Currently, check for LE-U only. For channels that we don't recognizes,
	// re-combine them anyway, and discard them later when we dispatch the PDU
	// according to CID.
	if p.cid() == cidLEAtt && p.dlen() > c.rxMPS {
		return fmt.Errorf("fragment size (%d) larger than rxMPS (%d)", p.dlen(), c.rxMPS)
	}
	// TODO check ACL-U packets length
	// not supporting Extended Flow Specification <= 48
	// supporting the Extended Flow Specification <= 672
	// [Vol 3, Part 4]

	// If this pkt is not a complete PDU, and we'll be receiving more
	// fragments, re-allocate the whole PDU (including Header).
	if len(p.payload()) < p.dlen() {
		p = make([]byte, 0, 4+p.dlen())
		p = append(p, pdu(pkt.data())...)
	}
	for len(p) < 4+p.dlen() {
		select {
		case <-c.hci.Closed():
			return io.EOF
		case <-c.Disconnected():
			return io.EOF
		case pkt, ok = <-c.chInPkt:
			if !ok || (pkt.pbf()&pbfContinuing) == 0 {
				return io.ErrUnexpectedEOF
			}
			p = append(p, pdu(pkt.data())...)
		}
	}

	cid := p.cid()
	switch {
	case cid == cidSignal:
		return c.handleSignal(p)
	case cid == cidLEAtt:
		return c.receivePDU(p)
	case cid == cidLESignal:
		return c.handleSignal(p)
	case cid == cidSMP:
		return c.handleSMP(p)
	case cid >= cidDynamicStart:
		return c.receivePDU(p)
	default:
		c.hci.log.Debug("recombine: unrecognized CID", slog.String("CID", fmt.Sprintf("%04X", cid)), log.Bytes("pdu", p))
	}
	return nil
}

func (c *Conn) receivePDU(p pdu) error {
	select {
	case c.chInPDU <- p:
		return nil
	case <-c.hci.Closed():
		return io.EOF
	case <-c.Disconnected():
		return io.EOF
	}
}

// Disconnected returns a receiving channel, which is closed when the connection disconnects.
func (c *Conn) Disconnected() <-chan struct{} {
	return c.chDone
}

// Close disconnects the connection by sending hci disconnect command to the device.
func (c *Conn) Close(ctx context.Context) error {
	select {
	case <-c.Disconnected():
		return nil
	default:
	}

	disconnectCMD := &cmd.Disconnect{
		ConnectionHandle: c.param.ConnectionHandle(),
		Reason:           0x13,
	}
	if err := c.hci.Send(ctx, disconnectCMD, nil); err != nil {
		return fmt.Errorf("unable to send disconnect command: %w", err)
	}

	return c.close(ctx)
}

func (c *Conn) close(ctx context.Context) error {
	// wait for disconnect complete
	select {
	case <-ctx.Done():
	case <-c.Disconnected():
		c.log.Debug("Already disconnected")
		return nil
	case <-c.chDisconnect:
	}

	// remote peripheral disconnected
	select {
	case <-c.chDone:
	default:
		c.log.Debug("Closing connection")
		close(c.chDone)
	}
	defer close(c.chDisconnect)

	c.log.Debug("Waiting for connection to clean up")

	c.Wait()

	c.log.Debug("Releasing for connection resources")

	close(c.chInPkt)
	close(c.cfgRequest)
	close(c.chInPDU)

	// When a connection disconnects, all the sent packets and weren't acked yet
	// will be recycled. [Vol2, Part E 4.1.1]
	//
	// must be done with the pool locked to avoid race conditions where
	// writePDU is in progress and does a Get from the pool after this completes,
	// leaking a buffer from the main pool.
	c.txBuffer.LockPool()
	c.txBuffer.PutAll()
	c.txBuffer.UnlockPool()

	c.log.Debug("Calling disconnect handler")
	c.disconnectHandler(c)

	if c.param.Role() == roleSlave {
		// Re-enable advertising, if it was advertising. Refer to the
		// handleLEConnectionComplete() for details.
		// This may fail with ErrCommandDisallowed, if the controller
		// was actually in advertising state. It does no harm though.
		c.hci.params.RLock()
		if c.hci.params.advEnable.AdvertisingEnable == 1 {
			c.hci.Add(1)
			go func() {
				defer c.hci.Done()
				if err := c.hci.Send(context.Background(), &c.hci.params.advEnable, nil); err != nil {
					c.hci.err = fmt.Errorf("unable to reenable advertising: %w", err)
				}
			}()
		}
		c.hci.params.RUnlock()
	}
	return nil
}

// LocalAddr returns local device's MAC address.
func (c *Conn) LocalAddr() ble.Addr { return c.hci.Addr() }

// RemoteAddr returns remote device's MAC address.
func (c *Conn) RemoteAddr() ble.Addr {
	a := c.param.PeerAddress()
	return net.HardwareAddr([]byte{a[5], a[4], a[3], a[2], a[1], a[0]})
}

// RxMTU returns the MTU which the upper layer is capable of accepting.
func (c *Conn) RxMTU() int { return c.rxMTU }

// SetRxMTU sets the MTU which the upper layer is capable of accepting.
func (c *Conn) SetRxMTU(mtu int) { c.rxMTU, c.rxMPS = mtu, mtu }

// TxMTU returns the MTU which the remote device is capable of accepting.
func (c *Conn) TxMTU() int { return c.txMTU }

// SetTxMTU sets the MTU which the remote device is capable of accepting.
func (c *Conn) SetTxMTU(mtu int) { c.txMTU = mtu }

// pkt implements HCI ACL Data Packet [Vol 2, Part E, 5.4.2]
// Packet boundary flags , bit[5:6] of handle field's MSB
// Broadcast flags. bit[7:8] of handle field's MSB
// Not used in LE-U. Leave it as 0x00 (Point-to-Point).
// Broadcasting in LE uses ADVB logical transport.
type packet []byte

func (a packet) handle() uint16 { return uint16(a[0]) | (uint16(a[1]&0x0f) << 8) }
func (a packet) pbf() int       { return (int(a[1]) >> 4) & 0x3 }
func (a packet) bcf() int       { return (int(a[1]) >> 6) & 0x3 }
func (a packet) dlen() int      { return int(a[2]) | (int(a[3]) << 8) }
func (a packet) data() []byte   { return a[4:] }

type pdu []byte

func (p pdu) dlen() int       { return int(binary.LittleEndian.Uint16(p[0:2])) }
func (p pdu) cid() uint16     { return binary.LittleEndian.Uint16(p[2:4]) }
func (p pdu) payload() []byte { return p[4:] }

type leFrameHdr pdu

func (f leFrameHdr) slen() int       { return int(binary.LittleEndian.Uint16(f[4:6])) }
func (f leFrameHdr) payload() []byte { return f[6:] }
