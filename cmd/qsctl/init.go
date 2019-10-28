package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qsctl/v2/constants"
)

// register available flag vars here
var (
	// bench will be set if bench flag was set
	bench bool
	// configPath will be set if config flag was set
	configPath string
)

// rootCmd is the main command of qsctl
var rootCmd = &cobra.Command{
	Use:     constants.Name,
	Long:    constants.Description,
	Version: constants.Version,
}

func init() {
	initGlobalFlag()
	// init flags for every single cmd
	initCpFlag()
	initLsFlag()
	initMbFlag()
	initPresignFlag()
	initRbFlag()
	initRmFlag()
	initStatFlag()
	initTeeFlag()

	// init config before command run
	rootCmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		return initConfig()
	}

	// add sub-command to rootCmd
	rootCmd.AddCommand(CatCommand)
	rootCmd.AddCommand(CpCommand)
	rootCmd.AddCommand(LsCommand)
	rootCmd.AddCommand(MbCommand)
	rootCmd.AddCommand(PresignCommand)
	rootCmd.AddCommand(RbCommand)
	rootCmd.AddCommand(RmCommand)
	rootCmd.AddCommand(StatCommand)
	rootCmd.AddCommand(TeeCommand)
}

func initConfig() (err error) {
	// Allow viper read from env.
	viper.SetEnvPrefix("qsctl")
	viper.AutomaticEnv()

	// Set default value for config.
	viper.SetDefault(constants.ConfigHost, constants.DefaultHost)
	viper.SetDefault(constants.ConfigPort, constants.DefaultPort)
	viper.SetDefault(constants.ConfigProtocol, constants.DefaultProtocol)
	viper.SetDefault(constants.ConfigConnectionRetries, constants.DefaultConnectionRetries)
	viper.SetDefault(constants.ConfigLogLevel, constants.DefaultLogLevel)

	// Load config from config file.
	if configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configPath)
	} else {
		// Search config in home directory with name ".qingstor" (without extension).
		viper.AddConfigPath("$HOME/.qingstor")
		// Search config in XDG style.
		viper.AddConfigPath("$HOME/.config/qingstor")
		// Read from "/etc/qingstor" instead if not found.
		viper.AddConfigPath("/etc/qingstor")
		viper.SetConfigName("config")
	}

	err = viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warnf("Config not loaded, use default and environment value instead.")
			if viper.GetString(constants.ConfigAccessKeyID) == "" && viper.GetString(constants.ConfigSecretAccessKey) == "" {
				log.Infof("AK and SK not found. Please setup your config now.")
				fileName, err := SetupConfigInteractive()
				if err != nil {
					return err
				}
				log.Infof("Your config has been set to <%v>.", fileName)
				viper.SetConfigFile(fileName)
				if err = viper.ReadInConfig(); err != nil {
					return err
				}
			}
			err = nil
		default:
			log.Errorf("Load config failed [%v]", err)
			return
		}
	}

	if viper.GetString(constants.ConfigLogLevel) != "" {
		lvl, err := log.ParseLevel(viper.GetString(constants.ConfigLogLevel))
		if err != nil {
			log.Errorf("Parse log level failed [%v]", err)
			return err
		}
		log.SetLevel(lvl)
	}

	return nil
}

func initGlobalFlag() {
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		"", "assign config path manually")
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().BoolVar(&bench, constants.BenchFlag,
		false, "enable benchmark or not")
	// Overwrite the default help flag to free -h shorthand.
	rootCmd.PersistentFlags().Bool("help", false, "help for this command")
}

type inputConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`

	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
	LogLevel string `yaml:"log_level"`
}

// SetupConfigInteractive setup input config interactively
func SetupConfigInteractive() (fileName string, err error) {
	in := newInputConfig()

	// multiple times to set config if not confirmed
	for {
		fmt.Printf("Access Key ID: ")
		if _, err = fmt.Scanf("%s", &in.AccessKeyID); err != nil {
			return "", err
		}

		fmt.Printf("Secret Access Key: ")
		if _, err = fmt.Scanf("%s", &in.SecretAccessKey); err != nil {
			return "", err
		}

		if !confirm("Do you apply qsctl for QingStor public cloud") {
			fmt.Printf("Host: ")
			if _, err = fmt.Scanf("%s", &in.Host); err != nil {
				return "", err
			}

			fmt.Printf("Protocol [http/https]: ")
			if _, err = fmt.Scanf("%s", &in.Protocol); err != nil {
				return "", err
			}

			fmt.Printf("Port: ")
			if _, err = fmt.Scanf("%s", &in.Port); err != nil {
				return "", err
			}
		}

		fmt.Printf("Log level [low to high: debug/info/warn/error/fatal]: ")
		if _, err = fmt.Scanf("%s", &in.LogLevel); err != nil {
			return "", err
		}

		b, err := yaml.Marshal(in)
		if err != nil {
			return "", err
		}

		fmt.Println("Your config listed below:")
		fmt.Printf("\n%v\n", string(b))

		if confirm("Confirm") {
			break
		}
		// reset the input
		in = newInputConfig()
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

// newInputConfig setup inputConfig and return the struct
func newInputConfig() inputConfig {
	return inputConfig{
		AccessKeyID:     "",
		SecretAccessKey: "",
		Host:            constants.DefaultHost,
		Port:            constants.DefaultPort,
		Protocol:        constants.DefaultProtocol,
		LogLevel:        constants.DefaultLogLevel,
	}
}

func confirm(tip string) bool {
	var out string
	fmt.Printf("%s? [yes/no] ", tip)

	if _, err := fmt.Scanf("%s", &out); err != nil {
		panic(err)
	}

	return strings.ToLower(out[:1]) == "y"
}
