package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	reItalic     	= regexp.MustCompile(`_(.+?)_`)
	reBold       	= regexp.MustCompile(`\*(.+?)\*`)
	reH6         	= regexp.MustCompile(`(?m)^ *h6\.\s*(.+)$`)
	reH5         	= regexp.MustCompile(`(?m)^ *h5\.\s*(.+)$`)
	reH4         	= regexp.MustCompile(`(?m)^ *h4\.\s*(.+)$`)
	reH3         	= regexp.MustCompile(`(?m)^ *h3\.\s*(.+)$`)
	reH2         	= regexp.MustCompile(`(?m)^ *h2\.\s*(.+)$`)
	reH1         	= regexp.MustCompile(`(?m)^ *h1\.\s*(.+)$`)
	reQuote      	= regexp.MustCompile(`(?ms)\{quote\}(.+?)\{quote\}`)
	reBulletNested	= regexp.MustCompile(`(?m)^ *(\*{2,})\s*(.+)$`)
	reBulletTop  	= regexp.MustCompile(`(?m)^ *\*\s+(.+)$`)
	reShorten		= regexp.MustCompile(`(\s+ *\s+) *\s+ *\s+`)
)

/** ParseJiraMarkdown converts Jira wiki markup to ANSI-formatted terminal output. */
func ParseJiraMarkdown(input string) string {
	s := input

	// Shorten the string by removing unecessary amounts of newlines
	s = reShorten.ReplaceAllString(s, "$1")

	// Quotes (before bold, since {quote} blocks may contain *)
	s = reQuote.ReplaceAllStringFunc(s, func(match string) string {
		inner := reQuote.FindStringSubmatch(match)[1]
		inner = strings.TrimSpace(inner)
		return Dim + inner + Reset
	})

	// Headlines h1-h5 → Bold
	for _, re := range []*regexp.Regexp{reH1, reH2, reH3, reH4, reH5} {
		s = re.ReplaceAllStringFunc(s, func(match string) string {
			inner := re.FindStringSubmatch(match)[1]
			return Bold + inner + Reset
		})
	}

	// h6 → Italic
	s = reH6.ReplaceAllStringFunc(s, func(match string) string {
		inner := reH6.FindStringSubmatch(match)[1]
		return Italic + inner + Reset
	})

	// Nested bullets (** and deeper) before top-level
	s = reBulletNested.ReplaceAllStringFunc(s, func(match string) string {
		groups := reBulletNested.FindStringSubmatch(match)
		depth := len(groups[1]) // number of * characters
		indent := strings.Repeat("  ", depth)
		return indent + "•" + " " + groups[2]
	})

	// Top-level bullets
	s = reBulletTop.ReplaceAllString(s, "  • $1")

	// Bold (before italic, to avoid * conflicts)
	s = reBold.ReplaceAllString(s, Bold+"$1"+Reset)

	// Italic
	s = reItalic.ReplaceAllString(s, Italic+"$1"+Reset)

	return s
}

func FormatTicket(ticket Ticket) FormattedTicket {
	config, err := GetConfig()
	if err != nil {
		log.Fatal("Could not load user configuration", err)
	}

	out := FormattedTicket{
		Headline: Ansii(
			"» ",
			Cyan,
			Underline,
			ticket.Key,
			Reset,
			": \"",
			Bold,
			ticket.Fields.Summary,
			NoBold,
			"\"",
		),
		Link: fmt.Sprintf(
			"╰ %s%s/%s%s",
			BBlue,
			config.JiraURL,
			ticket.Key,
			Reset,
		),
		Priority: fmt.Sprintf(
			"  • %s",
			ticket.Fields.Priority.Name,
		),
		IssueType: fmt.Sprintf(
			"  • %s %s(%s)%s",
			ticket.Fields.IssueType.Name,
			Italic,
			ticket.Fields.IssueType.Description,
			Reset,
		),
		TypeCombined: fmt.Sprintf(
			"  • %s %s(%s)%s",
			ticket.Fields.Priority.Name,
			Italic,
			ticket.Fields.IssueType.Name,
			Reset,
		),
		Status: fmt.Sprintf(
			"  • \"%s\"",
			ticket.Fields.Status.Name,
		),
		StatusLong: fmt.Sprintf(
			"  • \"%s\" (%s, %s)",
			ticket.Fields.Status.Name,
			ticket.Fields.Status.Description,
			ticket.Fields.Status.Category.Name,
		),
		Creator: fmt.Sprintf(
			"  %sC:%s %s %s(%s)%s",
			Cyan,
			Reset,
			ticket.Fields.Creator.DisplayName,
			Italic,
			ticket.Fields.Creator.Name,
			Reset,
		),
		Reporter: fmt.Sprintf(
			"  %sR:%s %s %s(%s)%s",
			Cyan,
			Reset,
			ticket.Fields.Reporter.DisplayName,
			Italic,
			ticket.Fields.Reporter.Name,
			Reset,
		),
		Assignee: fmt.Sprintf(
			"  %sA:%s %s %s(%s)%s",
			Cyan,
			Reset,
			ticket.Fields.Assignee.DisplayName,
			Italic,
			ticket.Fields.Assignee.Name,
			Reset,
		),
		Members: fmt.Sprintf(
			"  %sR%s: %s %s(%s)%s, %sA%s: %s %s(%s)%s",
			Bold,
			NoBold,
			ticket.Fields.Reporter.DisplayName,
			Italic,
			ticket.Fields.Reporter.Name,
			Reset,
			Bold,
			NoBold,
			ticket.Fields.Assignee.DisplayName,
			Italic,
			ticket.Fields.Assignee.Name,
			Reset,
		),
		TicketDetails: Ansii(
			Dim,
			fmt.Sprintf(
				"  Project: %s%s (%s)%s, Labels: %s%v",
				Italic,
				ticket.Fields.Project.Key,
				ticket.Fields.Project.Name,
				NoItalic,
				Italic,
				ticket.Fields.Labels,
			),
		),
		TicketDates: Ansii(
			Dim,
			fmt.Sprintf(
				"  Created: %s%v%s, Updated: %s%v%s",
				Italic,
				ticket.Fields.Created,
				NoItalic,
				Italic,
				ticket.Fields.Updated,
				NoItalic,
			),
		),
		Divider: "  ",
		Description:	ParseJiraMarkdown(ticket.Fields.Description),
		Comments:		FormatTicketComments(ticket.Fields.Comment),
	}

	out.Divider += strings.Repeat(
		"─",
		len(fmt.Sprintf( // len(..) w/o Ansii sequences
			"  Created: %v, Updated: %v",
			ticket.Fields.Created,
			ticket.Fields.Updated,
		)) - 2,
	)

	if len(ticket.Fields.Status.Description) > 0 {
		out.Status += " ("+ticket.Fields.Status.Description+")"
	}

	return out
}

func FormatTicketComments(commentSection TicketComments) string {
	if commentSection.Total == 0 || len(commentSection.Comments) == 0 {
		return ""
	}

	comments := Map(commentSection.Comments, func(comment JiraComment) string {
		return fmt.Sprintf(
			"┆\n╰ %s %s(%s), %s (%s)%s\n  %s",
			comment.Author.DisplayName,
			Dim,
			comment.Author.Name,
			comment.Created,
			comment.Updated,
			NoDim,
			ParseJiraMarkdown(comment.Body),
		)
	})

	return "» Comments:\n" + strings.Join(comments, "\n")
}
