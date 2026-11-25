package coll_test

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/luminyra/coll"
)

func ExampleSeq() {
	isEven := func(n int) bool { return n%2 == 0 }

	labels := coll.Of(1, 2, 3, 4, 5, 6).Filter(isEven).Map(strconv.Itoa).Collect()
	fmt.Println(labels)
	// Output: [2 4 6]
}

func ExampleSeq_FlatMap() {
	words := coll.Of("lazy sequences", "with generic methods").FlatMap(strings.FieldsSeq)
	fmt.Println(words.Count(), words.Collect())
	// Output: 5 [lazy sequences with generic methods]
}

func ExampleSeq_Fold() {
	total := coll.Of(1200, 300, 4500).Fold(0, func(sum, cents int) int { return sum + cents })
	fmt.Println(total)
	// Output: 6000
}

func ExampleSeq_Indexed() {
	for i, name := range coll.Of("alice", "bob").Indexed() {
		fmt.Println(i, name)
	}
	// Output:
	// 0 alice
	// 1 bob
}

func ExampleFromSeq() {
	ages := map[string]int{"carol": 41, "alice": 30, "bob": 25}

	adults := coll.FromSeq(maps.Keys(ages)).Filter(func(name string) bool { return ages[name] >= 30 }).Collect()
	slices.Sort(adults)
	fmt.Println(adults)
	// Output: [alice carol]
}

func ExamplePartition() {
	paid, pending := coll.Partition([]int{120, -30, 45, -5}, func(n int) bool { return n > 0 })
	fmt.Println(paid, pending)
	// Output: [120 45] [-30 -5]
}

func ExampleDifference() {
	lastWeek := []string{"alice", "bob", "carol"}
	thisWeek := []string{"bob", "carol", "dave"}

	fmt.Println("churned:", coll.Difference(lastWeek, thisWeek))
	fmt.Println("joined:", coll.Difference(thisWeek, lastWeek))
	// Output:
	// churned: [alice]
	// joined: [dave]
}
