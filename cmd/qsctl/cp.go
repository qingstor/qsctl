package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var cpInput struct {
	ExpectSize           string
	MaximumMemoryContent string
}

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

func initCpFlag() {
	CpCommand.PersistentFlags().StringVar(&cpInput.ExpectSize,
		constants.ExpectSizeFlag,
		"",
		"expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)",
	)
	CpCommand.PersistentFlags().StringVar(&cpInput.MaximumMemoryContent,
		constants.MaximumMemoryContentFlag,
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
}

func cpRun(_ *cobra.Command, args []string) (err error) {
	rootTask := &taskutils.BetweenStorageTask{}
	err = utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	t := task.NewCopyFile(rootTask)

	t.Run()
	if t.ValidateFault() {
		return t.GetFault()
	}
	return
}
