package collections

import "testing"

func TestQueueZeroValueAndFIFO(t *testing.T) {
	var q Queue[int]
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on an empty queue returned ok=true")
	}

	for _, v := range []int{1, 2, 3} {
		q.Enqueue(v)
	}
	if got, ok := q.Peek(); !ok || got != 1 {
		t.Fatalf("Peek = (%d, %v), want (1, true)", got, ok)
	}
	for _, want := range []int{1, 2, 3} {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Fatalf("Dequeue = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
	if !q.IsEmpty() || q.Len() != 0 {
		t.Fatalf("empty queue: IsEmpty=%v Len=%d", q.IsEmpty(), q.Len())
	}
}

func TestQueueWrapAndGrow(t *testing.T) {
	q := NewQueue[int](8)
	for i := 0; i < 8; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < 4; i++ {
		if got, _ := q.Dequeue(); got != i {
			t.Fatalf("initial Dequeue = %d, want %d", got, i)
		}
	}
	for i := 8; i < 17; i++ {
		q.Enqueue(i)
	}

	if q.Cap() != 16 {
		t.Fatalf("Cap after wrapped growth = %d, want 16", q.Cap())
	}
	for want := 4; want < 17; want++ {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Fatalf("wrapped Dequeue = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
}

func TestQueueClearsRemovedReferences(t *testing.T) {
	type payload struct{ data [64]byte }
	p := &payload{}
	q := NewQueue[*payload](1)
	q.Enqueue(p)
	if got, ok := q.Dequeue(); !ok || got != p {
		t.Fatalf("Dequeue = (%p, %v), want (%p, true)", got, ok, p)
	}
	for i, v := range q.buf {
		if v != nil {
			t.Fatalf("backing slot %d retains a dequeued pointer", i)
		}
	}
}

func TestQueueClearWrapped(t *testing.T) {
	q := NewQueue[*int](8)
	values := make([]int, 12)
	for i := 0; i < 8; i++ {
		q.Enqueue(&values[i])
	}
	for i := 0; i < 5; i++ {
		q.Dequeue()
	}
	for i := 8; i < 12; i++ {
		q.Enqueue(&values[i])
	}
	capacity := q.Cap()
	q.Clear()
	if !q.IsEmpty() || q.Cap() != capacity {
		t.Fatalf("after Clear: IsEmpty=%v Cap=%d, want Cap=%d", q.IsEmpty(), q.Cap(), capacity)
	}
	for i, v := range q.buf {
		if v != nil {
			t.Fatalf("backing slot %d was not cleared", i)
		}
	}
}

func TestNewQueueCapacity(t *testing.T) {
	if got := NewQueue[int](9).Cap(); got != 16 {
		t.Fatalf("NewQueue(9).Cap = %d, want 16", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("NewQueue with negative capacity did not panic")
		}
	}()
	NewQueue[int](-1)
}

func FuzzQueueMatchesSlice(f *testing.F) {
	f.Add([]byte{0, 2, 4, 1, 3, 5})
	f.Add([]byte{0, 0, 0, 0, 1, 1, 0, 1})

	f.Fuzz(func(t *testing.T, operations []byte) {
		var q Queue[byte]
		var reference []byte
		for _, operation := range operations {
			if operation&1 == 0 {
				q.Enqueue(operation)
				reference = append(reference, operation)
				continue
			}

			got, ok := q.Dequeue()
			if len(reference) == 0 {
				if ok {
					t.Fatalf("Dequeue on empty reference = (%d, true)", got)
				}
				continue
			}
			if !ok || got != reference[0] {
				t.Fatalf("Dequeue = (%d, %v), want (%d, true)", got, ok, reference[0])
			}
			reference = reference[1:]
		}

		if q.Len() != len(reference) {
			t.Fatalf("Len = %d, want %d", q.Len(), len(reference))
		}
		for _, want := range reference {
			got, ok := q.Dequeue()
			if !ok || got != want {
				t.Fatalf("drain Dequeue = (%d, %v), want (%d, true)", got, ok, want)
			}
		}
	})
}
