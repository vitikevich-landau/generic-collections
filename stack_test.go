package collections

import "testing"

func TestStackZeroValueAndLIFO(t *testing.T) {
	var s Stack[int]
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on an empty stack returned ok=true")
	}
	if got, ok := s.Peek(); ok || got != 0 {
		t.Fatalf("Peek on empty stack = (%d, %v), want (0, false)", got, ok)
	}

	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}
	if got, ok := s.Peek(); !ok || got != 3 {
		t.Fatalf("Peek = (%d, %v), want (3, true)", got, ok)
	}
	for _, want := range []int{3, 2, 1} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Fatalf("Pop = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatalf("empty stack: IsEmpty=%v Len=%d", s.IsEmpty(), s.Len())
	}
}

func TestStackClearsRemovedReferences(t *testing.T) {
	type payload struct{ data [64]byte }
	p := &payload{}

	var s Stack[*payload]
	s.Push(p)
	if got, ok := s.Pop(); !ok || got != p {
		t.Fatalf("Pop = (%p, %v), want (%p, true)", got, ok, p)
	}
	if retained := s.items[:cap(s.items)][0]; retained != nil {
		t.Fatal("popped pointer is still retained in the backing array")
	}
}

func TestStackCapacityAndClear(t *testing.T) {
	s := NewStack[*int](16)
	if s.Cap() != 16 {
		t.Fatalf("Cap = %d, want 16", s.Cap())
	}
	a, b := 1, 2
	s.Push(&a)
	s.Push(&b)
	s.Clear()
	if !s.IsEmpty() || s.Cap() != 16 {
		t.Fatalf("after Clear: IsEmpty=%v Cap=%d", s.IsEmpty(), s.Cap())
	}
	for i, v := range s.items[:cap(s.items)] {
		if v != nil {
			t.Fatalf("backing slot %d was not cleared", i)
		}
	}
}
