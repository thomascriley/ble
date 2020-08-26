package gatt

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/log"
	"github.com/thomascriley/ble/linux/att"
)

const (
	cccNotify   = 0x0001
	cccIndicate = 0x0002
)

// NewClient returns a GATT Client.
func NewClient(conn ble.Conn) (*Client, error) {
	p := &Client{
		subs: make(map[uint16]*sub),
		conn: conn,
		addr: conn.LocalAddr().String(),
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
	
	addr string
}

func (p *Client) String() string { return p.addr }

func (p *Client) RLock() {
	log.Printf("BLE GATT: %s: read lock", p)
	p.RWMutex.RLock()
}

func (p *Client) RUnlock() {
	log.Printf("BLE GATT: %s: read unlock", p)
	p.RWMutex.RUnlock()
}

func (p *Client) Lock() {
	log.Printf("BLE GATT: %s: lock", p)
	p.RWMutex.Lock()
}

func (p *Client) Unlock() {
	log.Printf("BLE GATT: %s: unlocking", p)
	p.RWMutex.Unlock()
}
// Addr returns the address of the client.
func (p *Client) Address() ble.Addr {
	log.Printf("BLE GATT: %s: getting address", p)
	defer log.Printf("BLE GATT: %s: got address", p)
	p.RLock()
	defer p.RUnlock()
	return p.conn.RemoteAddr()
}

// Name returns the name of the client.
func (p *Client) Name() string {
	log.Printf("BLE GATT: %s: getting name", p)
	defer log.Printf("BLE GATT: %s: got name", p)
	p.RLock()
	defer p.RUnlock()
	return p.name
}

func (p *Client) Connection() ble.Conn {
	return p.conn
}

// Profile returns the discovered profile.
func (p *Client) Profile() *ble.Profile {
	log.Printf("BLE GATT: %s: getting profile", p)
	defer log.Printf("BLE GATT: %s: got profile", p)
	p.RLock()
	defer p.RUnlock()
	return p.profile
}

// DiscoverProfile discovers the whole hierarchy of a server.
func (p *Client) DiscoverProfile(force bool) (*ble.Profile, error) {
	log.Printf("BLE GATT: %s: discovering profile", p)
	if p.profile != nil && !force {
		return p.profile, nil
	}
	ss, err := p.DiscoverServices(nil)
	if err != nil {
		return nil, fmt.Errorf("can't discover services: %s", err)
	}
	for _, s := range ss {
		cs, err := p.DiscoverCharacteristics(nil, s)
		if err != nil {
			return nil, fmt.Errorf("can't discover characteristics: %s", err)
		}
		for _, c := range cs {
			_, err := p.DiscoverDescriptors(nil, c)
			if err != nil {
				return nil, fmt.Errorf("can't discover descriptors: %s", err)
			}
		}
	}
	p.profile = &ble.Profile{Services: ss}
	log.Printf("BLE GATT: %s: discovered profile", p)
	return p.profile, nil
}

// DiscoverServices finds all the primary services on a server. [Vol 3, Part G, 4.4.1]
// If filter is specified, only filtered services are returned.
func (p *Client) DiscoverServices(filter []ble.UUID) ([]*ble.Service, error) {
	log.Printf("BLE GATT: %s: discovering services", p)
	p.Lock()
	defer p.Unlock()
	if p.profile == nil {
		p.profile = &ble.Profile{}
	}
	start := uint16(0x0001)
	for {
		length, b, err := p.ac.ReadByGroupType(start, 0xFFFF, ble.PrimaryServiceUUID)
		if err == ble.ErrAttrNotFound {
			log.Printf("BLE GATT: %s: discovered %d services", p, len(p.profile.Services))
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
				log.Printf("BLE GATT: %s: discovered %d services", p, len(p.profile.Services))
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
	log.Printf("BLE GATT: %s: discovering included services", p)
	p.Lock()
	defer p.Unlock()
	return nil, nil
}

// DiscoverCharacteristics finds all the characteristics within a service. [Vol 3, Part G, 4.6.1]
// If filter is specified, only filtered characteristics are returned.
func (p *Client) DiscoverCharacteristics(filter []ble.UUID, s *ble.Service) ([]*ble.Characteristic, error) {
	log.Printf("BLE GATT: %s: %s: discovering characteristics in service", p, s.UUID)
	p.Lock()
	defer p.Unlock()
	start := s.Handle
	var lastChar *ble.Characteristic
	for start <= s.EndHandle {
		length, b, err := p.ac.ReadByType(start, s.EndHandle, ble.CharacteristicUUID)
		if err == ble.ErrAttrNotFound {
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
	log.Printf("BLE GATT: %s: %s: discovered %d characteristics in service", p, s.UUID, len(s.Characteristics))
	return s.Characteristics, nil
}

// DiscoverDescriptors finds all the descriptors within a characteristic. [Vol 3, Part G, 4.7.1]
// If filter is specified, only filtered descriptors are returned.
func (p *Client) DiscoverDescriptors(filter []ble.UUID, c *ble.Characteristic) ([]*ble.Descriptor, error) {
	log.Printf("BLE GATT: %s: %s: discovering descriptors", p, c.UUID)
	p.Lock()
	defer p.Unlock()
	start := c.ValueHandle + 1
	for start <= c.EndHandle {
		found, b, err := p.ac.FindInformation(start, c.EndHandle)
		if err == ble.ErrAttrNotFound {
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
	log.Printf("BLE GATT: %s: %s: discovered %d descriptors in characteristic", p, c.UUID, len(c.Descriptors))
	return c.Descriptors, nil
}

// ReadCharacteristic reads a characteristic value from a server. [Vol 3, Part G, 4.8.1]
func (p *Client) ReadCharacteristic(c *ble.Characteristic) ([]byte, error) {
	log.Printf("BLE GATT: %s: %s: read characteristic", p, c.UUID)
	p.Lock()
	defer p.Unlock()
	val, err := p.ac.Read(c.ValueHandle)
	if err != nil {
		return nil, err
	}

	c.Value = val
	log.Printf("BLE GATT: %s: %s: read %d bytes from characteristic", p, c.UUID, len(val))
	return val, nil
}

// ReadLongCharacteristic reads a characteristic value which is longer than the MTU. [Vol 3, Part G, 4.8.3]
func (p *Client) ReadLongCharacteristic(c *ble.Characteristic) ([]byte, error) {
	log.Printf("BLE GATT: %s: read long characteristic", c.UUID)
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
	log.Printf("BLE GATT: %s: %s: read %d bytes from characteristic", p, c.UUID, len(buffer))
	return buffer, nil
}

// WriteCharacteristic writes a characteristic value to a server. [Vol 3, Part G, 4.9.3]
func (p *Client) WriteCharacteristic(c *ble.Characteristic, v []byte, noRsp bool) error {
	log.Printf("BLE GATT: %s: write %d bytes to characteristic (noRsp: %v)", c.UUID, len(v), noRsp)
	p.Lock()
	defer p.Unlock()
	if noRsp {
		return p.ac.WriteCommand(c.ValueHandle, v)
	}
	if err := p.ac.Write(c.ValueHandle, v); err != nil {
		return fmt.Errorf("failed to write to characterstic handle: %w", err)
	}
	log.Printf("BLE GATT: %s: %s: wrote %d bytes to characteristic", p, c.UUID, len(v))
	return nil
}

// ReadDescriptor reads a characteristic descriptor from a server. [Vol 3, Part G, 4.12.1]
func (p *Client) ReadDescriptor(d *ble.Descriptor) ([]byte, error) {
	log.Printf("BLE GATT: %s: read descriptor", d.UUID)
	p.Lock()
	defer p.Unlock()
	val, err := p.ac.Read(d.Handle)
	if err != nil {
		return nil, err
	}

	d.Value = val
	log.Printf("BLE GATT: %s: %s: read %d bytes from descriptor", p, d.UUID, len(val))
	return val, nil
}

// WriteDescriptor writes a characteristic descriptor to a server. [Vol 3, Part G, 4.12.3]
func (p *Client) WriteDescriptor(d *ble.Descriptor, v []byte) error {
	log.Printf("BLE GATT: %s: write descriptor", d.UUID)
	p.Lock()
	defer p.Unlock()
	if err := p.ac.Write(d.Handle, v); err != nil {
		return fmt.Errorf("failed to write to descriptor handle: %w", err)
	}
	log.Printf("BLE GATT: %s: %s: wrote %d bytes to descriptor", p, d.UUID, len(v))
	return nil
}

// ReadRSSI retrieves the current RSSI value of remote peripheral. [Vol 2, Part E, 7.5.4]
func (p *Client) ReadRSSI() int {
	log.Printf("BLE GATT: read RSSI", p)
	p.Lock()
	defer p.Unlock()
	// TODO:
	return 0
}

// ExchangeMTU informs the server of the client’s maximum receive MTU size and
// request the server to respond with its maximum receive MTU size. [Vol 3, Part F, 3.4.2.1]
func (p *Client) ExchangeMTU(mtu int) (int, error) {
	log.Printf("BLE GATT: %s: exchange mtu %d",p, mtu)
	p.Lock()
	defer p.Unlock()
	out, err := p.ac.ExchangeMTU(mtu)
	if err != nil {
		return 0, err
	}
	log.Printf("BLE GATT: %s: exchanged mtu of %d", p, mtu)
	return out, nil
}

// Subscribe subscribes to indication (if ind is set true), or notification of a
// characteristic value. [Vol 3, Part G, 4.10 & 4.11]
func (p *Client) Subscribe(c *ble.Characteristic, ind bool, h ble.NotificationHandler) error {
	log.Printf("BLE GATT: %s: subscribing to characteristic",p, c.UUID)
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
	log.Printf("BLE GATT: %s: %s: subscribed to characteristic", p, c.UUID)
	return nil
}

// Unsubscribe unsubscribes to indication (if ind is set true), or notification
// of a specified characteristic value. [Vol 3, Part G, 4.10 & 4.11]
func (p *Client) Unsubscribe(c *ble.Characteristic, ind bool) error {
	log.Printf("BLE GATT: %s: unsubscribing to characteristic",p, c.UUID)
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
	log.Printf("BLE GATT: %s: %s: unsubscribed to characteristic", p, c.UUID)
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
	log.Printf("BLE GATT: %s: clearing subscriptions", p)
	p.Lock()
	defer p.Unlock()
	zero := make([]byte, 2)
	for vh, s := range p.subs {
		if err := p.ac.Write(s.cccdh, zero); err != nil {
			return err
		}
		delete(p.subs, vh)
	}
	log.Printf("BLE GATT: %s: cleared subscriptions", p)
	return nil
}

// CancelConnection disconnects the connection.
func (p *Client) CancelConnection(ctx context.Context) error {
	log.Printf("BLE GATT: %s: canceling connection", p)
	defer p.Wait()
	p.Lock()
	defer p.Unlock()
	if err := p.conn.Close(ctx); err != nil {
		return err
	}
	log.Printf("BLE GATT: %s: cancel connection", p)
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
	log.Printf("BLE GATT: %s: handling notification", p)
	defer log.Printf("BLE GATT: %s: handled notification", p)
	fn, ok := p.getNotificationHandler(req)
	if !ok {
		// FIXME: disconnects and propagate an error to the user.
		log.Printf("BLE GATT: %s: received an unregistered notification", p)
		return
	}
	if fn == nil {
		return
	}
	out := make([]byte,len(req)-3)
	copy(out, req[3:])
	fn(out)
}

func (p *Client) getNotificationHandler(req []byte) (ble.NotificationHandler, bool){
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
