package cmd

import "fmt"

var registry Commands

func BuildCommandRegistry(
	args []string,
	details ProgramDetails,
) Commands {
	fmt.Println(
		Ansii(Bold, details.Name),
		details.Version,
		args,
	)
	fmt.Println() // nl

	registry = Commands{
		register_help(args, details),
		register_init(args, details),
		register_comment(args, details),
		register_default(args, details),
	}

    return registry
}
