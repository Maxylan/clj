package cmd

type ProgramDetails struct {
	Name	string
	Version	string
}

type CommandDetails struct {
	/** Name of this command */
	Name		string
	/** Usage notes */
	Usage		string
	/** Lengthy description */
	Description	string
	/** Available subcommands */
	Subcommands	[]CommandDetails
}

type Command struct {
	/** Full-name (prettyprint) */
	Name	string
	/** Did this {Command} match command-line arguments? */
	Match	bool
	/** Executes this command */
	Execute	func()
	/** Returns details about this command, used by "help" */
	Details	CommandDetails
}

type Commands [3]Command

type Config struct {
    JiraURL string `json:"jira_url"`
    Token   string `json:"token"`
}
