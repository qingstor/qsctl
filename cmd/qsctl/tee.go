package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/utils"
)

var teeInput struct {
	ExpectSize string
	MaxMemory  string
}

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

func teeRun(_ *cobra.Command, _ []string) (err error) {
	return nil
}

func initTeeFlag() {
	TeeCommand.PersistentFlags().StringVar(&teeInput.ExpectSize,
		constants.ExpectSizeFlag,
		"",
		"expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)",
	)
	TeeCommand.PersistentFlags().StringVar(&teeInput.MaxMemory,
		constants.MaximumMemoryContentFlag,
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
}

func validateTeeFlag(_ *cobra.Command, _ []string) (err error) {
	// TODO: parse should be moved into teeParse func
	if teeInput.ExpectSize != "" {
		teeExpectSize, err := utils.ParseByteSize(teeInput.ExpectSize)
		_ = teeExpectSize
		if err != nil {
			return err
		}
	}

	if teeInput.MaxMemory != "" {
		teeMaxMemory, err := utils.ParseByteSize(teeInput.MaxMemory)
		_ = teeMaxMemory
		if err != nil {
			return err
		}
	}
	return nil
}
