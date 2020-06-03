package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qingstor/qsctl/v2/cmd/utils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
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
	// noProgress will be set if no-progress flag was set
	noProgress bool
	// zone will be set if zone flag was set
	zone string
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

	rootCmd.SetVersionTemplate(i18n.Sprintf(
		`{{with .Name}}{{printf "%%s " .}}{{end}}{{printf "version %%s" .Version}}`,
	))
	rootCmd.SetHelpTemplate(i18n.Sprintf(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`))
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

	// if zone flag was set, overwrite the config
	if zone != "" {
		viper.Set(constants.ConfigZone, zone)
	}

	// try to read config from path set above
	err = viper.ReadInConfig()
	if err == nil {
		log.Debugf("Load config success from [%s]: %v", viper.ConfigFileUsed(), viper.AllSettings())
		return nil
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		i18n.Printf("Load config failed [%v]", err)
		return err
	}

	// if env set, get config from env
	if configuredByEnv() {
		i18n.Printf("Config not loaded, use default and environment value instead.")
		log.Debug("Config not loaded, use default and environment value instead.")
		return nil
	}

	// if env not set, try to start interactive setup
	// if not run interactively, return error
	if !utils.IsInteractiveEnable() {
		log.Errorf("qsctl not run interactively, and cannot load config with err: [%v]", err)
		return err
	}

	i18n.Printf("AccessKey and SecretKey not found. Please setup your config now, or exit and setup manually.")
	log.Debug("AccessKey and SecretKey not found. Ready to turn into setup config interactively.")
	var fileName string
	fileName, err = utils.SetupConfigInteractive()
	if err != nil {
		return fmt.Errorf("setup config failed [%v], please try again", err)
	}

	i18n.Printf("Your config has been set to <%v>. You can still modify it manually.", fileName)
	viper.SetConfigFile(fileName)
	log.Debugf("Config was set to [%s]", fileName)
	// read in config again after interactively setup config file
	if err = viper.ReadInConfig(); err != nil {
		log.Errorf("Read config after interactively setup failed: [%v]", err)
		return err
	}
	return nil
}

func initGlobalFlag() {
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		"", i18n.Sprintf("assign config path manually"))
	// Add zone flag which can be used in all sub commands.
	rootCmd.PersistentFlags().StringVarP(&zone, constants.ConfigZone, "z",
		"", i18n.Sprintf("in which zone to do the operation"))
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().BoolVar(&bench, constants.BenchFlag,
		false, i18n.Sprintf("enable benchmark or not"))
	rootCmd.PersistentFlags().BoolVar(&noProgress, constants.NoProgressFlag,
		false, i18n.Sprintf("disable progress bar display or not"))
	// Overwrite the default help flag to free -h shorthand.
	rootCmd.PersistentFlags().Bool("help", false, i18n.Sprintf("help for this command"))
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, i18n.Sprintf("print logs for debug"))
}

func silenceUsage(c *cobra.Command) {
	c.SilenceUsage = true
}

// configuredByEnv returns true if either ak or sk set
func configuredByEnv() bool {
	return viper.GetString(constants.ConfigAccessKeyID) != "" ||
		viper.GetString(constants.ConfigSecretAccessKey) != ""
}
