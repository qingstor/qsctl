package cmd

import (
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/utils"
)

var (
	expectSize           string
	maximumMemoryContent string
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
	PreRun: func(cmd *cobra.Command, args []string) {
		var v datasize.ByteSize
		err := v.UnmarshalText([]byte(expectSize))
		if err != nil {
			log.Fatalf("Input expect size <%s> is invalid", expectSize)
		}
		contexts.ExpectSize = int64(v)

		err = v.UnmarshalText([]byte(maximumMemoryContent))
		if err != nil {
			log.Fatalf("Input expect size <%s> is invalid", maximumMemoryContent)
		}
		contexts.MaximumMemoryContent = int64(v)
	},
	RunE: cpRun,
}

func cpRun(_ *cobra.Command, args []string) (err error) {
	return action.Copy(args[0], args[1])
}

func init() {
	CpCommand.PersistentFlags().StringVar(
		&expectSize,
		"expect-size",
		"",
		`expected size of the input file
accept: 100MB, 1.8G
(only used for input from stdin)`)
	CpCommand.MarkPersistentFlagRequired("expect-size")

	CpCommand.PersistentFlags().StringVar(
		&maximumMemoryContent,
		"maximum-memory-content",
		"",
		"maximum content loaded in memory \n (only used for input from stdin)")
}
