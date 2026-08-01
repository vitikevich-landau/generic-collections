package collections

import (
	"reflect"
	"testing"
)

// TestMapFilterReduce проверяет базовое поведение трёх основных функций над
// срезами, а также сохранение «нулевости» (nil на входе — nil на выходе).
func TestMapFilterReduce(t *testing.T) {
	numbers := []int{1, 2, 3, 4}

	// Map: каждый элемент возводится в квадрат, порядок сохраняется.
	if got, want := Map(numbers, func(v int) int { return v * v }), []int{1, 4, 9, 16}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Map = %v, ожидалось %v", got, want)
	}

	// Filter: остаются только чётные, порядок исходных элементов сохраняется.
	if got, want := Filter(numbers, func(v int) bool { return v%2 == 0 }), []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter = %v, ожидалось %v", got, want)
	}

	// Reduce: свёртка в сумму, 1+2+3+4 = 10.
	if got := Reduce(numbers, 0, func(acc, v int) int { return acc + v }); got != 10 {
		t.Fatalf("Reduce = %d, ожидалось 10", got)
	}

	// nil на входе должен давать именно nil, а не пустой срез: это позволяет
	// вызывающему коду отличить «данных не было» от «всё отфильтровалось».
	// Параметры типа приходится указывать явно — из nil их не вывести.
	if Map[int, int](nil, func(v int) int { return v }) != nil {
		t.Fatal("Map(nil) вернул не nil — «нулевость» не сохранена")
	}
	if Filter[int](nil, func(int) bool { return true }) != nil {
		t.Fatal("Filter(nil) вернул не nil — «нулевость» не сохранена")
	}
}

// TestAppendMapReusesNonOverlappingDestination проверяет главное обещание
// AppendMap: если у переданного буфера достаточно ёмкости, новая память не
// выделяется и результат дописывается прямо в него.
func TestAppendMapReusesNonOverlappingDestination(t *testing.T) {
	// Длина 1, ёмкость 4: хватит места ещё на три дописанных элемента.
	dst := make([]int, 1, 4)
	dst[0] = 10

	// Запоминаем адрес первого элемента: если AppendMap выделит новый массив,
	// адрес изменится — это и есть признак аллокации.
	first := &dst[0]

	got := AppendMap(dst, []int{1, 2, 3}, func(v int) int { return v * 2 })

	// Уже лежавшее в dst значение 10 должно остаться на месте, результаты
	// отображения — дописаться следом.
	if want := []int{10, 2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("AppendMap = %v, ожидалось %v", got, want)
	}
	if &got[0] != first {
		t.Fatal("AppendMap выделил новую память, хотя ёмкости приёмника хватало")
	}
}

// TestAppendFilterInPlace проверяет фильтрацию «на месте»: форма
// AppendFilter(in[:0], in, keep) должна переиспользовать тот же backing array.
func TestAppendFilterInPlace(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	first := &in[0]

	got := AppendFilter(in[:0], in, func(v int) bool { return v%2 == 0 })

	if want := []int{2, 4, 6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("AppendFilter = %v, ожидалось %v", got, want)
	}

	// Совпадение адресов означает, что результат записан поверх входа
	// и ни одной аллокации не произошло.
	if &got[0] != first {
		t.Fatal("AppendFilter не переиспользовал память входа при фильтрации на месте")
	}
}

// TestAppendFilterInPlaceClearsRejectedReferences проверяет, что при фильтрации
// на месте отброшенный хвост обнуляется.
//
// Это принципиально для ссылочных типов: отфильтрованные указатели физически
// остаются в backing array за пределами длины результата и, если их не обнулить,
// продолжают удерживать объекты от сборки мусора — незаметная утечка памяти.
func TestAppendFilterInPlaceClearsRejectedReferences(t *testing.T) {
	type payload struct {
		keep bool
	}

	first := &payload{keep: true}
	second := &payload{}
	third := &payload{keep: true}
	fourth := &payload{}
	in := []*payload{first, second, third, fourth}

	got := AppendFilter(in[:0], in, func(v *payload) bool { return v.keep })

	if want := []*payload{first, third}; !reflect.DeepEqual(got, want) {
		t.Fatalf("AppendFilter = %v, ожидалось %v", got, want)
	}

	// Заглядываем в хвост за длиной результата — туда, куда обычный код уже не
	// добирается, но где память по-прежнему жива.
	for i, v := range got[:cap(got)][len(got):] {
		if v != nil {
			t.Fatalf("отброшенный слот %d всё ещё удерживает указатель %p", len(got)+i, v)
		}
	}
}

// age — пользовательский тип с базовым типом int. Нужен, чтобы проверить, что
// ограничение Number написано через тильду (~int) и потому принимает не только
// сам int, но и типы, определённые на его основе.
type age int

// TestSumAllNumericFamilies проверяет, что Sum принимает представителей всех
// семейств типов, перечисленных в ограничении Number: пользовательский тип с
// базовым int, беззнаковый uintptr и комплексные числа.
func TestSumAllNumericFamilies(t *testing.T) {
	if got := Sum([]age{10, 20, 30}); got != 60 {
		t.Fatalf("Sum(age) = %d, ожидалось 60", got)
	}
	if got := Sum([]uintptr{1, 2, 3}); got != 6 {
		t.Fatalf("Sum(uintptr) = %d, ожидалось 6", got)
	}

	// Комплексные числа складываются покомпонентно и, в частности, показывают,
	// почему ограничение Number не сводится к cmp.Ordered: упорядочить их нельзя.
	if got := Sum([]complex128{1 + 2i, 3 + 4i}); got != 4+6i {
		t.Fatalf("Sum(complex128) = %v, ожидалось 4+6i", got)
	}
}

// TestCompatibilityHelpers проверяет функции, оставленные для обратной
// совместимости (Index и SortedKeys). Они помечены как Deprecated, но должны
// продолжать работать, пока не удалены.
func TestCompatibilityHelpers(t *testing.T) {
	if got := Index([]string{"a", "b", "c"}, "b"); got != 1 {
		t.Fatalf("Index = %d, ожидалось 1", got)
	}

	// Карта заполняется в «неправильном» порядке специально: SortedKeys обязан
	// вернуть ключи по возрастанию независимо от порядка обхода карты.
	m := map[string]int{"b": 2, "a": 1}
	if got, want := SortedKeys(m), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedKeys = %v, ожидалось %v", got, want)
	}
}

// TestKeysNilForEmptyMaps фиксирует контракт Keys и SortedKeys: и для nil-карты,
// и для пустой карты возвращается именно nil, а не пустой срез.
func TestKeysNilForEmptyMaps(t *testing.T) {
	if Keys[string, int](nil) != nil {
		t.Fatal("Keys(nil) вернул не nil")
	}
	if Keys(map[string]int{}) != nil {
		t.Fatal("Keys для пустой карты вернул не nil")
	}
	if SortedKeys(map[string]int{}) != nil {
		t.Fatal("SortedKeys для пустой карты вернул не nil")
	}
}
