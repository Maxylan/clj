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
					Name:			"Field: assignee",
					Usage:			fmt.Sprintf("%s set assignee <User> on <...>", program.Name),
					Description:	"Set 'Assignee' of ticket(s). Picks best user-match. See subcommand `users` for available users. Ex. `clj users`",
					Subcommands: 	[]CommandDetails{
						{
							Name:			"Exact / Strict matching",
							Usage:			fmt.Sprintf("%s <...> [-e|--exact]", program.Name),
							Description:	"Force exact / strict matching in your user search.",
							Subcommands: 	[]CommandDetails{},
						},
						{
							Name:			"Prompt to select",
							Usage:			fmt.Sprintf("%s <...> [-s|--select]", program.Name),
							Description:	"Instead of auto-picking best match, this will prompt you for input. Let's you pick the user.",
							Subcommands: 	[]CommandDetails{},
						},
					},
				},
				{
					Name:			"Field: reporter",
					Usage:			fmt.Sprintf("%s set reporter <User> on <...>", program.Name),
					Description:	"Set 'Reporter' of ticket(s). See subcommand `users` for available users. Ex. `clj users`",
					Subcommands: 	[]CommandDetails{
						{
							Name:			"Exact / Strict matching",
							Usage:			fmt.Sprintf("%s <...> [-e|--exact]", program.Name),
							Description:	"Force exact / strict matching in your user search.",
							Subcommands: 	[]CommandDetails{},
						},
						{
							Name:			"Prompt to select",
							Usage:			fmt.Sprintf("%s <...> [-s|--select]", program.Name),
							Description:	"Instead of auto-picking best match, this will prompt you for input. Let's you pick the user.",
							Subcommands: 	[]CommandDetails{},
						},
					},
				},
				{
					Name:			"Field: status",
					Usage:			fmt.Sprintf("%s set [status|transition] <Value> on <...>", program.Name),
					Description:	"Set 'status' ('transition') on ticket(s). See subcommand `stat` for available transitions. Ex. `clj stat <PROJ-1337>`",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Print Updated Tickets",
					Usage:			fmt.Sprintf("%s set <...> on <...> [-p|--print]", program.Name),
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
	case strings.EqualFold(chain.Keywords[1], "assignee") || strings.EqualFold(chain.Keywords[1], "reporter"):
		set_member(chain, chain.Keywords[2])
	case strings.EqualFold(chain.Keywords[1], "status") || strings.EqualFold(chain.Keywords[1], "transition"):
		set_status(chain, chain.Keywords[2])
	default:
		chain.Args = append(chain.Args, "--help")
		Help(chain, Ansii(Red, "(!)", NoColor, " Invalid field \"", chain.Keywords[1], "\""))
		return
	}

	if slices.Contains(chain.Args, "-p") || slices.Contains(chain.Args, "--print") {
		view_tickets(CommandArgChain{
			TicketIDs:	chain.TicketIDs,
		})
	}
}

func set_member(chain CommandArgChain, searchTerm string) {
	if len(chain.TicketIDs) < 1 {
		log.Fatal("No tickets given", chain)
	}

	var field string
	var value string
	switch {
	case strings.EqualFold(chain.Keywords[1], "assignee"):
		field = "assignee"
	case strings.EqualFold(chain.Keywords[1], "reporter"):
		field = "reporter"
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

	if len(searchTerm) < 2 || strings.EqualFold(searchTerm, "unassign") || strings.EqualFold(searchTerm, "unassigned") {
		value = ""
	} else {
		if searchTerm[0] == '@' {
			searchTerm = searchTerm[1:]
		} 

		// Search for users w/ same autocomplete as Jira's dropdown..
		users := get_matching_users(searchTerm)
		
		if len(users) == 0 {
			log.Fatalf(
				"%s(!)%s Found no users matching '%s'! %s%sMaybe try to simplify your search?%s\n",
				Red, Reset, searchTerm,
				Dim, Italic, Reset,
			)
		}
		
		users = filter_deleted_users(users)
		slices.Reverse(users)

		if slices.Contains(chain.Args, "-e") || slices.Contains(chain.Args, "--exact") {
			users = filter_users_exact_match(users, searchTerm)
		}

		picked := users[0]

		if slices.Contains(chain.Args, "-s") || slices.Contains(chain.Args, "--select") {
			picked = PromptSelect(users, func (u JiraUser) string {
				return Ansii(u.DisplayName, " ", Dim, Italic, "(", u.Name, ")\t", u.Key)
			})
		}

		value = picked.Name
	}

	bodyData := []byte(
		fmt.Sprintf(`{"update":{"%s": [{"set":{"name": "%s"}}]}}`, field, value),
	)

	if slices.Contains(chain.Args, "-v") || slices.Contains(chain.Args, "--verbose") {
		fmt.Println(Ansii(Dim, Italic, string(bodyData)))
	}

	if len(ticketIDs) == 1 {
		success, put_ticket_err := put_ticket_fields(ticketIDs[0], bodyData)

		if !success || put_ticket_err != nil {
			log.Fatal("Could not update ticket ", ticketIDs[0], " ", put_ticket_err.Error())
		}

		fmt.Println("» Successfully updated "+field+" of ticket", Ansii(Underline, Cyan, ticketIDs[0]))
		return
	}

	ch := make(chan *bool, len(ticketIDs))
	var wg sync.WaitGroup

	for _, id := range ticketIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			success, err := put_ticket_fields(id, bodyData)

			if err != nil {
				fmt.Println("Could not update ticket", id, err.Error())
			} else if success {
				fmt.Println("» Successfully updated "+field+" of ticket", Ansii(Underline, Cyan, id))
			}

			ch <- &success
		}(id)
	}

	wg.Wait()
	close(ch)

	for range ch {} // Just drains the channel buffer
}

func set_status(chain CommandArgChain, newStatus string) []string {
	if len(chain.TicketIDs) < 1 {
		log.Fatal("No tickets given", chain)
	}

	// Includes statuses..
	transitions := get_transitions_for_tickets(chain.TicketIDs)

	// TODO! Should use go channels like every other loop w/ multiple HTTP Calls!
	for _, ticket := range chain.TicketIDs {
		ticketTransitions, exists := transitions[ticket]
		transitionId := ""

		if !exists || len(ticketTransitions) == 0 {
			fmt.Println(ticketTransitions, transitions)
			fmt.Printf(
				"%s(!)%s Could not get available transitions %s%s(statuses)%s for ticket '%s%s%s'. %sSkipping..%s\n",
				Red, Reset, Dim, Italic, Reset, Cyan, ticket, Reset, Bold, Reset,
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
				func(t IssueTransition, _ int) string {
					return FormatTransition(t, false)
				},
			)

			log.Fatalf(
				"\n%s(!)%s Status \"%s\" is not a viable transition for '%s%s%s%s'.\n    %s%s(Tip: Run `clj stat <Tickets...>` beforehand to view available transitions)%s\n  » %sAvailable Transitions:%s\n%s\n",
				Red, Reset, newStatus,
				Cyan, Underline, ticket, Reset,
				Dim, Italic, Reset,
				Bold, Reset, strings.Join(options, "\n"),
			)
		}

		if slices.Contains(chain.Args, "-v") || slices.Contains(chain.Args, "--verbose") {
			fmt.Println(Ansii(Dim, Italic, "» Setting transition '", transitionId, "' on '", ticket, "'"))
		}

		success, postErr := post_ticket_transition(ticket, transitionId)
		if !success || postErr != nil {
			log.Fatalf(
				"%s(!)%s Could not POST status transition %s for '%s%s%s%s'.\n» %sError:%s\n\n%v\n",
				Red, Reset, newStatus,
				Cyan, Underline, ticket, Reset,
				Bold, Reset, postErr,
			)
		} else {
			fmt.Println("» Successfully updated status of ticket", Ansii(Underline, Cyan, ticket))
		}
	}

	return chain.TicketIDs
}
