package coll

import (
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func isEven(v int) bool { return v%2 == 0 }

func TestSeqChain(t *testing.T) {
	got := From([]int{1, 2, 3, 4, 5, 6}).Filter(isEven).Map(strconv.Itoa).Collect()
	if !slices.Equal(got, []string{"2", "4", "6"}) {
		t.Errorf("Filter.Map = %v", got)
	}

	words := Of("a b", "c", "d e f").FlatMap(strings.FieldsSeq).Collect()
	if !slices.Equal(words, []string{"a", "b", "c", "d", "e", "f"}) {
		t.Errorf("FlatMap = %v", words)
	}

	parsed := Of("1", "x", "3").FilterMap(func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		return n, err == nil
	}).Collect()
	if !slices.Equal(parsed, []int{1, 3}) {
		t.Errorf("FilterMap = %v", parsed)
	}

	if got := Of(1, 2, 3, 4).Skip(1).Take(2).Collect(); !slices.Equal(got, []int{2, 3}) {
		t.Errorf("Skip.Take = %v", got)
	}
	if got := Of(1, 2).Take(0).Collect(); len(got) != 0 {
		t.Errorf("Take(0) = %v", got)
	}
	if got := Of(3, 1, 3, 2, 1).UniqBy(func(v int) int { return v }).Collect(); !slices.Equal(got, []int{3, 1, 2}) {
		t.Errorf("UniqBy = %v", got)
	}
	if got := Of(1, 2, 3, 4, 5).Fold("", func(acc string, v int) string { return acc + strconv.Itoa(v) }); got != "12345" {
		t.Errorf("Fold = %q", got)
	}
}

func TestSeqIsLazy(t *testing.T) {
	pulled := 0
	source := FromSeq(func(yield func(int) bool) {
		for i := 1; ; i++ {
			pulled++
			if !yield(i) {
				return
			}
		}
	})

	got := source.Filter(isEven).Take(3).Collect()
	if !slices.Equal(got, []int{2, 4, 6}) || pulled != 6 {
		t.Errorf("got %v after pulling %d values, expected [2 4 6] after 6", got, pulled)
	}

	pulled = 0
	if v, ok := source.Find(func(v int) bool { return v > 2 }); !ok || v != 3 || pulled != 3 {
		t.Errorf("Find = %d, %t after pulling %d values", v, ok, pulled)
	}
}

func TestSeqStopsEarly(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6)
	seqs := map[string]Seq[int]{
		"Filter":    s.Filter(func(int) bool { return true }),
		"Map":       s.Map(func(v int) int { return v }),
		"FilterMap": s.FilterMap(func(v int) (int, bool) { return v, true }),
		"FlatMap":   s.FlatMap(func(v int) iter.Seq[int] { return Of(v, v).Iter() }),
		"Take":      s.Take(5),
		"Skip":      s.Skip(1),
		"UniqBy":    s.UniqBy(func(v int) int { return v }),
	}
	for name, seq := range seqs {
		t.Run(strings.ToLower(name), func(t *testing.T) {
			n := 0
			for range seq {
				if n++; n == 2 {
					break
				}
			}
			if n != 2 {
				t.Errorf("stopped after %d values", n)
			}
		})
	}

	chunks := 0
	for range s.Chunk(2) {
		chunks++
		break
	}
	if chunks != 1 {
		t.Errorf("Chunk did not stop: %d", chunks)
	}
	indexed := 0
	for i := range s.Indexed() {
		if indexed = i; i == 1 {
			break
		}
	}
	if indexed != 1 {
		t.Errorf("Indexed did not stop: %d", indexed)
	}
}

func TestSeqTerminals(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)

	if got := s.Filter(isEven).Count(); got != 2 {
		t.Errorf("Count = %d", got)
	}
	if v, ok := s.First(); !ok || v != 1 {
		t.Errorf("First = %d, %t", v, ok)
	}
	if _, ok := Of[int]().First(); ok {
		t.Error("First on empty Seq must report false")
	}
	if _, ok := s.Find(func(v int) bool { return v > 9 }); ok {
		t.Error("Find found a missing value")
	}
	if !s.Any(isEven) || s.Any(func(v int) bool { return v > 9 }) || Of[int]().Any(isEven) {
		t.Error("Any mismatch")
	}
	if !s.Every(func(v int) bool { return v > 0 }) || s.Every(isEven) || !Of[int]().Every(isEven) {
		t.Error("Every mismatch")
	}

	groups := s.GroupBy(isEven)
	if !slices.Equal(groups[true], []int{2, 4}) || !slices.Equal(groups[false], []int{1, 3, 5}) {
		t.Errorf("GroupBy = %v", groups)
	}
	byLen := Of("a", "bb", "cc").KeyBy(func(s string) int { return len(s) })
	if byLen[1] != "a" || byLen[2] != "cc" {
		t.Errorf("KeyBy = %v", byLen)
	}

	for i, v := range s.Indexed() {
		if v != i+1 {
			t.Errorf("Indexed yielded (%d, %d)", i, v)
		}
	}
	if got := slices.Collect(s.Iter()); len(got) != 5 {
		t.Errorf("Iter = %v", got)
	}
	if got := FromSeq(maps.Keys(map[string]int{"k": 1})).Collect(); !slices.Equal(got, []string{"k"}) {
		t.Errorf("FromSeq = %v", got)
	}
}

func TestSeqChunk(t *testing.T) {
	chunks := slices.Collect(Of(1, 2, 3, 4, 5).Chunk(2))
	if len(chunks) != 3 || !slices.Equal(chunks[0], []int{1, 2}) || !slices.Equal(chunks[2], []int{5}) {
		t.Fatalf("Chunk = %v", chunks)
	}
	// Chunks are independent: appending to one must not overwrite the next.
	chunks[0] = append(chunks[0], 99)
	if chunks[1][0] != 3 {
		t.Fatalf("append to chunk overwrote neighbor: %v", chunks)
	}
	if got := slices.Collect(Of[int]().Chunk(3)); len(got) != 0 {
		t.Fatalf("Chunk(empty) = %v", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Chunk(size 0) did not panic")
		}
	}()
	Of(1).Chunk(0)
}
