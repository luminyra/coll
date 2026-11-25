package coll

import (
	"slices"
	"strconv"
	"testing"
)

func TestSliceHelpers(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}

	if got := Map(in, strconv.Itoa); !slices.Equal(got, []string{"1", "2", "3", "4", "5", "6"}) {
		t.Errorf("Map = %v", got)
	}
	if got := Filter(in, isEven); !slices.Equal(got, []int{2, 4, 6}) {
		t.Errorf("Filter = %v", got)
	}
	matched, rest := Partition(in, isEven)
	if !slices.Equal(matched, []int{2, 4, 6}) || !slices.Equal(rest, []int{1, 3, 5}) {
		t.Errorf("Partition = %v, %v", matched, rest)
	}
	if got := Uniq([]int{3, 1, 3, 2, 1}); !slices.Equal(got, []int{3, 1, 2}) {
		t.Errorf("Uniq = %v", got)
	}
	if got := UniqBy(in, isEven); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("UniqBy = %v", got)
	}
	if got := Counts([]string{"a", "b", "a"}); got["a"] != 2 || got["b"] != 1 {
		t.Errorf("Counts = %v", got)
	}

	groups := GroupBy(in, isEven)
	if !slices.Equal(groups[true], []int{2, 4, 6}) || !slices.Equal(groups[false], []int{1, 3, 5}) {
		t.Errorf("GroupBy = %v", groups)
	}
	if got := KeyBy([]string{"a", "bb"}, func(s string) int { return len(s) }); got[1] != "a" || got[2] != "bb" {
		t.Errorf("KeyBy = %v", got)
	}
	if got := ToMap([]string{"x", "yy"}, func(s string) (string, int) { return s, len(s) }); got["yy"] != 2 {
		t.Errorf("ToMap = %v", got)
	}
}

func TestHelpersKeepNamedSliceType(t *testing.T) {
	type IDs []int
	var got IDs = Filter(IDs{1, 2, 3}, isEven)
	if !slices.Equal(got, IDs{2}) {
		t.Errorf("Filter(IDs) = %v", got)
	}
}

func TestHelpersDoNotModifyInput(t *testing.T) {
	in := []int{1, 2, 3}
	out := Filter(in, func(v int) bool { return v != 2 })
	out[0] = 100
	if !slices.Equal(in, []int{1, 2, 3}) {
		t.Fatalf("input modified: %v", in)
	}
}
