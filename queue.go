package collections

const minQueueCapacity = 8

// Queue is a dynamically growing FIFO ring buffer. Its zero value is ready to
// use. A Queue must not be copied after first use and is not safe for
// concurrent use.
type Queue[T any] struct {
	noCopy noCopy
	buf    []T
	head   int
	size   int
}

// NewQueue returns an empty queue with capacity rounded up to a power of two.
func NewQueue[T any](capacity int) *Queue[T] {
	capacity = queueCapacity(capacity)
	if capacity == 0 {
		return &Queue[T]{}
	}
	return &Queue[T]{buf: make([]T, capacity)}
}

// Enqueue adds v to the back of the queue.
func (q *Queue[T]) Enqueue(v T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	tail := (q.head + q.size) & (len(q.buf) - 1)
	q.buf[tail] = v
	q.size++
}

// Dequeue removes and returns the front item. The boolean result is false when
// the queue is empty. The vacated slot is cleared so it cannot retain pointers.
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	v := q.buf[q.head]
	var zero T
	q.buf[q.head] = zero
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.size--
	if q.size == 0 {
		q.head = 0
	}
	return v, true
}

// Peek returns the front item without removing it. The boolean result is false
// when the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}
	return q.buf[q.head], true
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	return q.size
}

// Cap returns the queue's current storage capacity.
func (q *Queue[T]) Cap() int {
	return len(q.buf)
}

// IsEmpty reports whether the queue contains no items.
func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

// Clear removes all items while retaining allocated storage for reuse.
func (q *Queue[T]) Clear() {
	if q.size == 0 {
		return
	}

	end := q.head + q.size
	if end <= len(q.buf) {
		clear(q.buf[q.head:end])
	} else {
		clear(q.buf[q.head:])
		clear(q.buf[:end-len(q.buf)])
	}
	q.head = 0
	q.size = 0
}

func (q *Queue[T]) grow() {
	capacity := minQueueCapacity
	if len(q.buf) != 0 {
		maxInt := int(^uint(0) >> 1)
		if len(q.buf) > maxInt/2 {
			panic("collections: queue capacity overflow")
		}
		capacity = len(q.buf) * 2
	}

	buf := make([]T, capacity)
	if q.size != 0 {
		if q.head+q.size <= len(q.buf) {
			copy(buf, q.buf[q.head:q.head+q.size])
		} else {
			n := copy(buf, q.buf[q.head:])
			copy(buf[n:], q.buf[:q.size-n])
		}
	}
	q.buf = buf
	q.head = 0
}

func queueCapacity(requested int) int {
	if requested < 0 {
		panic("collections: negative queue capacity")
	}
	if requested == 0 {
		return 0
	}

	capacity := minQueueCapacity
	maxInt := int(^uint(0) >> 1)
	for capacity < requested {
		if capacity > maxInt/2 {
			panic("collections: queue capacity overflow")
		}
		capacity *= 2
	}
	return capacity
}
