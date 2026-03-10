package cmd

import "regexp"

type RegexPatterns struct {
	TicketId		*regexp.Regexp
	DoubleQuoted	*regexp.Regexp
	Marshaling 		MarshalingRegexPatterns
	Unmarshaling 	UnmarshalingRegexPatterns
}
type UnmarshalingRegexPatterns struct {
	Italic     		*regexp.Regexp
	Bold       		*regexp.Regexp
	Mono			*regexp.Regexp
	Mention       	*regexp.Regexp
	H6         		*regexp.Regexp
	H5         		*regexp.Regexp
	H4         		*regexp.Regexp
	H3         		*regexp.Regexp
	H2         		*regexp.Regexp
	H1         		*regexp.Regexp
	Quote      		*regexp.Regexp
	BulletNested	*regexp.Regexp
	BulletTop  		*regexp.Regexp
	Shorten			*regexp.Regexp
}
type MarshalingRegexPatterns struct {
	Mono			*regexp.Regexp
	Mention       	*regexp.Regexp
}

var patterns RegexPatterns = RegexPatterns{
	TicketId:		regexp.MustCompile(`\w+-\d+`),
	DoubleQuoted:	regexp.MustCompile(`^"(.+?)"$`),
	Unmarshaling:	UnmarshalingRegexPatterns{
		Italic:     	regexp.MustCompile(`_(.+?)_`),
		Bold:       	regexp.MustCompile(`\*(.+?)\*`),
		Mono:			regexp.MustCompile(`\{+}*([^\{\}]+)\{*\}+`),
		Mention:       	regexp.MustCompile(`\[~(\w[\w- ]+\w)\]`),
		H6:         	regexp.MustCompile(`(?m)^ *h6\.\s*(.+)$`),
		H5:         	regexp.MustCompile(`(?m)^ *h5\.\s*(.+)$`),
		H4:         	regexp.MustCompile(`(?m)^ *h4\.\s*(.+)$`),
		H3:         	regexp.MustCompile(`(?m)^ *h3\.\s*(.+)$`),
		H2:         	regexp.MustCompile(`(?m)^ *h2\.\s*(.+)$`),
		H1:         	regexp.MustCompile(`(?m)^ *h1\.\s*(.+)$`),
		Quote:      	regexp.MustCompile(`(?ms)\{quote\}(.+?)\{quote\}`),
		BulletNested:	regexp.MustCompile(`(?m)^ *(\*{2,})\s*(.+)$`),
		BulletTop:  	regexp.MustCompile(`(?m)^ *\*\s+(.+)$`),
		Shorten:		regexp.MustCompile(`(\s+ *\s+) *\s+ *\s+`),
	},
	Marshaling: MarshalingRegexPatterns{
		Mono:			regexp.MustCompile("`(.+?)`"),
		Mention:       	regexp.MustCompile(`@(.{3})`),
	},
}

/** Matches given name(s) against simple RegEx pattern to determine if it/they are valid Jira ticket names */
func IsValidTicketName(names ...string) bool {
	if len(names) == 0 {
		return false
	}

	for _, n := range names {
		if !patterns.TicketId.Match([]byte(n)) {
			return false;
		}
	}

	return true
}

func isStringArg(arg string) bool {
	if len(arg) < 3 {
		return false
	}

	// Ensure non-empty double-quoted string
	return patterns.DoubleQuoted.Match([]byte(arg))
}
