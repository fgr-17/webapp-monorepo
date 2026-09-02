// Package history keeps an in-memory ring of recent successful calculator operations.
package history

import (
	"sync"
	"time"
)

// DefaultCap is the maximum number of entries retained (oldest dropped first).
const DefaultCap = 10

// Entry is one successful calculator operation.
type Entry struct {
	Op        string    `json:"op"`
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Result    float64   `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

// Store is a thread-safe, capped history. List() returns newest first.
type Store struct {
	mu      sync.Mutex
	cap     int
	entries []Entry
}

// New returns a store that keeps at most cap entries. Cap < 1 uses DefaultCap.
func New(cap int) *Store {
	if cap < 1 {
		cap = DefaultCap
	}
	return &Store{cap: cap, entries: make([]Entry, 0, cap)}
}

// Add appends a successful operation. Evicts the oldest entry when over capacity.
func (s *Store) Add(op string, a, b, result float64, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if at.IsZero() {
		at = time.Now().UTC()
	} else {
		at = at.UTC()
	}

	s.entries = append(s.entries, Entry{
		Op:        op,
		A:         a,
		B:         b,
		Result:    result,
		Timestamp: at,
	})
	if len(s.entries) > s.cap {
		// Drop oldest (front of slice).
		s.entries = append([]Entry(nil), s.entries[len(s.entries)-s.cap:]...)
	}
}

// List returns a copy of entries, newest first.
func (s *Store) List() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := len(s.entries)
	out := make([]Entry, n)
	for i := range s.entries {
		out[i] = s.entries[n-1-i]
	}
	return out
}

// Len returns the current number of entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}
