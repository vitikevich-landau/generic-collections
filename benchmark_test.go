package collections

import "testing"

// Переменные-«стоки» (sink). Результаты операций записываются в них, чтобы
// компилятор не счёл вычисления бесполезными и не выбросил их целиком:
// присваивание в пакетную переменную он оптимизировать не вправе.
var (
	benchmarkIntSink  int
	benchmarkSetSink  Set[int]
	benchmarkByteSink byte
)

// benchmarkLargeValue — намеренно крупное значение (256 байт). Нужно, чтобы
// показать стоимость копирования самих элементов, а не только служебных
// операций очереди.
type benchmarkLargeValue [256]byte

// BenchmarkQueueChurnInt измеряет установившийся режим очереди на маленьких
// значениях: длина очереди не меняется, буфер не растёт, и стоимость пары
// Enqueue/Dequeue сводится к записи в слот и сдвигу индексов.
//
// ReportAllocs добавляет в отчёт колонки с аллокациями — здесь ожидается ноль.
func BenchmarkQueueChurnInt(b *testing.B) {
	q := NewQueue[int](1024)

	// Прогрев: заполняем буфер наполовину. Заполнять его до предела нельзя —
	// тогда первый же Enqueue в измеряемом цикле вызвал бы grow с копированием
	// всего буфера, и замер перестал бы быть установившимся режимом. На длинном
	// прогоне этот единственный рост размазался бы по b.N и округлился до нуля,
	// а на коротком (-benchtime=1x) исказил бы результат на три порядка.
	for i := 0; i < 512; i++ {
		q.Enqueue(i)
	}

	b.ReportAllocs()

	// Сбрасываем таймер, чтобы время прогрева не попало в результат.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		v, _ := q.Dequeue()
		benchmarkIntSink = v
	}
}

// BenchmarkQueueChurnLargeValue — тот же сценарий, что и BenchmarkQueueChurnInt,
// но с 256-байтными элементами. Сравнение двух бенчмарков показывает, какая
// часть стоимости приходится на копирование значений, а какая — на саму очередь.
func BenchmarkQueueChurnLargeValue(b *testing.B) {
	q := NewQueue[benchmarkLargeValue](1024)

	// Как и выше, заполняем лишь половину буфера, чтобы в измеряемом цикле
	// не происходило роста.
	for i := 0; i < 512; i++ {
		q.Enqueue(benchmarkLargeValue{byte(i)})
	}

	// Значение готовится заранее, чтобы его создание не попало в измерения.
	value := benchmarkLargeValue{42}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(value)
		v, _ := q.Dequeue()
		benchmarkByteSink = v[0]
	}
}

// BenchmarkSetIntersectSkewed измеряет пересечение множеств сильно разного
// размера: 100 000 против 64.
//
// Intersect обходит меньшее множество, поэтому время работы должно определяться
// числом 64, а не 100 000. Именно этот перекос («skew») и проверяет бенчмарк:
// наивная реализация, обходящая ресивер, была бы здесь на три порядка медленнее.
func BenchmarkSetIntersectSkewed(b *testing.B) {
	large := NewSetWithCapacity[int](100_000)
	for i := 0; i < 100_000; i++ {
		large.Add(i)
	}
	small := NewSetWithCapacity[int](64)
	for i := 0; i < 64; i++ {
		small.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSetSink = large.Intersect(small)
	}
}
