package ble

import (
	"log/slog"
)

// NewService creates and initialize a new Service using u as it's UUID.
func NewService(u UUID) *Service {
	return &Service{UUID: u}
}

// NewDescriptor creates and returns a Descriptor.
func NewDescriptor(u UUID) *Descriptor {
	return &Descriptor{UUID: u}
}

// NewCharacteristic creates and returns a Characteristic.
func NewCharacteristic(u UUID) *Characteristic {
	return &Characteristic{UUID: u}
}

// Property ...
type Property int

// Characteristic property flags (spec 3.3.3.1)
const (
	CharBroadcast   Property = 0x01 // may be brocasted
	CharRead        Property = 0x02 // may be read
	CharWriteNR     Property = 0x04 // may be written to, with no reply
	CharWrite       Property = 0x08 // may be written to, with a reply
	CharNotify      Property = 0x10 // supports notifications
	CharIndicate    Property = 0x20 // supports Indications
	CharSignedWrite Property = 0x40 // supports signed write
	CharExtended    Property = 0x80 // supports extended properties
)

// A Profile is composed of one or more services necessary to fulfill a use case.
type Profile struct {
	Services []*Service `json:"services"`
}

// Find searches discovered profile for the specified target's type and UUID.
// The target must has the type of *Service, *Characteristic, or *Descriptor.
func (p *Profile) Find(target interface{}) interface{} {
	switch t := target.(type) {
	case *Service:
		return p.FindService(t)
	case *Characteristic:
		return p.FindCharacteristic(t)
	case *Descriptor:
		return p.FindDescriptor(t)
	default:
		return nil
	}
}

// FindService searches discoverd profile for the specified service and UUID
func (p *Profile) FindService(service *Service) *Service {
	for _, s := range p.Services {
		if s.UUID.Equal(service.UUID) {
			return s
		}
	}
	return nil
}

// FindCharacteristic searches discoverd profile for the specified characteristic and UUID
func (p *Profile) FindCharacteristic(char *Characteristic) *Characteristic {
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			if c.UUID.Equal(char.UUID) {
				return c
			}
		}
	}
	return nil
}

// FindDescriptor searches discoverd profile for the specified descriptor and UUID
func (p *Profile) FindDescriptor(desc *Descriptor) *Descriptor {
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			for _, d := range c.Descriptors {
				if d.UUID.Equal(desc.UUID) {
					return d
				}
			}
		}
	}
	return nil
}

// A Service is a BLE service.
type Service struct {
	UUID            UUID              `json:"uuid"`
	Characteristics []*Characteristic `json:"characteristics"`

	Handle    uint16 `json:"handle"`
	EndHandle uint16 `json:"endHandle"`

	log *slog.Logger
}

// AddCharacteristic adds a characteristic to a service.
// AddCharacteristic panics if the service already contains another characteristic with the same UUID.
func (s *Service) AddCharacteristic(c *Characteristic) *Characteristic {
	for _, x := range s.Characteristics {
		if x.UUID.Equal(c.UUID) {
			panic("service already contains a characteristic with UUID " + c.UUID.String())
		}
	}
	c.log = s.log.With(slog.String("characteristic", c.UUID.String()))
	s.Characteristics = append(s.Characteristics, c)
	return c
}

// NewCharacteristic adds a characteristic to a service.
// NewCharacteristic panics if the service already contains another characteristic with the same UUID.
func (s *Service) NewCharacteristic(u UUID) *Characteristic {
	return s.AddCharacteristic(&Characteristic{
		UUID: u,
		log:  s.log.With(slog.String("characteristic", u.String())),
	})
}

// A Characteristic is a BLE characteristic.
type Characteristic struct {
	UUID        UUID          `json:"uuid"`
	Property    Property      `json:"property"`
	Secure      Property      `json:"secure"`
	Descriptors []*Descriptor `json:"descriptors"`
	CCCD        *Descriptor   `json:"cccd"`

	Value []byte `json:"value"`

	ReadHandler     ReadHandler   `json:"-"`
	WriteHandler    WriteHandler  `json:"-"`
	NotifyHandler   NotifyHandler `json:"-"`
	IndicateHandler NotifyHandler `json:"-"`

	Handle      uint16 `json:"handle"`
	ValueHandle uint16 `json:"valueHandle"`
	EndHandle   uint16 `json:"endHandle"`

	log *slog.Logger
}

// AddDescriptor adds a descriptor to a characteristic.
// AddDescriptor panics if the characteristic already contains another descriptor with the same UUID.
func (c *Characteristic) AddDescriptor(d *Descriptor) *Descriptor {
	for _, x := range c.Descriptors {
		c.log.Debug("existing descriptor", slog.String("descriptor", d.UUID.String()))
		if x.UUID.Equal(d.UUID) {
			panic("characteristic already contains a descriptor with UUID " + d.UUID.String())
		}
	}
	c.log.Debug("adding descriptor", slog.String("descriptor", d.UUID.String()))
	d.log = c.log.With(slog.String("descriptor", d.UUID.String()))
	c.Descriptors = append(c.Descriptors, d)
	return d
}

// NewDescriptor adds a descriptor to a characteristic.
// NewDescriptor panics if the characteristic already contains another descriptor with the same UUID.
func (c *Characteristic) NewDescriptor(u UUID) *Descriptor {
	return c.AddDescriptor(&Descriptor{
		UUID: u,
		log:  c.log.With(slog.String("descriptor", u.String())),
	})
}

// SetValue makes the characteristic support read requests, and returns a static value.
// SetValue must be called before the containing service is added to a server.
// SetValue panics if the characteristic has been configured with a ReadHandler.
func (c *Characteristic) SetValue(b []byte) {
	if c.ReadHandler != nil {
		panic("characteristic has been configured with a read handler")
	}
	c.Property |= CharRead
	c.Value = make([]byte, len(b))
	copy(c.Value, b)
}

// HandleRead makes the characteristic support read requests, and routes read requests to h.
// HandleRead must be called before the containing service is added to a server.
// HandleRead panics if the characteristic has been configured with a static value.
func (c *Characteristic) HandleRead(h ReadHandler) {
	if c.Value != nil {
		panic("characteristic has been configured with a static value")
	}
	c.Property |= CharRead
	c.ReadHandler = h
}

// HandleWrite makes the characteristic support write and write-no-response requests, and routes write requests to h.
// The WriteHandler does not differentiate between write and write-no-response requests; it is handled automatically.
// HandleWrite must be called before the containing service is added to a server.
func (c *Characteristic) HandleWrite(h WriteHandler) {
	c.Property |= CharWrite | CharWriteNR
	c.WriteHandler = h
}

// HandleNotify makes the characteristic support notify requests, and routes notification requests to h.
// HandleNotify must be called before the containing service is added to a server.
func (c *Characteristic) HandleNotify(h NotifyHandler) {
	c.Property |= CharNotify
	c.NotifyHandler = h
}

// HandleIndicate makes the characteristic support indicate requests, and routes notification requests to h.
// HandleIndicate must be called before the containing service is added to a server.
func (c *Characteristic) HandleIndicate(h NotifyHandler) {
	c.Property |= CharIndicate
	c.IndicateHandler = h
}

// Descriptor is a BLE descriptor
type Descriptor struct {
	UUID     UUID     `json:"uuid"`
	Property Property `json:"property"`

	Handle uint16 `json:"handle"`
	Value  []byte `json:"value"`

	ReadHandler  ReadHandler  `json:"-"`
	WriteHandler WriteHandler `json:"-"`

	log *slog.Logger
}

// SetValue makes the descriptor support read requests, and returns a static value.
// SetValue must be called before the containing service is added to a server.
// SetValue panics if the descriptor has already configured with a ReadHandler.
func (d *Descriptor) SetValue(b []byte) {
	if d.ReadHandler != nil {
		panic("descriptor has been configured with a read handler")
	}
	d.Property |= CharRead
	d.Value = make([]byte, len(b))
	copy(d.Value, b)
}

// HandleRead makes the descriptor support read requests, and routes read requests to h.
// HandleRead must be called before the containing service is added to a server.
// HandleRead panics if the descriptor has been configured with a static value.
func (d *Descriptor) HandleRead(h ReadHandler) {
	if d.Value != nil {
		panic("descriptor has been configured with a static value")
	}
	d.Property |= CharRead
	d.ReadHandler = h
}

// HandleWrite makes the descriptor support write and write-no-response requests, and routes write requests to h.
// The WriteHandler does not differentiate between write and write-no-response requests; it is handled automatically.
// HandleWrite must be called before the containing service is added to a server.
func (d *Descriptor) HandleWrite(h WriteHandler) {
	d.Property |= CharWrite | CharWriteNR
	d.WriteHandler = h
}
