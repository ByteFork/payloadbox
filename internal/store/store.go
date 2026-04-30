// Package store provides a thread-safe fixed-capacity ring buffer.
// When the buffer is full, Add evicts the oldest item. A capacity of 0
// means unbounded growth via append.
package store

import (
	"slices"
	"sync"
)

type Store[T any] struct {
	mu       sync.RWMutex
	capacity int
	data     []T
	head     int
	size     int
}

func NewStore[T any](capacity int) *Store[T] {
	if capacity < 0 {
		capacity = 0
	}

	return &Store[T]{capacity: capacity, data: make([]T, capacity)}
}

func (s *Store[T]) Add(item T) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.capacity == 0 {
		s.data = append(s.data, item)
		s.size++

		return
	}

	if s.size < s.capacity {
		pos := (s.head + s.size) % s.capacity
		s.data[pos] = item
		s.size++
	} else {
		s.head = (s.head + 1) % s.capacity
		pos := (s.head + s.size - 1) % s.capacity
		s.data[pos] = item
	}
}

func (s *Store[T]) List() []T {
	if s == nil {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.size == 0 {
		return make([]T, 0)
	}

	if s.capacity == 0 {
		return slices.Clone(s.data[:s.size])
	}

	out := make([]T, s.size)
	first := min(s.size, s.capacity-s.head)
	copy(out, s.data[s.head:s.head+first])

	if first < s.size {
		copy(out[first:], s.data[:s.size-first])
	}

	return out
}

// Find returns the first item for which match returns true, or the zero value
// of T and false if none match. The scan walks oldest-first.
func (s *Store[T]) Find(match func(T) bool) (T, bool) {
	var zero T
	if s == nil {
		return zero, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.size {
		var item T
		if s.capacity == 0 {
			item = s.data[i]
		} else {
			item = s.data[(s.head+i)%s.capacity]
		}

		if match(item) {
			return item, true
		}
	}

	return zero, false
}

func (s *Store[T]) Len() int {
	if s == nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.size
}

func (s *Store[T]) Clear() {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Zero used slots so pointer-bearing Ts can be GC'd, but keep the
	// backing array so future Adds do not need to reallocate.
	if s.capacity == 0 {
		clear(s.data)
		s.data = s.data[:0]
	} else {
		clear(s.data[:s.size])
	}

	s.head = 0
	s.size = 0
}
