package collections

import "testing"

func TestSetBasics(t *testing.T) {
	s := NewSet("a", "b", "a")
	if s.Len() != 2 || !s.Has("a") || !s.Has("b") {
		t.Fatalf("NewSet = %v, want {a,b}", s)
	}
	s.Remove("a")
	if s.Has("a") || s.Len() != 1 {
		t.Fatalf("after Remove: Has(a)=%v Len=%d", s.Has("a"), s.Len())
	}
	s.Clear()
	if !s.IsEmpty() {
		t.Fatalf("set is not empty after Clear: %v", s)
	}
}

func TestSetOperations(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(2, 3, 4)

	assertSet(t, a.Union(b), 1, 2, 3, 4)
	assertSet(t, a.Intersect(b), 2, 3)
	assertSet(t, a.Difference(b), 1)

	clone := a.Clone()
	clone.Add(9)
	if a.Has(9) {
		t.Fatal("Clone shares map storage with the source")
	}
	assertSet(t, a, 1, 2, 3)
}

func TestSetIntersectWalksSmallerInput(t *testing.T) {
	large := NewSetWithCapacity[int](10_000)
	for i := 0; i < 10_000; i++ {
		large.Add(i)
	}
	small := NewSet(1, 9_999, 20_000)
	assertSet(t, large.Intersect(small), 1, 9_999)
	assertSet(t, small.Intersect(large), 1, 9_999)
}

func TestNilSetReadsAndRemoveAreSafeButAddPanics(t *testing.T) {
	var s Set[int]
	if s.Has(1) || !s.IsEmpty() {
		t.Fatal("nil Set has unexpected values")
	}
	s.Remove(1)
	s.Clear()
	if s.Clone() != nil {
		t.Fatal("Clone of a nil Set did not preserve nilness")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("Add on a nil Set did not panic")
		}
	}()
	s.Add(1)
}

func assertSet[T comparable](t *testing.T, got Set[T], want ...T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("set length = %d, want %d; set=%v", len(got), len(want), got)
	}
	for _, v := range want {
		if !got.Has(v) {
			t.Fatalf("set %v does not contain %v", got, v)
		}
	}
}
