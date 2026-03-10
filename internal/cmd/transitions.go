package cmd

import (
	"fmt"
	"log"
	"sync"
	"slices"
	"strings"
)

const cmd_transitions = "Available Ticket Transitions (statuses)"

func register_transitions(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_transitions,
			program.Name,
			program.Version,
		),
		Match: match_transitions,
		Execute: view_transitions,
		Details: CommandDetails{
			Name:			cmd_transitions,
			Usage:			fmt.Sprintf("%s [stat|statuses|transitions] on <Tickets...>", program.Name),
			Description:	"Retrieves information about each available transition (status), for each given ticket.",
			Subcommands: []CommandDetails{
				{
					Name:			"Detailed view",
					Usage:			fmt.Sprintf("%s <...> [-d|--detailed]", program.Name),
					Description:	"Include as many details as possible",
					Subcommands: 	[]CommandDetails{},
				},
				/*{ // Not implemented
					Name:			"Verbose",
					Usage:			fmt.Sprintf("%s <Tickets...> [-v|--verbose]", program.Name),
					Description:	"Give a more verbose output, useful for debugging",
					Subcommands: 	[]CommandDetails{},
				},*/
			},
		},
	};
}

func match_transitions(chain CommandArgChain) bool { 
	if len(chain.Keywords) == 0 || len(chain.TicketIDs) == 0 {
		return false
	}

	alias1 := strings.EqualFold(chain.Keywords[0], "stat")
	alias2 := alias1 || strings.EqualFold(chain.Keywords[0], "statuses")
	return alias2 || strings.EqualFold(chain.Keywords[0], "transitions")
}

func view_transitions(chain CommandArgChain) {
	if len(chain.TicketIDs) < 1 {
		fmt.Println("No tickets given", chain.Args)
		return
	}

	available_transitions := get_transitions_for_tickets(chain.TicketIDs)

	for ticket, transitions := range available_transitions {
		render_transitions(ticket, transitions, chain.Args)
		fmt.Println("")
	}
}

func get_transitions_for_tickets(ticketIDs []string) TicketTransitionsMap {
	if len(ticketIDs) < 1 {
		log.Fatal("No Ticket IDs given", ticketIDs)
	}

	if len(ticketIDs) == 1 {
		transitions, get_transitions_err := get_issue_transitions(ticketIDs[0])

		if get_transitions_err != nil {
			log.Fatal("Could not get ticket '", ticketIDs[0], "' transitions. ", get_transitions_err.Error())
		}

		return TicketTransitionsMap{
			ticketIDs[0]: transitions.Transitions,
		}
	}

	ch := make(chan *TicketTransitions, len(ticketIDs))
	ticket_transitions := TicketTransitionsMap{}
	var wg sync.WaitGroup

	for _, ticketID := range ticketIDs {
		if _, exists := ticket_transitions[ticketID]; exists {
			continue
		}

		ticket_transitions[ticketID] = []IssueTransition{}
		wg.Add(1)

		go func(ticketID string) {
			defer wg.Done()
			transitions, err := get_issue_transitions(ticketID)

			if err != nil {
				fmt.Println("Could not get Ticket", ticketID, "Transitions.", err.Error())
			}

			ch <- transitions
		}(ticketID)
	}

	wg.Wait()
	close(ch)

	for result := range ch {
		if result != nil {
			ticket_transitions[result.TicketID] = result.Transitions
		}
	}

	return ticket_transitions
}

func render_transitions(ticketID string, transitions []IssueTransition, args []string) {
	// [-d|--detailed] - Include more details
	detailed := slices.Contains(args, "-d") || slices.Contains(args, "--detailed")

	options := Map(
		transitions,
		func(t IssueTransition, _ int) string {
			return FormatTransition(t, detailed)
		},
	)

	fmt.Println("» Available Transitions -", Ansii(Cyan, Underline, ticketID, Reset, ":"))
	fmt.Println(strings.Join(options, "\n"))
}
