package coll

// Intersect returns the distinct items present in both a and b, in the order
// they appear in a.
func Intersect[S ~[]T, T comparable](a, b S) S {
	inB := toSet(b)
	return Uniq(Filter(a, func(v T) bool {
		_, ok := inB[v]
		return ok
	}))
}

// Difference returns the items of a that are not in b, keeping their order
// and duplicates. Call Difference(b, a) for the other direction.
func Difference[S ~[]T, T comparable](a, b S) S {
	inB := toSet(b)
	return Filter(a, func(v T) bool {
		_, ok := inB[v]
		return !ok
	})
}

// Union returns the distinct items of all slices, in order of first
// appearance.
func Union[S ~[]T, T comparable](sets ...S) S {
	total := 0
	for _, s := range sets {
		total += len(s)
	}
	all := make(S, 0, total)
	for _, s := range sets {
		all = append(all, s...)
	}
	return Uniq(all)
}

func toSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, v := range items {
		set[v] = struct{}{}
	}
	return set
}
