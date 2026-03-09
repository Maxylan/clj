package cmd

import (
	"fmt"
	"log"
	"sync"
	"slices"
	"strings"
)

const cmd_default = "View Ticket(s)"

func register_default(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_default,
			program.Name,
			program.Version,
		),
		Match: match_default,
		Execute: view_tickets,
		Details: CommandDetails{
			Name:			cmd_default,
			Usage:			fmt.Sprintf("%s <Tickets...>", program.Name),
			Description:	"Retrieves each given ticket, printing their title + description.",
			Subcommands: []CommandDetails{
				{
					Name:			"Detailed view",
					Usage:			fmt.Sprintf("%s <Tickets...> [-d|--detailed]", program.Name),
					Description:	"Include as many details as possible",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Include comments",
					Usage:			fmt.Sprintf("%s <Tickets...> [-c|--comments]", program.Name),
					Description:	"Render the whole comment section",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Only comments",
					Usage:			fmt.Sprintf("%s <Tickets...> [-o|--only-comments]", program.Name),
					Description:	"Render *only* the comment section",
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

func match_default(chain CommandArgChain) bool { 
	return len(chain.Keywords) == 0 && len(chain.TicketIDs) > 0
}

func view_tickets(chain CommandArgChain) {
	if len(chain.TicketIDs) < 1 {
		fmt.Println("No tickets given", chain.Args)
		return
	}

	tickets := get_tickets(chain.TicketIDs)

	for _, ticket := range tickets {
		render_ticket(&ticket, chain.Args)
		fmt.Println("")
	}
}

func get_tickets(ticketIDs []string) []Ticket {
	if len(ticketIDs) < 1 {
		log.Fatal("No tickets given", ticketIDs)
	}

	if len(ticketIDs) == 1 {
		ticket, get_ticket_err := get_ticket(ticketIDs[0])

		if get_ticket_err != nil {
			log.Fatal("Could not get ticket ", ticketIDs[0], " ", get_ticket_err.Error())
		}

		return []Ticket{
			*ticket,
		}
	}

	ch := make(chan *Ticket, len(ticketIDs))
	var wg sync.WaitGroup
	var tickets []Ticket

	for _, id := range ticketIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ticket, err := get_ticket(id)

			if err != nil {
				fmt.Println("Could not get ticket", id, err.Error())
			}

			ch <- ticket
		}(id)
	}

	wg.Wait()
	close(ch)

	for ticket := range ch {
		if ticket != nil {
			tickets = append(tickets, *ticket)
		}
	}

	return tickets
}

func render_ticket(ticket *Ticket, args []string) {
	if ticket == nil {
		fmt.Println("Could not view nil ticket")
		return
	}

	formatted := FormatTicket(*ticket)

	// [-o|--only-comments] - Render only comments
	if slices.Contains(args, "-o") || slices.Contains(args, "--only-comments") {
		fmt.Printf(
			"%s\n%s\n%s\n%s",
			formatted.Headline,
			formatted.Link,
			formatted.DividerShort,
			formatted.Comments,
		)
		return 
	}

	out := fmt.Sprintf(
		"%s\n%s",
		formatted.Headline,
		formatted.Link,
	)

	// [-d|--detailed] - Include more details
	if slices.Contains(args, "-d") || slices.Contains(args, "--detailed") {
		out += "\n" + strings.Join(
			[]string{
				formatted.StatusLong,
				formatted.Priority,
				formatted.IssueType,
				formatted.Creator,
				formatted.Reporter,
				formatted.Assignee,
				formatted.TicketDetails,
				formatted.TicketDates,
				formatted.Divider,
				formatted.Description,
			},
			"\n",
		)
	} else {
		out += "\n" + strings.Join(
			[]string{
				formatted.Status,
				formatted.TypeCombined,
				formatted.Members,
				formatted.TicketDates,
				formatted.Divider,
				formatted.Description,
			},
			"\n",
		)
	}

	// [-c|--comments] - Append comment section
	if slices.Contains(args, "-c") || slices.Contains(args, "--comments") {
		out += "\n" + formatted.Divider + "\n" + formatted.Comments
	}

	fmt.Println(out)
}
