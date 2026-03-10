package cmd

import (
	"fmt"
	"strings"
)

var registry Commands

/**
 * Register each defined command, passing details about the program like version & name.
 * 
 * Each registered command will provide..
 *  - Details about itself, for help / documentation.
 *  - Expose a `Match(..)` function returns if an arg-chain matches (should invoke) this command
 *  - Expose a `Execute(..)` function to invoke this command
 */
func BuildCommandRegistry(
	details	ProgramDetails,
) Commands {
	SetConfigPath(details.Name)

    registry = Commands{
		register_help(details),
		register_init(details),
		register_set(details),
		register_transitions(details),
		register_comment(details),
		register_default(details),
	}

	return registry
}

func CommandNotFound(message string) {
	Help(
		CommandArgChain{
			TicketIDs: []string{},
			Args: []string{ "--minimal" },
		},
		Ansii(Red, "(!)", NoColor, " ", message),
	)
}

/**
 * Parses slice of strings (like `os.Args`) into "chains" of arguments,
 * each "chain" representing a single operation.
 *
 * "Arguments" are sorted into buckets of..
 *  - Arguments		(ex. --only-comments)
 *  - Ticket IDs	(ex. PROJ-1337)
 *  - Keywords		(..rest)
 *  ..in the order they're given in the terminal.
 *
 *  Note: '-v|--verbose' argument is ignored. It's expected that the caller
 *  checks if its present and passes the boolean result parameter `verbose`
 */
func ParseArgs(args []string, verbose bool) []CommandArgChain {
	cur := 0
	setTickets := false
	out := []CommandArgChain{
		{
			TicketIDs:	[]string{},
			Keywords:	[]string{},
			Args:		[]string{},
		},
	}

	if verbose {
		out[0].Args = append(out[0].Args, "--verbose")
	}

	for _, arg := range args[1:] {
		isValidTicketName := IsValidTicketName(arg)

		switch {
		case strings.EqualFold(arg, "and"):
			cur++
			setTickets = false
			out = append(out, CommandArgChain{
				TicketIDs:	[]string{},
				Keywords:	[]string{},
				Args:		[]string{},
			})
			if verbose {
				out[cur].Args = append(out[cur].Args, "--verbose")
			}
		case arg[0] == '-':
			if arg != "-v" && arg != "--verbose" {
				out[cur].Args = append(out[cur].Args, arg)
			}
		case setTickets || isValidTicketName:
			if !isValidTicketName {
				fmt.Println(Ansii(
					Red, "(!)", Reset, " ", Italic,
					"Ticket ", Cyan, Underline, arg, NoUnderline,
					NoColor, " potentially poorly formatted. ", Bold, "Skipped!",
				))
				break
			}

			out[cur].TicketIDs = append(out[cur].TicketIDs, arg)
		case strings.EqualFold(arg, "on"):
			setTickets = true
		default:
			out[cur].Keywords = append(out[cur].Keywords, arg)
		}
	}

	return out
}
