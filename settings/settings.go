package settings

import (
	"autogit/semanticgit/git/gitraw"
	"autogit/settings/envs"
	"autogit/settings/logus"
	"autogit/settings/types"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

const ToolName = "autogit"

var HookFolderName = fmt.Sprintf(".%s-hooks", ToolName)

//go:embed version.txt
var Version string

type ConfigScheme struct {
	Changelog  ChangelogScheme  `yaml:"changelog"`
	Regex      RegexScheme      `yaml:"regex"`
	Validation ValidationScheme `yaml:"validation"`
	Git        struct {
		SSHPath string `yaml:"sshPath"`
	} `yaml:"git"`
}

var GlobalConfigPath types.ConfigPath

var ProjectConfigPath types.ConfigPath

var ProjectPath types.FilePath

const ConfigFileName types.ConfigPath = "autogit.yml"

func init() {
	GlobalConfigPath = types.ConfigPath(filepath.Join(string(envs.PathUserHome), string(ConfigFileName)))

	g := gitraw.NewGitRepo()
	w, err := g.Worktree()
	logus.CheckFatal(err, "are we not in git repo folder?")
	ProjectPath = types.FilePath(w.Filesystem.Root())

	ProjectConfigPath = types.ConfigPath(filepath.Join(string(ProjectPath), string(ConfigFileName)))
}

var cachedConfigFile []byte = []byte{}

func is_file_missing(err error) bool {
	_, ok := err.(*fs.PathError)
	return ok
}

func readSettingsfile(configPath types.ConfigPath) []byte {
	// TODO You could have written less config readings across your code.
	// Caching for the purpose of rendering logging records only once
	if len(cachedConfigFile) != 0 {
		return cachedConfigFile
	}

	file, err := os.ReadFile(string(configPath))
	local_file_is_not_found := false
	if err != nil {
		if is_file_missing(err) {
			logus.Debug("not found at path repository local file with config. Fallback to global config", logus.FilePath(configPath.ToFilePath()), logus.OptError(err))
			local_file_is_not_found = true
		} else {
			logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), logus.OptError(err))
		}
	} else {
		logus.Debug("succesfuly read config from local repository project path", logus.ConfigPath(configPath))
	}

	global_file_is_not_found := false
	if local_file_is_not_found {
		file, err = os.ReadFile(string(GlobalConfigPath))
		if err != nil {
			if is_file_missing(err) {
				logus.Debug("not found at path repository global file with config. Fallback to other in memory config", logus.FilePath(configPath.ToFilePath()), logus.OptError(err))
				global_file_is_not_found = true
			} else {
				logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), logus.OptError(err))
			}
		}
	}

	if local_file_is_not_found && global_file_is_not_found {
		logus.Debug("fallback to memory settings file")
		file = []byte(ConfigExample)
	}

	cachedConfigFile = file
	return file
}

func merge_maps[T comparable](result map[T]interface{}, additions map[T]interface{}) map[T]interface{} {
	/*
		for key value in result {
			if key is present in additions {

				if value is not map[string]interface{} {
					result[key] = additions[key]
				} else {
					merge_maps(result[key], additions[key])
				}

			}
		}
	*/

	// For all keys and values in resulting hashmap
	for result_key, result_value := range result {

		// if key is present in additions hasmap
		if additions_value, is_present := additions[result_key]; is_present {

			// if value is not map[string]{interface}
			if asserted_addition_value, ok := additions_value.(map[string]interface{}); !ok {
				// just override with value from additions to resulting hashmap
				result[result_key] = additions_value
			} else {

				// otherwise try to merge recursively values from additions nested hashmaps to resulting key
				if asserted_result_value, ok := result_value.(map[string]interface{}); ok {
					merge_maps(asserted_result_value, asserted_addition_value)
				} else {
					logus.Fatal(`
						failed to assert value of config in memory being of same type as value in input config
						potentially broken config
					`)
				}
			}
		}
	}

	return result
}

func configRead(file []byte) *ConfigScheme {
	config := ConfigScheme{}

	file_config := make(map[interface{}]interface{})
	err := yaml.Unmarshal(file, &file_config)
	logus.CheckFatal(err, "unable to unmarshal input config")

	memory_config := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(ConfigExample), &memory_config)
	logus.CheckFatal(err, "unable to unmrashal memory config")

	// merged file content onto memory config
	merged_config := merge_maps(memory_config, file_config)

	merged_config_as_bytes, err := yaml.Marshal(&merged_config)
	logus.CheckFatal(err, "unable to marshal merged config")

	err = yaml.Unmarshal(merged_config_as_bytes, &config)
	logus.CheckFatal(err, "unable to unmarshal merged config")

	// Config overrides for dev env purposes
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_SSH_PATH"); ok {
		config.Git.SSHPath = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_COMMIT_URL"); ok {
		config.Changelog.CommitURL = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_COMMIT_RANGE_URL"); ok {
		config.Changelog.CommitRangeURL = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_ISSUE_URL"); ok {
		config.Changelog.IssueURL = value
	}

	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_VALIDATION_RULES_HEADER_SUBJECT_MIN_WORDS"); ok {
		res, err := strconv.Atoi(value)
		logus.CheckFatal(err, "crashed when trying to atoi min words env value")
		config.Validation.Rules.Header.Subject.MinWords = res
	}

	return &config
}

// yml package has no way to validate that there is no unknown undeclared fields
func validateSettingsScheme(configPath types.ConfigPath) {
	var config ConfigScheme
	var err error

	file := readSettingsfile(configPath)
	// Marshal file to struct
	err = yaml.Unmarshal(file, &config)
	logus.CheckFatal(err, "failed to unmarshal config")

	// Unmarshal struct to bytes
	m, err := yaml.Marshal(&config)
	logus.CheckFatal(err, "unable to marshal settings")

	// Marshal bytes to map
	a := make(map[interface{}]interface{})
	err = yaml.Unmarshal(m, &a)
	logus.CheckFatal(err, "failed unmarshaling to yaml")

	// compare with file marshaled to map
	b := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &b)
	logus.CheckFatal(err, "failed unmarshaling to yaml again")

	if !reflect.DeepEqual(a, b) {
		logus.Fatal(`
		setting file contains not registered keys.
		Check your version of autogit, and documentation related to settings
		`, logus.Expected(a), logus.Actual(b))
	}
}

func LoadSettings(configPath types.ConfigPath) *ConfigScheme {
	file := readSettingsfile(configPath)
	config := configRead(file)
	ChangelogInit(*config)
	RegexInit(config)
	validateSettingsScheme(configPath)

	return config
}

var config *ConfigScheme

func GetConfig() ConfigScheme {
	if config == nil {
		config = LoadSettings(ProjectConfigPath)
	}
	return *config
}

//go:embed autogit.example.yml
var ConfigExample string
