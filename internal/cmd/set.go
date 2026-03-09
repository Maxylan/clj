package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

const cmd_set = "Set some field(s) on Ticket(s)"

func register_set(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_set,
			program.Name,
			program.Version,
		),
		Match: match_set,
		Execute: set_field_on_tickets,
		Details: CommandDetails{
			Name:			cmd_set,
			Usage:			fmt.Sprintf("%s set <Field> <Value> on <Tickets...>", program.Name),
			Description:	"Update some field on ticket(s), see subcommands. Ex. `clj set status \"Done\" on PROJ-1337`",
			Subcommands: []CommandDetails{
				{
					Name:			"Field: status",
					Usage:			fmt.Sprintf("%s set status <Value> on <...>", program.Name),
					Description:	"Set 'status' field on ticket(s). Statuses defined in a project's workflow, see `clj statuses <PROJ-1337>`",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"View Updated Tickets",
					Usage:			fmt.Sprintf("%s set <...> on <...> [-t|--view-tickets]", program.Name),
					Description:	"Print updated ticket(s).",
					Subcommands: 	[]CommandDetails{},
				},
				/*{ // Not implemented
					Name:			"Verbose",
					Usage:			fmt.Sprintf("%s <...> [-v|--verbose]", program.Name),
					Description:	"Give a more verbose output, useful for debugging",
					Subcommands: 	[]CommandDetails{},
				},*/
			},
		},
	};
}

func match_set(chain CommandArgChain) bool {
	if len(chain.Keywords) < 3 || len(chain.TicketIDs) == 0 {
		return false
	}

	return strings.EqualFold(chain.Keywords[0], "s") || strings.EqualFold(chain.Keywords[0], "set")
}

func set_field_on_tickets(chain CommandArgChain) {
	if len(chain.Keywords) < 3 || len(chain.TicketIDs) == 0 {
		log.Fatal("Insufficient num. of arguments provided to `set_field_on_tickets`", chain)
	}

	switch {
	case strings.EqualFold(chain.Keywords[1], "status"):
		set_status(chain, chain.Keywords[2])
	default:
		chain.Args = append(chain.Args, "--help")
		Help(chain, Ansii(Red, "(!)", NoColor, " Invalid field \"", chain.Keywords[1], "\""))
		return
	}

	if slices.Contains(chain.Args, "-t") || slices.Contains(chain.Args, "--view-tickets") {
		view_tickets(CommandArgChain{
			TicketIDs:	chain.TicketIDs,
		})
	}
}

func set_status(chain CommandArgChain, status string) {
	fmt.Println(status)
}
