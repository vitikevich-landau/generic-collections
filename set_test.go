package collections

import "testing"

// TestSetBasics проверяет базовый жизненный цикл множества: создание с
// отбрасыванием дубликатов, удаление и полную очистку.
func TestSetBasics(t *testing.T) {
	// "a" передано дважды — во множестве должно остаться два элемента, не три.
	s := NewSet("a", "b", "a")
	if s.Len() != 2 || !s.Has("a") || !s.Has("b") {
		t.Fatalf("NewSet = %v, ожидалось {a,b}", s)
	}

	s.Remove("a")
	if s.Has("a") || s.Len() != 1 {
		t.Fatalf("после Remove: Has(a)=%v Len=%d, ожидалось Has(a)=false Len=1", s.Has("a"), s.Len())
	}

	s.Clear()
	if !s.IsEmpty() {
		t.Fatalf("множество не пусто после Clear: %v", s)
	}
}

// TestSetOperations проверяет теоретико-множественные операции и то, что Clone
// действительно создаёт независимую копию, а не второе имя для той же карты.
func TestSetOperations(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(2, 3, 4)

	assertSet(t, a.Union(b), 1, 2, 3, 4) // объединение
	assertSet(t, a.Intersect(b), 2, 3)   // пересечение
	assertSet(t, a.Difference(b), 1)     // разность: есть в a, нет в b

	// Карта — ссылочный тип, поэтому «копия» без явного копирования данных
	// делила бы хранилище с оригиналом. Проверяем, что это не так.
	clone := a.Clone()
	clone.Add(9)
	if a.Has(9) {
		t.Fatal("Clone разделяет хранилище с исходным множеством")
	}

	// Заодно убеждаемся, что операции выше не изменили исходное множество.
	assertSet(t, a, 1, 2, 3)
}

// TestSetIntersectWalksSmallerInput проверяет, что Intersect даёт одинаковый
// результат независимо от порядка операндов.
//
// Внутри Intersect всегда обходит меньшее множество — на сильно разных размерах
// это разница между 3 и 10 000 проверок. Тест фиксирует, что оптимизация не
// ломает корректность и что операция симметрична.
func TestSetIntersectWalksSmallerInput(t *testing.T) {
	large := NewSetWithCapacity[int](10_000)
	for i := 0; i < 10_000; i++ {
		large.Add(i)
	}

	// 20_000 в large отсутствует и в пересечение попасть не должен.
	small := NewSet(1, 9_999, 20_000)

	assertSet(t, large.Intersect(small), 1, 9_999)
	assertSet(t, small.Intersect(large), 1, 9_999)
}

// TestNilSetReadsAndRemoveAreSafeButAddPanics фиксирует поведение нулевого
// значения Set[T] — то есть nil-карты.
//
// Чтение, Remove, Clear и Clone на ней безопасны (так работают встроенные
// операции над nil-картами), а вот Add обязан паниковать: запись в nil-карту
// в Go невозможна.
func TestNilSetReadsAndRemoveAreSafeButAddPanics(t *testing.T) {
	var s Set[int]

	if s.Has(1) || !s.IsEmpty() {
		t.Fatal("nil-множество ведёт себя как непустое")
	}

	// Ни один из этих вызовов не должен паниковать.
	s.Remove(1)
	s.Clear()
	if s.Clone() != nil {
		t.Fatal("Clone для nil-множества вернул не nil — «нулевость» не сохранена")
	}

	// Отложенный recover ловит ожидаемую панику. Если её не случилось,
	// recover() вернёт nil — значит, контракт нарушен.
	defer func() {
		if recover() == nil {
			t.Fatal("Add на nil-множестве не вызвал панику")
		}
	}()
	s.Add(1)
}

// assertSet — вспомогательная функция: проверяет, что множество got состоит
// ровно из значений want (порядок для множества не определён, поэтому
// сравниваются длина и наличие каждого элемента).
//
// Вызов t.Helper() заставляет testing показывать в отчёте строку вызывающего
// теста, а не строку внутри самой assertSet.
func assertSet[T comparable](t *testing.T, got Set[T], want ...T) {
	t.Helper()

	// Проверка длины нужна, чтобы поймать лишние элементы: без неё
	// множество-надмножество прошло бы проверку наличия целиком.
	if len(got) != len(want) {
		t.Fatalf("размер множества = %d, ожидалось %d; множество=%v", len(got), len(want), got)
	}
	for _, v := range want {
		if !got.Has(v) {
			t.Fatalf("множество %v не содержит %v", got, v)
		}
	}
}
