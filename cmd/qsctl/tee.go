package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/constants"

	"github.com/yunify/qsctl/v2/utils"
)

var (
	teeExpectSize int64
	teeMaxMemory  int64
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
	return nil
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
		constants.MaximumMemoryContentFlag,
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
}

func validateTeeFlag(_ *cobra.Command, _ []string) (err error) {
	if expectSize != "" {
		teeExpectSize, err = utils.ParseByteSize(expectSize)
		if err != nil {
			return err
		}
	}

	if maximumMemoryContent != "" {
		teeMaxMemory, err = utils.ParseByteSize(maximumMemoryContent)
		if err != nil {
			return err
		}
	}
	return nil
}
