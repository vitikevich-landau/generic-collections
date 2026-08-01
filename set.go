package collections

// Set is a hash set for comparable values. A nil Set supports reads and
// removals, but Add requires a Set created with NewSet, NewSetWithCapacity, or
// make. Set is not safe for concurrent use.
type Set[T comparable] map[T]struct{}

// NewSet creates a set containing values.
func NewSet[T comparable](values ...T) Set[T] {
	s := make(Set[T], len(values))
	for _, v := range values {
		s[v] = struct{}{}
	}
	return s
}

// NewSetWithCapacity creates an empty set preallocated for capacity values.
func NewSetWithCapacity[T comparable](capacity int) Set[T] {
	return make(Set[T], capacity)
}

// Add inserts v into the set.
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// Has reports whether v is in the set.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Remove deletes v from the set.
func (s Set[T]) Remove(v T) {
	delete(s, v)
}

// Len returns the number of values in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty reports whether the set contains no values.
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all values while retaining allocated storage for reuse.
func (s Set[T]) Clear() {
	clear(s)
}

// Clone returns a shallow copy of s.
func (s Set[T]) Clone() Set[T] {
	if s == nil {
		return nil
	}
	out := make(Set[T], len(s))
	for v := range s {
		out[v] = struct{}{}
	}
	return out
}

// Union returns a new set containing values present in either set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(other))
	for v := range s {
		out[v] = struct{}{}
	}
	for v := range other {
		out[v] = struct{}{}
	}
	return out
}

// Intersect returns a new set containing values present in both sets. It walks
// the smaller input so highly skewed set sizes remain efficient.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	if len(s) > len(other) {
		s, other = other, s
	}

	out := make(Set[T], len(s))
	for v := range s {
		if _, ok := other[v]; ok {
			out[v] = struct{}{}
		}
	}
	return out
}

// Difference returns a new set containing values present in s but not other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	out := make(Set[T], len(s))
	for v := range s {
		if _, ok := other[v]; !ok {
			out[v] = struct{}{}
		}
	}
	return out
}
