package cmd

import (
	"fmt"
	"log"
	"sync"
	"slices"
	"strings"
)

const cmd_comment = "Comment on Ticket(s)"

func register_comment(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_comment,
			program.Name,
			program.Version,
		),
		Match: match_comment,
		Execute: func(chain CommandArgChain) { comment_on_tickets(chain) },
		Details: CommandDetails{
			Name:			cmd_comment,
			Usage:			fmt.Sprintf("%s comment \"Lorem ipsum dolor..\" on <Tickets...>", program.Name),
			Description:	"Creates a new comment on each given ticket. Outputs their comment sections.",
			Subcommands: []CommandDetails{
				{
					Name:			"Oldest First",
					Usage:			fmt.Sprintf("%s <...> [--oldest-first]", program.Name),
					Description:	"Change default sort-order to show the oldest comments first.",
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

func match_comment(chain CommandArgChain) bool {
	if len(chain.Keywords) < 2 || len(chain.TicketIDs) == 0 {
		return false
	}

	return strings.EqualFold(chain.Keywords[0], "c") || strings.EqualFold(chain.Keywords[0], "comment")
}

func comment_on_tickets(chain CommandArgChain) {
	if len(chain.Keywords) < 2 {
		log.Fatal("Insufficient num. of arguments provided to `comment_on_tickets`", chain)
	}
	if len(chain.TicketIDs) < 1 {
		fmt.Println("No tickets given", chain.Args)
		return
	}

	comment := NewComment{
		Body: MarshalJiraMarkdown(chain.Keywords[1]),
	}

	if slices.Contains(chain.Args, "-v") || slices.Contains(chain.Args, "--verbose") {
		fmt.Println(Ansii(Dim, Italic, "» Posting comment '", comment.Body, "' on ticket(s) ", strings.Join(chain.TicketIDs, ", ")))
	}

	if len(chain.TicketIDs) == 1 {
		success, postCommentErr := post_ticket_comment(chain.TicketIDs[0], comment)

		if !success || postCommentErr != nil {
			log.Fatal(Ansii(Bold, Red, "(!)", NoBold, " Failed", Reset, " to posted comment on ticket ", Underline, Cyan, chain.TicketIDs[0]), postCommentErr)
		} else {
			fmt.Println("» Successfully posted comment on ticket", Ansii(Underline, Cyan, chain.TicketIDs[0]))
		}
	} else {
		ch := make(chan string, len(chain.TicketIDs))
		var wg sync.WaitGroup

		for _, id := range chain.TicketIDs {
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

	viewTicketArgs := []string{ "--only-comments" }

	if slices.Contains(chain.Args, "--oldest-first") {
		viewTicketArgs = append(viewTicketArgs, "--oldest-first")
	}

	view_tickets(CommandArgChain{
		TicketIDs:	chain.TicketIDs,
		Args: 		viewTicketArgs,
	})
}
