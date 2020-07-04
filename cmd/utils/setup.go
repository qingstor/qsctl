package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Xuanwo/storage"
	"github.com/qingstor/noah/task"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

// CredentialConfig is the struct for setup config input
type CredentialConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`

	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`

	Inputers    []StringInputer `yaml:"-"`
	Checker     InputChecker    `yaml:"-"`
	WriteCloser io.WriteCloser  `yaml:"-"`
}

// InputChecker indicates what check input need
type InputChecker struct {
	Confirm  BoolInputer
	Test     BoolInputer
	TestFunc func() error
}

// String implement Stringer
func (c CredentialConfig) String() string {
	return fmt.Sprintf("ak: %s, sk: %s, url: %s://%s:%s",
		c.AccessKeyID, c.SecretAccessKey, c.Protocol, c.Host, c.Port)
}

// Format credential config
func (c CredentialConfig) Format() ([]byte, error) {
	return yaml.Marshal(c)
}

// ConfigOption is opts to init config
type ConfigOption func(config *CredentialConfig)

// newCredentialConfig init a credential config
func newCredentialConfig(opts ...ConfigOption) CredentialConfig {
	c := DefaultCredentialConfig()
	c.addOption(opts...)
	return c
}

// withAccessKeyID set config with given ak
func withAccessKeyID(ak string) ConfigOption {
	return func(c *CredentialConfig) {
		if ak != "" {
			c.AccessKeyID = ak
		}
	}
}

// withSecretAccessKey set config with given sk
func withSecretAccessKey(sk string) ConfigOption {
	return func(c *CredentialConfig) {
		if sk != "" {
			c.SecretAccessKey = sk
		}
	}
}

// withHost set config with given host
func withHost(h string) ConfigOption {
	return func(c *CredentialConfig) {
		if h != "" {
			c.Host = h
		}
	}
}

// withPort set config with given port
func withPort(p string) ConfigOption {
	return func(c *CredentialConfig) {
		if p != "" {
			c.Port = p
		}
	}
}

// withProtocol set config with given protocol
func withProtocol(p string) ConfigOption {
	return func(c *CredentialConfig) {
		if p != "" {
			c.Protocol = p
		}
	}
}

// withConfirmChecker set config with given confirm checker
func withConfirmChecker(confirm BoolInputer) ConfigOption {
	return func(c *CredentialConfig) {
		c.Checker.Confirm = confirm
	}
}

// withTestChecker set config with given test checker
func withTestChecker(t BoolInputer) ConfigOption {
	return func(c *CredentialConfig) {
		c.Checker.Test = t
	}
}

// withTestFunc set config with given test func
func withTestFunc(f func() error) ConfigOption {
	return func(c *CredentialConfig) {
		c.Checker.TestFunc = f
	}
}

// withWriteCloser set config with given io.WriteCloser
func withWriteCloser(w io.WriteCloser) ConfigOption {
	return func(c *CredentialConfig) {
		c.WriteCloser = w
	}
}

// withInputers set config with given StringInputers
func withInputers(inputers ...StringInputer) ConfigOption {
	return func(c *CredentialConfig) {
		c.Inputers = append(c.Inputers, inputers...)
	}
}

// addOption exec config options to set up config
func (c *CredentialConfig) addOption(opts ...ConfigOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// DefaultCredentialConfig setup CredentialConfig and return the struct
func DefaultCredentialConfig() CredentialConfig {
	return CredentialConfig{
		AccessKeyID:     "",
		SecretAccessKey: "",
		Host:            constants.DefaultHost,
		Port:            constants.DefaultPort,
		Protocol:        constants.DefaultProtocol,

		Inputers: make([]StringInputer, 0),
		Checker:  InputChecker{},
	}
}

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

// setupConfig runs the whole process to set up credential config
func (c *CredentialConfig) setupConfig() (err error) {
	var confirm, test bool
	for !confirm {
		for _, inputer := range c.Inputers {
			if err = inputer.SetStringP(inputer.GetStringP()); err != nil {
				return
			}
		}

		if c.Checker.Test != nil {
			test, err = c.Checker.Test.GetBool()
			if err != nil {
				return err
			}
			if test && c.Checker.TestFunc != nil {
				var out []byte
				out, err = yaml.Marshal(c)
				if err != nil {
					return fmt.Errorf("marshal config failed: [%w]", err)
				}

				viper.SetConfigType("yaml")
				if err = viper.ReadConfig(bytes.NewBuffer(out)); err != nil {
					return fmt.Errorf("read config from struct failed: [%w]", err)
				}

				if err = c.Checker.TestFunc(); err != nil {
					i18n.Sprintf("test failed with error: [%s]", err.Error())
					err = nil
				} else {
					i18n.Sprintf("test succeed")
				}
			}
		}

		// if no confirm checker set, treat as confirm
		if c.Checker.Confirm == nil {
			confirm = true
			break
		}
		confirm, err = c.Checker.Confirm.GetBool()
		if err != nil {
			return err
		}
	}

	b, err := c.Format()
	if err != nil {
		return
	}
	defer c.WriteCloser.Close()
	if _, err = c.WriteCloser.Write(b); err != nil {
		return
	}
	return
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

// SetupConfigInteractive do the config initial, and then call setup config to set it up
func SetupConfigInteractive() (fileName string, err error) {
	in := newCredentialConfig(
		withAccessKeyID(viper.GetString(constants.ConfigAccessKeyID)),
		withSecretAccessKey(viper.GetString(constants.ConfigSecretAccessKey)),
		withProtocol(viper.GetString(constants.ConfigProtocol)),
		withHost(viper.GetString(constants.ConfigHost)),
		withPort(viper.GetString(constants.ConfigPort)),
	)

	// prepare for the writer
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// write config into preference file
	fileName = filepath.Join(homeDir, ".qingstor", "config.yaml")
	if err = os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return "", err
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}

	// add more options
	in.addOption(
		withInputers(
			PromptInput{
				Key:      constants.ConfigAccessKeyID,
				Value:    &in.AccessKeyID,
				MsgTitle: "AccessKey",
			},
			// ask for secret access key
			PromptInput{
				Key:      constants.ConfigSecretAccessKey,
				Value:    &in.SecretAccessKey,
				MsgTitle: "SecretAccessKey",
			},
			// ask for host
			PromptInput{
				Key:      constants.ConfigHost,
				Value:    &in.Host,
				MsgTitle: "Host",
			},
			// ask for port
			PromptInput{
				Key:      constants.ConfigPort,
				Value:    &in.Port,
				MsgTitle: "Port",
			},
			// ask for protocol
			PromptInput{
				Key:      constants.ConfigProtocol,
				Value:    &in.Protocol,
				MsgTitle: "Protocol",
			},
		),
		withConfirmChecker(
			PromptConfirm{
				Msg: "Confirm your config?",
			},
		),
		withTestChecker(
			PromptConfirm{
				Msg: "Test access with your config?",
			},
		),
		withWriteCloser(f),
		withTestFunc(func() error {
			rootTask := taskutils.NewAtServiceTask(10)
			err = utils.ParseAtServiceInput(rootTask)
			if err != nil {
				return err
			}

			t := task.NewListStorage(rootTask)
			t.SetZone("")
			t.SetStoragerFunc(func(storage.Storager) {})
			t.Run()

			if t.GetFault().HasError() {
				return t.GetFault()
			}
			return nil
		}),
	)
	i18n.Printf("Enter new values or accept defaults in brackets with Enter.\n")
	i18n.Printf("Access key and Secret key are your identifiers for QingStor. Leave them empty if you want to use the env variables.\n")
	if err = in.setupConfig(); err != nil {
		return
	}
	return
}

// IsInteractiveEnable checks whether qsctl run interactively by
// checking stdin and stdout is terminal or not
func IsInteractiveEnable() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd())) && terminal.IsTerminal(int(os.Stdin.Fd()))
}
