package cmd

import "regexp"

var (
	reItalic     	= regexp.MustCompile(`_(.+?)_`)
	reBold       	= regexp.MustCompile(`\*(.+?)\*`)
	reMono			= regexp.MustCompile(`\{+}*([^\{\}]+)\{*\}+`)
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
	reDoubleQuoted	= regexp.MustCompile(`^"(.+?)"$`)
	reTicketId		= regexp.MustCompile(`\w+-\d+`)
)
