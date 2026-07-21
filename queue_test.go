package main

import "testing"

// ─────────────────────────────────────────────────────────────────────────────
// queue_test.go — проверяем порядок FIFO и поведение на пустой очереди. В
// отличие от стека, очередь отдаёт элементы в том же порядке, в каком их клали.
// ─────────────────────────────────────────────────────────────────────────────

// Пустая очередь: Dequeue → (нулевое значение, false), без паники.
func TestQueueEmptyDequeue(t *testing.T) {
	var q Queue[string]
	if v, ok := q.Dequeue(); ok || v != "" {
		t.Fatalf("пустая очередь: получили (%q, %v), ожидали (\"\", false)", v, ok)
	}
}

// Кладём три, снимаем три — порядок должен совпасть (FIFO).
func TestQueueFIFO(t *testing.T) {
	var q Queue[int]
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Len() != 3 {
		t.Fatalf("Len после трёх Enqueue = %d, ожидали 3", q.Len())
	}
	for _, want := range []int{1, 2, 3} { // берём с начала → прямой порядок
		if v, ok := q.Dequeue(); !ok || v != want {
			t.Fatalf("Dequeue = (%v, %v), ожидали (%d, true)", v, ok, want)
		}
	}
	if q.Len() != 0 {
		t.Fatalf("после опустошения Len = %d, ожидали 0", q.Len())
	}
}
