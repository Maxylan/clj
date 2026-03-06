package cmd

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const cmd_comment = "Comment on Ticket(s)"

func register_comment(
	args	[]string,
	program	ProgramDetails,
) Command {
	var (
		hasArgs		= len(args) > 3
		lastArg		= args[len(args)-1]
		isComment	= false
		ticketIDs	[]string
	)

	if hasArgs {
		isComment = strings.EqualFold(args[1], "c") || strings.EqualFold(args[1], "comment")
		ticketIDs = Filter(args[2:], func(arg string, i int) bool {
			// Prevents last arg from being counted as a ticket ID..
			return (i != len(args) - 3) && arg[0] != '-'
		})
	}

	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_comment,
			program.Name,
			program.Version,
		),
		Match: hasArgs && isComment && IsValidTicketName(ticketIDs...),
		Execute: func() { comment_on_tickets(ticketIDs, NewComment{ Body: lastArg }, args) },
		Details: CommandDetails{
			Name:			cmd_comment,
			Usage:			fmt.Sprintf("%s [c|comment] <PROJ-1337> <...> \"Lorem ipsum dolor..\"", program.Name),
			Description:	"\"Creates a new comment on each given ticket, outputs their comment sections.\"",
			Subcommands: []CommandDetails{
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

func comment_on_tickets(ticketIDs []string, comment NewComment, args []string) {
	if len(ticketIDs) < 1 {
		fmt.Println("No tickets given", args)
		return
	}

	if len(ticketIDs) == 1 {
		success, postCommentErr := post_ticket_comment(ticketIDs[0], comment)

		if !success || postCommentErr != nil {
			log.Fatal(Ansii(Bold, Red, "(!)", NoBold, " Failed", Reset, " to posted comment on ticket ", Underline, Cyan, ticketIDs[0]), postCommentErr)
		} else {
			fmt.Println("» Successfully posted comment on ticket", Ansii(Underline, Cyan, ticketIDs[0]))
		}
	} else {
		ch := make(chan string, len(ticketIDs))
		var wg sync.WaitGroup

		for _, id := range ticketIDs {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				success, err := post_ticket_comment(id, comment)

				if !success || err != nil {
					fmt.Println(Ansii(Bold, Red, "(!)", NoBold, " Failed", Reset, " to posted comment on ticket ", Underline, Cyan, id), err)
					ch <- ""
					return
				}

				ch <- id
			}(id)
		}

		wg.Wait()
		close(ch)

		for ticketID := range ch {
			if len(ticketID) > 0 {
				fmt.Println("» Successfully posted comment on ticket", Ansii(Underline, Cyan, ticketID))
			}
		}
	}

	view_tickets(ticketIDs, []string{ "--only-comments" })
}

