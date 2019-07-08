package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
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
