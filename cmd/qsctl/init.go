package main

import (
	"github.com/Xuanwo/storage/services/qingstor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/constants"
)

var (
	// register available flag vars here
	bench                bool
	expectSize           string
	humanReadable        bool
	longFormat           bool
	maximumMemoryContent string
	recursive            bool
	reverse              bool
	zone                 string
)

var (
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
	viper.SetDefault(constants.ConfigHost, "qingstor.com")
	viper.SetDefault(constants.ConfigPort, 443)
	viper.SetDefault(constants.ConfigProtocol, "https")
	viper.SetDefault(constants.ConfigConnectionRetries, 3)
	viper.SetDefault(constants.ConfigLogLevel, "info")

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

// StorageOption is the alias for config register func
type StorageOption func(*qingstor.Config)

// NewQingstorConfig will collect option func and conduct a pointer to qingstor.Config
func NewQingstorConfig(opt ...StorageOption) *qingstor.Config {
	r := new(qingstor.Config)
	for _, o := range opt {
		o(r)
	}
	return r
}

// WriteBase will return a conduct func for base info
func WriteBase() StorageOption {
	return func(c *qingstor.Config) {
		c.AccessKeyID = viper.GetString(constants.ConfigAccessKeyID)
		c.SecretAccessKey = viper.GetString(constants.ConfigSecretAccessKey)
		c.Host = viper.GetString(constants.ConfigHost)
		c.Port = viper.GetInt(constants.ConfigPort)
		c.Protocol = viper.GetString(constants.ConfigProtocol)
	}
}

// WriteBucketName will return a conduct func for bucket name
func WriteBucketName(name string) StorageOption {
	return func(c *qingstor.Config) {
		c.BucketName = name
	}
}

// WriteZone will return a conduct func for zone
func WriteZone(z string) StorageOption {
	return func(c *qingstor.Config) {
		c.Zone = z
	}
}
