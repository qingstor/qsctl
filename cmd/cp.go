package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/action"
)

// CpCommand will handle copy command.
var CpCommand = &cobra.Command{
	Use:   "cp <source-path> <dest-path>",
	Short: "copy from/to qingstor",
	Long:  "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout",
	Example: `  Copy file:       qsctl cp /path/to/file qs://prefix/a
  Copy folder:     qsctl cp qs://prefix/a /path/to/folder -r
  Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin
  Write to stdout: qsctl cp qs://prefix/b - > /path/to/file
`,
	Args: cobra.ExactArgs(2),
	RunE: run,
}

func run(cmd *cobra.Command, args []string) (err error) {
	err = action.Copy(args[0], args[1])
	if err != nil {
		panic(err)
	}
	return
}
