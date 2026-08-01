package collections

import "testing"

// These tests enforce the allocation invariants documented in the README:
// a warmed-up queue's Enqueue/Dequeue pairs are allocation-free, and Map and
// Filter make at most one allocation for the result.

func TestQueueSteadyStateDoesNotAllocate(t *testing.T) {
	q := NewQueue[int](64)
	for i := 0; i < 32; i++ {
		q.Enqueue(i)
	}
	allocs := testing.AllocsPerRun(1000, func() {
		q.Enqueue(1)
		q.Dequeue()
	})
	if allocs != 0 {
		t.Fatalf("steady-state Enqueue/Dequeue = %.1f allocs per pair, want 0", allocs)
	}
}

func TestMapFilterAllocateAtMostOnce(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8}
	double := func(v int) int { return v * 2 }
	even := func(v int) bool { return v%2 == 0 }

	if allocs := testing.AllocsPerRun(100, func() { Map(in, double) }); allocs > 1 {
		t.Fatalf("Map = %.1f allocs, want at most 1", allocs)
	}
	if allocs := testing.AllocsPerRun(100, func() { Filter(in, even) }); allocs > 1 {
		t.Fatalf("Filter = %.1f allocs, want at most 1", allocs)
	}
}
