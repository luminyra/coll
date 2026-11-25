package coll

// Map transforms every item of s with f and returns the results in order.
func Map[T, R any](s []T, f func(T) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Filter returns the items of s for which keep reports true, in order.
func Filter[S ~[]T, T any](s S, keep func(T) bool) S {
	result := make(S, 0, len(s))
	for _, v := range s {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

// Partition splits s into the items for which match reports true and the
// rest. Both results keep the input order.
func Partition[S ~[]T, T any](s S, match func(T) bool) (matched, rest S) {
	for _, v := range s {
		if match(v) {
			matched = append(matched, v)
		} else {
			rest = append(rest, v)
		}
	}
	return matched, rest
}

// GroupBy groups the items of s by the key computed by key. Items keep their
// order within each group.
func GroupBy[S ~[]T, T any, K comparable](s S, key func(T) K) map[K]S {
	groups := make(map[K]S)
	for _, v := range s {
		k := key(v)
		groups[k] = append(groups[k], v)
	}
	return groups
}

// KeyBy builds a map from the key computed by key to the item. Later items
// overwrite earlier ones with the same key.
func KeyBy[T any, K comparable](s []T, key func(T) K) map[K]T {
	result := make(map[K]T, len(s))
	for _, v := range s {
		result[key(v)] = v
	}
	return result
}

// ToMap builds a map from the key-value pairs returned by f. Later pairs
// overwrite earlier ones with the same key.
func ToMap[T any, K comparable, V any](s []T, f func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(s))
	for _, item := range s {
		k, v := f(item)
		result[k] = v
	}
	return result
}

// Uniq returns the items of s with duplicates removed, keeping the first
// occurrence of each value.
func Uniq[S ~[]T, T comparable](s S) S {
	return UniqBy(s, func(v T) T { return v })
}

// UniqBy removes items whose key, computed by key, has already been seen,
// keeping the first occurrence.
func UniqBy[S ~[]T, T any, K comparable](s S, key func(T) K) S {
	seen := make(map[K]struct{}, len(s))
	result := make(S, 0, len(s))
	for _, v := range s {
		k := key(v)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, v)
	}
	return result
}

// Counts returns how many times each value appears in s.
func Counts[T comparable](s []T) map[T]int {
	result := make(map[T]int)
	for _, v := range s {
		result[v]++
	}
	return result
}
