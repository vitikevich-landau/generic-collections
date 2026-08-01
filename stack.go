package collections

// Stack is a LIFO collection. Its zero value is ready to use. A Stack must not
// be copied after first use and is not safe for concurrent use.
type Stack[T any] struct {
	noCopy noCopy
	items  []T
}

// NewStack returns an empty stack with space preallocated for capacity items.
func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{items: make([]T, 0, capacity)}
}

// Push adds v to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top item. The boolean result is false when the
// stack is empty. The vacated slot is cleared so it cannot keep references
// alive through the backing array.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	last := len(s.items) - 1
	v := s.items[last]
	var zero T
	s.items[last] = zero
	s.items = s.items[:last]
	return v, true
}

// Peek returns the top item without removing it. The boolean result is false
// when the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Cap returns the capacity of the stack's backing storage.
func (s *Stack[T]) Cap() int {
	return cap(s.items)
}

// IsEmpty reports whether the stack contains no items.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all items while retaining allocated storage for reuse.
func (s *Stack[T]) Clear() {
	clear(s.items)
	s.items = s.items[:0]
}
