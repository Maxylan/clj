package cmd

import (
	"os"
	"fmt"
	"log"
	"encoding/json"
	"path/filepath"
)

const configFile = "config.json" 
var configDirPath string
var configPath string

func SetConfigPath(programName string) {
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		// Sets config path based on program name.
		configDirPath = filepath.Join(
			userConfigDir,
			"." + programName,
		)
		configPath = filepath.Join(
			configDirPath,
			configFile,
		)
	} else {
		log.Fatal("Cannot get `UserConfigDir`", err)
	}
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
