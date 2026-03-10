package cmd

import "strings"

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
