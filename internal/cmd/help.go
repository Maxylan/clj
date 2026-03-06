package cmd

import (
	"fmt"
	"slices"
	"strings"
)

const cmd_help = "Help"

func register_help(
	args	[]string,
	program	ProgramDetails,
) Command {

	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_help,
			program.Name,
			program.Version,
		),
		Match: len(args) < 2 || strings.ToLower(args[1]) == "help",
		Execute: func() { Help(args, "") },
		Details: CommandDetails{
			Name:			cmd_help,
			Usage:			fmt.Sprintf("%s help | alt. [-h|--help] after any command.", program.Name),
			Description:	"\"Prints usage & detailed explanations of each subcommand of their accepted arguments\"",
			Subcommands: []CommandDetails{{
				Name:			"Minimal",
				Usage:			fmt.Sprintf("%s help minimal|[-m|--minimal]", program.Name),
				Description:	"Omit lengthy descriptions, only printing command name + usage.",
				Subcommands: 	[]CommandDetails{},
			}},
		},
	};
}

func Help(args []string, msg string) {
	if len(msg) > 0 {
		fmt.Println(msg + "\n")
	}

	minimal :=
		slices.Contains(args, "-m") ||
		slices.Contains(args, "--minimal") ||
		(len(args) > 2 && strings.EqualFold(args[2], "minimal"))

	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		// If appended as argument (i.e -h|--help)
		// ..for-each registered command *matching* given args, print help. 
		matches := FilterMap(
			registry[:],
			func(cmd Command, _ int) bool { return cmd.Name != "Help" && cmd.Match },
			func(cmd Command, _ int) CommandDetails { return cmd.Details },
		)

		for i := range matches {
			PrintDetails(matches[i], minimal)
		}

		return
	}

	// <name> help - Print help for each registered command 
	for i := range len(registry) {
		PrintDetails(registry[i].Details, minimal)
	}
}

func FormatDetails(
	details CommandDetails,
	minimal bool,
	level int,
	prefix string,
) string {
	indent := strings.Repeat("\t", level)
	nl := "\n" + indent

	var sb strings.Builder;

	sb.WriteString(Ansii(
		indent,
		prefix,
		Bold,
		details.Name,
		NoBold,
		nl,
	))

	if len(details.Subcommands) > 0 {
		sb.WriteString(Ansii(
			"├ ",
			Dim,
			details.Usage,
			NoDim,
			nl,
		))

		if !minimal {
			sb.WriteString(Ansii(
				"├ ",
				Dim,
				details.Description,
				NoDim,
				nl,
			))
		}

		sb.WriteString("╰ Subcommands:\n")

		for _, s := range details.Subcommands {
			sb.WriteString(FormatDetails(s, minimal, level + 1, "• "))
		}
	} else if minimal {
		sb.WriteString(Ansii(
			"╰ ",
			Dim,
			details.Usage,
			NoDim,
		))
	} else {
		sb.WriteString(Ansii(
			"├ ",
			Dim,
			details.Usage,
			NoDim,
			nl,
			"╰ ",
			Dim,
			details.Description,
			NoDim,
		))
	}

	return sb.String() + "\n";
}

func PrintDetails(details CommandDetails, minimal bool) {
	fmt.Println(FormatDetails(details, minimal, 0, "» "))
}
