package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qsctl/v2/constants"
)

// InputConfig is the struct for setup config input
type InputConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`

	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
	LogLevel string `yaml:"log_level"`
}

// NewInputConfig setup InputConfig and return the struct
func NewInputConfig() InputConfig {
	return InputConfig{
		AccessKeyID:     "",
		SecretAccessKey: "",
		Host:            constants.DefaultHost,
		Port:            constants.DefaultPort,
		Protocol:        constants.DefaultProtocol,
		LogLevel:        constants.DefaultLogLevel,
	}
}

var keyPrompt = []*survey.Question{
	{
		Name:     "AccessKeyID",
		Prompt:   &survey.Input{Message: "AccessKeyID:"},
		Validate: survey.Required,
	},
	{
		Name:     "SecretAccessKey",
		Prompt:   &survey.Password{Message: "SecretAccessKey:"},
		Validate: survey.Required,
	},
}

var isPublicCloud = true
var publicCloudPrompt = &survey.Confirm{
	Message: "Apply qsctl for QingStor public cloud?",
}

var privatePrompt = []*survey.Question{
	{
		Name:     "Host",
		Prompt:   &survey.Input{Message: "Host:"},
		Validate: survey.Required,
	},
	{
		Name:   "Port",
		Prompt: &survey.Input{Message: "Port:"},
		Validate: func(ans interface{}) error {
			if v, ok := ans.(string); ok {
				if _, err := strconv.Atoi(v); err != nil {
					return fmt.Errorf("cannot parse port from your input <%v>: [%w]", ans, err)
				}
				return nil
			}
			return fmt.Errorf("cannot transfer port from non-string input, please check your input")
		},
	},
	{
		Name: "Protocol",
		Prompt: &survey.Select{
			Message: "Protocol:",
			Options: []string{"http", "https"},
		},
		Validate: survey.Required,
	},
}

var logLevelPrompt = &survey.Select{
	Message: "Log level:",
	Options: []string{"debug", "info", "warn", "error", "fatal"},
}

var confirm = false
var confirmPrompt = &survey.Confirm{
	Message: "Confirm your config?",
}

// SetupConfigInteractive setup input config interactively
func SetupConfigInteractive() (fileName string, err error) {
	in := NewInputConfig()

	// multiple times to set config if not confirmed
	for {
		if err = survey.Ask(keyPrompt, &in); err != nil {
			return "", err
		}

		if err = survey.AskOne(publicCloudPrompt, &isPublicCloud); err != nil {
			return "", err
		}

		if !isPublicCloud {
			if err = survey.Ask(privatePrompt, &in); err != nil {
				return "", err
			}
		}

		if err = survey.AskOne(logLevelPrompt, &in.LogLevel); err != nil {
			return "", err
		}

		if err = survey.AskOne(confirmPrompt, &confirm); err != nil {
			return "", err
		}

		if confirm {
			break
		}
		// reset the input
		in = NewInputConfig()
	}

	b, err := yaml.Marshal(in)
	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fileName = filepath.Join(homeDir, ".qingstor/config.yaml")
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write(b); err != nil {
		return "", err
	}
	return fileName, nil
}
