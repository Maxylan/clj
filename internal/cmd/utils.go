package cmd

import (
	"strings"
)

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

const (
    // Reset
    Reset = "\033[0m"

    // Regular colors
    Black   = "\033[30m"
    Red     = "\033[31m"
    Green   = "\033[32m"
    Yellow  = "\033[33m"
    Blue    = "\033[34m"
	BBlue	= "\033[94m"
    Magenta = "\033[35m"
    Cyan    = "\033[36m"
    White   = "\033[37m"
    Gray	= "\033[90m"
    NoColor	= "\033[39m"

    // Bold
    BoldBlack   = "\033[1;30m"
    BoldRed     = "\033[1;31m"
    BoldGreen   = "\033[1;32m"
    BoldYellow  = "\033[1;33m"
    BoldBlue    = "\033[1;34m"
    BoldMagenta = "\033[1;35m"
    BoldCyan    = "\033[1;36m"
    BoldWhite   = "\033[1;37m"
    BoldGray	= "\033[1;90m"
    BoldNoColor	= "\033[1;39m"

    // Background
    BgBlack   = "\033[40m"
    BgRed     = "\033[41m"
    BgGreen   = "\033[42m"
    BgYellow  = "\033[43m"
    BgBlue    = "\033[44m"
    BgMagenta = "\033[45m"
    BgCyan    = "\033[46m"
    BgWhite   = "\033[47m"

    // Styles
    Bold      = "\033[1m"
    Dim       = "\033[2m"
    Italic    = "\033[3m"
    Underline = "\033[4m"
    Blink     = "\033[5m"
    Reversed  = "\033[7m"
	NoBold      = "\033[22m"
    NoDim       = "\033[22m"
    NoItalic    = "\033[23m"
    NoUnderline = "\033[24m"
    NoBlink     = "\033[25m"
    NoReversed  = "\033[27m"
)

func Ansii(strArr ...string) string {
	return strings.Join(strArr, "") + Reset;
}

/** Matches given name(s) against simple RegEx pattern to determine if it/they are valid Jira ticket names */
func IsValidTicketName(names ...string) bool {
	if len(names) == 0 {
		return false
	}

	for _, n := range names {
		if !reTicketId.Match([]byte(n)) {
			return false;
		}
	}

	return true
}

func isStringArg(arg string) bool {
	if len(arg) < 3 {
		return false
	}

	// Ensure non-empty double-quoted string
	return reDoubleQuoted.Match([]byte(arg))
}
