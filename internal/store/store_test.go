package store

import (
	"slices"
	"sync"
	"testing"
)

func TestEmpty(t *testing.T) {
	s := NewStore[int](3)

	if got := s.Len(); got != 0 {
		t.Fatalf("Len on empty = %d, want 0", got)
	}

	got := s.List()
	if got == nil {
		t.Fatalf("List on empty returned nil; want empty non-nil slice")
	}

	if len(got) != 0 {
		t.Fatalf("List on empty = %v, want []", got)
	}
}

func TestAddUpToCapacity(t *testing.T) {
	s := NewStore[int](3)
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if got, want := s.Len(), 3; got != want {
		t.Fatalf("Len = %d, want %d", got, want)
	}

	if got, want := s.List(), []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Fatalf("List = %v, want %v", got, want)
	}
}

func TestOverflowEvictsOldest(t *testing.T) {
	tests := []struct {
		name string
		cap  int
		in   []int
		want []int
	}{
		{"exactly capacity", 3, []int{1, 2, 3}, []int{1, 2, 3}},
		{"overflow by one", 3, []int{1, 2, 3, 4}, []int{2, 3, 4}},
		{"overflow by many", 3, []int{1, 2, 3, 4, 5, 6, 7}, []int{5, 6, 7}},
		{"deep overflow (3×)", 2, []int{1, 2, 3, 4, 5, 6}, []int{5, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewStore[int](tc.cap)
			for _, v := range tc.in {
				s.Add(v)
			}

			if got := s.List(); !slices.Equal(got, tc.want) {
				t.Fatalf("List = %v, want %v", got, tc.want)
			}

			if s.Len() != len(tc.want) {
				t.Fatalf("Len = %d, want %d", s.Len(), len(tc.want))
			}
		})
	}
}

func TestUnboundedCapacity(t *testing.T) {
	s := NewStore[int](0)
	for i := range 1000 {
		s.Add(i)
	}

	if got, want := s.Len(), 1000; got != want {
		t.Fatalf("Len = %d, want %d", got, want)
	}

	got := s.List()
	if len(got) != 1000 || got[0] != 0 || got[999] != 999 {
		t.Fatalf("List[0]=%d List[999]=%d, want 0 and 999", got[0], got[999])
	}
}

func TestClear(t *testing.T) {
	s := NewStore[int](3)
	s.Add(1)
	s.Add(2)
	s.Clear()

	if got := s.Len(); got != 0 {
		t.Fatalf("Len after Clear = %d, want 0", got)
	}

	if got := s.List(); len(got) != 0 {
		t.Fatalf("List after Clear = %v, want []", got)
	}

	// Can still Add after Clear.
	s.Add(10)

	if got, want := s.List(), []int{10}; !slices.Equal(got, want) {
		t.Fatalf("List after Clear+Add = %v, want %v", got, want)
	}
}

func TestWrapAroundThenClearThenRefill(t *testing.T) {
	s := NewStore[int](3)
	for _, v := range []int{1, 2, 3, 4, 5} {
		s.Add(v)
	}

	s.Clear()

	for _, v := range []int{10, 11} {
		s.Add(v)
	}

	if got, want := s.List(), []int{10, 11}; !slices.Equal(got, want) {
		t.Fatalf("List = %v, want %v", got, want)
	}
}

func TestNilReceiver(t *testing.T) {
	var s *Store[int] // nil pointer

	// None of these should panic.
	s.Add(1)

	if got := s.Len(); got != 0 {
		t.Fatalf("nil Len = %d, want 0", got)
	}

	if got := s.List(); got != nil {
		t.Fatalf("nil List = %v, want nil", got)
	}

	s.Clear()
}

func TestConcurrentAdds(t *testing.T) {
	const (
		writers   = 8
		perWriter = 500
		capacity  = 64
	)

	s := NewStore[int](capacity)

	var wg sync.WaitGroup
	for range writers {
		wg.Go(func() {
			for i := range perWriter {
				s.Add(i)
			}
		})
	}

	wg.Wait()

	if got := s.Len(); got != capacity {
		t.Fatalf("Len after concurrent writes = %d, want %d", got, capacity)
	}

	if got := len(s.List()); got != capacity {
		t.Fatalf("len(List()) = %d, want %d", got, capacity)
	}
}

func TestConcurrentReadersAndWriters(_ *testing.T) {
	s := NewStore[int](32)
	done := make(chan struct{})

	go func() {
		for i := range 10_000 {
			s.Add(i)
		}

		close(done)
	}()

	for range 4 {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					_ = s.Len()
					_ = s.List()
				}
			}
		}()
	}

	<-done
}
