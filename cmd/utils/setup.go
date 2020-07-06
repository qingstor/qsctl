package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

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

	inputers   []StringInputer `yaml:"-"`
	checker    InputChecker    `yaml:"-"`
	writerFunc func() error    `yaml:"-"`
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

// FormatType return type of credential config when format
func (CredentialConfig) FormatType() string {
	return "yaml"
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
		c.checker.Confirm = confirm
	}
}

// withTestChecker set config with given test checker
func withTestChecker(t BoolInputer) ConfigOption {
	return func(c *CredentialConfig) {
		c.checker.Test = t
	}
}

// withTestFunc set config with given test func
func withTestFunc(f func() error) ConfigOption {
	return func(c *CredentialConfig) {
		c.checker.TestFunc = f
	}
}

// withWriteFunc set config with given func
func withWriteFunc(w func() error) ConfigOption {
	return func(c *CredentialConfig) {
		c.writerFunc = w
	}
}

// withInputers set config with given StringInputers
func withInputers(inputers ...StringInputer) ConfigOption {
	return func(c *CredentialConfig) {
		c.inputers = append(c.inputers, inputers...)
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

		inputers: make([]StringInputer, 0),
		checker:  InputChecker{},
		writerFunc: func() error {
			return nil
		},
	}
}

// setupConfig runs the whole process to set up credential config
// 1. call all inpupters to input config attr
// 2. if test check was set, try to ask to test. Then if get true, load config and run test func.
// 3. if confirm check was not set, handle as always confirm. Otherwise, ask to confirm.
// 4. if not confirmed, back to step 1.
// 5. if confirmed, format config and call write func
func (c *CredentialConfig) setupConfig() (err error) {
	var confirm, test bool
	for !confirm {
		for _, inputer := range c.inputers {
			if err = inputer.SetStringP(inputer.GetStringP()); err != nil {
				return fmt.Errorf("set string for %s failed: [%w]", inputer.GetName(), err)
			}
		}

		if c.checker.Test != nil {
			test, err = c.checker.Test.GetBool()
			if err != nil {
				return fmt.Errorf("get bool for %s failed: [%w]", c.checker.Test.GetName(), err)
			}
			if test && c.checker.TestFunc != nil {
				var out []byte
				out, err = c.Format()
				if err != nil {
					return fmt.Errorf("marshal config failed when test: [%w]", err)
				}

				viper.SetConfigType(c.FormatType())
				if err = viper.ReadConfig(bytes.NewBuffer(out)); err != nil {
					return fmt.Errorf("read config from struct failed when test: [%w]", err)
				}

				if err = c.checker.TestFunc(); err != nil {
					i18n.Printf("Test failed with error: [%s]\n", err.Error())
					err = nil
				} else {
					i18n.Printf("Test connection succeed\n")
				}
			}
		}

		// if no confirm checker set, treat as confirm
		if c.checker.Confirm == nil {
			confirm = true
			break
		}
		confirm, err = c.checker.Confirm.GetBool()
		if err != nil {
			return fmt.Errorf("get bool for %s failed: [%w]", c.checker.Confirm.GetName(), err)
		}
	}

	if err = c.writerFunc(); err != nil {
		return fmt.Errorf("call write func failed: [%w]", err)
	}
	return
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
				Key: "confirm",
				Msg: "Confirm your config?",
			},
		),
		withTestChecker(
			PromptConfirm{
				Key: "test",
				Msg: "Test access with your config?",
			},
		),
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
		withWriteFunc(func() (err error) {
			// prepare for the writer
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("get home dir failed: [%w]", err)
			}

			// write config into preference file
			fileName = filepath.Join(homeDir, ".qingstor", "config.yaml")
			if err = os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
				return fmt.Errorf("make dir for <%s> failed: [%w]", fileName, err)
			}

			f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("open file for <%s> failed: [%w]", fileName, err)
			}
			defer f.Close()

			b, err := in.Format()
			if err != nil {
				return fmt.Errorf("format config failed: [%w]", err)
			}
			if _, err = f.Write(b); err != nil {
				return fmt.Errorf("write config failed: [%w]", err)
			}
			return nil
		}),
	)
	i18n.Printf("Enter new values or accept defaults in brackets with Enter.\n")
	i18n.Printf("Access key and Secret key are your identifiers for QingStor. Leave them empty if you want to use the env variables.\n")
	if err = in.setupConfig(); err != nil {
		return fileName, fmt.Errorf("setup config failed: [%w]", err)
	}
	return
}

// IsInteractiveEnable checks whether qsctl run interactively by
// checking stdin and stdout is terminal or not
func IsInteractiveEnable() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd())) && terminal.IsTerminal(int(os.Stdin.Fd()))
}
