package hci

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux/hci/cmd"
	"github.com/currantlabs/ble/linux/l2cap"
	"github.com/currantlabs/ble/linux/rfcomm"
	"github.com/pkg/errors"
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

func (h *HCI) DialRFCOMM(ctx context.Context, a ble.Addr, clockOffset uint16, pageScanRepetitionMode uint8) (ble.RFCOMMClient, error) {
	b, err := net.ParseMAC(a.String())
	if err != nil {
		return nil, ErrInvalidAddr
	}

	h.params.Lock()
	h.params.connBREDRParams.BDADDR = [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}
	h.params.connBREDRParams.ClockOffset = clockOffset
	h.params.connBREDRParams.PageScanRepetitionMode = pageScanRepetitionMode
	err = h.Send(&h.params.connBREDRParams, nil)
	h.params.Unlock()

	if err != nil {
		return nil, err
	}

	var tmo <-chan time.Time
	if h.dialerTmo != time.Duration(0) {
		tmo = time.After(h.dialerTmo)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-h.done:
		return nil, h.err
	case c := <-h.chMasterSPPConn:
		// increment to the next available dynamicCID channel id
		c.SourceID = h.dynamicCID
		if h.dynamicCID == 0xFFFF {
			h.dynamicCID = cidDynamicStart
		} else {
			h.dynamicCID++
		}

		timeout := time.Duration(15 * time.Second)

		if err = c.InformationRequest(l2cap.InfoTypeConnectionlessMTU, timeout); err != nil {
			return nil, err
		}
		if err = c.InformationRequest(l2cap.InfoTypeExtendedFeatures, timeout); err != nil {
			return nil, err
		}
		// 1.2 - 2.1 + EDR will return not supported
		c.InformationRequest(l2cap.InfoTypeFixedChannels, timeout)

		if err = c.ConnectionRequest(psmRFCOMM, timeout); err != nil {
			return nil, err
		}

		// Even if all default values are acceptable, a Configuration Request
		// packet with no options shall be sent. [Vol 3, Part A, 4.4]
		// TODO: make this non-static
		options := []l2cap.Option{&l2cap.MTUOption{MTU: 0x03f5}}
		if err = c.ConfigurationRequest(options, timeout); err != nil {
			return nil, err
		}

		return rfcomm.NewClient(c)
	case <-tmo:
		h.params.Lock()
		h.params.connCancelBREDR.BDADDR = [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}
		err = h.Send(&h.params.connCancelBREDR, nil)
		h.params.Unlock()

		if err == nil {
			// The pending connection was canceled successfully.
			return nil, fmt.Errorf("connection timed out")
		}
		// The connection has been established, the cancel command
		// failed with ErrDisallowed.
		if err == ErrDisallowed {
			return rfcomm.NewClient(<-h.chMasterSPPConn)
		}
		return nil, errors.Wrap(err, "cancel connection failed")
	}
}
