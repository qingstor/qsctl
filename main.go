package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/cmd"
	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
)

var application = &cobra.Command{
	Use:  constants.Name,
	Long: constants.Description,
}

var (
	configPath string
)

func init() {
	cobra.OnInitialize(initConfig)

	application.AddCommand(cmd.CatCommand)
	application.AddCommand(cmd.CpCommand)
	application.AddCommand(cmd.TeeCommand)

	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().BoolVar(&contexts.Bench, "bench", false, "enable benchmark or not")
}

func initConfig() {
	// Allow viper read from env.
	viper.SetEnvPrefix("qsctl")
	viper.AutomaticEnv()

	// Load config from config file.
	if configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configPath)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".qingstor" (without extension).
		viper.AddConfigPath(home + "/" + ".qingstor")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Set default value for config.
	viper.SetDefault(constants.ConfigHost, "qingstor.com")
	viper.SetDefault(constants.ConfigPort, 443)
	viper.SetDefault(constants.ConfigProtocol, "https")
	viper.SetDefault(constants.ConfigConnectionRetries, 3)
	viper.SetDefault(constants.ConfigLogLevel, "info")
}

func main() {
	err := application.Execute()
	if err != nil {
		os.Exit(1)
	}
}
