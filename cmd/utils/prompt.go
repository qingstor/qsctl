package utils

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
)

// StringInputer is interface to set string from input
type StringInputer interface {
	SetStringP(*string) error
	GetStringP() *string
}

// BoolInputer is interface to get a bool variable
type BoolInputer interface {
	GetBool() (bool, error)
}

// PromptInput used to set up a prompt input
type PromptInput struct {
	Key      string
	Value    *string
	MsgTitle string
}

// SetStringP set value to given string pointer
func (p PromptInput) SetStringP(s *string) error {
	if s == nil {
		return fmt.Errorf("invalid value to set, value cannot be nil")
	}
	var msg string
	defaultValue := viper.GetString(p.Key)

	if defaultValue != "" {
		msg = " [" + defaultValue + "]"
	}
	prompt := &survey.Input{
		Message: p.MsgTitle + msg + ":",
	}

	var tmp string
	if err := survey.AskOne(prompt, &tmp); err != nil {
		return err
	}
	if tmp != "" {
		*s = tmp
	}
	return nil
}

// GetStringP returns a pointer to string to set
func (p PromptInput) GetStringP() *string {
	return p.Value
}

// PromptConfirm used for prompt confirm
type PromptConfirm struct {
	Msg string
}

// GetBool returns a bool
func (p PromptConfirm) GetBool() (bool, error) {
	var res bool
	var prompt = &survey.Confirm{
		Message: p.Msg,
	}
	if err := survey.AskOne(prompt, &res); err != nil {
		return false, err
	}
	return res, nil
}
