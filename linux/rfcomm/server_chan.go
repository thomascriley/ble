package rfcomm

import (
	"errors"
	"sync"
)

// MaxServerChannels The server channel range is from 1 to 30
const maxServerChannels uint8 = 30

type serverChannels struct {
	sync.RWMutex

	channels map[uint8]*Client
}

func newServerChannels() *serverChannels {
	return &serverChannels{
		channels: make(map[uint8]*Client, maxServerChannels)}
}

func (s *serverChannels) Add(c *Client) (uint8, error) {
	s.Lock()
	defer s.Unlock()
	var i uint8
	for i = 1; i <= maxServerChannels; i++ {
		if _, ok := s.channels[i]; !ok {
			s.channels[i] = c
			return i, nil
		}
	}
	return 0, errors.New("There is no more room to add another server channel")
}

func (s *serverChannels) Remove(serverChannel uint8) {
	s.Lock()
	defer s.Unlock()
	delete(s.channels, serverChannel)
}
