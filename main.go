package main

import (
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/cmd"
	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
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
	cobra.OnInitialize(initConfig)

	application.PersistentPreRunE = cmd.ParseFlagIntoContexts

	application.AddCommand(cmd.CatCommand)
	application.AddCommand(cmd.CpCommand)
	application.AddCommand(cmd.StatCommand)
	application.AddCommand(cmd.TeeCommand)
	application.AddCommand(cmd.MbCommand)
	application.AddCommand(cmd.RmCommand)

	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().BoolVar(&contexts.Bench, "bench", false, "enable benchmark or not")
}

func initConfig() {
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
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Errorf("Get homedir failed [%v]", err)
			return
		}

		// Search config in home directory with name ".qingstor" (without extension).
		viper.AddConfigPath(home + "/" + ".qingstor")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Load config failed [%v]", err)
		return
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if viper.GetString(constants.ConfigLogLevel) != "" {
		lvl, err := log.ParseLevel(viper.GetString(constants.ConfigLogLevel))
		if err != nil {
			log.Errorf("Parse log level failed [%v]", err)
			return
		}
		log.SetLevel(lvl)
	}
	// Setup global qs service
	if err := contexts.SetupServices(); err != nil {
		log.Errorf("Setup up service failed [%v]", err)
		return
	}
}

func main() {
	err := application.Execute()
	if err != nil {
		os.Exit(1)
	}
}
