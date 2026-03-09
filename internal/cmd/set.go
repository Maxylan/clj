package cmd

import (
	"fmt"
	"log"
	"sync"
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

	if slices.Contains(chain.Args, "-t") || slices.Contains(chain.Args, "--view-ticket") || slices.Contains(chain.Args, "--view-tickets") {
		view_tickets(CommandArgChain{
			TicketIDs:	chain.TicketIDs,
		})
	}
}

func set_status(chain CommandArgChain, newStatus string) []Ticket {
	if len(chain.TicketIDs) < 1 {
		log.Fatal("No tickets given", chain)
	}

	tickets := get_tickets(chain.TicketIDs)
	if len(tickets) < 1 {
		log.Fatal("Could not get tickets", chain)
	}

	ticketIDs := Map(
		tickets,
		func(ticket Ticket, _ int) string {
			return ticket.Key
		},
	)

	// Includes statuses..
	transitions := get_transitions_for_tickets(ticketIDs)

	for _, ticket := range tickets {
		ticketTransitions, exists := transitions[ticket.Key]
		transitionId := ""

		if !exists || len(ticketTransitions) == 0 {
			fmt.Println(ticketTransitions, transitions)
			fmt.Printf(
				"%s(!)%s Could not get available transitions %s%s(statuses)%s for ticket '%s%s%s'. %sSkipping..%s\n",
				Red, Reset, Dim, Italic, Reset, Cyan, ticket.Key, Reset, Bold, Reset,
			)
			continue
		}

		for _, t := range ticketTransitions {
			if strings.EqualFold(t.Name, newStatus) || strings.EqualFold(t.Id, newStatus) || strings.EqualFold(t.Status.Name, newStatus) {
				transitionId = t.Id
				break;
			}
		}

		if len(transitionId) == 0 {
			options := Map(
				ticketTransitions,
				func(t IssueTransition, i int) [3]string {
					return [3]string{ fmt.Sprintf("%d", i), t.Name, t.Description }
				},
			)

			log.Fatalf(
				"%s(!)%s Status %s is not a viable transition for '%s%s%s%s'.\n    %s%s(Tip: Run `clj statuses <Tickets...>` to get available statuses)%s\n» %sAvailable Options:%s\n\n%v\n",
				Red, Reset, newStatus,
				Cyan, Underline, ticket.Key, Reset,
				Dim, Italic, Reset,
				Bold, Reset, options,
			)
		}

		success, postErr := post_ticket_transition(ticket.Key, transitionId)
		if !success || postErr != nil {
			log.Fatalf(
				"%s(!)%s Could not POST status transition %s for '%s%s%s%s'.\n» %sError:%s\n\n%v\n",
				Red, Reset, newStatus,
				Cyan, Underline, ticket.Key, Reset,
				Bold, Reset, postErr,
			)
		} else {
			fmt.Println("» Successfully updated status of ticket", Ansii(Underline, Cyan, ticket.Key))
		}
	}

	return tickets
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
