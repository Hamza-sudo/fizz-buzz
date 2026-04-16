package stats

import (
	"sync"

	"fizz-buzz/internal/service"
)

// Entry represents a counted request and its hit count.
type Entry struct {
	Params service.FizzBuzzParams `json:"params"`
	Hits   int                    `json:"hits"`
}

// Store keeps FizzBuzz request statistics in memory.
type Store struct {
	mu     sync.RWMutex
	counts map[service.FizzBuzzParams]int
	top    Entry
	hasTop bool
}

// NewStore creates a ready-to-use statistics store.
func NewStore() *Store {
	return &Store{
		counts: make(map[service.FizzBuzzParams]int),
	}
}

// Record increments the hit counter for a request.
func (s *Store) Record(params service.FizzBuzzParams) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hits := s.counts[params] + 1
	s.counts[params] = hits

	if !s.hasTop || hits > s.top.Hits {
		s.top = Entry{
			Params: params,
			Hits:   hits,
		}
		s.hasTop = true
	}
}

// Top returns the most frequently requested parameters.
func (s *Store) Top() (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.hasTop {
		return Entry{}, false
	}

	return s.top, true
}
