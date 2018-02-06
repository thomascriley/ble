package hci

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
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
func (h *HCI) Inquire(length int, numResponses int) error {
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
	return h.Send(&h.params.inquiry, nil)
}

// StopInquiry stops inquiring for BR/EDR devices by sending an InquiryCancel command
func (h *HCI) StopInquiry() error {
	return h.Send(&h.params.inquiryCancel, nil)
}

func (h *HCI) RequestRemoteName(a ble.Addr) (string, error) {
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
	h.Send(req, nil)

	select {
	case nameEvent := <-ch:
		return nameEvent.name, nameEvent.err
	case <-time.After(60 * time.Second):
		return "", errors.New("Timed out waiting for remote name response")
	}
}

func (h *HCI) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode, channel uint8) (cli ble.RFCOMMClient, err error) {
	b, err := net.ParseMAC(a.String())
	if err != nil {
		return nil, ErrInvalidAddr
	}
	addr := [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}

	h.params.Lock()
	h.params.connBREDRParams.BDADDR = addr
	h.params.connBREDRParams.ClockOffset = clockOffset
	h.params.connBREDRParams.PageScanRepetitionMode = pageScanRepetitionMode
	h.params.connBREDRParams.AllowRoleSwitch = roleMaster
	err = h.Send(&h.params.connBREDRParams, nil)
	if err != nil && err.Error() == ErrACLConnExists.Error() {
		h.params.connCancelBREDR.BDADDR = [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}
		err = h.Send(&h.params.connCancelBREDR, nil)
		err = h.Send(&h.params.connBREDRParams, nil)
	}
	h.params.Unlock()

	if err != nil {
		return nil, err
	}

	var tmo <-chan time.Time
	if h.dialerTmo != time.Duration(0) {
		tmo = time.After(h.dialerTmo)
	}
	select {
	case <-h.done:
		return nil, h.err
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

		timeout := time.Duration(15 * time.Second)

		if err = c.InformationRequest(l2cap.InfoTypeConnectionlessMTU, timeout); err != nil {
			fmt.Printf("Error: %s", err.Error())
			//c.Close()
			//return nil, err
		}

		if err = c.InformationRequest(l2cap.InfoTypeExtendedFeatures, timeout); err != nil {
			fmt.Printf("Error: %s", err.Error())
			//c.Close()
			//return nil, err
		}

		// 1.2 - 2.1 + EDR will return not supported
		c.InformationRequest(l2cap.InfoTypeFixedChannels, timeout)

		if err = c.ConnectionRequest(psmRFCOMM, timeout); err != nil {
			c.Close()
			return nil, err
		}

		// Even if all default values are acceptable, a Configuration Request
		// packet with no options shall be sent. [Vol 3, Part A, 4.4]
		// TODO: make this non-static

		mtuOption := &l2cap.MTUOption{MTU: 0x03f5}
		//mtuOption.SetHint(0x01)
		if err = c.ConfigurationRequest([]l2cap.Option{mtuOption}, timeout); err != nil {
			c.Close()
			return nil, err
		}

		select {
		case <-ctx.Done():
			c.Close()
			return nil, ctx.Err()
		case <-c.cfgRequest:
		}

		if cli, err = rfcomm.NewClient(ctx, c, channel); err != nil {
			c.Close()
		}
		return cli, err

	case <-ctx.Done():
		return h.cancelConnectionBREDR(addr, ctx.Err())
	case <-tmo:
		return h.cancelConnectionBREDR(addr, fmt.Errorf("connection timed out"))
	}
}

func (h *HCI) cancelConnectionBREDR(addr [6]byte, connErr error) (ble.RFCOMMClient, error) {
	h.params.Lock()
	h.params.connCancelBREDR.BDADDR = addr
	err := h.Send(&h.params.connCancelBREDR, nil)
	h.params.Unlock()

	if err == nil {
		// The pending connection was canceled successfully.
		return nil, connErr
	}
	// The connection has been established, the cancel command
	// failed with ErrDisallowed.
	if err == ErrDisallowed {
		c := <-h.chMasterBREDRConn
		c.Close()
	}
	return nil, errors.Wrap(err, "cancel connection failed")
}
