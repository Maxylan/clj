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

	if slices.Contains(chain.Args, "-t") || slices.Contains(chain.Args, "--view-tickets") {
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

	projectIDs := Map(
		tickets,
		func(ticket Ticket, _ int) string {
			return ticket.Fields.Project.Id
		},
	)

	// Includes statuses..
	proj_issue_types := get_issue_types_for_projects(projectIDs)

	for _, ticket := range tickets {
		if len(ticket.Fields.Project.Id) == 0 {
			fmt.Println(Ansii(Red, "(!)"), "Ticket", ticket.Key, "malformed, skipping..")
			continue
		}

		var pickedStatus *TicketStatus
		var available []TicketStatus

		if issueTypes, exists := proj_issue_types[ticket.Fields.Project.Id]; exists {
			for _, issueType := range issueTypes {
				if issueType.Id == ticket.Fields.IssueType.Id {
					available = issueType.Statuses
					break
				}
			}
		} else {
			log.Fatal(
				"Could not get project ",
				ticket.Fields.Project.Id,
				". This should not be possible! ",
				ticket.Key,
			)
		}

		if len(available) == 0 {
			fmt.Println(Ansii(Red, "(!)"), "Could not determine available statuses for Ticket", Ansii(Cyan, ticket.Key), ", skipping..")
			continue;
		}

		for _, status := range available {
			if strings.EqualFold(status.Name, newStatus) || strings.EqualFold(status.Category.Name, newStatus) {
				pickedStatus = &status
				break;
			}
		}

		if pickedStatus == nil {
			log.Fatalf(
				"%s(!)%s Status %s is not a viable option for '%s%s%s%s' %s%s(Project '%s').\n    (Tip: Run `clj statuses <Tickets...>` to get available statuses)%s\n» %sAvailable Options:%s\n\n%v",
				Red, Reset, newStatus,
				Cyan, Underline, ticket.Key, Reset,
				Dim, Italic, ticket.Fields.Project.Id, Reset,
				Bold, Reset, available,
			)
		}
	}

	return tickets
}

func get_issue_types_for_projects(projectIDs []string) ProjectsWithIssueTypesMap {
	if len(projectIDs) < 1 {
		log.Fatal("No Project IDs given", projectIDs)
	}

	if len(projectIDs) == 1 {
		projectIssueTypes, get_statuses_err := get_project_issue_types(projectIDs[0])

		if get_statuses_err != nil {
			log.Fatal("Could not get project '", projectIDs[0], "' statuses. ", get_statuses_err.Error())
		}

		return ProjectsWithIssueTypesMap{
			projectIDs[0]: projectIssueTypes.IssueTypes,
		}
	}

	ch := make(chan *ProjectIssueTypes, len(projectIDs))
	var project_statuses ProjectsWithIssueTypesMap
	var wg sync.WaitGroup

	for _, projectId := range projectIDs {
		if _, exists := project_statuses[projectId]; exists {
			continue
		}

		project_statuses[projectId] = []JiraIssueType{}
		wg.Add(1)

		go func(projectId string) {
			defer wg.Done()
			projectIssueTypes, err := get_project_issue_types(projectId)

			if err != nil {
				fmt.Println("Could not get Project", projectId, "Issue Types.", err.Error())
			}

			ch <- projectIssueTypes
		}(projectId)
	}

	wg.Wait()
	close(ch)

	for result := range ch {
		if result != nil {
			project_statuses[result.ProjectID] = result.IssueTypes
		}
	}

	return project_statuses
}
