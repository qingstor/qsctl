package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

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
	Args:    cobra.ExactArgs(1),
	RunE:    teeRun,
	PreRunE: validateTeeFlag,
}

func teeRun(_ *cobra.Command, args []string) (err error) {
	return action.Copy("-", args[0])
}

func initTeeFlag() {
	TeeCommand.PersistentFlags().StringVar(&expectSize,
		"expect-size",
		"",
		"expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)",
	)
	TeeCommand.PersistentFlags().StringVar(&maximumMemoryContent,
		"maximum-memory-content",
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
}

func validateTeeFlag(_ *cobra.Command, _ []string) (err error) {
	if expectSize != "" {
		contexts.ExpectSize, err = utils.ParseByteSize(expectSize)
		if err != nil {
			return
		}
	}

	if maximumMemoryContent != "" {
		contexts.MaximumMemoryContent, err = utils.ParseByteSize(maximumMemoryContent)
		if err != nil {
			return
		}
	}

	return nil
}
