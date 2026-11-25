// Package coll provides collection helpers for Go 1.27 that complement the
// standard [slices], [maps] and [iter] packages.
//
// # Lazy sequences
//
// The centerpiece is [Seq], an [iter.Seq] with chainable methods. Operations
// such as [Seq.Filter] and [Seq.Map] do no work until a terminal operation,
// such as [Seq.Collect] or a range loop, pulls values. Map, [Seq.FilterMap],
// [Seq.Fold] and [Seq.GroupBy] are generic methods, so the element type can
// change in the middle of a chain:
//
//	names := coll.From(users).Filter(isActive).Map(userName).Take(10).Collect()
//
// # Eager helpers
//
// [Map], [Filter], [Partition], [GroupBy], [KeyBy], [ToMap], [Uniq], [UniqBy]
// and [Counts] cover one-shot slice transforms. [Intersect], [Difference] and
// [Union] treat slices as sets, and [MapKeys], [MapValues], [Invert] and
// [Merge] transform maps.
//
// # Conventions
//
// Callbacks receive only the element; use [Seq.Indexed] when the position
// matters. Map callbacks receive only the key or only the value. Functions
// never modify their input.
//
// coll deliberately leaves out what the standard library already provides:
// use [slices.Concat], [slices.Chunk], [slices.ContainsFunc],
// [slices.IndexFunc], [maps.Keys], [maps.Values], [cmp.Or] and new(v)
// directly.
//
// The package is inspired by Lodash, samber/lo and Rust iterators.
package coll
