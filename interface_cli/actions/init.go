package actions

import (
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/settings/types"
	"autogit/settings/utils"
	"fmt"
	"strings"
)

const (
	InitAdvice string = "activate hook with `autogit hook activate [--global]`"
)

func init_write_config(config_path types.ConfigPath) {
	if utils.FileExists(config_path.ToFilePath()) {
		logus.Fatal("file with settings already exists", logus.ConfigPath(config_path))
		return
	}

	file := utils.NewFile(config_path.ToFilePath()).CreateToWriteF()
	defer file.Close()

	config_lines := strings.Split(settings.ConfigExample, "\n")
	for i, line := range config_lines {
		config_lines[i] = fmt.Sprintf("# %s", line)
	}
	commented_out_config := strings.Join(config_lines, "\n")

	file.WritelnF(commented_out_config)

	logus.Info("Succesfully created settings in location", logus.ConfigPath(config_path))
	logus.Info("Try to " + InitAdvice + ". It will automatically verify committs for you!")
}

func Initialize(init_globally bool) {
	config_path := settings.ProjectConfigPath
	if init_globally {
		config_path = settings.GlobalConfigPath
	}
	init_write_config(config_path)
}