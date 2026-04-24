package hub

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingleSubscriberReceives(t *testing.T) {
	h := NewHub[int]()

	ch := h.Subscribe()
	defer h.Unsubscribe(ch)

	h.Publish(42)

	select {
	case got := <-ch:
		if got != 42 {
			t.Fatalf("received %d, want 42", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("subscriber did not receive within 100ms")
	}
}

func TestAllSubscribersReceive(t *testing.T) {
	h := NewHub[string]()

	subs := make([]chan string, 3)
	for i := range subs {
		subs[i] = h.Subscribe()
	}

	defer func() {
		for _, c := range subs {
			h.Unsubscribe(c)
		}
	}()

	h.Publish("hi")

	for i, c := range subs {
		select {
		case got := <-c:
			if got != "hi" {
				t.Errorf("sub %d got %q, want %q", i, got, "hi")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("sub %d did not receive", i)
		}
	}
}

func TestSlowSubscriberDropsNotBlocks(t *testing.T) {
	h := NewHub[int]()

	slow := h.Subscribe()

	fast := h.Subscribe()
	defer h.Unsubscribe(slow)
	defer h.Unsubscribe(fast)

	for i := range 1000 {
		h.Publish(i)
	}

	select {
	case <-fast:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("fast subscriber starved by slow one")
	}
}

func TestUnsubscribeClosesChannel(t *testing.T) {
	h := NewHub[int]()
	ch := h.Subscribe()

	h.Unsubscribe(ch)

	_, ok := <-ch
	if ok {
		t.Fatal("channel still open after Unsubscribe")
	}
}

func TestPublishWithNoSubscribersIsNoOp(_ *testing.T) {
	h := NewHub[int]()
	h.Publish(1)
	h.Publish(2)
}

func TestPublishAfterUnsubscribe(t *testing.T) {
	h := NewHub[int]()
	a := h.Subscribe()
	b := h.Subscribe()

	h.Unsubscribe(a)
	defer h.Unsubscribe(b)

	h.Publish(7)

	select {
	case got := <-b:
		if got != 7 {
			t.Fatalf("got %d, want 7", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("remaining subscriber did not receive")
	}
}

func TestConcurrentSubUnsubPublish(t *testing.T) {
	h := NewHub[int]()

	var (
		publishes int64
		wg        sync.WaitGroup
	)

	for range 4 {
		wg.Go(func() {
			for i := range 2_000 {
				h.Publish(i)
				atomic.AddInt64(&publishes, 1)
			}
		})
	}

	for range 4 {
		wg.Go(func() {
			for range 200 {
				ch := h.Subscribe()

				for range 5 {
					select {
					case <-ch:
					case <-time.After(5 * time.Millisecond):
					}
				}

				h.Unsubscribe(ch)
			}
		})
	}

	wg.Wait()

	if publishes == 0 {
		t.Fatal("no publishes observed")
	}
}
