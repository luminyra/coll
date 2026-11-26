package coll

import (
	"iter"
	"slices"
)

// Seq is a lazy sequence of values with chainable operations. It is an
// [iter.Seq] with methods, so it can be ranged over directly:
//
//	for name := range coll.From(users).Filter(isActive).Map(userName) { ... }
//
// Intermediate operations such as Filter and Map do no work until a terminal
// operation such as Collect or a range loop pulls values. A Seq built from a
// slice can be iterated any number of times; a Seq built from a single-use
// iterator inherits that restriction.
type Seq[T any] iter.Seq[T]

// From returns a Seq over the items of s, in order.
func From[T any](s []T) Seq[T] {
	return Seq[T](slices.Values(s))
}

// Of returns a Seq over the given items, in order.
func Of[T any](items ...T) Seq[T] {
	return From(items)
}

// FromSeq wraps an iterator, such as maps.Keys(m) or strings.SplitSeq(s, ","),
// as a Seq.
func FromSeq[T any](seq iter.Seq[T]) Seq[T] {
	return Seq[T](seq)
}

// Iter returns s as a plain [iter.Seq], for APIs such as [slices.Collect] or
// [maps.Collect] that take one.
func (s Seq[T]) Iter() iter.Seq[T] {
	return iter.Seq[T](s)
}

// Indexed returns an iterator over the values of s paired with their
// zero-based position.
func (s Seq[T]) Indexed() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range s {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Filter returns the values for which keep reports true.
func (s Seq[T]) Filter(keep func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if keep(v) && !yield(v) {
				return
			}
		}
	}
}

// Map transforms every value with f. The result type may differ from T.
func (s Seq[T]) Map[R any](f func(T) R) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// FilterMap transforms values with f and keeps only the results for which f
// reports true.
func (s Seq[T]) FilterMap[R any](f func(T) (R, bool)) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if r, ok := f(v); ok && !yield(r) {
				return
			}
		}
	}
}

// FlatMap transforms every value into an iterator with f and yields the
// values of each iterator in turn. Iterator helpers such as
// [strings.FieldsSeq] and [slices.Values] fit directly.
func (s Seq[T]) FlatMap[R any](f func(T) iter.Seq[R]) Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			for r := range f(v) {
				if !yield(r) {
					return
				}
			}
		}
	}
}

// Take returns at most the first n values. It stops pulling from s once n
// values have been yielded.
func (s Seq[T]) Take(n int) Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		taken := 0
		for v := range s {
			if !yield(v) {
				return
			}
			if taken++; taken == n {
				return
			}
		}
	}
}

// Skip drops the first n values and returns the rest.
func (s Seq[T]) Skip(n int) Seq[T] {
	return func(yield func(T) bool) {
		skipped := 0
		for v := range s {
			if skipped < n {
				skipped++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// UniqBy drops values whose key, computed by key, has already been seen,
// keeping the first occurrence.
func (s Seq[T]) UniqBy[K comparable](key func(T) K) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for v := range s {
			k := key(v)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

// Chunk groups consecutive values into slices of size values; the last slice
// may be shorter. Each slice is newly allocated. Chunk panics if size is not
// positive. For a slice that is already in memory, [slices.Chunk] avoids the
// copies.
//
// Chunk returns a plain [iter.Seq] because a method of Seq[T] cannot return
// Seq[[]T] (that would instantiate Seq infinitely); wrap it with FromSeq to
// keep chaining.
func (s Seq[T]) Chunk(size int) iter.Seq[[]T] {
	if size <= 0 {
		panic("coll: Chunk size must be positive")
	}
	return func(yield func([]T) bool) {
		chunk := make([]T, 0, size)
		for v := range s {
			chunk = append(chunk, v)
			if len(chunk) == size {
				if !yield(chunk) {
					return
				}
				chunk = make([]T, 0, size)
			}
		}
		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// Collect returns the values of s in a new slice.
func (s Seq[T]) Collect() []T {
	// Start with room for a few values: most pipelines yield at least some,
	// and growing from zero costs several small allocations.
	return slices.AppendSeq(make([]T, 0, 8), s.Iter())
}

// Fold combines the values of s from left to right with f, starting from
// init.
func (s Seq[T]) Fold[R any](init R, f func(acc R, v T) R) R {
	acc := init
	for v := range s {
		acc = f(acc, v)
	}
	return acc
}

// Count returns the number of values. Combine it with Filter to count
// matches.
func (s Seq[T]) Count() int {
	n := 0
	for range s {
		n++
	}
	return n
}

// First returns the first value. The boolean is false when s is empty.
func (s Seq[T]) First() (T, bool) {
	for v := range s {
		return v, true
	}
	var zero T
	return zero, false
}

// Find returns the first value for which match reports true. The boolean is
// false when no value matches.
func (s Seq[T]) Find(match func(T) bool) (T, bool) {
	return s.Filter(match).First()
}

// Any reports whether match is true for at least one value. It stops at the
// first match and returns false for an empty sequence.
func (s Seq[T]) Any(match func(T) bool) bool {
	_, ok := s.Find(match)
	return ok
}

// Every reports whether match is true for every value. It stops at the first
// mismatch and returns true for an empty sequence.
func (s Seq[T]) Every(match func(T) bool) bool {
	return !s.Any(func(v T) bool { return !match(v) })
}

// GroupBy groups the values by the key computed by key. Values keep their
// order within each group.
func (s Seq[T]) GroupBy[K comparable](key func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for v := range s {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

// KeyBy builds a map from the key computed by key to the value. Later values
// overwrite earlier ones with the same key.
func (s Seq[T]) KeyBy[K comparable](key func(T) K) map[K]T {
	result := make(map[K]T)
	for v := range s {
		result[key(v)] = v
	}
	return result
}
