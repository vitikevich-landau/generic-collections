package main

import "testing"

// ─────────────────────────────────────────────────────────────────────────────
// set_test.go — проверяем уникальность, наличие, удаление и операции над
// множествами. Отдельно — главная грабля темы: Add на nil-карте ПАНИКует, а
// созданное через NewSet множество работает.
// ─────────────────────────────────────────────────────────────────────────────

// NewSet отсеивает дубликаты; Has/Remove/Len работают как ожидается.
func TestSetBasics(t *testing.T) {
	s := NewSet("a", "b", "a") // дубликат "a" должен схлопнуться
	if s.Len() != 2 {
		t.Fatalf("Len = %d, ожидали 2 (дубликат отброшен)", s.Len())
	}
	if !s.Has("a") || !s.Has("b") {
		t.Fatal("Has должен находить оба элемента")
	}
	if s.Has("z") {
		t.Fatal("Has(\"z\") = true, ожидали false")
	}
	s.Remove("a")
	if s.Has("a") || s.Len() != 1 {
		t.Fatalf("после Remove(\"a\"): Has=%v Len=%d, ожидали false и 1", s.Has("a"), s.Len())
	}
	s.Remove("нет-такого") // Remove отсутствующего — не ошибка
}

// Union и Intersect возвращают новое множество и не трогают исходные.
func TestSetUnionIntersect(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(2, 3, 4)

	u := a.Union(b)
	if u.Len() != 4 || !u.Has(1) || !u.Has(4) {
		t.Fatalf("Union = %v (Len %d), ожидали {1,2,3,4}", SortedKeys(u), u.Len())
	}
	i := a.Intersect(b)
	if i.Len() != 2 || !i.Has(2) || !i.Has(3) || i.Has(1) {
		t.Fatalf("Intersect = %v (Len %d), ожидали {2,3}", SortedKeys(i), i.Len())
	}
	// Исходные множества не изменились.
	if a.Len() != 3 || b.Len() != 3 {
		t.Fatalf("исходные множества изменились: a=%d b=%d", a.Len(), b.Len())
	}
}

// nil-множество: Has/Remove безопасны, а вот Add ПАНИКует — это и показываем
// через recover. Именно из-за этого множество создают через NewSet/make.
func TestNilSetAddPanics(t *testing.T) {
	var s Set[int] // nil-карта

	if s.Has(1) {
		t.Fatal("Has на nil-множестве должен вернуть false, а не паниковать")
	}
	s.Remove(1) // no-op, без паники

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("ожидали панику при Add на nil-множестве, её не было")
		}
	}()
	s.Add(1) // ПАНИКА: assignment to entry in nil map
	t.Fatal("до этой строки дойти не должны — Add обязан был паникнуть")
}
