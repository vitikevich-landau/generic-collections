package collections

import (
	"math"
	"testing"
)

// TestQueueZeroValueAndFIFO проверяет, что нулевое значение Queue[T] готово к
// работе (буфер выделяется при первой вставке) и что порядок извлечения — FIFO.
func TestQueueZeroValueAndFIFO(t *testing.T) {
	// Никакого NewQueue — работаем прямо с нулевым значением.
	var q Queue[int]

	// На пустой очереди Dequeue возвращает (нулевое значение, false).
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue на пустой очереди вернул ok=true")
	}

	for _, v := range []int{1, 2, 3} {
		q.Enqueue(v)
	}

	// Peek показывает начало очереди, не извлекая элемент.
	if got, ok := q.Peek(); !ok || got != 1 {
		t.Fatalf("Peek = (%d, %v), ожидалось (1, true)", got, ok)
	}

	// FIFO: порядок извлечения совпадает с порядком вставки — в отличие от
	// стека, где он обратный.
	for _, want := range []int{1, 2, 3} {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Fatalf("Dequeue = (%d, %v), ожидалось (%d, true)", got, ok, want)
		}
	}

	if !q.IsEmpty() || q.Len() != 0 {
		t.Fatalf("опустошённая очередь: IsEmpty=%v Len=%d, ожидалось IsEmpty=true Len=0", q.IsEmpty(), q.Len())
	}
}

// TestQueueWrapAndGrow проверяет самый каверзный случай кольцевого буфера:
// рост, когда данные «завёрнуты» через конец буфера.
//
// Сценарий: заполняем 8 слотов, извлекаем 4 (голова уезжает на индекс 4),
// доливаем ещё 9 — хвост заворачивается в начало буфера, а затем буфер должен
// удвоиться. При росте элементы обязаны быть перенесены в правильном порядке:
// сначала кусок от головы до конца буфера, потом кусок из его начала.
func TestQueueWrapAndGrow(t *testing.T) {
	q := NewQueue[int](8)
	for i := 0; i < 8; i++ {
		q.Enqueue(i)
	}

	// Освобождаем место в начале буфера и сдвигаем голову.
	for i := 0; i < 4; i++ {
		if got, _ := q.Dequeue(); got != i {
			t.Fatalf("первичный Dequeue = %d, ожидалось %d", got, i)
		}
	}

	// Эти вставки сначала заворачиваются в начало буфера, а последняя
	// переполняет его и вызывает удвоение.
	for i := 8; i < 17; i++ {
		q.Enqueue(i)
	}

	if q.Cap() != 16 {
		t.Fatalf("Cap после роста из завёрнутого состояния = %d, ожидалось 16", q.Cap())
	}

	// Главная проверка: после переноса порядок FIFO не нарушен.
	for want := 4; want < 17; want++ {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Fatalf("Dequeue после заворота = (%d, %v), ожидалось (%d, true)", got, ok, want)
		}
	}
}

// TestQueueGrowContiguous проверяет вторую ветку роста — когда элементы лежат
// в буфере непрерывно и переносятся одним копированием.
//
// 17 вставок в нулевую очередь: буфер создаётся размером 8, затем удваивается
// до 16 и до 32 (ёмкость всегда степень двойки).
func TestQueueGrowContiguous(t *testing.T) {
	var q Queue[int]
	for i := 0; i < 17; i++ {
		q.Enqueue(i)
	}

	if q.Cap() != 32 {
		t.Fatalf("Cap после непрерывного роста = %d, ожидалось 32", q.Cap())
	}
	for want := 0; want < 17; want++ {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Fatalf("Dequeue = (%d, %v), ожидалось (%d, true)", got, ok, want)
		}
	}
}

// TestQueuePeekEmptyAndWrapped проверяет Peek в двух пограничных состояниях:
// на пустой очереди и на очереди, завёрнутой через конец буфера.
func TestQueuePeekEmptyAndWrapped(t *testing.T) {
	var q Queue[int]
	if got, ok := q.Peek(); ok || got != 0 {
		t.Fatalf("Peek на пустой очереди = (%d, %v), ожидалось (0, false)", got, ok)
	}

	// Приводим очередь в завёрнутое состояние: 8 вставок, 5 извлечений
	// (голова на индексе 5), затем ещё 3 вставки — они лягут в начало буфера.
	for i := 0; i < 8; i++ {
		q.Enqueue(i)
	}
	for i := 0; i < 5; i++ {
		q.Dequeue()
	}
	for i := 8; i < 11; i++ {
		q.Enqueue(i)
	}

	// Начало очереди определяется головой, а не физическим началом буфера.
	if got, ok := q.Peek(); !ok || got != 5 {
		t.Fatalf("Peek на завёрнутой очереди = (%d, %v), ожидалось (5, true)", got, ok)
	}
}

// TestQueueClearsRemovedReferences проверяет, что Dequeue обнуляет
// освободившийся слот и буфер не удерживает извлечённый указатель.
func TestQueueClearsRemovedReferences(t *testing.T) {
	type payload struct{ data [64]byte }
	p := &payload{}

	// Запрошена ёмкость 1, но фактическая будет minQueueCapacity (8).
	q := NewQueue[*payload](1)
	q.Enqueue(p)
	if got, ok := q.Dequeue(); !ok || got != p {
		t.Fatalf("Dequeue = (%p, %v), ожидалось (%p, true)", got, ok, p)
	}

	for i, v := range q.buf {
		if v != nil {
			t.Fatalf("слот %d буфера удерживает извлечённый указатель", i)
		}
	}
}

// TestQueueClearWrapped проверяет Clear на завёрнутой очереди — случай, когда
// занятые слоты образуют два участка (хвост буфера и его начало) и обнулять
// нужно оба.
func TestQueueClearWrapped(t *testing.T) {
	q := NewQueue[*int](8)

	// Отдельный массив значений — чтобы в очередь попали 12 разных указателей.
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

	// Clear очищает содержимое, но сохраняет буфер для переиспользования.
	if !q.IsEmpty() || q.Cap() != capacity {
		t.Fatalf("после Clear: IsEmpty=%v Cap=%d, ожидалось IsEmpty=true Cap=%d", q.IsEmpty(), q.Cap(), capacity)
	}
	for i, v := range q.buf {
		if v != nil {
			t.Fatalf("слот %d буфера не был обнулён", i)
		}
	}
}

// TestQueueClearContiguous проверяет вторую ветку Clear — когда занятые слоты
// лежат одним непрерывным участком. Заодно проверяется, что повторный Clear на
// уже пустой очереди безопасен и ничего не делает.
func TestQueueClearContiguous(t *testing.T) {
	q := NewQueue[*int](8)
	values := make([]int, 8)
	for i := range values {
		q.Enqueue(&values[i])
	}

	// Голова сдвигается, но заворота не происходит: участок остаётся непрерывным.
	for i := 0; i < 3; i++ {
		q.Dequeue()
	}
	q.Clear()

	if !q.IsEmpty() || q.Cap() != 8 {
		t.Fatalf("после Clear: IsEmpty=%v Cap=%d, ожидалось IsEmpty=true Cap=8", q.IsEmpty(), q.Cap())
	}
	for i, v := range q.buf {
		if v != nil {
			t.Fatalf("слот %d буфера не был обнулён", i)
		}
	}

	// Clear на пустой очереди — не ошибка и не должен ничего ломать.
	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("Clear на пустой очереди не является пустой операцией")
	}
}

// TestNewQueueCapacity проверяет правила округления запрошенной ёмкости:
// ноль означает «не выделять буфер», всё от 1 до 8 округляется до
// minQueueCapacity, дальше — до ближайшей степени двойки вверх. Отрицательная
// ёмкость — паника.
func TestNewQueueCapacity(t *testing.T) {
	if got := NewQueue[int](0).Cap(); got != 0 {
		t.Fatalf("NewQueue(0).Cap = %d, ожидалось 0", got)
	}

	// Всё, что не превышает минимум, поднимается ровно до минимума.
	for _, requested := range []int{1, 2, 4, 8} {
		if got := NewQueue[int](requested).Cap(); got != minQueueCapacity {
			t.Fatalf("NewQueue(%d).Cap = %d, ожидалось %d", requested, got, minQueueCapacity)
		}
	}

	// 9 не является степенью двойки — округляется вверх до 16.
	if got := NewQueue[int](9).Cap(); got != 16 {
		t.Fatalf("NewQueue(9).Cap = %d, ожидалось 16", got)
	}

	// Проверку паники держим последней: после неё выполнение теста прервётся,
	// и отложенный recover доведёт его до конца.
	defer func() {
		if recover() == nil {
			t.Fatal("NewQueue с отрицательной ёмкостью не вызвал панику")
		}
	}()
	NewQueue[int](-1)
}

// TestNewQueueCapacityOverflowPanics проверяет защиту от переполнения: для
// math.MaxInt ближайшей степени двойки вверх в пределах int не существует,
// поэтому вместо бесконечного цикла или отрицательной ёмкости ожидается паника.
func TestNewQueueCapacityOverflowPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewQueue с непредставимой ёмкостью не вызвал панику")
		}
	}()
	NewQueue[int](math.MaxInt)
}

// FuzzQueueMatchesSlice — фаззинг-тест: сравнивает поведение очереди с заведомо
// корректной эталонной моделью (обычным срезом).
//
// Каждый байт входных данных трактуется как операция: чётный — Enqueue этого
// байта, нечётный — Dequeue. Ту же последовательность применяем к срезу и после
// каждого шага сверяем результаты. Такой приём (model-based testing) находит
// редкие комбинации заворота и роста буфера, которые трудно придумать вручную.
func FuzzQueueMatchesSlice(f *testing.F) {
	// Затравочные данные (seed corpus): фаззер стартует с них и мутирует дальше.
	f.Add([]byte{0, 2, 4, 1, 3, 5})       // вставки вперемешку с извлечениями
	f.Add([]byte{0, 0, 0, 0, 1, 1, 0, 1}) // серия вставок, затем извлечения

	f.Fuzz(func(t *testing.T, operations []byte) {
		var q Queue[byte]

		// Эталон: срез, у которого начало — это начало очереди.
		var reference []byte

		for _, operation := range operations {
			if operation&1 == 0 {
				// Чётный байт — вставка. Кладём сам байт как значение.
				q.Enqueue(operation)
				reference = append(reference, operation)
			} else {
				// Нечётный байт — извлечение.
				got, ok := q.Dequeue()
				if len(reference) == 0 {
					// Эталон пуст — очередь тоже обязана вернуть ok=false.
					if ok {
						t.Fatalf("Dequeue при пустом эталоне = (%d, true), ожидалось ok=false", got)
					}
				} else {
					if !ok || got != reference[0] {
						t.Fatalf("Dequeue = (%d, %v), ожидалось (%d, true)", got, ok, reference[0])
					}
					reference = reference[1:]
				}
			}

			// Инвариант, проверяемый после каждой операции: Peek показывает то
			// же начало, что и эталон, и согласован с его пустотой.
			front, ok := q.Peek()
			if ok != (len(reference) > 0) || (ok && front != reference[0]) {
				t.Fatalf("Peek = (%d, %v), ожидалось начало %v", front, ok, reference)
			}
		}

		// Финальная сверка: длина и полное содержимое до опустошения.
		if q.Len() != len(reference) {
			t.Fatalf("Len = %d, ожидалось %d", q.Len(), len(reference))
		}
		for _, want := range reference {
			got, ok := q.Dequeue()
			if !ok || got != want {
				t.Fatalf("Dequeue при опустошении = (%d, %v), ожидалось (%d, true)", got, ok, want)
			}
		}
	})
}
