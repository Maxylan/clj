package main

import (
	"os"
	"fmt"
	"slices"
	"strings"
	"github.com/Maxylan/clj/internal/cmd"
)

func main() {
	program := cmd.ProgramDetails{
		Name: "clj",
		Version: "v1.0.0",
	}

	var verbose = slices.Contains(os.Args, "--verbose") || slices.Contains(os.Args, "-v")

	if verbose {
		fmt.Println(
			cmd.Ansii(cmd.Bold, program.Name),
			program.Version,
		)
	}

	var registry = cmd.BuildCommandRegistry(program)

	// Parses `os.Args` to into argument "chains".
	// Each chain represents its own operation / command, w/ its own list of Ticket IDs.
	// Chains are separated by the keyword "and".
	argChains := cmd.ParseArgs(os.Args)

	for i, chain := range argChains {
		var selected *cmd.Command // Pick first match..

		if verbose {
			fmt.Printf("%s[%d] %v%s\n\n", cmd.Dim, i, chain, cmd.Reset)
		}

		for _, cmd := range registry {
			if cmd.Match(chain) {
				if strings.Contains(cmd.Name, "Init ") {
					cmd.Execute(chain);
					continue
				}

				selected = &cmd
				break
			}
		}

		if selected != nil {
			selected.Execute(chain);
		} else {
			cmd.CommandNotFound("Found no command matching given list of arguments.")
		}
	}
}
