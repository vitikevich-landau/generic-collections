package main

import (
	"reflect"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// slices_test.go — проверяем универсальные функции и ограничения. Важные точки:
// Map умеет МЕНЯТЬ тип (int→string), Sum работает с любым числовым типом (в том
// числе объявленным через ~), а SortedKeys даёт детерминированный порядок.
// ─────────────────────────────────────────────────────────────────────────────

func TestMapFilterReduce(t *testing.T) {
	nums := []int{1, 2, 3, 4}

	got := Map(nums, func(n int) int { return n * n })
	if want := []int{1, 4, 9, 16}; !reflect.DeepEqual(got, want) {
		t.Errorf("Map: got %v, want %v", got, want)
	}

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4}; !reflect.DeepEqual(evens, want) {
		t.Errorf("Filter: got %v, want %v", evens, want)
	}

	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	if sum != 10 {
		t.Errorf("Reduce: got %d, want 10", sum)
	}
}

// Главная фишка Map — второй параметр типа U: тип на выходе может отличаться.
func TestMapChangesType(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(n int) string {
		return string(rune('a' + n - 1))
	})
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Map int→string: got %v, want %v", got, want)
	}
}

// Age — свой тип с базой int. Тильда ~int в ограничении Number разрешает
// использовать Sum и на нём: это проверяем компиляцией + результатом.
type Age int

func TestSumConstraints(t *testing.T) {
	if got := Sum([]int{1, 2, 3, 4}); got != 10 {
		t.Errorf("Sum[int]: got %d, want 10", got)
	}
	if got := Sum([]float64{0.5, 1.5, 2}); got != 4.0 {
		t.Errorf("Sum[float64]: got %v, want 4.0", got)
	}
	if got := Sum([]Age{10, 20, 30}); got != Age(60) { // ~int в действии
		t.Errorf("Sum[Age]: got %d, want 60", got)
	}
}

func TestIndex(t *testing.T) {
	in := []string{"a", "b", "c"}
	if got := Index(in, "b"); got != 1 {
		t.Errorf("Index(b): got %d, want 1", got)
	}
	if got := Index(in, "z"); got != -1 {
		t.Errorf("Index(z): got %d, want -1", got)
	}
}

// SortedKeys обязан давать один и тот же порядок независимо от случайного обхода
// карты. Заодно проверяем, что он принимает Set[T] (это тоже map).
func TestSortedKeys(t *testing.T) {
	s := NewSet("груша", "яблоко", "слива")
	got := SortedKeys(s)
	if want := []string{"груша", "слива", "яблоко"}; !reflect.DeepEqual(got, want) {
		t.Errorf("SortedKeys: got %v, want %v", got, want)
	}
}
