package coll

import (
	"maps"
	"strings"
	"testing"
)

func TestMapHelpers(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if got := MapValues(m, func(v int) int { return v * 2 }); !maps.Equal(got, map[string]int{"a": 2, "b": 4, "c": 6}) {
		t.Errorf("MapValues = %v", got)
	}
	if got := MapKeys(m, strings.ToUpper); !maps.Equal(got, map[string]int{"A": 1, "B": 2, "C": 3}) {
		t.Errorf("MapKeys = %v", got)
	}
	if got := Invert(m); !maps.Equal(got, map[int]string{1: "a", 2: "b", 3: "c"}) {
		t.Errorf("Invert = %v", got)
	}
	merged := Merge(m, map[string]int{"a": 10, "d": 4})
	if !maps.Equal(merged, map[string]int{"a": 10, "b": 2, "c": 3, "d": 4}) {
		t.Errorf("Merge = %v", merged)
	}
	if m["a"] != 1 || len(m) != 3 {
		t.Errorf("Merge modified its input: %v", m)
	}
}
