package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/utils"
)

// TeeCommandFlags records all flags for TeeCommand
var TeeCommandFlags = FlagSet{}

// TeeCommand will handle tee command.
var TeeCommand = &cobra.Command{
	Use:   "tee qs://<bucket_name>/<object_key>",
	Short: "tee a remote object from stdin",
	Long: `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`,
	Example: utils.AlignPrintWithColon(
		"Tee object: qsctl tee qs://prefix/a",
	),
	Args: cobra.ExactArgs(1),
	RunE: teeRun,
}

func teeRun(_ *cobra.Command, args []string) (err error) {
	return action.Copy("-", args[0])
}

func initTeeCommandFlag() {
	addFlagToTeeCommand()
	if flag, ok := TeeCommandFlags[expectSizeFlag]; ok {
		TeeCommand.PersistentFlags().StringVar(flag.(StringCtlFlag).StringVar(&expectSize))
	}

	if flag, ok := TeeCommandFlags[maximumMemoryContentFlag]; ok {
		TeeCommand.PersistentFlags().StringVar(flag.(StringCtlFlag).StringVar(&maximumMemoryContent))
	}
	// register TeeCommandFlags to cmd-flag map
	cmdToFlagSet.AddFlagSet(TeeCommand.Name(), &TeeCommandFlags)
}

func addFlagToTeeCommand() {
	TeeCommandFlags.AddFlag(expectSizeFlag, expectSizeFlagInfo.SetRequired())
	TeeCommandFlags.AddFlag(maximumMemoryContentFlag, maximumMemoryContentFlagInfo)
}
