package coll

import "maps"

// MapKeys transforms every key with f and keeps the values. When several keys
// map to the same result, which value survives is unspecified.
func MapKeys[M ~map[K]V, K comparable, V any, R comparable](m M, f func(K) R) map[R]V {
	result := make(map[R]V, len(m))
	for k, v := range m {
		result[f(k)] = v
	}
	return result
}

// MapValues transforms every value with f and keeps the keys.
func MapValues[M ~map[K]V, K comparable, V, R any](m M, f func(V) R) map[K]R {
	result := make(map[K]R, len(m))
	for k, v := range m {
		result[k] = f(v)
	}
	return result
}

// Invert swaps keys and values. When several keys share a value, which key
// survives is unspecified.
func Invert[M ~map[K]V, K, V comparable](m M) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Merge copies the maps from left to right into a new map. Later maps
// overwrite earlier keys. The inputs are not modified.
func Merge[M ~map[K]V, K comparable, V any](ms ...M) M {
	n := 0
	for _, m := range ms {
		n += len(m)
	}
	result := make(M, n)
	for _, m := range ms {
		maps.Copy(result, m)
	}
	return result
}
