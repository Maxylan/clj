package cmd

import (
	"time"
	"strings"
)

type ProgramDetails struct {
	Name	string
	Version	string
}

type CommandArgChain struct {
	TicketIDs 	[]string
	Keywords	[]string
	Args		[]string
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
	Match	func(CommandArgChain) bool
	/** Executes this command */
	Execute	func(CommandArgChain)
	/** Returns details about this command, used by "help" */
	Details	CommandDetails
}

type Commands [7]Command

type Config struct {
    JiraURL string `json:"jira_url"`
    Token   string `json:"token"`
}

// API Types

type Headers map[string]string

type Ticket struct {
	Key		string			`json:"key"`
	Fields	TicketFields	`json:"fields"`
}

type TicketFields struct {
	Priority	TicketPriority	`json:"priority"`
	Labels		[]string		`json:"labels"`
	Status		TicketStatus	`json:"status"`
	Creator		JiraMember		`json:"creator"`
	Reporter	JiraMember		`json:"reporter"`
	Assignee	JiraMember		`json:"assignee"`
	IssueType	TicketIssueType	`json:"issuetype"`
	Project		JiraProject		`json:"project"`
	Summary		string			`json:"summary"`
	Description	string			`json:"description"`
	Comment		TicketComments	`json:"comment"`
	Created		JiraTime		`json:"created"`
	Updated		JiraTime		`json:"updated"`
}

type TicketPriority struct {
	Name		string		`json:"name"`
}

type TicketStatus struct {
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	Category	TicketStatusCategory	`json:"statusCategory"`
}

type TicketStatusCategory struct {
	Id			int			`json:"id"`
	Name		string		`json:"name"`
}

type JiraMember struct {
	Name		string		`json:"name"`
	DisplayName	string		`json:"displayName"`
}

type TicketIssueType struct {
	Id			string		`json:"id"`
	Name		string		`json:"name"`
	Subtask		bool		`json:"subtask"`
	Description	string		`json:"description"`
}

type JiraIssueType struct {
	Id			string		`json:"id"`
	Name		string		`json:"name"`
	Subtask		string		`json:"subtask"`
	Description	string		`json:"description"`
	Statuses	[]TicketStatus	`json:"statuses"`
}

type JiraProject struct {
	Id			string		`json:"id"`
	Key			string		`json:"key"`
	Name		string		`json:"name"`
}

type TicketComments struct {
	Comments	[]JiraComment	`json:"comments"`
	Total		int				`json:"total"`
}

type JiraComment struct {
	Author		JiraMember	`json:"author"`
	Body		string		`json:"body"`
	Created		JiraTime	`json:"created"`
	Updated		JiraTime	`json:"updated"`
}

type NewComment struct {
	Body		string		`json:"body"`
}

type JiraTime struct {
    time.Time
}

func (t *JiraTime) UnmarshalJSON(data []byte) error {
    str := strings.Trim(string(data), `"`)
    parsed, err := time.Parse("2006-01-02T15:04:05.000-0700", str)
    if err != nil {
        return err
    }
    t.Time = parsed
    return nil
}

type FormattedTicket struct {
	Headline		string
	Link			string
	Priority		string
	IssueType		string
	TypeCombined	string
	Status			string
	StatusLong		string
	Creator			string
	Reporter		string
	Assignee		string
	Members			string
	TicketDetails	string
	TicketDates		string
	Divider 		string
	DividerShort 	string
	Description		string
	Comments		string
}

type TicketTransitions struct {
	TicketID		string
	Transitions		[]IssueTransition
}

type JiraIssueTransitions struct {
	Transitions		[]IssueTransition `json:"transitions"`
}

type IssueTransition struct {
	Id				string			`json:"id"`
	Name			string			`json:"name"`
	Description		string			`json:"description"`
	OpsbarSequence	int				`json:"opsbarSequence"`
	Status			TicketStatus	`json:"to"`
}

type TicketTransitionsMap map[string][]IssueTransition

type JiraUser struct {
	Key				string			`json:"key"`
	Name			string			`json:"name"`
	DisplayName		string			`json:"displayName"`
	Active			bool			`json:"active"`
	Deleted			bool			`json:"deleted"`
	TimeZone		string			`json:"timeZone"`
}
