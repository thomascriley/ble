package gatt

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/att"
	"log/slog"
	"sync"
)

const (
	cccNotify   = 0x0001
	cccIndicate = 0x0002
)

// NewClient returns a GATT Client.
func NewClient(log *slog.Logger, conn ble.Conn) (*Client, error) {
	p := &Client{
		subs: make(map[uint16]*sub),
		conn: conn,
		addr: conn.LocalAddr().String(),
		log: log.With(
			slog.String("peripheral", conn.LocalAddr().String()),
			slog.String("proto", "GATT"),
		),
	}
	p.ac = att.NewClient(conn, p)
	p.Add(1)
	go func() {
		defer p.Done()
		p.ac.Loop()
	}()
	return p, nil
}

// A Client is a GATT Client.
type Client struct {
	sync.RWMutex
	sync.WaitGroup

	profile *ble.Profile
	name    string
	subs    map[uint16]*sub

	ac   *att.Client
	conn ble.Conn

	log *slog.Logger

	addr string
}

func (p *Client) String() string { return p.addr }

func (p *Client) RLock() {
	p.log.Debug("read lock")
	p.RWMutex.RLock()
}

func (p *Client) RUnlock() {
	p.log.Debug("read unlock")
	p.RWMutex.RUnlock()
}

func (p *Client) Lock() {
	p.log.Debug("lock")
	p.RWMutex.Lock()
}

func (p *Client) Unlock() {
	p.log.Debug("unlocking")
	p.RWMutex.Unlock()
}

// Addr returns the address of the client.
func (p *Client) Address() ble.Addr {
	p.log.Debug("getting address")
	defer p.log.Debug("got address")
	p.RLock()
	defer p.RUnlock()
	return p.conn.RemoteAddr()
}

// Name returns the name of the client.
func (p *Client) Name() string {
	p.log.Debug("getting name")
	defer p.log.Debug("got name")
	p.RLock()
	defer p.RUnlock()
	return p.name
}

func (p *Client) Connection() ble.Conn {
	return p.conn
}

// Profile returns the discovered profile.
func (p *Client) Profile() *ble.Profile {
	p.log.Debug("getting profile")
	defer p.log.Debug("got profile")
	p.RLock()
	defer p.RUnlock()
	return p.profile
}

// DiscoverProfile discovers the whole hierarchy of a server.
func (p *Client) DiscoverProfile(force bool) (*ble.Profile, error) {
	p.log.Debug("discovering profile")
	if p.profile != nil && !force {
		return p.profile, nil
	}
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		return nil, fmt.Errorf("can't discover services: %s", err)
	}

	var cs []*ble.Characteristic
	var c *ble.Characteristic
	for _, s := range ss {
		if cs, err = p.DiscoverCharacteristics(nil, s); err != nil {
			return nil, fmt.Errorf("can't discover characteristics: %s", err)
		}
		for _, c = range cs {
			if _, err = p.DiscoverDescriptors(nil, c); err != nil {
				return nil, fmt.Errorf("can't discover descriptors: %s", err)
			}
		}
	}
	p.profile = &ble.Profile{Services: ss}
	p.log.Debug("discovered profile")
	return p.profile, nil
}

// DiscoverServices finds all the primary services on a server. [Vol 3, Part G, 4.4.1]
// If filter is specified, only filtered services are returned.
func (p *Client) DiscoverServices(filter []ble.UUID) ([]*ble.Service, error) {
	p.log.Debug("discovering services")
	p.Lock()
	defer p.Unlock()
	if p.profile == nil {
		p.profile = &ble.Profile{}
	}
	start := uint16(0x0001)
	for {
		length, b, err := p.ac.ReadByGroupType(start, 0xFFFF, ble.PrimaryServiceUUID)
		if errors.Is(err, ble.ErrAttrNotFound) {
			p.log.Debug("discovered services", slog.Int("count", len(p.profile.Services)))
			return p.profile.Services, nil
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read by group type: %w", err)
		}
		for len(b) != 0 {
			h := binary.LittleEndian.Uint16(b[:2])
			endh := binary.LittleEndian.Uint16(b[2:4])
			u := ble.UUID(b[4:length])
			if filter == nil || ble.Contains(filter, u) {
				s := &ble.Service{
					UUID:      u,
					Handle:    h,
					EndHandle: endh,
				}
				p.profile.Services = append(p.profile.Services, s)
			}
			if endh == 0xFFFF {
				p.log.Debug("discovered services", slog.Int("count", len(p.profile.Services)))
				return p.profile.Services, nil
			}
			start = endh + 1
			b = b[length:]
		}
	}
}

// DiscoverIncludedServices finds the included services of a service. [Vol 3, Part G, 4.5.1]
// If filter is specified, only filtered services are returned.
func (p *Client) DiscoverIncludedServices(_ []ble.UUID, _ *ble.Service) ([]*ble.Service, error) {
	p.log.Debug("discovering included services")
	p.Lock()
	defer p.Unlock()
	return nil, nil
}

// DiscoverCharacteristics finds all the characteristics within a service. [Vol 3, Part G, 4.6.1]
// If filter is specified, only filtered characteristics are returned.
func (p *Client) DiscoverCharacteristics(filter []ble.UUID, s *ble.Service) ([]*ble.Characteristic, error) {
	p.log.Debug("discovering characteristics in service", slog.String("uuid", s.UUID.String()))
	p.Lock()
	defer p.Unlock()
	start := s.Handle
	var lastChar *ble.Characteristic
	for start <= s.EndHandle {
		length, b, err := p.ac.ReadByType(start, s.EndHandle, ble.CharacteristicUUID)
		if errors.Is(err, ble.ErrAttrNotFound) {
			break
		} else if err != nil {
			return nil, err
		}
		for len(b) != 0 {
			h := binary.LittleEndian.Uint16(b[:2])
			p := ble.Property(b[2])
			vh := binary.LittleEndian.Uint16(b[3:5])
			u := ble.UUID(b[5:length])
			c := &ble.Characteristic{
				UUID:        u,
				Property:    p,
				Handle:      h,
				ValueHandle: vh,
				EndHandle:   s.EndHandle,
			}
			if filter == nil || ble.Contains(filter, u) {
				s.Characteristics = append(s.Characteristics, c)
			}
			if lastChar != nil {
				lastChar.EndHandle = c.Handle - 1
			}
			lastChar = c
			start = vh + 1
			b = b[length:]
		}
	}
	p.log.Debug("discovered characteristics in service", slog.String("uuid", s.UUID.String()), slog.Int("count", len(s.Characteristics)))
	return s.Characteristics, nil
}

// DiscoverDescriptors finds all the descriptors within a characteristic. [Vol 3, Part G, 4.7.1]
// If filter is specified, only filtered descriptors are returned.
func (p *Client) DiscoverDescriptors(filter []ble.UUID, c *ble.Characteristic) ([]*ble.Descriptor, error) {
	p.log.Debug("discovering descriptors", slog.String("uuid", c.UUID.String()))
	p.Lock()
	defer p.Unlock()
	start := c.ValueHandle + 1
	for start <= c.EndHandle {
		found, b, err := p.ac.FindInformation(start, c.EndHandle)
		if errors.Is(err, ble.ErrAttrNotFound) {
			break
		} else if err != nil {
			return nil, err
		}
		length := 2 + 2
		if found == 0x02 {
			length = 2 + 16
		}
		for len(b) != 0 {
			h := binary.LittleEndian.Uint16(b[:2])
			u := ble.UUID(b[2:length])
			d := &ble.Descriptor{UUID: u, Handle: h}
			if filter == nil || ble.Contains(filter, u) {
				c.Descriptors = append(c.Descriptors, d)
			}
			if u.Equal(ble.ClientCharacteristicConfigUUID) {
				c.CCCD = d
			}
			start = h + 1
			b = b[length:]
		}
	}
	p.log.Debug("discovered descriptors in characteristic", slog.String("uuid", c.UUID.String()), slog.Int("count", len(c.Descriptors)))
	return c.Descriptors, nil
}

// ReadCharacteristic reads a characteristic value from a server. [Vol 3, Part G, 4.8.1]
func (p *Client) ReadCharacteristic(c *ble.Characteristic) ([]byte, error) {
	p.log.Debug("read characteristic", slog.String("uuid", c.UUID.String()))
	p.Lock()
	defer p.Unlock()
	val, err := p.ac.Read(c.ValueHandle)
	if err != nil {
		return nil, err
	}

	c.Value = val
	p.log.Debug("read bytes from characteristic", slog.String("uuid", c.UUID.String()), slog.Int("count", len(val)))
	return val, nil
}

// ReadLongCharacteristic reads a characteristic value which is longer than the MTU. [Vol 3, Part G, 4.8.3]
func (p *Client) ReadLongCharacteristic(c *ble.Characteristic) ([]byte, error) {
	p.log.Debug("read long characteristic", slog.String("uuid", c.UUID.String()))
	p.Lock()
	defer p.Unlock()

	// The maximum length of an attribute value shall be 512 octects [Vol 3, 3.2.9]
	buffer := make([]byte, 0, 512)

	read, err := p.ac.Read(c.ValueHandle)
	if err != nil {
		return nil, err
	}
	buffer = append(buffer, read...)

	for len(read) >= p.conn.TxMTU()-1 {
		if read, err = p.ac.ReadBlob(c.ValueHandle, uint16(len(buffer))); err != nil {
			return nil, err
		}
		buffer = append(buffer, read...)
	}

	c.Value = buffer
	p.log.Debug("read bytes from characteristic", slog.String("uuid", c.UUID.String()), slog.Int("count", len(buffer)))
	return buffer, nil
}

// WriteCharacteristic writes a characteristic value to a server. [Vol 3, Part G, 4.9.3]
func (p *Client) WriteCharacteristic(c *ble.Characteristic, v []byte, noRsp bool) error {
	p.log.Debug("write bytes to characteristic", slog.String("uuid", c.UUID.String()), slog.Int("count", len(v)), slog.Bool("noRsp", noRsp))
	p.Lock()
	defer p.Unlock()
	if noRsp {
		return p.ac.WriteCommand(c.ValueHandle, v)
	}
	if err := p.ac.Write(c.ValueHandle, v); err != nil {
		return fmt.Errorf("failed to write to characterstic handle: %w", err)
	}
	p.log.Debug("wrote bytes to characteristic", slog.String("uuid", c.UUID.String()), slog.Int("count", len(v)))
	return nil
}

// ReadDescriptor reads a characteristic descriptor from a server. [Vol 3, Part G, 4.12.1]
func (p *Client) ReadDescriptor(d *ble.Descriptor) ([]byte, error) {
	p.log.Debug("read descriptor", slog.String("uuid", d.UUID.String()))
	p.Lock()
	defer p.Unlock()
	val, err := p.ac.Read(d.Handle)
	if err != nil {
		return nil, err
	}

	d.Value = val
	p.log.Debug("read bytes from descriptor", slog.String("uuid", d.UUID.String()), slog.Int("count", len(val)))
	return val, nil
}

// WriteDescriptor writes a characteristic descriptor to a server. [Vol 3, Part G, 4.12.3]
func (p *Client) WriteDescriptor(d *ble.Descriptor, v []byte) error {
	p.log.Debug("write descriptor", slog.String("uuid", d.UUID.String()))
	p.Lock()
	defer p.Unlock()
	if err := p.ac.Write(d.Handle, v); err != nil {
		return fmt.Errorf("failed to write to descriptor handle: %w", err)
	}
	p.log.Debug("wrote bytes to descriptor", slog.String("uuid", d.UUID.String()), slog.Int("count", len(v)))
	return nil
}

// ReadRSSI retrieves the current RSSI value of remote peripheral. [Vol 2, Part E, 7.5.4]
func (p *Client) ReadRSSI() int {
	p.log.Debug("read RSSI")
	p.Lock()
	defer p.Unlock()
	// TODO:
	return 0
}

// ExchangeMTU informs the server of the client’s maximum receive MTU size and
// request the server to respond with its maximum receive MTU size. [Vol 3, Part F, 3.4.2.1]
func (p *Client) ExchangeMTU(mtu int) (int, error) {
	p.log.Debug("exchange mtu", slog.Int("mtu", mtu))
	p.Lock()
	defer p.Unlock()
	out, err := p.ac.ExchangeMTU(mtu)
	if err != nil {
		return 0, err
	}
	p.log.Debug("exchanged mtu", slog.Int("mtu", mtu))
	return out, nil
}

// Subscribe subscribes to indication (if ind is set true), or notification of a
// characteristic value. [Vol 3, Part G, 4.10 & 4.11]
func (p *Client) Subscribe(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	p.log.Debug("subscribing to characteristic", slog.String("uuid", c.UUID.String()))
	p.Lock()
	defer p.Unlock()

	if c.CCCD == nil {
		return fmt.Errorf("CCCD not found")
	}
	if ind {
		if err := p.setHandlers(c.CCCD.Handle, c.ValueHandle, cccIndicate, h); err != nil {
			return err
		}
	} else if err := p.setHandlers(c.CCCD.Handle, c.ValueHandle, cccNotify, h); err != nil {
		return err
	}
	p.log.Debug("subscribed to characteristic", slog.String("uuid", c.UUID.String()))
	return nil
}

// Unsubscribe unsubscribes to indication (if ind is set true), or notification
// of a specified characteristic value. [Vol 3, Part G, 4.10 & 4.11]
func (p *Client) Unsubscribe(c *ble.Characteristic, ind bool) error {
	p.log.Debug("unsubscribing to characteristic", slog.String("uuid", c.UUID.String()))
	p.Lock()
	defer p.Unlock()
	if c.CCCD == nil {
		return fmt.Errorf("CCCD not found")
	}
	if ind {
		if err := p.setHandlers(c.CCCD.Handle, c.ValueHandle, cccIndicate, nil); err != nil {
			return err
		}
	} else if err := p.setHandlers(c.CCCD.Handle, c.ValueHandle, cccNotify, nil); err != nil {
		return err
	}
	p.log.Debug("unsubscribed to characteristic", slog.String("uuid", c.UUID.String()))
	return nil
}

func (p *Client) setHandlers(cccdh, vh, flag uint16, h ble.NotificationHandler) error {
	s, ok := p.subs[vh]
	if !ok {
		s = &sub{cccdh, 0x0000, nil, nil}
		p.subs[vh] = s
	}
	switch {
	case h == nil && (s.ccc&flag) == 0:
		return nil
	case h != nil && (s.ccc&flag) != 0:
		return nil
	case h == nil && (s.ccc&flag) != 0:
		s.ccc &= ^flag
	case h != nil && (s.ccc&flag) == 0:
		s.ccc |= flag
	}

	v := make([]byte, 2)
	binary.LittleEndian.PutUint16(v, s.ccc)
	if flag == cccNotify {
		s.nHandler = h
	} else {
		s.iHandler = h
	}
	return p.ac.Write(s.cccdh, v)
}

// ClearSubscriptions clears all subscriptions to notifications and indications.
func (p *Client) ClearSubscriptions() error {
	p.log.Debug("clearing subscriptions")
	p.Lock()
	defer p.Unlock()
	zero := make([]byte, 2)
	for vh, s := range p.subs {
		if err := p.ac.Write(s.cccdh, zero); err != nil {
			return err
		}
		delete(p.subs, vh)
	}
	p.log.Debug("cleared subscriptions")
	return nil
}

// CancelConnection disconnects the connection.
func (p *Client) CancelConnection(ctx context.Context) error {
	p.log.Debug("canceling connection")
	defer p.Wait()
	p.Lock()
	defer p.Unlock()
	if err := p.conn.Close(ctx); err != nil {
		return err
	}
	p.log.Debug("cancel connection")
	return nil
}

// Disconnected returns a receiving channel, which is closed when the client disconnects.
func (p *Client) Disconnected() <-chan struct{} {
	p.Lock()
	defer p.Unlock()
	return p.conn.Disconnected()
}

// Conn returns the client's current connection.
func (p *Client) Conn() ble.Conn {
	return p.conn
}

// HandleNotification ...
func (p *Client) HandleNotification(req []byte) {
	p.log.Debug("handling notification")
	defer p.log.Debug("handled notification")
	fn, ok := p.getNotificationHandler(req)
	if !ok {
		// FIXME: disconnects and propagate an error to the user.
		p.log.Debug("received an unregistered notification")
		return
	}
	if fn == nil {
		return
	}
	out := make([]byte, len(req)-3)
	copy(out, req[3:])
	fn(out)
}

func (p *Client) getNotificationHandler(req []byte) (ble.NotificationHandler, bool) {
	p.Lock()
	defer p.Unlock()
	vh := att.HandleValueIndication(req).AttributeHandle()
	sub, ok := p.subs[vh]
	if !ok {
		return nil, false
	}
	fn := sub.nHandler
	if req[0] == att.HandleValueIndicationCode {
		fn = sub.iHandler
	}
	return fn, true
}

type sub struct {
	cccdh    uint16
	ccc      uint16
	nHandler ble.NotificationHandler
	iHandler ble.NotificationHandler
}
