package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/utils"
)

// MbCommandFlags records all flags for MbCommand
var MbCommandFlags = FlagSet{}

var (
	zoneFlagInfo = NewStringCtlFlag(
		zoneFlag,
		"z",
		"In which zone to do the operation",
		"",
	)
)

// MbCommand will handle make bucket command.
var MbCommand = &cobra.Command{
	Use:   "mb [qs://]<bucket-name>",
	Short: "make a new bucket",
	Long: `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`,
	Example: utils.AlignPrintWithColon(
		"Make bucket: qsctl mb bucket-name",
	),
	Args: cobra.ExactArgs(1),
	RunE: mbRun,
}

func mbRun(_ *cobra.Command, args []string) (err error) {
	return action.MakeBucket(args[0])
}

func initMbCommandFlag() {
	addFlagToMbCommand()
	if flag, ok := MbCommandFlags[zoneFlag]; ok {
		MbCommand.PersistentFlags().StringVarP(flag.(StringCtlFlag).StringVarP(&zone))
	}
	// register MbCommandFlags to cmd-flag map
	cmdToFlagSet.AddFlagSet(MbCommand.Name(), &MbCommandFlags)
}

func addFlagToMbCommand() {
	MbCommandFlags.AddFlag(zoneFlag, zoneFlagInfo.SetRequired())
}
