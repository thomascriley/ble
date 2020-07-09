package hci

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/adv"
	"github.com/thomascriley/ble/linux/gatt"
)

// Addr ...
func (h *HCI) Addr() ble.Addr { return h.addr }

// SetAdvHandler ...
func (h *HCI) SetAdvHandler(ah ble.AdvHandler) error {
	h.advHandler = ah
	return nil
}

// Scan starts scanning.
func (h *HCI) Scan(ctx context.Context, allowDup bool) error {
	h.params.scanEnable.FilterDuplicates = 1
	if allowDup {
		h.params.scanEnable.FilterDuplicates = 0
	}
	h.params.scanEnable.LEScanEnable = 1
	h.adHist = make(map[string]*Advertisement, 128)
	return h.Send(ctx, &h.params.scanEnable, nil)
}

// StopScanning stops scanning.
func (h *HCI) StopScanning(ctx context.Context) error {
	h.params.scanEnable.LEScanEnable = 0
	return h.Send(ctx, &h.params.scanEnable, nil)
}

// AdvertiseAdv advertises a given Advertisement, context is used for timing out long running send command to the hci
// device in case the device does not respond as expected
func (h *HCI) AdvertiseAdv(ctx context.Context, a ble.Advertisement) error {
	ad, err := adv.NewPacket(adv.Flags(adv.FlagGeneralDiscoverable | adv.FlagLEOnly))
	if err != nil {
		return err
	}
	f := adv.AllUUID

	// Current length of ad packet plus two bytes of length and tag.
	l := ad.Len() + 1 + 1
	for _, u := range a.Services() {
		l += u.Len()
	}
	if l > adv.MaxEIRPacketLength {
		f = adv.SomeUUID
	}
	for _, u := range a.Services() {
		if err := ad.Append(f(u)); err != nil {
			if err == adv.ErrNotFit {
				break
			}
			return err
		}
	}
	sr, _ := adv.NewPacket()
	switch {
	case ad.Append(adv.CompleteName(a.LocalName())) == nil:
	case sr.Append(adv.CompleteName(a.LocalName())) == nil:
	case sr.Append(adv.ShortName(a.LocalName())) == nil:
	}

	if a.ManufacturerData() != nil {
		manufacuturerData := adv.ManufacturerData(1337, a.ManufacturerData())
		switch {
		case ad.Append(manufacuturerData) == nil:
		case sr.Append(manufacuturerData) == nil:
		}
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), sr.Bytes()); err != nil {
		return nil
	}
	return h.Advertise(ctx)

}

// AdvertiseNameAndServices advertises device name, and specified service UUIDs.
// It tries to fit the UUIDs in the advertising data as much as possible.
// If name doesn't fit in the advertising data, it will be put in scan response.
func (h *HCI) AdvertiseNameAndServices(ctx context.Context, name string, uuids ...ble.UUID) error {
	ad, err := adv.NewPacket(adv.Flags(adv.FlagGeneralDiscoverable | adv.FlagLEOnly))
	if err != nil {
		return err
	}
	f := adv.AllUUID

	// Current length of ad packet plus two bytes of length and tag.
	l := ad.Len() + 1 + 1
	for _, u := range uuids {
		l += u.Len()
	}
	if l > adv.MaxEIRPacketLength {
		f = adv.SomeUUID
	}
	for _, u := range uuids {
		if err := ad.Append(f(u)); err != nil {
			if err == adv.ErrNotFit {
				break
			}
			return err
		}
	}
	sr, _ := adv.NewPacket()
	switch {
	case ad.Append(adv.CompleteName(name)) == nil:
	case sr.Append(adv.CompleteName(name)) == nil:
	case sr.Append(adv.ShortName(name)) == nil:
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), sr.Bytes()); err != nil {
		return fmt.Errorf("unable to set advertisement: %w", err)
	}
	return h.Advertise(ctx)
}

// AdvertiseMfgData avertises the given manufacturer data.
func (h *HCI) AdvertiseMfgData(ctx context.Context, id uint16, md []byte) error {
	ad, err := adv.NewPacket(adv.ManufacturerData(id, md))
	if err != nil {
		return err
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), nil); err != nil {
		return nil
	}
	return h.Advertise(ctx)
}

// AdvertiseServiceData16 advertises data associated with a 16bit service uuid
func (h *HCI) AdvertiseServiceData16(ctx context.Context,id uint16, b []byte) error {
	ad, err := adv.NewPacket(adv.ServiceData16(id, b))
	if err != nil {
		return err
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), nil); err != nil {
		return nil
	}
	return h.Advertise(ctx)
}

// AdvertiseIBeaconData advertise iBeacon with given manufacturer data.
func (h *HCI) AdvertiseIBeaconData(ctx context.Context,md []byte) error {
	ad, err := adv.NewPacket(adv.IBeaconData(md))
	if err != nil {
		return err
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), nil); err != nil {
		return nil
	}
	return h.Advertise(ctx)
}

// AdvertiseIBeacon advertises iBeacon with specified parameters.
func (h *HCI) AdvertiseIBeacon(ctx context.Context,u ble.UUID, major, minor uint16, pwr int8) error {
	ad, err := adv.NewPacket(adv.IBeacon(u, major, minor, pwr))
	if err != nil {
		return err
	}
	if err := h.SetAdvertisement(ctx, ad.Bytes(), nil); err != nil {
		return nil
	}
	return h.Advertise(ctx)
}

// StopAdvertising stops advertising.
func (h *HCI) StopAdvertising(ctx context.Context) error {
	h.params.advEnable.AdvertisingEnable = 0
	return h.Send(ctx, &h.params.advEnable, nil)
}

// Accept starts advertising and accepts connection.
func (h *HCI) Accept() (ble.Conn, error) {
	select {
	case <-h.Closed():
		return nil, h.err
	case c := <-h.chSlaveConn:
		return c, nil
	}
}

// Dial ...
func (h *HCI) Dial(ctx context.Context, a ble.Addr, addressType ble.AddressType) (ble.ClientBLE, error) {
	b, err := net.ParseMAC(a.String())
	if err != nil {
		return nil, ErrInvalidAddr
	}

	h.params.Lock()
	h.params.connParams.PeerAddress = [6]byte{b[5], b[4], b[3], b[2], b[1], b[0]}
	if addressType == ble.AddressTypeRandom {
		h.params.connParams.PeerAddressType = 1
	} else {
		h.params.connParams.PeerAddressType = 0
	}
	err = h.Send(ctx, &h.params.connParams, nil)
	h.params.Unlock()

	if err != nil {
		return nil, err
	}

	cancelCTX, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return h.cancelDial(cancelCTX)
	case <-h.Closed():
		return nil, h.err
	case c := <-h.chMasterConn:
		c.SourceID = cidLEAtt
		c.DestinationID = cidLEAtt
		return gatt.NewClient(c)
	}
}

// cancelDial cancels the Dialing
func (h *HCI) cancelDial(ctx context.Context) (ble.ClientBLE, error) {
	err := h.Send(ctx, &h.params.connCancel, nil)
	if err == nil {
		// The pending connection was canceled successfully.
		return nil, fmt.Errorf("connection canceled")
	}
	// The connection has been established, the cancel command
	// failed with ErrDisallowed.
	if err == ErrDisallowed {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-h.Closed():
			return nil, h.err
		case ch := <- h.chMasterConn:
			return gatt.NewClient(ch)
		}

	}
	return nil, fmt.Errorf( "cancel connection failed: %w", err)
}

// Advertise starts advertising.
func (h *HCI) Advertise(ctx context.Context) error {
	h.params.advEnable.AdvertisingEnable = 1
	return h.Send(ctx, &h.params.advEnable, nil)
}

// SetAdvertisement sets advertising data and scanResp.
func (h *HCI) SetAdvertisement(ctx context.Context, ad []byte, sr []byte) error {
	if len(ad) > adv.MaxEIRPacketLength || len(sr) > adv.MaxEIRPacketLength {
		return ble.ErrEIRPacketTooLong
	}

	h.params.advData.AdvertisingDataLength = uint8(len(ad))
	copy(h.params.advData.AdvertisingData[:], ad)
	if err := h.Send(ctx, &h.params.advData, nil); err != nil {
		return err
	}

	h.params.scanResp.ScanResponseDataLength = uint8(len(sr))
	copy(h.params.scanResp.ScanResponseData[:], sr)
	if err := h.Send(ctx, &h.params.scanResp, nil); err != nil {
		return err
	}
	return nil
}

func (h *HCI) cancelConnection(ctx context.Context, connErr error) (ble.Client, error) {
	h.params.Lock()
	err := h.Send(ctx, &h.params.connCancel, nil)
	h.params.Unlock()

	if err == nil {
		// The pending connection was canceled successfully.
		return nil, connErr
	}
	// The connection has been established, the cancel command
	// failed with ErrDisallowed.
	if err == ErrDisallowed {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-h.Closed():
			return nil, h.err
		case ch := <- h.chMasterConn:
			return gatt.NewClient(ch)
		}
	}
	return nil, fmt.Errorf( "cancel connection failed: %w", err)
}
