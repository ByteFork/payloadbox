// Package hub provides a generic in-memory pub/sub fan-out for typed events.
// Subscribers receive events on a buffered channel. A subscriber whose buffer
// is full will drop events rather than block the publisher.
package hub

import "sync"

// subscriberBufferSize is the per-subscriber channel capacity.
const subscriberBufferSize = 32

type Hub[T any] struct {
	mu   sync.Mutex
	subs map[chan T]struct{}
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{subs: make(map[chan T]struct{})}
}

func (h *Hub[T]) Subscribe() chan T {
	ch := make(chan T, subscriberBufferSize)

	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()

	return ch
}

func (h *Hub[T]) Unsubscribe(ch chan T) {
	h.mu.Lock()
	delete(h.subs, ch)
	close(ch)
	h.mu.Unlock()
}

// Publish delivers event to every current subscriber. A subscriber with a
// full buffer has the event dropped rather than blocking Publish.
//
// The mutex is held across the fan-out so concurrent Unsubscribe calls cannot
// close a channel mid-send (which would panic). Since each send is
// non-blocking (select with default), the hold time remains bounded by the
// current subscriber count, not by any subscriber's receive speed.
func (h *Hub[T]) Publish(event T) {
	h.mu.Lock()
	for ch := range h.subs {
		select {
		case ch <- event:
		default:
		}
	}
	h.mu.Unlock()
}
