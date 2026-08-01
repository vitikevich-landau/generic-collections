package collections

import (
	"reflect"
	"testing"
)

func TestMapFilterReduce(t *testing.T) {
	numbers := []int{1, 2, 3, 4}

	if got, want := Map(numbers, func(v int) int { return v * v }), []int{1, 4, 9, 16}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Map = %v, want %v", got, want)
	}
	if got, want := Filter(numbers, func(v int) bool { return v%2 == 0 }), []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter = %v, want %v", got, want)
	}
	if got := Reduce(numbers, 0, func(acc, v int) int { return acc + v }); got != 10 {
		t.Fatalf("Reduce = %d, want 10", got)
	}
	if Map[int, int](nil, func(v int) int { return v }) != nil {
		t.Fatal("Map(nil) did not preserve nilness")
	}
	if Filter[int](nil, func(int) bool { return true }) != nil {
		t.Fatal("Filter(nil) did not preserve nilness")
	}
}

func TestAppendMapReusesDestination(t *testing.T) {
	dst := make([]int, 1, 4)
	dst[0] = 10
	first := &dst[0]
	got := AppendMap(dst, []int{1, 2, 3}, func(v int) int { return v * 2 })
	if want := []int{10, 2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("AppendMap = %v, want %v", got, want)
	}
	if &got[0] != first {
		t.Fatal("AppendMap allocated despite sufficient destination capacity")
	}
}

func TestAppendFilterInPlace(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	first := &in[0]
	got := AppendFilter(in[:0], in, func(v int) bool { return v%2 == 0 })
	if want := []int{2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("AppendFilter = %v, want %v", got, want)
	}
	if &got[0] != first {
		t.Fatal("AppendFilter did not reuse in-place storage")
	}
}

type age int

func TestSumAllNumericFamilies(t *testing.T) {
	if got := Sum([]age{10, 20, 30}); got != 60 {
		t.Fatalf("Sum(age) = %d, want 60", got)
	}
	if got := Sum([]uintptr{1, 2, 3}); got != 6 {
		t.Fatalf("Sum(uintptr) = %d, want 6", got)
	}
	if got := Sum([]complex128{1 + 2i, 3 + 4i}); got != 4+6i {
		t.Fatalf("Sum(complex128) = %v, want 4+6i", got)
	}
}

func TestCompatibilityHelpers(t *testing.T) {
	if got := Index([]string{"a", "b", "c"}, "b"); got != 1 {
		t.Fatalf("Index = %d, want 1", got)
	}
	m := map[string]int{"b": 2, "a": 1}
	if got, want := SortedKeys(m), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedKeys = %v, want %v", got, want)
	}
}
