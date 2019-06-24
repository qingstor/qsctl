package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yunify/qsctl/constants"
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

	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config path")
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
}

func main() {
	application.Execute()
}
