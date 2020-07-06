package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cutil "github.com/qingstor/qsctl/v2/cmd/utils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

// ConfigureCommand will handle configure command.
var ConfigureCommand = &cobra.Command{
	Use:   "configure",
	Short: i18n.Sprintf("configure interactively"),
	Long:  i18n.Sprintf("qsctl configure can invoke interactive (re)configuration tool"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Set config: qsctl configure"),
	),
	Args: cobra.ExactArgs(0),
	RunE: configureRun,
}

func configureRun(c *cobra.Command, _ []string) error {
	silenceUsage(c)
	// if not run interactively, return error
	if !cutil.IsInteractiveEnable() {
		log.Errorf("qsctl not run interactively")
		return fmt.Errorf(i18n.Sprintf("qsctl not run interactively"))
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			i18n.Printf("Try to load config failed [%v]\n", err)
			return err
		}
		err = nil // ignore config file not found error
	}

	filename, err := cutil.SetupConfigInteractive()
	if err != nil {
		return err
	}
	fmt.Println(filename)

	// i18n.Printf("Your config has been set to <%v>. You can still modify it manually.", configFile)
	return nil
}
