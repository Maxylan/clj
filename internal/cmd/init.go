package cmd

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
)

const cmd_init = "Init " + Dim + "(initial setup / configuration)" + NoDim

func register_init(
	program	ProgramDetails,
) Command {
	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_init,
			program.Name,
			program.Version,
		),
		Match: match_init,
		Execute: Setup,
		Details: CommandDetails{
			Name:			cmd_init,
			Usage:			fmt.Sprintf("%s init", program.Name),
			Description:	"Performs initial setup / configuration. Can be re-run later to reconfigure this utility.",
			Subcommands:	[]CommandDetails{},
		},
	};
}

func match_init(chain CommandArgChain) bool {
	return !IsConfigured() || (len(chain.Keywords) > 0 && strings.EqualFold(chain.Keywords[0], "init"))
}


func Setup(chain CommandArgChain) {
    fmt.Println(cmd_init + "\n ")

	config, _ := GetConfig()

	question := Ansii("» Enter the ", Underline, "URL", NoUnderline, " to your Jira installation:")
	if len(config.JiraURL) > 0 {
		question += Ansii(Dim, Italic, "current (", config.JiraURL, ")");
	}

    fmt.Println(question)

    reader := bufio.NewReader(os.Stdin)

    line, readError := reader.ReadString('\n')
    if readError != nil {
        log.Fatal(readError)
    }

	config.JiraURL = line[:len(line) - 2] // gets rid of '\n' delim..

	question = Ansii("» Enter your ", Underline, "\"PAT\"")
	question += Ansii(" ", Dim, Italic, "(Personal Access Token)")
	question += Ansii(" to your Jira account:")

	if len(config.Token) > 16 {
		question += Ansii(Dim, Italic, "current ([...]", config.Token[16:], ")");
	}

    fmt.Println(question)

    line, readError = reader.ReadString('\n')
    if readError != nil {
        log.Fatal(readError)
    }

	config.Token = line[:len(line) - 2]

    if saveError := SaveConfig(config); saveError != nil {
        log.Fatal("Failed to save configuration. ", saveError.Error())
    }
}
