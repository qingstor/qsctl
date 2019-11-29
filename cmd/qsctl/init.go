package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yunify/qsctl/v2/cmd/utils"
	"github.com/yunify/qsctl/v2/constants"
)

//go:generate go run ../../internal/cmd/generator/i18nextract
//go:generate go run ../../internal/cmd/generator/i18ngenerator

// register available flag vars here
var (
	// bench will be set if bench flag was set
	bench bool
	// configPath will be set if config flag was set
	configPath string
	debug      bool
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
	initMvFlag()
	initPresignFlag()
	initRbFlag()
	initRmFlag()
	initStatFlag()
	initSyncFlag()
	initTeeFlag()

	// init config before command run
	rootCmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.PanicLevel)
		}

		return initConfig()
	}

	// add sub-command to rootCmd
	rootCmd.AddCommand(CatCommand)
	rootCmd.AddCommand(CpCommand)
	rootCmd.AddCommand(LsCommand)
	rootCmd.AddCommand(MbCommand)
	rootCmd.AddCommand(MvCommand)
	rootCmd.AddCommand(PresignCommand)
	rootCmd.AddCommand(RbCommand)
	rootCmd.AddCommand(RmCommand)
	rootCmd.AddCommand(StatCommand)
	rootCmd.AddCommand(SyncCommand)
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

	// try to read config from path set above
	err = viper.ReadInConfig()
	if err == nil {
		return
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		fmt.Printf("Load config failed [%v]", err)
		return
	}

	// if env not set, start interactive setup
	if viper.GetString(constants.ConfigAccessKeyID) == "" && viper.GetString(constants.ConfigSecretAccessKey) == "" {
		fmt.Printf("AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.")
		fileName, err := utils.SetupConfigInteractive()
		if err != nil {
			return fmt.Errorf("setup config failed [%v], please try again", err)
		}
		fmt.Printf("Your config has been set to <%v>. You can still modify it manually.", fileName)
		viper.SetConfigFile(fileName)
		if err = viper.ReadInConfig(); err != nil {
			return err
		}
	} else {
		fmt.Printf("Config not loaded, use default and environment value instead.")
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
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print logs for debug")
}
