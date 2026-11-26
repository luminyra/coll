<img align="right" src="https://images.luminyra.dev/luminyra/github/coll/logo.webp" alt="An ochre ticket press pulling ID stubs from an order ticket." width="96" height="96">

# coll

[![Go Reference](https://pkg.go.dev/badge/github.com/luminyra/coll.svg)](https://pkg.go.dev/github.com/luminyra/coll)
[![CI](https://github.com/luminyra/coll/actions/workflows/ci.yml/badge.svg)](https://github.com/luminyra/coll/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/luminyra/coll)](https://goreportcard.com/report/github.com/luminyra/coll)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/doc/devel/release)

coll turns multi-step slice and map processing into short, readable pipelines, such as "the IDs of the first ten paid orders" or "orders grouped by customer", without hand-writing the loops each time.

Its core is `Seq`, a lazy sequence: each step runs only when a value is pulled, so a pipeline that needs ten results stops after finding ten. The standard `slices` and `maps` packages cover the basics; coll adds what they lack.

- **Lazy pipelines.** `Filter`, `Map`, `Take` and friends do no work until `Collect` or a `range` loop asks for values.
- **Type-changing chains.** `Map`, `FilterMap` and `Fold` are Go 1.27 generic methods, so a pipeline can turn orders into IDs without leaving the chain.
- **Works with any iterator.** Anything that produces an `iter.Seq`, such as `maps.Keys` or `strings.FieldsSeq`, can start a pipeline, and a `Seq` can be ranged over directly.
- **One callback shape.** Callbacks receive just the element. No unused `index` parameters.
- **Close to a hand-written loop.** A pipeline that stops early runs within 1.5× of the equivalent loop.

<p align="center">
  <img src="https://images.luminyra.dev/luminyra/github/coll/cover.webp" alt="A ticket press separates numbered stubs from a strip of illustrated tickets." width="640">
</p>

## Contents

- [When to use it](#when-to-use-it) · [Install](#install) · [Getting started](#getting-started)
- [Common tasks](#common-tasks) · [How it works](#how-it-works) · [Performance](#performance) · [API at a glance](#api-at-a-glance)

## When to use it

coll is a good fit when a transformation has several steps that would otherwise need a loop with nested `if`s and helper variables: filtering and mapping records, grouping, de-duplicating, set operations, or pulling the first few matches out of a large input.

Keep a plain `for` loop for one-step transformations and for code on a hot path; see [Performance](#performance). And use the standard library first: `slices.Sort`, `slices.Contains`, `slices.Chunk` or `maps.Keys` need no help.

## Install

```sh
go get github.com/luminyra/coll
```

coll requires Go 1.27 or later and has no dependencies.

## Getting started

Copy the program into `main.go` and run it with `go run .`:

```go
package main

import (
	"fmt"

	"github.com/luminyra/coll"
)

type Order struct {
	ID       int
	Customer string
	Cents    int
	Paid     bool
}

func main() {
	orders := []Order{
		{ID: 1, Customer: "alice", Cents: 1200, Paid: true},
		{ID: 2, Customer: "bob", Cents: 800, Paid: false},
		{ID: 3, Customer: "alice", Cents: 300, Paid: true},
		{ID: 4, Customer: "carol", Cents: 4500, Paid: true},
	}

	// A lazy pipeline: nothing runs until Collect pulls the values.
	paid := coll.From(orders).Filter(func(o Order) bool { return o.Paid })

	ids := paid.Map(func(o Order) int { return o.ID }).Collect()
	total := paid.Fold(0, func(sum int, o Order) int { return sum + o.Cents })

	// Eager helpers for one-shot transforms.
	byCustomer := coll.GroupBy(orders, func(o Order) string { return o.Customer })

	fmt.Println("paid orders:", ids)
	fmt.Println("total cents:", total)
	fmt.Println("alice has", len(byCustomer["alice"]), "orders")
}
```

```text
paid orders: [1 3 4]
total cents: 6000
alice has 2 orders
```

`paid` is a recipe, not a result: `Collect` and `Fold` each run it once. A `Seq` built from a slice can be run any number of times.

<p align="center">
  <a href="https://images.luminyra.dev/luminyra/github/coll/demo.webp"><img src="https://images.luminyra.dev/luminyra/github/coll/demo.webp" alt="From four orders, Collect retrieves IDs one and three and stops before pulling the fourth order." width="760"></a>
</p>

## Common tasks

### Filter and transform in one pass

```go
emails := coll.From(users).
	Filter(func(u User) bool { return u.Active }).
	Map(func(u User) string { return u.Email }).
	Collect()
```

For a single step on a whole slice, the eager functions read more simply: `coll.Filter(users, isActive)` or `coll.Map(users, userEmail)`.

### Stop as soon as you have enough

`Take`, `First`, `Find`, `Any` and `Every` stop pulling values once they know the answer, even from an infinite source:

```go
firstTen := coll.From(orders).Filter(isLarge).Take(10).Collect()
admin, found := coll.From(users).Find(func(u User) bool { return u.Admin })
```

### Start from any iterator

```go
words := coll.FromSeq(strings.FieldsSeq(text)).UniqBy(strings.ToLower).Collect()

for name := range coll.FromSeq(maps.Keys(scores)).Filter(isVIP) {
	fmt.Println(name)
}
```

`FlatMap` expands each value into an iterator, so `coll.From(lines).FlatMap(strings.FieldsSeq)` yields every word of every line.

### Use the position of each value

Callbacks receive only the element. When the position matters, iterate with `Indexed`:

```go
for i, step := range coll.From(steps).Indexed() {
	fmt.Printf("%d. %s\n", i+1, step)
}
```

### Group, index and count

```go
byTeam := coll.GroupBy(users, func(u User) string { return u.Team }) // map[string][]User
byID := coll.KeyBy(users, func(u User) int { return u.ID })          // map[int]User
votes := coll.Counts(ballots)                                        // map[string]int
active, inactive := coll.Partition(users, func(u User) bool { return u.Active })
```

### Treat slices as sets

```go
retained := coll.Intersect(lastWeek, thisWeek) // in both, without duplicates
churned := coll.Difference(lastWeek, thisWeek) // in lastWeek only
everyone := coll.Union(lastWeek, thisWeek)     // all, in order of first appearance
```

### Transform maps

```go
prices := coll.MapValues(cents, func(c int) string { return fmt.Sprintf("$%d.%02d", c/100, c%100) })
upper := coll.MapKeys(scores, strings.ToUpper)
merged := coll.Merge(defaults, overrides) // later maps win; inputs are not modified
```

## How it works

- `Seq[T]` is defined as `iter.Seq[T]` with methods, so a `Seq` works in `range` loops directly, and `Iter` converts it for functions that expect an `iter.Seq`, such as `slices.Collect`.
- Intermediate steps (`Filter`, `Map`, `Take`, `Skip`, `UniqBy`…) return a new `Seq` and do no work. Terminal steps (`Collect`, `Fold`, `Count`, `First`, `Find`, `Any`, `Every`, `GroupBy`, `KeyBy`) run the pipeline.
- A `Seq` from a slice can be run repeatedly; a `Seq` from a single-use iterator, such as one reading from a network stream, can be run once.
- No function modifies its input. Eager functions keep named slice and map types: filtering an `IDs` slice returns an `IDs`.

## Performance

A lazy pipeline calls a function per element and step, so it costs a little more than a hand-written loop; unlike eager helpers, it stops as soon as it has enough.

Finding the IDs of the first 100 matching orders among 10,000 (`BenchmarkFirstPaidIDs` in [`seq_test.go`](seq_test.go)). Measured with the benchmarks in this repository (`make bench`; Go 1.27.4, linux/arm64, 5 CPUs, median of 6 runs). Run them on your hardware for numbers that matter to you.

| Approach | Time | Allocated |
| --- | --- | --- |
| Hand-written loop | 0.36 µs | 896 B in 1 allocation |
| `coll.From(orders).Filter(…).Map(…).Take(100).Collect()` | 0.52 µs | 1.9 KB in 5 allocations |
| `coll.Map(coll.Filter(orders, …), …)[:100]` (eager) | 29 µs | 267 KB in 2 allocations |

`Collect` starts with room for a few values, so short results do not reallocate repeatedly. Eager helpers process the whole input; prefer a `Seq` when you only need part of the result.

## API at a glance

| Area | API |
| --- | --- |
| Start a pipeline | `From`, `Of`, `FromSeq` |
| Transform lazily | `Filter`, `Map`, `FilterMap`, `FlatMap`, `Take`, `Skip`, `UniqBy`, `Chunk`, `Indexed` |
| Finish a pipeline | `Collect`, `Fold`, `Count`, `First`, `Find`, `Any`, `Every`, `GroupBy`, `KeyBy`, `Iter` |
| Slices | `Map`, `Filter`, `Partition`, `GroupBy`, `KeyBy`, `ToMap`, `Uniq`, `UniqBy`, `Counts` |
| Sets | `Intersect`, `Difference`, `Union` |
| Maps | `MapKeys`, `MapValues`, `Invert`, `Merge` |

The full reference with runnable examples is on [pkg.go.dev](https://pkg.go.dev/github.com/luminyra/coll). Complete programs live in [`examples/`](examples): [`slices`](examples/slices/main.go) and [`maps-sets`](examples/maps-sets/main.go).

## Versioning

coll follows [semantic versioning](https://semver.org). Until v1.0.0, a minor release may change the API; every change is listed in the [changelog](CHANGELOG.md). It supports the two most recent Go releases, starting with Go 1.27.

## Contributing

Questions, bug reports and pull requests are welcome. See the [contributing guide](https://github.com/luminyra/.github/blob/main/CONTRIBUTING.md), and report security issues as described in the [security policy](https://github.com/luminyra/.github/blob/main/SECURITY.md).

coll was inspired by Lodash, [samber/lo](https://github.com/samber/lo) and Rust iterators. It is released under the [MIT License](LICENSE).
