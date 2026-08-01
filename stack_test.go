package collections

import "testing"

// TestStackZeroValueAndLIFO проверяет два свойства сразу: нулевое значение
// Stack[T] готово к работе без конструктора, и порядок извлечения — LIFO.
func TestStackZeroValueAndLIFO(t *testing.T) {
	// Никакого NewStack — работаем прямо с нулевым значением.
	var s Stack[int]

	// На пустом стеке Pop и Peek обязаны вернуть (нулевое значение, false),
	// а не паниковать.
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop на пустом стеке вернул ok=true")
	}
	if got, ok := s.Peek(); ok || got != 0 {
		t.Fatalf("Peek на пустом стеке = (%d, %v), ожидалось (0, false)", got, ok)
	}

	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}

	// Peek показывает вершину, но не снимает её.
	if got, ok := s.Peek(); !ok || got != 3 {
		t.Fatalf("Peek = (%d, %v), ожидалось (3, true)", got, ok)
	}

	// LIFO: положили 1, 2, 3 — снимаем в обратном порядке 3, 2, 1.
	for _, want := range []int{3, 2, 1} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Fatalf("Pop = (%d, %v), ожидалось (%d, true)", got, ok, want)
		}
	}

	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatalf("опустошённый стек: IsEmpty=%v Len=%d, ожидалось IsEmpty=true Len=0", s.IsEmpty(), s.Len())
	}
}

// TestStackClearsRemovedReferences проверяет, что Pop обнуляет освободившийся
// слот и backing array не удерживает снятый указатель.
//
// Без этого объект оставался бы живым для сборщика мусора сколь угодно долго,
// хотя логически он из стека уже удалён.
func TestStackClearsRemovedReferences(t *testing.T) {
	type payload struct{ data [64]byte }
	p := &payload{}

	var s Stack[*payload]
	s.Push(p)
	if got, ok := s.Pop(); !ok || got != p {
		t.Fatalf("Pop = (%p, %v), ожидалось (%p, true)", got, ok, p)
	}

	// Расширяем срез до полной ёмкости, чтобы заглянуть за его длину — туда,
	// где физически остался снятый элемент.
	if retained := s.items[:cap(s.items)][0]; retained != nil {
		t.Fatal("снятый указатель всё ещё удерживается в backing array")
	}
}

// TestStackCapacityAndClear проверяет, что NewStack резервирует запрошенную
// ёмкость, а Clear очищает содержимое, но сохраняет уже выделенную память.
func TestStackCapacityAndClear(t *testing.T) {
	s := NewStack[*int](16)
	if s.Cap() != 16 {
		t.Fatalf("Cap = %d, ожидалось 16", s.Cap())
	}

	a, b := 1, 2
	s.Push(&a)
	s.Push(&b)
	s.Clear()

	// Смысл Clear — именно переиспользование: стек пуст, но память осталась.
	if !s.IsEmpty() || s.Cap() != 16 {
		t.Fatalf("после Clear: IsEmpty=%v Cap=%d, ожидалось IsEmpty=true Cap=16", s.IsEmpty(), s.Cap())
	}

	// При этом ни один слот не должен удерживать указатель.
	for i, v := range s.items[:cap(s.items)] {
		if v != nil {
			t.Fatalf("слот %d backing array не был обнулён", i)
		}
	}
}
