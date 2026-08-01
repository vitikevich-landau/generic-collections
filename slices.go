package collections

import (
	"cmp"
	"slices"
)

// Map applies f to every item in in and returns the results. A nil input
// produces a nil result.
func Map[T, U any](in []T, f func(T) U) []U {
	if in == nil {
		return nil
	}
	return AppendMap(make([]U, 0, len(in)), in, f)
}

// AppendMap applies f to every item in in and appends the results to dst. It is
// useful on hot paths where the caller can reuse destination storage. The
// append region of dst must not overlap in; supporting arbitrary overlap would
// require preserving the input in an additional allocation.
func AppendMap[T, U any](dst []U, in []T, f func(T) U) []U {
	dst = slices.Grow(dst, len(in))
	for _, v := range in {
		dst = append(dst, f(v))
	}
	return dst
}

// Filter returns the items for which keep reports true. A nil input produces a
// nil result. The result preallocates for the worst case to guarantee at most
// one allocation.
func Filter[T any](in []T, keep func(T) bool) []T {
	if in == nil {
		return nil
	}
	return AppendFilter(make([]T, 0, len(in)), in, keep)
}

// AppendFilter appends the items accepted by keep to dst. Passing in[:0] as dst
// enables allocation-free in-place filtering and clears the rejected tail so it
// cannot retain references. Other overlapping layouts of dst and in are
// unsupported because appends may overwrite unread input values.
func AppendFilter[T any](dst []T, in []T, keep func(T) bool) []T {
	inPlace := len(in) > 0 && len(dst) == 0 && cap(dst) >= len(in) && &dst[:1][0] == &in[0]
	dst = slices.Grow(dst, len(in))
	for _, v := range in {
		if keep(v) {
			dst = append(dst, v)
		}
	}
	if inPlace {
		clear(in[len(dst):])
	}
	return dst
}

// Reduce folds in into an accumulator, starting with init.
func Reduce[T, U any](in []T, init U, f func(U, T) U) U {
	acc := init
	for _, v := range in {
		acc = f(acc, v)
	}
	return acc
}

// Number contains all predeclared numeric types and user-defined types with the
// same underlying types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~complex64 | ~complex128
}

// Sum returns the sum of in, or the zero value for an empty input.
func Sum[T Number](in []T) T {
	var total T
	for _, v := range in {
		total += v
	}
	return total
}

// Index returns the first index of target, or -1 if target is not present.
//
// Deprecated: use slices.Index directly in new code.
func Index[T comparable](in []T, target T) int {
	return slices.Index(in, target)
}

// Keys returns the keys of m in unspecified order. A nil or empty map produces
// a nil slice.
//
// Deprecated: use maps.Keys with slices.Collect in new code.
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// SortedKeys returns the keys of m in ascending order.
//
// Deprecated: use slices.Sorted(maps.Keys(m)) in new code.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	out := Keys(m)
	slices.Sort(out)
	return out
}
