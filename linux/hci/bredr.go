package hci

import (
	"context"
	"fmt"
	"github.com/thomascriley/ble/log"
	"net"
	"time"
	"errors"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/hci/cmd"
	"github.com/thomascriley/ble/linux/l2cap"
	"github.com/thomascriley/ble/linux/rfcomm"
)

type nameEvent struct {
	name string
	err  error
}

// SetAdvHandler ...
func (h *HCI) SetInqHandler(ah ble.InqHandler) error {
	h.inqHandler = ah
	return nil
}

// Inquire starts inquiring for BR/EDR devices.
func (h *HCI) Inquire(ctx context.Context, length int, numResponses int) error {
	if numResponses > 255 || numResponses < 0 {
		numResponses = 0x00 // unlimited
	}
	if length > 0x30 || length <= 0 {
		length = 0x30 // max length (61.44 sec)
	}

	// General/Unlimited Inquiry Access Code (GIAC)
	// see: https://www.bluetooth.com/specifications/assigned-numbers/baseband
	h.params.inquiry.LAP = [3]byte{0x33, 0x8b, 0x9e}
	h.params.inquiry.InquiryLength = uint8(length)
	h.params.inquiry.NumResponses = uint8(numResponses)
	return h.Send(ctx, &h.params.inquiry, nil)
}

// StopInquiry stops inquiring for BR/EDR devices by sending an InquiryCancel command
func (h *HCI) StopInquiry(ctx context.Context) error {
	return h.Send(ctx, &h.params.inquiryCancel, nil)
}

func (h *HCI) RequestRemoteName(ctx context.Context, a ble.Addr) (string, error) {
	bdaddr := a.(net.HardwareAddr)

	ch := make(chan *nameEvent)
	defer func() {
		h.nameHandlers.Lock()
		close(ch)
		delete(h.nameHandlers.handlers, a)
		h.nameHandlers.Unlock()
	}()

	h.nameHandlers.Lock()
	h.nameHandlers.handlers[a] = ch
	h.nameHandlers.Unlock()

	req := &cmd.RemoteNameRequest{
		ClockOffset:          0x0000,
		PageScanRepitionMode: 0x00}
	copy(req.BDADDR[:], bdaddr)
	if err := h.Send(ctx, req, nil); err != nil {
		return "", fmt.Errorf("remote name request: %s", err)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case nameEvent := <-ch:
		return nameEvent.name, nameEvent.err
	case <-h.Closed():
		return "", fmt.Errorf("disconnected: %w", h.err)
	}
}

func (h *HCI) sendBREDRParams(ctx context.Context, addr [6]byte, clockOffset uint16, pageScanRepetitionMode uint8) error {
	h.params.Lock()
	defer h.params.Unlock()
	h.params.connBREDRParams.BDADDR = addr
	h.params.connBREDRParams.ClockOffset = clockOffset
	h.params.connBREDRParams.PageScanRepetitionMode = pageScanRepetitionMode
	h.params.connBREDRParams.AllowRoleSwitch = roleMaster
	err := h.Send(ctx, &h.params.connBREDRParams, nil)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrACLConnExists):
		h.params.connCancelBREDR.BDADDR = addr
		if err = h.Send(ctx, &h.params.connCancelBREDR, nil); err != nil {
			return fmt.Errorf("unable to cancel existing connection: %w", err)
		}
		if err = h.Send(ctx, &h.params.connBREDRParams, nil); err != nil {
			return fmt.Errorf("unable to reestabolish existing connection: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unable to send BREDR params: %w", err)
	}
}

func (h *HCI) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (cli ble.ClientRFCOMM, err error) {
	b, err := net.ParseMAC(a.String())
	if err != nil {
		return nil, ErrInvalidAddr
	}
	addr := [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}

	if err := h.sendBREDRParams(ctx, addr, clockOffset,pageScanRepetitionMode); err != nil {
		return nil, fmt.Errorf("unable to send BREDR params: %w", err)
	}

	select {
	case <-h.Closed():
		return nil, fmt.Errorf("hci device closed: %w", h.err)
	case <-ctx.Done():
		return nil, h.cancelConnectionBREDR(ctx, addr, fmt.Errorf("connection timed out"))
	case c := <-h.chMasterBREDRConn:

		// increment to the next available dynamicCID channel id
		c.SourceID = h.dynamicCID
		if h.dynamicCID == 0xFFFF {
			h.dynamicCID = cidDynamicStart
		} else {
			// this doesn't work for some reason, maybe they need to be reused after
			// connection closes
			//h.dynamicCID++
		}

		timeout := 15 * time.Second

		if err = c.InformationRequest(ctx, l2cap.InfoTypeConnectionlessMTU, timeout); err != nil {
			log.Printf("Warning: unable to make information request for connectionless mtu: %s\n", err)
			//c.Close()
			//return nil, err
		}

		if err = c.InformationRequest(ctx, l2cap.InfoTypeExtendedFeatures, timeout); err != nil {
			log.Printf("Warning: unable to make information request for extended features: %s\n", err)
			//c.Close()
			//return nil, err
		}

		// 1.2 - 2.1 + EDR will return not supported
		if err = c.InformationRequest(ctx, l2cap.InfoTypeFixedChannels, timeout); err != nil {
			log.Printf("Warning: unable to make information request for fixed channels: %s\n", err)
		}

		if err = c.ConnectionRequest(ctx, psmRFCOMM, timeout); err != nil {
			_ = c.Close(ctx)
			return nil, fmt.Errorf("unable to make connection request: %w", err)
		}

		// Even if all default values are acceptable, a Configuration Request
		// packet with no options shall be sent. [Vol 3, Part A, 4.4]
		// TODO: make this non-static

		mtuOption := &l2cap.MTUOption{MTU: 0x03f5}
		//mtuOption.SetHint(0x01)
		if err = c.ConfigurationRequest(ctx, []l2cap.Option{mtuOption}, timeout); err != nil {
			_ = c.Close(ctx)
			return nil, fmt.Errorf("unable to make configuration request: %w", err)
		}

		select {
		case <-h.Closed():
			_ = c.Close(ctx)
			return nil, fmt.Errorf("hci device closed: %w", h.err)
		case <-ctx.Done():
			_ = c.Close(ctx)
			return nil, errors.New("timed out waiting for cfgRequest")
		case <-c.cfgRequest:
		}

		cli := rfcomm.NewClient(c, channel)
		if err := cli.DialContext(ctx); err != nil {
			_ = c.Close(ctx)
			return nil, fmt.Errorf("unable to dial client: %w", err)
		}
		return cli, err
	}
}

func (h *HCI) cancelConnectionBREDR(ctx context.Context, addr [6]byte, connErr error) error {
	h.params.Lock()
	h.params.connCancelBREDR.BDADDR = addr
	err := h.Send(ctx, &h.params.connCancelBREDR, nil)
	h.params.Unlock()

	if err == nil {
		// The pending connection was canceled successfully.
		return connErr
	}
	// The connection has been established, the cancel command
	// failed with ErrDisallowed.
	if err == ErrDisallowed {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out canceling connection: %w", ctx.Err())
		case <-h.Closed():
			return fmt.Errorf("hci device closed: %w", h.err)
		case c := <-h.chMasterBREDRConn:
			_ = c.Close(ctx)
		}
	}
	return fmt.Errorf( "cancel connection failed: %w", err)
}
