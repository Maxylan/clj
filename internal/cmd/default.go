package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

const cmd_default = "View Ticket(s)"

func register_default(
	args	[]string,
	program	ProgramDetails,
) Command {
	var input []string

	if len(args) > 0 {
		input = Filter(args[1:], func(arg string) bool {
			return strings.Contains(arg, "--")
		})
	}

	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_default,
			program.Name,
			program.Version,
		),
		Match: len(input) > 1 && IsValidTicketName(input...),
		Execute: func() { View(args, input) },
		Details: CommandDetails{
			Name:			cmd_default,
			Usage:			fmt.Sprintf("%s <PROJ-1337> <...>", program.Name),
			Description:	"\"Retrieves each given ticket, printing their title + description.\"",
			Subcommands: []CommandDetails{
				{
					Name:			"Detailed view",
					Usage:			fmt.Sprintf("%s <...> [-d|--detailed]", program.Name),
					Description:	"Include as many details as possible",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Verbose",
					Usage:			fmt.Sprintf("%s <...> [-v|--verbose]", program.Name),
					Description:	"Give a more verbose output, useful for debugging",
					Subcommands: 	[]CommandDetails{},
				},
			},
		},
	};
}

/** Matches given name(s) against simple RegEx pattern to determine if it/they are valid Jira ticket names */
func IsValidTicketName(names ...string) bool {
	pattern := regexp.MustCompile(`\w+-\d+`)

	for _, n := range names {
		if !pattern.Match([]byte(n)) {
			return false;
		}
	}

	return true
}

func View(args []string, tickets []string) {
	if len(tickets) < 1 {
		fmt.Println("No tickets given", args)
	}
}
