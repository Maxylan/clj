package cmd

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"encoding/json"
	"path/filepath"
	"strings"
)

const cmd_init = "Init " + Dim + "(initial setup / configuration)" + NoDim
const configFile = "config.json" 
var configDirPath string
var configPath string

func register_init(
	args	[]string,
	program	ProgramDetails,
) Command {
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		// Sets config path based on program name.
		configDirPath = filepath.Join(
			userConfigDir,
			"." + program.Name,
		)
		configPath = filepath.Join(
			configDirPath,
			configFile,
		)
	} else {
		log.Fatal("Cannot get `UserConfigDir`", err)
	}

	return Command{
		Name: fmt.Sprintf(
			"%s (%s, %s)",
			cmd_init,
			program.Name,
			program.Version,
		),
		Match: !IsConfigured() || (len(args) > 1 && strings.EqualFold(args[1], "init")),
		Execute: func() { Setup(args) },
		Details: CommandDetails{
			Name:			cmd_init,
			Usage:			fmt.Sprintf("%s init", program.Name),
			Description:	"\"Performs initial setup / configuration. Can be re-run later to reconfigure this utility.\"",
			Subcommands:	[]CommandDetails{},
		},
	};
}

func ConfigExists() (bool, error) {
	_, err := os.Stat(configPath)
	return err == nil || !os.IsNotExist(err), err
}

func GetConfig() (Config, error) {
	if len(configPath) < 1 {
		log.Fatal("No config path set")
	}

	if exists, statError := ConfigExists(); !exists {
		// Check if the config's directory exists, creates it if it doesn't.
		_, dirStatError := os.Stat(configDirPath)
		if os.IsNotExist(dirStatError) {
			mkdirError := os.MkdirAll(configDirPath, 0755)
			return Config{}, mkdirError
		}

		return Config{}, statError
	}

	// Decode (read)
	data, fileReadError := os.ReadFile(configPath)

	if fileReadError != nil {
		if os.IsNotExist(fileReadError) {
			return Config{}, fileReadError
		}

		log.Fatal("Could not read config at " + configPath)
	}

	var config Config
	unmarshalError := json.Unmarshal(data, &config)
	return config, unmarshalError
}

func SaveConfig(config Config) error {
	if len(config.JiraURL) < 1 && len(config.Token) < 8 {
		fmt.Println(Ansii(Bold, "(!)", NoBold, " Refusing to save invalid configuration") + "\n")
		return nil
	}

	if len(configPath) < 1 {
		log.Fatal("No config path set")
	}

	// Encode (write)
	data, marshalError := json.MarshalIndent(config, "", "  ")
	if marshalError != nil {
		return marshalError
	}

	fmt.Println(Ansii("» Saving configuration", Dim, " to ", Italic, "'", configPath, "'") + "\n")
	return os.WriteFile(configPath, data, 0600)
}

func IsConfigured() bool {
	config, err := GetConfig()
	if err != nil {
		return false
	}

	return len(config.JiraURL) > 0 && len(config.Token) > 0
}

func Setup(args []string) {
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

	config.JiraURL = line

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

	config.Token = line

    if saveError := SaveConfig(config); saveError != nil {
        log.Fatal("Failed to save configuration. ", saveError.Error())
    }
}
