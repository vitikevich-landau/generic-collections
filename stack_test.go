package main

import "testing"

// ─────────────────────────────────────────────────────────────────────────────
// stack_test.go — проверяем порядок LIFO, поведение на пустом стеке и то, что
// параметр типа «донёсся»: Stack[int].Pop отдаёт именно int, Stack[string] —
// string, без всяких приведений.
// ─────────────────────────────────────────────────────────────────────────────

// Пустой стек не должен ничего отдавать — только (нулевое значение, false).
func TestStackEmptyPop(t *testing.T) {
	var s Stack[int]
	if v, ok := s.Pop(); ok || v != 0 {
		t.Fatalf("пустой стек: получили (%v, %v), ожидали (0, false)", v, ok)
	}
	if s.Len() != 0 {
		t.Fatalf("Len пустого стека = %d, ожидали 0", s.Len())
	}
}

// Кладём три, снимаем три — порядок должен быть обратным (LIFO).
func TestStackLIFO(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Len() != 3 {
		t.Fatalf("Len после трёх Push = %d, ожидали 3", s.Len())
	}
	for _, want := range []int{3, 2, 1} { // снимаем сверху → обратный порядок
		if v, ok := s.Pop(); !ok || v != want {
			t.Fatalf("Pop = (%v, %v), ожидали (%d, true)", v, ok, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Fatal("после опустошения Pop должен вернуть ok=false")
	}
}

// Peek показывает вершину, не снимая её.
func TestStackPeek(t *testing.T) {
	var s Stack[string]
	if _, ok := s.Peek(); ok {
		t.Fatal("Peek пустого стека должен вернуть ok=false")
	}
	s.Push("низ")
	s.Push("верх")
	if v, ok := s.Peek(); !ok || v != "верх" {
		t.Fatalf("Peek = (%q, %v), ожидали (\"верх\", true)", v, ok)
	}
	if s.Len() != 2 {
		t.Fatalf("Peek не должен менять размер: Len = %d, ожидали 2", s.Len())
	}
}

// Stack работает и с пользовательскими структурами — тип сохраняется.
func TestStackStructType(t *testing.T) {
	var s Stack[Point]
	s.Push(Point{1, 2})
	if v, ok := s.Pop(); !ok || v != (Point{1, 2}) {
		t.Fatalf("Pop = (%v, %v), ожидали ((1,2), true)", v, ok)
	}
}
