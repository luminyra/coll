package coll

import (
	"slices"
	"testing"
)

func TestSets(t *testing.T) {
	a, b := []int{1, 2, 2, 3, 1}, []int{2, 3, 4}

	tests := []struct {
		name     string
		got      []int
		expected []int
	}{
		{name: "intersect", got: Intersect(a, b), expected: []int{2, 3}},
		{name: "difference keeps order and duplicates", got: Difference(a, b), expected: []int{1, 1}},
		{name: "difference other direction", got: Difference(b, a), expected: []int{4}},
		{name: "union", got: Union(a, b), expected: []int{1, 2, 3, 4}},
		{name: "union of nothing", got: Union[[]int](), expected: []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !slices.Equal(tt.got, tt.expected) {
				t.Errorf("got %v, expected %v", tt.got, tt.expected)
			}
		})
	}
}
