package cmd

import "strings"

// Array Utilities

/** Filters `source` {[]T}, returning a new slice of items matching predicate `pred` */
func Filter[T any](source []T, pred func(T, int) bool) []T {
	var out []T

	for i, e := range source {
		if pred(e, i) {
			out = append(out, e)
		}
	}

	return out
}

/** Filters `source` {[]T}, appending items matching predicate `pred` to destination `dest` {[]T} */
func AppendFilter[T any](source []T, dest *[]T, pred func(T, int) bool) {
	for i, e := range source {
		if pred(e, i) {
			*dest = append(*dest, e)
		}
	}
}

/** Maps `source` {[]TIn}, returning a new slice of items w/ predicate `pred` applied {[]TOut} */
func Map[TIn any, TOut any](source []TIn, pred func(TIn, int) TOut) []TOut {
	var out []TOut

	for i, e := range source {
		out = append(out, pred(e, i))
	}

	return out
}

/** Maps `source` {[]TIn}, returning a new slice of items w/ predicate `pred` applied {[]TOut} */
func FlatMap[TIn any, TOut any](source []TIn, pred func(TIn, int) []TOut) []TOut {
	var out []TOut

	for i, e := range source {
		out = append(out, pred(e, i)...)
	}

	return out
}

/** Maps `source` {[]TIn}, appending items w/ predicate `pred` applied to destination `dest` {[]TOut} */
func AppendMap[TIn any, TOut any](source []TIn, dest *[]TOut, pred func(TIn, int) TOut) {
	for i, e := range source {
		*dest = append(*dest, pred(e, i))
	}
}

/** Maps `source` {[]TIn}, returning a new slice of items matching predicate `p1` w/ predicate `p2` applied {[]TOut} */
func FilterMap[TIn any, TOut any](source []TIn, p1 func(TIn, int) bool, p2 func(TIn, int) TOut) []TOut {
	var out []TOut

	for i, e := range source {
		if p1(e, i) {
			out = append(out, p2(e, i))
		}
	}

	return out
}

/** Same as `slices.Index`, but takes a slice of `string` and compares w/ `strings.EqualFold` */
func IndexEqualFold(source []string, needle string) int {
	for i, v := range source {
		if strings.EqualFold(v, needle) {
			return i
		}
	}

	return -1
}
