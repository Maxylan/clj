package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

const cmd_users = "Search / Find matching users"

func register_users(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_users,
			program.Name,
			program.Version,
		),
		Match: match_users,
		Execute: view_users,
		Details: CommandDetails{
			Name:			cmd_users,
			Usage:			fmt.Sprintf("%s [u|user|users] \"<Value>\"", program.Name),
			Description:	"Search for users matching <Value>, remember to use quotes if your search contains whitespace characters.",
			Subcommands: []CommandDetails{
				{
					Name:			"Detailed view",
					Usage:			fmt.Sprintf("%s <...> [-d|--detailed]", program.Name),
					Description:	"Include as many details as possible",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Exact / Strict matching",
					Usage:			fmt.Sprintf("%s <...> [-e|--exact]", program.Name),
					Description:	"Force exact / strict matching.",
					Subcommands: 	[]CommandDetails{},
				},
				{
					Name:			"Include Deleted Users",
					Usage:			fmt.Sprintf("%s <...> [-i|--include-deleted]", program.Name),
					Description:	"Override default behaviour that filters out deleted users.",
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

func match_users(chain CommandArgChain) bool { 
	if len(chain.Keywords) < 2 {
		return false
	}

	alias1 := strings.EqualFold(chain.Keywords[0], "u")
	alias2 := alias1 || strings.EqualFold(chain.Keywords[0], "user")
	return alias2 || strings.EqualFold(chain.Keywords[0], "users")
}

func view_users(chain CommandArgChain) {
	searchTerm := chain.Keywords[1]

	if searchTerm[0] == '@' {
		searchTerm = searchTerm[1:]
	}

	users := get_matching_users(searchTerm)

	if slices.Contains(chain.Args, "-e") || slices.Contains(chain.Args, "--exact") {
		users = filter_users_exact_match(users, searchTerm)
	}

	if !slices.Contains(chain.Args, "-i") && !slices.Contains(chain.Args, "--include-deleted") {
		users = filter_deleted_users(users)
	}

	slices.Reverse(users)

	render_users(users, searchTerm, chain.Args)
}

func get_matching_users(searchTerm string) []JiraUser {
	if len(searchTerm) < 1 {
		log.Fatal("No search query given", searchTerm)
	}

	users, get_users_err := get_matching_users_search(searchTerm)

	if get_users_err != nil {
		log.Fatal("Failed to search for users matching ", searchTerm, " ", get_users_err.Error())
	} else if users == nil || len(*users) == 0 {
		log.Fatalf(
			"%s(!)%s Found no users matching '%s'! %s%sMaybe try to simplify your search?%s\n",
			Red, Reset, searchTerm,
			Dim, Italic, Reset,
		)
	}

	return *users
}

func filter_deleted_users(users []JiraUser) []JiraUser {
	return Filter(users, func(u JiraUser, _ int) bool {
		return !u.Deleted
	})
}

func filter_users_exact_match(users []JiraUser, searchTerm string) []JiraUser {
	return Filter(users, func(u JiraUser, _ int) bool {
		return u.Name == searchTerm || strings.EqualFold(searchTerm, u.DisplayName) || u.Key == searchTerm
	})
}

func render_users(users []JiraUser, searchTerm string, args []string) {
	// [-d|--detailed] - Include more details
	detailed := slices.Contains(args, "-d") || slices.Contains(args, "--detailed")

	maxLen := 0

	for _, u := range users {
		if nl := len(u.DisplayName); nl > maxLen {
			maxLen = nl
		}
	}

	options := Map(
		users,
		func(u JiraUser, _ int) string {
			if pad := maxLen - len(u.DisplayName); pad > 0 {
				u.DisplayName += strings.Repeat(" ", pad)
			}

			return FormatUser(u, detailed)
		},
	)

	fmt.Println("» Active Users", Ansii(Dim, Italic, "matching '", searchTerm, "'", Reset, ":"))
	fmt.Println(strings.Join(options, "\n"))
}
