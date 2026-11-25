// Command slices filters, transforms, groups and batches a slice of orders,
// first with a lazy coll.Seq chain and then with eager helpers.
package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/luminyra/coll"
)

type order struct {
	ID       int
	Customer string
	Amount   int // in cents
	Paid     bool
}

func main() {
	orders := []order{
		{ID: 1, Customer: "alice", Amount: 1200, Paid: true},
		{ID: 2, Customer: "bob", Amount: 800, Paid: false},
		{ID: 3, Customer: "alice", Amount: 300, Paid: true},
		{ID: 4, Customer: "carol", Amount: 4500, Paid: true},
		{ID: 5, Customer: "bob", Amount: 2600, Paid: true},
	}
	isPaid := func(o order) bool { return o.Paid }

	// region: seq
	// A Seq is lazy: each step runs only when Collect or a range loop pulls
	// values. Map is a Go 1.27 generic method, so the chain can turn orders
	// into ints or strings.
	paid := coll.From(orders).Filter(isPaid)

	ids := paid.Map(func(o order) int { return o.ID }).Collect()
	total := paid.Fold(0, func(sum int, o order) int { return sum + o.Amount })
	customers := paid.UniqBy(func(o order) string { return o.Customer }).
		Map(func(o order) string { return o.Customer }).
		Collect()
	// endregion

	fmt.Println("paid orders:", ids)
	fmt.Println("total (cents):", total)
	fmt.Println("customers:", customers)

	// region: group
	// GroupBy returns a map, whose iteration order is unspecified; sort the
	// keys for stable output.
	byCustomer := coll.GroupBy(orders, func(o order) string { return o.Customer })
	for _, name := range slices.Sorted(maps.Keys(byCustomer)) {
		fmt.Printf("%s: %d orders\n", name, len(byCustomer[name]))
	}

	settled, pending := coll.Partition(orders, isPaid)
	fmt.Printf("%d paid, %d pending\n", len(settled), len(pending))
	// endregion

	// region: chunk
	// For a slice that is already in memory, the standard slices.Chunk
	// batches it without copying.
	for batch := range slices.Chunk(ids, 2) {
		fmt.Println("batch:", batch)
	}
	// endregion
}
