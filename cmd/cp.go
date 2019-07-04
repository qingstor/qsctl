package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/utils"
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

func init() {
	// TODO: support input x MB
	CpCommand.PersistentFlags().Int64Var(
		&contexts.ExpectSize,
		"expect-size",
		0,
		"expected size of the input file \n (only used for input from stdin)")
	CpCommand.PersistentFlags().Int64Var(
		&contexts.MaximumMemoryContent,
		"maximum-memory-content",
		constants.DefaultMaximumMemoryContent,
		"maximum content loaded in memory \n (only used for input from stdin)")
}
