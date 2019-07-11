package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/utils"
)

// CpCommandFlags records all flags for CpCommand
var CpCommandFlags = FlagSet{}

var (
	expectSizeFlagInfo = NewStringCtlFlag(
		expectSizeFlag,
		"",
		`expected size of the input file
accept: 100MB, 1.8G
(only used for input from stdin)`,
		"",
	)

	maximumMemoryContentFlagInfo = NewStringCtlFlag(
		maximumMemoryContentFlag,
		"",
		"maximum content loaded in memory \n (only used for input from stdin)",
		"",
	)
)

// CpCommand will handle copy command.
var CpCommand = &cobra.Command{
	Use:   "cp <source-path> <dest-path>",
	Short: "copy from/to qingstor",
	Long:  "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout",
	Example: utils.AlignPrintWithColon(
		"Copy file: qsctl cp /path/to/file qs://prefix/a",
		"Copy folder: qsctl cp qs://prefix/a /path/to/folder -r",
		"Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin",
		"Write to stdout: qsctl cp qs://prefix/b - > /path/to/file",
	),
	Args: cobra.ExactArgs(2),
	RunE: cpRun,
}

func cpRun(_ *cobra.Command, args []string) (err error) {
	return action.Copy(args[0], args[1])
}

func initCpCommandFlag() {
	addFlagToCpCommand()
	if flag, ok := CpCommandFlags[expectSizeFlag]; ok {
		CpCommand.PersistentFlags().StringVar(flag.(StringCtlFlag).StringVar(&expectSize))
	}

	if flag, ok := CpCommandFlags[maximumMemoryContentFlag]; ok {
		CpCommand.PersistentFlags().StringVar(flag.(StringCtlFlag).StringVar(&maximumMemoryContent))
	}
	// register CpCommandFlags to cmd-flag map
	cmdToFlagSet.AddFlagSet(CpCommand.Name(), &CpCommandFlags)
}

func addFlagToCpCommand() {
	CpCommandFlags.AddFlag(expectSizeFlag, expectSizeFlagInfo.SetRequired())
	CpCommandFlags.AddFlag(maximumMemoryContentFlag, maximumMemoryContentFlagInfo)
}
