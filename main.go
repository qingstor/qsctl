package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/v2/cmd"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

var application = &cobra.Command{
	Use:     constants.Name,
	Long:    constants.Description,
	Version: constants.Version,
}

var (
	configPath string
)

func init() {
	// Set log formatter firstly.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	application.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		err := initConfig()
		if err != nil {
			return err
		}

		return cmd.ParseFlagIntoContexts(c, args)
	}

	application.AddCommand(cmd.CatCommand)
	application.AddCommand(cmd.CpCommand)
	application.AddCommand(cmd.StatCommand)
	application.AddCommand(cmd.TeeCommand)
	application.AddCommand(cmd.MbCommand)
	application.AddCommand(cmd.RbCommand)
	application.AddCommand(cmd.RmCommand)

	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().BoolVar(&contexts.Bench, "bench", false, "enable benchmark or not")
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
	// Setup global qs service
	if err := contexts.SetupServices(); err != nil {
		log.Errorf("Setup up service failed [%v]", err)
		return err
	}

	return nil
}

func main() {
	err := application.Execute()
	if err != nil {
		os.Exit(1)
	}
}
