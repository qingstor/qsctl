package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		setLogLevel()
		setupViper()
		// if execute configure cmd, not read config into viper (also not read from env)
		if c.Name() == ConfigureCommand.Name() {
			return nil
		}
		setupEnvForViper()
		return initConfig()
	}

	// add sub-command to rootCmd
	rootCmd.AddCommand(CatCommand)
	rootCmd.AddCommand(ConfigureCommand)
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
		`{{with .Name}}{{printf "%%s " .}}{{end}}{{printf "version %%s\n" .Version}}`,
	))
	rootCmd.SetHelpTemplate(i18n.Sprintf(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`))
}

func initConfig() (err error) {
	// overwrite zone config by flag if set
	if zone != "" {
		viper.Set(constants.ConfigZone, zone)
	}

	// try to read config from path set in setupViper()
	err = viper.ReadInConfig()
	if err == nil {
		if configuredByEnv() {
			log.Debug("Ak sk loaded from env.")
		}
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
	return err
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
	return getEnvWithPrefix(constants.EnvPrefix, constants.ConfigAccessKeyID) != "" &&
		getEnvWithPrefix(constants.EnvPrefix, constants.ConfigSecretAccessKey) != ""
}

func setLogLevel() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.PanicLevel)
	}
}

func setupViper() {
	// Set default value for config.
	viper.SetDefault(constants.ConfigHost, constants.DefaultHost)
	viper.SetDefault(constants.ConfigPort, constants.DefaultPort)
	viper.SetDefault(constants.ConfigProtocol, constants.DefaultProtocol)

	// Load config from config file.
	if configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configPath)
	} else {
		// Search config in home directory with name ".qingstor".
		viper.AddConfigPath("$HOME/.qingstor")
		// Search config in XDG style.
		viper.AddConfigPath("$HOME/.config/qingstor")
		// Read from "/etc/qingstor" instead if not found.
		viper.AddConfigPath("/etc/qingstor")
		// Read config from configPath/config.xxx (without extension)
		viper.SetConfigName("config")
	}
}

func setupEnvForViper() {
	// Allow viper read from env.
	viper.SetEnvPrefix("qsctl")
	viper.AutomaticEnv()
}

func getEnvWithPrefix(prefix, key string) string {
	return os.Getenv(strings.ToUpper(prefix + "_" + key))
}
