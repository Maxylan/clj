package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func FixAnsiiOverlap(input string, codes []string, insert string) string {
	str := input

	for _, code := range codes {
		if i := strings.Index(input, code); i > -1 {
			offset := i + len(code)
			str = input[:offset] + insert

			if len(input) > offset + 1 {
				str += FixAnsiiOverlap(input[offset:], []string{ code }, insert)
			}
		}
	}

	return str
}

/** ParseJiraMarkdown converts Jira wiki markup to ANSI-formatted terminal output. */
func ParseJiraMarkdown(input string) string {
	s := input

	// Bold (before italic, to avoid * conflicts)
	s = reBold.ReplaceAllString(s, Bold+"$1"+NoBold)

	// Italic
	s = reItalic.ReplaceAllString(s, Italic+"$1"+NoItalic)

	// Shorten the string by removing unecessary amounts of newlines
	s = reShorten.ReplaceAllString(s, "$1")

	// Quotes
	s = reQuote.ReplaceAllStringFunc(s, func(match string) string {
		inner := reQuote.FindStringSubmatch(match)[1]
		inner = strings.TrimSpace(inner)
		inner = FixAnsiiOverlap(inner, []string{ NoDim, Reset }, Dim);

		return "\n> \"" + Dim + inner + NoDim + "\""
	})

	// Headlines h1-h5 → Bold
	for _, re := range []*regexp.Regexp{reH1, reH2, reH3, reH4, reH5} {
		s = re.ReplaceAllStringFunc(s, func(match string) string {
			inner := re.FindStringSubmatch(match)[1]
			inner = FixAnsiiOverlap(inner, []string{ NoBold, Reset }, Bold);

			return Bold + inner + NoBold
		})
	}

	// h6 → Italic
	s = reH6.ReplaceAllStringFunc(s, func(match string) string {
		inner := reH6.FindStringSubmatch(match)[1]
		inner = FixAnsiiOverlap(inner, []string{ NoItalic, Reset }, Italic);
		return Italic + inner + NoItalic
	})

	// Mono
	s = reMono.ReplaceAllString(s, Red+"$1"+NoColor)

	// Nested bullets (** and deeper) before top-level
	s = reBulletNested.ReplaceAllStringFunc(s, func(match string) string {
		groups := reBulletNested.FindStringSubmatch(match)
		depth := len(groups[1]) // number of * characters
		indent := strings.Repeat("  ", depth)
		return indent + "•" + " " + groups[2]
	})

	// Top-level bullets
	s = reBulletTop.ReplaceAllString(s, "  • $1")

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
		Divider: 		"  ",
		DividerShort: 	"  ",
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

	out.DividerShort += strings.Repeat(
		"─",
		len(fmt.Sprintf( // len(..) w/o Ansii sequences
			"%s/%s",
			config.JiraURL,
			ticket.Key,
		)),
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

	comments := Map(commentSection.Comments, func(comment JiraComment, i int) string {
		if i > 0 && commentSection.Comments[i - 1].Author.Name == comment.Author.Name {
			return fmt.Sprintf(
				":\n┆╴%s%s%s (%s)%s\n%s",
				Dim,
				Italic,
				comment.Created,
				comment.Updated,
				Reset,
				ParseJiraMarkdown(comment.Body),
			)
		}
		return fmt.Sprintf(
			"┆\n╰ %s %s(%s)\n  %s%s (%s)%s\n%s",
			comment.Author.DisplayName,
			Dim,
			comment.Author.Name,
			Italic,
			comment.Created,
			comment.Updated,
			Reset,
			ParseJiraMarkdown(comment.Body),
		)
	})

	return "» Comments:\n" + strings.Join(comments, "\n") + "\n"
}
