package rfcomm

import (
	"errors"
	"sync"
)

// MaxServerChannels The server channel range is from 1 to 30
const maxServerChannels int = 30

type serverChannels struct {
	sync.RWMutex

	channels map[int]*Client
}

func newServerChannels() *serverChannels {
	return &serverChannels{
		channels: make(map[int]*Client, maxServerChannels)}
}

func (s *serverChannels) Add(c *Client) (int, error) {
	s.RWLock()
	defer s.RWUnlock()
	for i := 1; i <= MaxServerChannels; i++ {
		if _, ok := s.channels[i]; !ok {
			s.channels[i] = c
			return i, nil
		}
	}
	return 0, errors.New("There is no more room to add another server channel")
}

func (s *serverChannels) Remove(serverChannel int) {
	s.RWLock()
	defer s.RWLock()
	delete(s.channels, serverChannel)
}
