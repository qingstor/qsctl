package main

import (
	"context"
	"encoding/json"

	"github.com/qingstor/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/internal/pkg/ilog"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
)

//go:generate go run ../../internal/cmd/generator/i18nextract
//go:generate go run ../../internal/cmd/generator/i18ngenerator

var globalFlag = globalFlags{}

// register available flag vars here
type globalFlags struct {
	// bench will be set if bench flag was set
	bench bool
	// configPath will be set if config flag was set
	configPath string
	debug      bool
	// zone will be set if zone flag was set
	zone string
}

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
		logger := ilog.InitLoggerWithDebug(globalFlag.debug)
		ctx := log.ContextWithLogger(c.Context(), logger)
		c.SetContext(ctx)
		return initConfig(ctx)
	}

	rootCmd.AddCommand(ShellCommand)
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

	rootCmd.SetVersionTemplate(i18n.Sprintf(constants.VersionTemplate))
	rootCmd.SetUsageTemplate(i18n.Sprintf(constants.UsageTemplate))
	rootCmd.SetHelpTemplate(i18n.Sprintf(constants.HelpTemplate))
}

func initConfig(ctx context.Context) (err error) {
	logger := log.FromContext(ctx)

	// Allow viper read from env.
	viper.SetEnvPrefix("qsctl")
	viper.AutomaticEnv()

	// Set default value for config.
	viper.SetDefault(constants.ConfigHost, constants.DefaultHost)
	viper.SetDefault(constants.ConfigPort, constants.DefaultPort)
	viper.SetDefault(constants.ConfigProtocol, constants.DefaultProtocol)

	// Load config from config file.
	if globalFlag.configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(globalFlag.configPath)
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
	if globalFlag.zone != "" {
		viper.Set(constants.ConfigZone, globalFlag.zone)
	}

	// try to read config from path set above
	err = viper.ReadInConfig()
	if err == nil {
		content, _ := json.Marshal(viper.AllSettings())
		logger.Debug(
			log.String("action", "load config success"),
			log.String("config_path", viper.ConfigFileUsed()),
			log.Bytes("content", content),
		)
		return nil
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		i18n.Printf("Load config failed [%v]", err)
		return err
	}

	// if env set, get config from env
	if configuredByEnv() {
		i18n.Printf("Config not loaded, use default and environment value instead.")
		logger.Debug(
			log.String("action", "load config from env"),
		)
		return nil
	}

	return
}

func initGlobalFlag() {
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().StringVarP(&globalFlag.configPath, "config", "c",
		"", i18n.Sprintf("assign config path manually"))
	// Add zone flag which can be used in all sub commands.
	rootCmd.PersistentFlags().StringVarP(&globalFlag.zone, constants.ConfigZone, "z",
		"", i18n.Sprintf("in which zone to do the operation"))
	// Add config flag which can be used in all sub commands.
	rootCmd.PersistentFlags().BoolVar(&globalFlag.bench, constants.BenchFlag,
		false, i18n.Sprintf("enable benchmark or not"))
	// Overwrite the default help flag to free -h shorthand.
	rootCmd.PersistentFlags().Bool("help", false, i18n.Sprintf("help for this command"))
	rootCmd.PersistentFlags().BoolVar(&globalFlag.debug, "debug", false, i18n.Sprintf("print logs for debug"))
}

func silenceUsage(c *cobra.Command) {
	c.SilenceUsage = true
}

// configuredByEnv returns true if either ak or sk set
func configuredByEnv() bool {
	return viper.GetString(constants.ConfigAccessKeyID) != "" &&
		viper.GetString(constants.ConfigSecretAccessKey) != ""
}

// resetGlobalFlags used to reset global flags before each sub-command run in shell,
// to avoid flags from last call is still working
func resetGlobalFlags() {
	globalFlag = globalFlags{}
}
