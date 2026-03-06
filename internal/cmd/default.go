package cmd

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"
)

const cmd_default = "View Ticket(s)"

func register_default(
	args	[]string,
	program	ProgramDetails,
) Command {
	var ticketIds []string

	if len(args) > 0 {
		ticketIds = Filter(args[1:], func(arg string) bool {
			return arg[0] != '-' && !strings.Contains(arg, "--")
		})
	}

	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_default,
			program.Name,
			program.Version,
		),
		Match: len(ticketIds) > 0 && IsValidTicketName(ticketIds...),
		Execute: func() { ViewTickets(ticketIds, args) },
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
					Name:			"Include comments",
					Usage:			fmt.Sprintf("%s <...> [-c|--comments]", program.Name),
					Description:	"Render the whole comment section",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Only comments",
					Usage:			fmt.Sprintf("%s <...> [-o|--only-comments]", program.Name),
					Description:	"Render *only* the comment section",
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

func ViewTickets(ticketIds []string, args []string) {
	if len(ticketIds) < 1 {
		fmt.Println("No tickets given", args)
	}

	if len(ticketIds) == 1 {
		ticket, getTicketErr := getTicket(ticketIds[0])

		if getTicketErr != nil {
			log.Fatal("Could not get ticket ", ticketIds[0], " ", getTicketErr.Error())
		}

		Render(ticket, args)
		return
	}

	ch := make(chan *Ticket, len(ticketIds))
	var wg sync.WaitGroup

	for _, id := range ticketIds {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ticket, err := getTicket(id)

			if err != nil {
				fmt.Println("Could not get ticket", id, err.Error())
			}

			ch <- ticket
		}(id)
	}

	wg.Wait()
	close(ch)

	for ticket := range ch {
		Render(ticket, args)
		fmt.Println("")
	}
}

func Render(ticket *Ticket, args []string) {
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
			"  " + Ansii(Cyan, strings.Repeat("─", len(formatted.Link) - 2)),
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
