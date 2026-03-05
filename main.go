package main

import (
	"os"
	"strings"
	"github.com/Maxylan/clj/internal/cmd"
)

func main() {
	program := cmd.ProgramDetails{
		Name: "clj",
		Version: "v1.0.0",
	}

	var registry = cmd.BuildCommandRegistry(os.Args, program)
	var selected *cmd.Command // Pick first match..

	for _, cmd := range registry {
		if cmd.Match {
			if strings.Contains(cmd.Name, "Init ") {
				cmd.Execute();
				continue
			}

			selected = &cmd
			break
		}
	}

	if selected != nil {
		selected.Execute();
	} else {
		cmd.Help(os.Args, "Found no command matching given list of arguments.")
	}
}
