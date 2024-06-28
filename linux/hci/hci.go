package hci

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/hci/cmd"
	"github.com/thomascriley/ble/linux/hci/evt"
	"github.com/thomascriley/ble/linux/hci/socket"
	"github.com/thomascriley/ble/linux/smp"
	"github.com/thomascriley/ble/log"
)

// Command ...
type Command interface {
	OpCode() int
	Len() int
	Marshal([]byte) error
}

// CommandRP ...
type CommandRP interface {
	Unmarshal(b []byte) error
}

type handlerFn func(b []byte) error

type pkt struct {
	cmd  Command
	done chan []byte
}

type nameHandlers struct {
	sync.Mutex
	handlers map[ble.Addr]chan *nameEvent
}

// NewHCI returns a hci device.
func NewHCI(log *slog.Logger) *HCI {
	h := &HCI{
		id: -1,

		chCmdPkt:  make(chan *pkt),
		chCmdBufs: make(chan []byte, 16),
		sent:      make(map[int]*pkt),
		sentMutex: &sync.RWMutex{},

		evth:     map[int]handlerFn{},
		evtMutex: &sync.RWMutex{},
		subh:     map[int]handlerFn{},
		subMutex: &sync.RWMutex{},

		adHist:  make(map[string]*Advertisement, 0),
		adTimes: make(map[string]time.Time, 0),

		dynamicCID: cidDynamicStart,

		muConns:           &sync.Mutex{},
		conns:             make(map[uint16]*Conn),
		chMasterConn:      make(chan *Conn),
		chMasterBREDRConn: make(chan *Conn),
		chSlaveConn:       make(chan *Conn),

		log: log.With("device", "hci"),

		//done: make(chan bool),
	}
	return h
}

// HCI ...
type HCI struct {
	sync.Mutex
	sync.WaitGroup

	params params

	skt socket.Closer
	id  int

	// Host to Controller command flow control [Vol 2, Part E, 4.4]
	chCmdPkt  chan *pkt
	chCmdBufs chan []byte
	sent      map[int]*pkt
	sentMutex *sync.RWMutex

	// evtHub
	evth     map[int]handlerFn
	evtMutex *sync.RWMutex
	subh     map[int]handlerFn
	subMutex *sync.RWMutex

	// aclHandler
	bufSize int
	bufCnt  int

	// Device information or status.
	addr    net.HardwareAddr
	txPwrLv int

	// adHist tracks the history of past scannable advertising packets.
	// Controller delivers AD(Advertising Data) and SR(Scan Response) separately
	// through HCI. Upon receiving an AD, no matter it's scannable or not, we
	// pass a Advertisement (AD only) to advHandler immediately.
	// Upon receiving a SR, we search the AD history for the AD from the same
	// device, and pass the Advertisiement (AD+SR) to advHandler.
	// The adHist and adLast are allocated in the Scan().
	advHandler ble.AdvHandler
	adHist     map[string]*Advertisement
	adTimes    map[string]time.Time
	adMutex    sync.RWMutex

	// Inquiry scan handler
	inqHandler ble.InqHandler

	// Outstanding Remote Name Requests
	nameHandlers *nameHandlers

	// dynamicCID pointer to next available channel for LMP connections for BR/EDR
	dynamicCID uint16

	// Host to Controller Data Flow Control Packet-based Data flow control for LE-U [Vol 2, Part E, 4.1.1]
	// Minimum 27 bytes. 4 bytes of L2CAP Header, and 23 bytes Payload from upper layer (ATT)
	pool *Pool

	// L2CAP connections
	muConns           *sync.Mutex
	conns             map[uint16]*Conn
	chMasterConn      chan *Conn // Dial returns master BLE connections.
	chMasterBREDRConn chan *Conn // DialBREDR returns master BREDR connections.
	chSlaveConn       chan *Conn // Peripheral accept slave connections.

	connectedHandler    func(evt.LEConnectionComplete)
	disconnectedHandler func(evt.DisconnectionComplete)

	// SMP capabilities
	smpCapabilites smp.Capabilities

	// holds socket errors before closing
	err error

	initialized bool
	log         *slog.Logger
	//done chan bool
}

// Init ...
func (h *HCI) Init(ctx context.Context) (err error) {
	if h.initialized {
		return ble.ErrAlreadyInitialized
	}
	h.initialized = true

	h.evtMutex.Lock()

	h.evth[0x3E] = h.handleLEMeta
	h.evth[evt.CommandCompleteCode] = h.handleCommandComplete
	h.evth[evt.CommandStatusCode] = h.handleCommandStatus
	h.evth[evt.DisconnectionCompleteCode] = h.handleDisconnectionComplete
	h.evth[evt.NumberOfCompletedPacketsCode] = h.handleNumberOfCompletedPackets

	// evt.EncryptionChangeCode:                     todo),
	// evt.ReadRemoteVersionInformationCompleteCode: todo),
	// evt.HardwareErrorCode:                        todo),
	// evt.DataBufferOverflowCode:                   todo),
	// evt.EncryptionKeyRefreshCompleteCode:         todo),
	// evt.AuthenticatedPayloadTimeoutExpiredCode:   todo),
	// evt.LEReadRemoteUsedFeaturesCompleteSubCode:   todo),
	// evt.LERemoteConnectionParameterRequestSubCode: todo),

	// BD/EDR
	h.evth[evt.InquiryCompleteCode] = h.handleInquiryComplete
	h.evth[evt.InquiryResultCode] = h.handleInquiryResult
	h.evth[evt.InquiryResultwithRSSICode] = h.handleInquiryWithRSSI
	h.evth[evt.ExtendedInquiryCode] = h.handleExtendedInquiry
	h.evth[evt.ConnectionCompleteCode] = h.handleConnectionComplete
	h.evth[evt.PageScanRepetitionModeChangeCode] = h.handlePageScanRepetitionModeChange
	h.evth[evt.ReadRemoteSupportedFeaturesCompleteCode] = h.handleReadRemoteSupportedFeaturesComplete
	h.evth[evt.MaxSlotsChangeCode] = h.handleMaxSlotsChange

	h.evtMutex.Unlock()

	h.subMutex.Lock()

	h.subh[evt.LEAdvertisingReportSubCode] = h.handleLEAdvertisingReport
	h.subh[evt.LEConnectionCompleteSubCode] = h.handleLEConnectionComplete
	h.subh[evt.LEConnectionUpdateCompleteSubCode] = h.handleLEConnectionUpdateComplete
	h.subh[evt.LELongTermKeyRequestSubCode] = h.handleLELongTermKeyRequest

	h.subMutex.Unlock()

	h.nameHandlers = &nameHandlers{handlers: make(map[ble.Addr]chan *nameEvent, 0)}

	if h.skt, err = socket.NewSocket(h.id); err != nil {
		return fmt.Errorf("unable to create new socket: %w", err)
	}
	h.Add(1)
	go func() {
		defer h.Close()
		defer h.Done()
		h.sktLoop()
	}()

	if err = h.setAllowedCommands(1); err != nil {
		return fmt.Errorf("unable to set allowed commands: %w", err)
	}

	h.log.Debug("hci init")
	if err = h.init(ctx); err != nil {
		return err
	}

	// Pre-allocate buffers with additional head room for lower layer headers.
	// HCI header (1 Byte) + ACL Data Header (4 bytes) + L2CAP PDU (or fragment)
	h.pool = NewPool(1+4+h.bufSize, h.bufCnt-1)

	h.params.RLock()
	defer h.params.RUnlock()

	if h.params.advEnable.AdvertisingEnable == 1 {
		h.log.Debug("send adv params")
		if err := h.Send(ctx, &h.params.advParams, nil); err != nil {
			return fmt.Errorf("unable to send advertising params: %w", err)
		}
	}

	h.log.Debug("send scan params")
	if err = h.Send(ctx, &h.params.scanParams, nil); err != nil {
		return fmt.Errorf("unable to send scan params: %w", err)
	}
	return nil
}

// Close ...
func (h *HCI) Close() error {
	defer h.Wait()

	errored := make([]string, 0)
	if h.skt != nil {
		// 2 seconds to close all connections
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		h.log.Debug("closing connections")
		// close connections to all peripherals
		h.muConns.Lock()
		for _, c := range h.conns {
			// 200 milliseconds to close each connection
			subCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
			if err := c.Close(subCtx); err != nil {
				errored = append(errored, err.Error())
			}
			cancel()
		}
		h.muConns.Unlock()

		h.log.Debug("closing socket")
		if err := h.skt.Close(); err != nil {
			errored = append(errored, err.Error())
		}
	}
	h.log.Debug("closed")
	if len(errored) > 0 {
		return fmt.Errorf("unable to nicely close: %s", strings.Join(errored, ", "))
	}
	return nil
}

// Closed ...
func (h *HCI) Closed() chan struct{} {
	if h.skt == nil {
		h.err = errors.New("socket has not been initialized")
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	return h.skt.Closed()
}

func (h *HCI) init(ctx context.Context) error {
	h.log.Debug("reseting")
	if err := h.Send(ctx, &cmd.Reset{}, nil); err != nil {
		return fmt.Errorf("unable to reset: %w", err)
	}

	h.log.Debug("read db addr")
	ReadBDADDRRP := cmd.ReadBDADDRRP{}
	if err := h.Send(ctx, &cmd.ReadBDADDR{}, &ReadBDADDRRP); err != nil {
		return fmt.Errorf("unable to read BDADDR: %w", err)
	}

	a := ReadBDADDRRP.BDADDR
	h.addr = []byte{a[5], a[4], a[3], a[2], a[1], a[0]}

	h.log.Debug("read buffer size")
	ReadBufferSizeRP := cmd.ReadBufferSizeRP{}
	if err := h.Send(ctx, &cmd.ReadBufferSize{}, &ReadBufferSizeRP); err != nil {
		return fmt.Errorf("unable to read buffer size: %w", err)
	}

	// Assume the buffers are shared between ACL-U and LE-U.
	h.bufCnt = int(ReadBufferSizeRP.HCTotalNumACLDataPackets)
	h.bufSize = int(ReadBufferSizeRP.HCACLDataPacketLength)

	h.log.Debug("le read buffer size")
	LEReadBufferSizeRP := cmd.LEReadBufferSizeRP{}
	if err := h.Send(ctx, &cmd.LEReadBufferSize{}, &LEReadBufferSizeRP); err != nil {
		return fmt.Errorf("unable to read le buffer size: %w", err)
	}

	if LEReadBufferSizeRP.HCTotalNumLEDataPackets != 0 {
		// Okay, LE-U do have their own buffers.
		h.bufCnt = int(LEReadBufferSizeRP.HCTotalNumLEDataPackets)
		h.bufSize = int(LEReadBufferSizeRP.HCLEDataPacketLength)
	}

	h.log.Debug("le read advertising channel")

	LEReadAdvertisingChannelTxPowerRP := cmd.LEReadAdvertisingChannelTxPowerRP{}
	if err := h.Send(ctx, &cmd.LEReadAdvertisingChannelTxPower{}, &LEReadAdvertisingChannelTxPowerRP); err != nil {
		return fmt.Errorf("unable to read advertising channel tx power: %w", err)
	}

	h.txPwrLv = int(LEReadAdvertisingChannelTxPowerRP.TransmitPowerLevel)

	h.log.Debug("le set event mask")
	LESetEventMaskRP := cmd.LESetEventMaskRP{}
	if err := h.Send(ctx, &cmd.LESetEventMask{LEEventMask: 0x000000000000001F}, &LESetEventMaskRP); err != nil {
		return fmt.Errorf("unable to set le event mask: %w", err)
	}

	h.log.Debug("set event mask")
	SetEventMaskRP := cmd.SetEventMaskRP{}
	if err := h.Send(ctx, &cmd.SetEventMask{EventMask: 0x3dbff807fffbffff}, &SetEventMaskRP); err != nil {
		return fmt.Errorf("unable to set event mask: %w", err)
	}

	h.log.Debug("write le host support")
	WriteLEHostSupportRP := cmd.WriteLEHostSupportRP{}
	if err := h.Send(ctx, &cmd.WriteLEHostSupport{LESupportedHost: 1, SimultaneousLEHost: 0}, &WriteLEHostSupportRP); err != nil {
		return fmt.Errorf("unable to write le host support: %w", err)
	}

	h.log.Debug("done init")
	return nil
}

// Send ...
func (h *HCI) Send(ctx context.Context, c Command, r CommandRP) error {
	h.log.Debug("sending command", slog.Int("opcode", c.OpCode()), slog.String("response", fmt.Sprintf("%T", r)))
	// reuse the byte array, marshalling the new command into it
	var b []byte
	select {
	case b = <-h.chCmdBufs:
	case <-h.Closed():
		if h.err == nil {
			return fmt.Errorf("hci device disconnected")
		}
		return fmt.Errorf("hci device disconnected: %w", h.err)
	case <-ctx.Done():
		return ctx.Err()
	}

	b[0] = pktTypeCommand // HCI header
	b[1] = byte(c.OpCode())
	b[2] = byte(c.OpCode() >> 8)
	b[3] = byte(c.Len())
	if err := c.Marshal(b[4:]); err != nil {
		return fmt.Errorf("hci: failed to marshal cmd: %w", err)
	}

	// compose the packet
	p := &pkt{c, make(chan []byte)}

	// keep track of sent packets awaiting responses
	h.sentMutex.Lock()
	h.sent[c.OpCode()] = p
	h.sentMutex.Unlock()

	// write the packet to the socket and check for errors
	if n, err := h.skt.Write(b[:4+c.Len()]); err != nil {
		return fmt.Errorf("hci: failed to send cmd: %w", err)
	} else if n != 4+c.Len() {
		return errors.New("hci: failed to send whole cmd pkt to hci socket")
	}

	// clear sent table when done, we sometimes get command complete or
	// command status messages with no matching send, which can attempt to
	// access stale packets in sent and fail or lock up.
	defer func() {
		h.sentMutex.Lock()
		delete(h.sent, c.OpCode())
		h.sentMutex.Unlock()
	}()

	// emergency timeout to prevent calls from locking up if the HCI
	// interface doesn't respond.  Response here should normally be fast
	// a timeout indicates a major problem with HCI.
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("hci: no response to command: %w", ctx.Err())
		case <-h.Closed():
			if h.err == nil {
				return errors.New("hci: no response to command: disconnected")
			}
			return fmt.Errorf("hci: no response to command: disconnected: %w", h.err)
		case b := <-p.done:
			if len(b) > 0 && b[0] != 0x00 {
				return ErrCommand(b[0])
			}
			if r != nil {
				if err := r.Unmarshal(b); err != nil {
					// assume this is due to receiving a response of previous command that timed out
					h.log.Warn("could not unmarshal bytes",
						slog.String("response", fmt.Sprintf("%T", r)),
						slog.String("bytes", fmt.Sprintf("%02X", b)),
						log.Error(err))
					continue
				}
			}
			return nil
		}
	}
}

func (h *HCI) sentPkt(code int) (*pkt, bool) {
	h.sentMutex.RLock()
	p, ok := h.sent[code]
	h.sentMutex.RUnlock()
	return p, ok
}

func (h *HCI) evtHandler(code int) (handlerFn, bool) {
	h.evtMutex.RLock()
	eventH, ok := h.evth[code]
	h.evtMutex.RUnlock()
	return eventH, ok
}

func (h *HCI) subHandler(code int) (handlerFn, bool) {
	h.subMutex.RLock()
	subH, ok := h.subh[code]
	h.subMutex.RUnlock()
	return subH, ok
}

func (h *HCI) sktLoop() {
	var (
		n   int
		err error
		b   = make([]byte, 4096)
	)
	for {
		// wait for read packet or for the socket to close
		if n, err = h.skt.Read(b); n == 0 || err != nil {
			h.err = fmt.Errorf("skt: %w", err)
			return
		}

		// handle the packet, copy bytes to prevent mangling in threads
		p := make([]byte, n)
		copy(p, b)
		err = h.handlePkt(p)

		// ignore some select errors, others should be printed for future work
		switch {
		case err == nil:
		case errors.Is(err, ErrUnknownCommand), errors.Is(err, ErrUnsupportedCommand):
			if err2 := h.setAllowedCommands(1); err2 != nil {
				h.log.Warn("failed to set allowed commands after read error", slog.String("allowed error", err2.Error()), log.Error(err))
			}
		case errors.Is(err, ErrUnsupportedVendorPacket):
		case errors.Is(err, ErrUnsupportedScoPacket):
		default:
			h.log.Warn("read error", log.Error(err))
		}
	}
}

func (h *HCI) handlePkt(b []byte) error {
	// Strip the 1-byte HCI header and pass down the rest of the packet.
	t, b := b[0], b[1:]
	switch t {
	case pktTypeCommand:
		return fmt.Errorf("%w: % X", ErrUnsupportedCommand, b)
	case pktTypeACLData:
		return h.handleACL(b)
	case pktTypeSCOData:
		return fmt.Errorf("%w: % X", ErrUnsupportedScoPacket, b)
	case pktTypeEvent:
		return h.handleEvt(b)
	case pktTypeVendor:
		return fmt.Errorf("%w: % X", ErrUnsupportedVendorPacket, b)
	default:
		return fmt.Errorf("%w: 0x%02X % X", ErrInvalidPacket, t, b)
	}
}

func (h *HCI) handleACL(b []byte) error {
	h.log.Debug("handling packet acl")

	handle := packet(b).handle()
	h.muConns.Lock()
	c, ok := h.conns[handle]
	h.muConns.Unlock()
	if !ok {
		return nil
	}
	select {
	case c.chInPkt <- b:
		return nil
	case <-h.Closed():
		if h.err == nil {
			return errors.New("hci: no response to command: disconnected")
		}
		return fmt.Errorf("hci device disconnected: %w", h.err)
	}
}

func (h *HCI) handleEvt(b []byte) error {
	code, plen := int(b[0]), int(b[1])
	//h.log.Debug("handling packet event", slog.Int("code", code))
	if plen != len(b[2:]) {
		return fmt.Errorf("invalid event packet: % X", b)
	}
	if code == evt.CommandCompleteCode || code == evt.CommandStatusCode {
		if f, found := h.evtHandler(code); found {
			return f(b[2:])
		}
	}
	if plen != len(b[2:]) {
		h.err = fmt.Errorf("invalid event packet: % X", b)
	}

	if f, found := h.evtHandler(code); found {
		h.err = f(b[2:])
		return nil
	}
	if code == 0xff { // Ignore vendor events
		return nil
	}
	return fmt.Errorf("unsupported event packet: % X", b)
}

func (h *HCI) handleLEMeta(b []byte) error {
	subcode := int(b[0])
	if f, found := h.subHandler(subcode); found {
		err := f(b)
		switch subcode {
		case evt.LEAdvertisingReportSubCode:
			return nil
		default:
			return err
		}
	}
	return fmt.Errorf("unsupported LE event: % X", b)
}

func (h *HCI) handleLEAdvertisingReport(b []byte) error {
	if h.advHandler == nil {
		return nil
	}

	e := evt.LEAdvertisingReport(b)
	for i := 0; i < int(e.NumReports()); i++ {
		var a *Advertisement
		switch e.EventType(i) {
		case evtTypAdvScanInd, evtTypAdvInd:
			a = newAdvertisement(e, i)
			h.adMutex.Lock()
			h.adHist[a.Address().String()] = a
			h.adTimes[a.Address().String()] = time.Now()
			for addr, stamp := range h.adTimes {
				if stamp.Add(5 * time.Minute).Before(time.Now()) {
					delete(h.adTimes, addr)
					delete(h.adHist, addr)
				}
			}
			h.adMutex.Unlock()
		case evtTypScanRsp:
			var ok bool
			sr := newAdvertisement(e, i)

			// Got a SR without having received an associated AD before?
			if a, ok = h.adHist[sr.Address().String()]; !ok {
				return fmt.Errorf("received scan response %s with no associated Advertising Data packet", sr.Address())
			}
			a.setScanResponse(sr)
		default:
			a = newAdvertisement(e, i)
		}

		h.advHandler(a)
	}

	return nil
}

func (h *HCI) handleCommandComplete(b []byte) error {
	e := evt.CommandComplete(b)
	if err := h.setAllowedCommands(int(e.NumHCICommandPackets())); err != nil {
		return fmt.Errorf("unable to set allowed commands: %w", err)
	}

	// NOP command, used for flow control purpose [Vol 2, Part E, 4.4]
	// no handling other than setAllowedCommands needed
	if e.CommandOpcode() == 0x0000 {
		return nil
	}
	p, found := h.sentPkt(int(e.CommandOpcode()))
	if !found {
		return fmt.Errorf("can't find the cmd for CommandCompleteEP (%w): % X", ErrUnknownCommand, e)
	}
	select {
	case p.done <- e.ReturnParameters():
		return nil
	case <-h.Closed():
		return fmt.Errorf("hci device closed: %w", h.err)
	}
}

func (h *HCI) handleCommandStatus(b []byte) error {
	e := evt.CommandStatus(b)
	if err := h.setAllowedCommands(int(e.NumHCICommandPackets())); err != nil {
		return fmt.Errorf("unable to set allowed commands: %w", err)
	}

	p, found := h.sentPkt(int(e.CommandOpcode()))
	if !found {
		return fmt.Errorf("can't find the cmd for CommandStatusEP: % X", e)
	}
	select {
	case p.done <- []byte{e.Status()}:
		return nil
	case <-h.Closed():
		return fmt.Errorf("hci device closed: %w", h.err)
	}
}

func (h *HCI) handleLEConnectionComplete(b []byte) error {
	e := evt.LEConnectionComplete(b)
	c := newConn(h, e)
	h.muConns.Lock()
	h.conns[e.ConnectionHandle()] = c
	h.muConns.Unlock()
	if e.Role() == roleMaster {
		if e.Status() == 0x00 {
			select {
			case h.chMasterConn <- c:
				// sent connection back to dialer
				return nil
			case <-h.Closed():
				// socket closed before connection made
				return fmt.Errorf("hci device closed: %w", h.err)
			}
		}
		if ErrCommand(e.Status()) == ErrConnID {
			// The connection was canceled successfully.
			return nil
		}
		return nil
	}
	if e.Status() == 0x00 {
		select {
		case h.chSlaveConn <- c:
		case <-h.Closed():
			return fmt.Errorf("hci device closed: %w", h.err)
		}
		// When a controller accepts a connection, it moves from advertising
		// state to idle/ready state. Host needs to explicitly ask the
		// controller to re-enable advertising. Note that the host was most
		// likely in advertising state. Otherwise it couldn't accept the
		// connection in the first place. The only exception is that user
		// asked the host to stop advertising during this tiny window.
		// The re-enabling might failed or ignored by the controller, if
		// it had reached the maximum number of concurrent connections.
		// So we also re-enable the advertising when a connection disconnected
		h.params.RLock()
		enabled := h.params.advEnable.AdvertisingEnable
		h.params.RUnlock()

		if enabled == 1 {
			if err := h.Send(context.Background(), &cmd.LESetAdvertiseEnable{AdvertisingEnable: 0}, nil); err != nil {
				_ = h.Close()
				return fmt.Errorf("unable to renable advertising: %w", err)
			}
		}
	}
	if h.connectedHandler != nil {
		h.connectedHandler(e)
	}
	return nil
}

func (h *HCI) handleLEConnectionUpdateComplete(b []byte) error {
	e := evt.LEConnectionUpdateComplete(b)
	h.muConns.Lock()
	c, ok := h.conns[e.ConnectionHandle()]
	h.muConns.Unlock()
	if !ok {
		return fmt.Errorf("le connection update complete has invalid connection handle %04X", e.ConnectionHandle())
	}
	return c.handleLEConnectionUpdateComplete(e)
}

func (h *HCI) handleDisconnectionComplete(b []byte) error {
	e := evt.DisconnectionComplete(b)
	h.muConns.Lock()
	c, found := h.conns[e.ConnectionHandle()]
	delete(h.conns, e.ConnectionHandle())
	h.muConns.Unlock()

	if !found {
		return fmt.Errorf("disconnecting an invalid handle %04X", e.ConnectionHandle())
	}

	close(c.chInPkt)

	if c.param.Role() == roleSlave {
		// Re-enable advertising, if it was advertising. Refer to the
		// handleLEConnectionComplete() for details.
		// This may failed with ErrCommandDisallowed, if the controller
		// was actually in advertising state. It does no harm though.
		h.params.RLock()
		if h.params.advEnable.AdvertisingEnable == 1 {
			h.Add(1)
			go func() {
				defer h.Done()
				if err := h.Send(context.Background(), &h.params.advEnable, nil); err != nil {
					h.err = fmt.Errorf("unable to reenable advertising: %w", err)
				}
			}()
		}
		h.params.RUnlock()
	} else {
		// remote peripheral disconnected
		select {
		case <-c.chDone:
		default:
			close(c.chDone)
		}
	}

	// When a connection disconnects, all the sent packets and weren't acked yet
	// will be recycled. [Vol2, Part E 4.1.1]
	//
	// must be done with the pool locked to avoid race conditions where
	// writePDU is in progress and does a Get from the pool after this completes,
	// leaking a buffer from the main pool.
	c.txBuffer.LockPool()
	c.txBuffer.PutAll()
	c.txBuffer.UnlockPool()
	if h.disconnectedHandler != nil {
		h.disconnectedHandler(e)
	}
	return nil
}

func (h *HCI) handleNumberOfCompletedPackets(b []byte) error {
	e := evt.NumberOfCompletedPackets(b)
	h.muConns.Lock()
	defer h.muConns.Unlock()
	for i := 0; i < int(e.NumberOfHandles()); i++ {
		c, found := h.conns[e.ConnectionHandle(i)]
		if !found {
			continue
		}

		// Put the delivered buffers back to the pool.
		for j := 0; j < int(e.HCNumOfCompletedPackets(i)); j++ {
			c.txBuffer.Put()
		}
	}
	return nil
}

func (h *HCI) handleLELongTermKeyRequest(b []byte) error {
	e := evt.LELongTermKeyRequest(b)
	return h.Send(context.Background(), &cmd.LELongTermKeyRequestNegativeReply{
		ConnectionHandle: e.ConnectionHandle(),
	}, nil)
}

func (h *HCI) handleInquiryComplete(_ []byte) error {
	return nil
}

func (h *HCI) handleExtendedInquiry(b []byte) error {
	if h.inqHandler == nil {
		return nil
	}
	// always a single response [Vol2, 7.7.38]
	h.Add(1)
	go func() {
		defer h.Done()
		h.inqHandler(newInquiry(evt.ExtendedInquiry(b), 0))
	}()

	return nil
}

func (h *HCI) handleInquiryResult(b []byte) error {
	if h.inqHandler == nil {
		return nil
	}

	e := evt.InquiryResult(b)
	for i := 0; i < int(e.NumResponses()); i++ {
		h.Add(1)
		go func(e evt.InquiryResult, i int) {
			defer h.Done()
			h.inqHandler(newInquiry(e, i))
		}(e, i)
	}
	return nil
}

func (h *HCI) handleInquiryWithRSSI(b []byte) error {
	if h.inqHandler == nil {
		return nil
	}

	e := evt.InquiryResultwithRSSI(b)
	for i := 0; i < int(e.NumResponses()); i++ {
		h.Add(1)
		go func(e evt.InquiryResultwithRSSI, i int) {
			defer h.Done()
			h.inqHandler(newInquiry(e, i))
		}(e, i)
	}

	return nil
}

func (h *HCI) handleConnectionComplete(b []byte) error {

	e := evt.ConnectionComplete(b)
	c := newConn(h, e)

	h.muConns.Lock()
	h.conns[e.ConnectionHandle()] = c
	h.muConns.Unlock()

	if e.Status() == 0x00 {
		select {
		case h.chMasterBREDRConn <- c:
			return nil
		case <-h.Closed():
			return fmt.Errorf("hci device closed: %w", h.err)
		}
	}
	if ErrCommand(e.Status()) == ErrConnID {
		// The connection was canceled successfully.
		return nil
	}
	return nil
}

func (h *HCI) handleReadRemoteSupportedFeaturesComplete(b []byte) error {
	e := evt.ReadRemoteSupportedFeaturesComplete(b)
	if e.Status() == 0x00 {
		h.muConns.Lock()
		h.conns[e.ConnectionHandle()].lmpFeatures = e.LMPFeatures()
		h.muConns.Unlock()
	}

	p, found := h.sentPkt(evt.ReadRemoteSupportedFeaturesCompleteCode)
	if !found {
		return fmt.Errorf("can't find the cmd for CommandReadRemoteSupportedFeatureEP: % X", e)
	}
	select {
	case p.done <- []byte{e.Status()}:
		return nil
	case <-h.Closed():
		return fmt.Errorf("hci device closed: %w", h.err)
	}
}

func (h *HCI) handlePageScanRepetitionModeChange(_ []byte) error {
	//e := evt.PageScanRepetitionModeChange(b)

	// remote controller has successfully changed the page
	// scan repetition mode [ Vol 2, 3.7, Table 3.8 ]
	return nil
}

func (h *HCI) handleMaxSlotsChange(_ []byte) error {
	//e := evt.MaxSlotsChange(b)

	// remote controller has successfully changed the page
	// scan repetition mode [ Vol 2, 3.7, Table 3.8 ]
	return nil
}

// [ Vol 2, 7.7.7 ]
func (h *HCI) handleReadRemoteNameRequestCompleteEvent(b []byte) error {
	e := evt.RemoteNameRequestComplete(b)

	nameEvent := &nameEvent{}
	if e.Status() == 0x00 {
		name := e.RemoteName()
		i := bytes.IndexByte(name[:], 0x00)
		if i == -1 {
			i = len(name)
		}
		nameEvent.name = string(name[0:i])
	} else {
		nameEvent.err = fmt.Errorf("the remote name request command failed: %X", e.Status())
	}

	h.nameHandlers.Lock()
	ch, ok := h.nameHandlers.handlers[ble.NewAddr(fmt.Sprintf("%X", e.BDADDR()))]
	h.nameHandlers.Unlock()
	if !ok {
		return fmt.Errorf("received remote name request complete from unknown address: %X", e.BDADDR())
	}
	select {
	case ch <- nameEvent:
		return nil
	case <-h.Closed():
		return fmt.Errorf("hci device closed: %w", h.err)
	}
}

func (h *HCI) setAllowedCommands(n int) error {

	//hard-coded limit to command queue depth
	//matches make(chan []byte, 16) in NewHCI
	// TODO make this a constant, decide correct size
	if n > 16 {
		n = 16
	}

	for len(h.chCmdBufs) < n {
		select {
		case h.chCmdBufs <- make([]byte, 64): // TODO make buffer size a constant
		case <-h.Closed():
			return fmt.Errorf("hci device closed: %w", h.err)
		}
	}
	return nil
}
