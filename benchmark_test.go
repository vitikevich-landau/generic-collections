package collections

import "testing"

var (
	benchmarkIntSink  int
	benchmarkSetSink  Set[int]
	benchmarkByteSink byte
)

type benchmarkLargeValue [256]byte

func BenchmarkQueueChurnInt(b *testing.B) {
	q := NewQueue[int](1024)
	for i := 0; i < 1024; i++ {
		q.Enqueue(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
		v, _ := q.Dequeue()
		benchmarkIntSink = v
	}
}

func BenchmarkQueueChurnLargeValue(b *testing.B) {
	q := NewQueue[benchmarkLargeValue](1024)
	for i := 0; i < 1024; i++ {
		q.Enqueue(benchmarkLargeValue{byte(i)})
	}
	value := benchmarkLargeValue{42}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enqueue(value)
		v, _ := q.Dequeue()
		benchmarkByteSink = v[0]
	}
}

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
