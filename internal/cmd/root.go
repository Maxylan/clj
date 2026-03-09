package cmd

var registry Commands

/**
 * Register each defined command, passing details about the program like version & name.
 * 
 * Each registered command will provide..
 *  - Details about itself, for help / documentation.
 *  - Expose a `Match(..)` function returns if an arg-chain matches (should invoke) this command
 *  - Expose a `Execute(..)` function to invoke this command
 */
func BuildCommandRegistry(
	details	ProgramDetails,
) Commands {
	SetConfigPath(details.Name)

    registry = Commands{
		register_help(details),
		register_init(details),
		register_set(details),
		register_comment(details),
		register_default(details),
	}

	return registry
}

func CommandNotFound(message string) {
	Help(
		CommandArgChain{
			TicketIDs: []string{},
			Args: []string{ "--minimal" },
		},
		Ansii(Red, "(!)", NoColor, " ", message),
	)
}
