// Command maps-sets uses coll's map helpers and set operations on slices
// together with the standard maps package.
package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/luminyra/coll"
)

func main() {
	stock := map[string]int{"apple": 12, "banana": 0, "cherry": 5, "durian": 0}

	// region: maps
	// Filtering a map is covered by the standard library: clone it, then
	// remove entries with maps.DeleteFunc.
	inStock := maps.Clone(stock)
	maps.DeleteFunc(inStock, func(_ string, n int) bool { return n == 0 })

	// MapValues passes only the value to the callback and keeps the keys.
	labels := coll.MapValues(inStock, func(n int) string { return fmt.Sprintf("%d in stock", n) })
	// endregion

	for _, name := range slices.Sorted(maps.Keys(labels)) {
		fmt.Println(name, labels[name])
	}

	// region: sets
	lastWeek := []string{"alice", "bob", "carol", "bob"}
	thisWeek := []string{"bob", "dave", "carol", "erin"}

	retained := coll.Intersect(lastWeek, thisWeek) // active both weeks, deduplicated
	churned := coll.Difference(lastWeek, thisWeek) // active last week only
	joined := coll.Difference(thisWeek, lastWeek)  // new this week
	everyone := coll.Union(lastWeek, thisWeek)     // everyone, in order of first appearance
	// endregion

	fmt.Println("retained:", retained)
	fmt.Println("churned: ", churned)
	fmt.Println("joined:  ", joined)
	fmt.Println("everyone:", everyone)
}
