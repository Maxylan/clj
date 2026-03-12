package cmd

import (
	"fmt"
	"log"
	"regexp"
	"slices"
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

/**
 * `MarshalJiraMarkdown(..)` converts *some* supported text-input to Jira wiki markup output.
 *
 * *Bold* and _Italic_ remains untouched, should work as-is.
 *
 * See `MarshalingRegexPatterns` for a comprehensive list of supported formattings.  
 */
func MarshalJiraMarkdown(input string) string {
	patterns := &patterns.Marshaling
	s := input

	// Mono
	s = patterns.Mono.ReplaceAllString(s, "{{{}$1{}}}")

	// Mentions
	s = patterns.Mention.ReplaceAllString(s, "[~$1]")

	return s
}

/**
 * `UnmarshalJiraMarkdown(..)` converts Jira wiki markup to ANSI-formatted terminal output.
 */
func UnmarshalJiraMarkdown(input string) string {
	patterns := &patterns.Unmarshaling
	s := input

	// Bold (before italic, to avoid * conflicts)
	s = patterns.Bold.ReplaceAllString(s, Bold+"$1"+NoBold)

	// Italic
	s = patterns.Italic.ReplaceAllString(s, Italic+"$1"+NoItalic)

	// Shorten the string by removing unecessary amounts of newlines
	s = patterns.Shorten.ReplaceAllString(s, "$1")

	// Quotes
	s = patterns.Quote.ReplaceAllStringFunc(s, func(match string) string {
		inner := patterns.Quote.FindStringSubmatch(match)[1]
		inner = strings.TrimSpace(inner)
		inner = FixAnsiiOverlap(inner, []string{ NoDim, Reset }, Dim);

		return "\n> \"" + Dim + inner + NoDim + "\""
	})

	// Headlines h1-h5 → Bold
	for _, re := range []*regexp.Regexp{patterns.H1, patterns.H2, patterns.H3, patterns.H4, patterns.H5} {
		s = re.ReplaceAllStringFunc(s, func(match string) string {
			inner := re.FindStringSubmatch(match)[1]
			inner = FixAnsiiOverlap(inner, []string{ NoBold, Reset }, Bold);

			return Bold + inner + NoBold
		})
	}

	// h6 → Italic
	s = patterns.H6.ReplaceAllStringFunc(s, func(match string) string {
		inner := patterns.H6.FindStringSubmatch(match)[1]
		inner = FixAnsiiOverlap(inner, []string{ NoItalic, Reset }, Italic);
		return Italic + inner + NoItalic
	})

	// Mono
	s = patterns.Mono.ReplaceAllString(s, Red+"$1"+NoColor)

	// Nested bullets (** and deeper) before top-level
	s = patterns.BulletNested.ReplaceAllStringFunc(s, func(match string) string {
		groups := patterns.BulletNested.FindStringSubmatch(match)
		depth := len(groups[1]) // number of * characters
		indent := strings.Repeat("  ", depth)
		return indent + "•" + " " + groups[2]
	})

	// Top-level bullets
	s = patterns.BulletTop.ReplaceAllString(s, "  • $1")

	// Mentions
	s = patterns.Mention.ReplaceAllString(s, "@$1")

	return s
}

func FormatTicket(ticket Ticket, recentCommentsFirst bool) FormattedTicket {
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
		Description:	UnmarshalJiraMarkdown(ticket.Fields.Description),
		Comments:		FormatTicketComments(ticket.Fields.Comment, recentCommentsFirst),
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

func FormatTicketComments(commentSection TicketComments, recentFirst bool) string {
	if commentSection.Total == 0 || len(commentSection.Comments) == 0 {
		return ""
	}

	if recentFirst {
		slices.Reverse(commentSection.Comments)
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
				UnmarshalJiraMarkdown(comment.Body),
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
			UnmarshalJiraMarkdown(comment.Body),
		)
	})

	return "» Comments:\n" + strings.Join(comments, "\n") + "\n"
}

func FormatTransition(transition IssueTransition, detailed bool) string {
	out := fmt.Sprintf("  • %s\t%s#%s%s", transition.Name, Dim, transition.Id, Italic)

	if detailed && len(transition.Status.Category.Name) > 0 {
		out += fmt.Sprintf("\tCategory \"%s\"", transition.Status.Category.Name)
	}

	if len(transition.Description) > 0 {
		out += fmt.Sprintf("\t(%s)", transition.Description)
	}

	return out + Reset
}

func FormatUser(user JiraUser, detailed bool) string {
	out := fmt.Sprintf("  • %s\t%s@%s%s", user.DisplayName, Dim, user.Name, Italic)

	if detailed {
		if len(user.TimeZone) > 0 {
			out += fmt.Sprintf("\tTZ \"%s\"", user.TimeZone)
		}

		if len(user.Key) > 0 {
			out += fmt.Sprintf("\t(%s)", user.Key)
		}
	}

	return out + Reset
}
