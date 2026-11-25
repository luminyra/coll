# coll

[![Go Reference](https://pkg.go.dev/badge/github.com/luminyra/coll.svg)](https://pkg.go.dev/github.com/luminyra/coll)
[![CI](https://github.com/luminyra/coll/actions/workflows/ci.yml/badge.svg)](https://github.com/luminyra/coll/actions/workflows/ci.yml)

Lazy, chainable sequences and a small set of slice, set and map helpers for Go 1.27, built to complement `slices`, `maps` and `iter`.

```go
names := coll.From(users).
	Filter(func(u User) bool { return u.Active }).
	Map(func(u User) string { return u.Name }). // Seq[User] becomes Seq[string]
	Take(10).
	Collect()
```

## Install

```sh
go get github.com/luminyra/coll
```

Requires Go 1.27 or later. The module has no dependencies.

## Highlights

- **Lazy sequences.** `Seq[T]` is an `iter.Seq[T]` with methods. Nothing runs until `Collect`, a terminal operation or a `range` loop pulls values, and `Take`, `Find` and `Any` stop pulling as soon as they can.
- **Type-changing chains.** `Map`, `FilterMap`, `FlatMap`, `Fold`, `GroupBy` and `KeyBy` are Go 1.27 generic methods, so a chain can turn users into names without leaving the pipeline.
- **One callback shape.** Callbacks receive only the element; use `Seq.Indexed` when the position matters. Map helpers receive only the key or only the value.
- **Standard library first.** coll does not duplicate `slices.Concat`, `slices.Chunk`, `slices.ContainsFunc`, `maps.Keys`, `cmp.Or` or `new(v)`; it adds only what they lack.

## Usage

```go
// Any iterator becomes a Seq.
words := coll.FromSeq(strings.FieldsSeq(text)).
	Filter(func(w string) bool { return len(w) > 3 }).
	UniqBy(strings.ToLower)

for i, w := range words.Indexed() {
	fmt.Println(i, w)
}

// Eager helpers for one-shot transforms.
byTeam := coll.GroupBy(users, func(u User) string { return u.Team })
active, inactive := coll.Partition(users, func(u User) bool { return u.Active })

// Slices as sets.
churned := coll.Difference(lastWeek, thisWeek)
```

| Area | API |
| --- | --- |
| Lazy | `From`, `Of`, `FromSeq`; `Filter`, `Map`, `FilterMap`, `FlatMap`, `Take`, `Skip`, `UniqBy`, `Chunk`, `Indexed`; `Collect`, `Fold`, `Count`, `First`, `Find`, `Any`, `Every`, `GroupBy`, `KeyBy` |
| Slices | `Map`, `Filter`, `Partition`, `GroupBy`, `KeyBy`, `ToMap`, `Uniq`, `UniqBy`, `Counts` |
| Sets | `Intersect`, `Difference`, `Union` |
| Maps | `MapKeys`, `MapValues`, `Invert`, `Merge` |

## Examples

Runnable programs live in [`examples/`](examples):

| Example | Shows |
| --- | --- |
| [`slices`](examples/slices/main.go) | a lazy pipeline over orders, grouping, batching |
| [`maps-sets`](examples/maps-sets/main.go) | map transforms and set operations |

```sh
go run ./examples/slices
```

## Design

coll is inspired by Lodash, [samber/lo](https://github.com/samber/lo) and Rust iterators. It targets Go 1.27 only and is built around iterators and generic methods instead of eager, index-aware helpers: it stays small by leaving everything the standard library already covers to `slices`, `maps` and `cmp`.

## License

[MIT](LICENSE)
