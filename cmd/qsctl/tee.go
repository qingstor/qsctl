package main

import (
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var teeInput struct {
	ExpectSize string
	MaxMemory  string
}

// TeeCommand will handle tee command.
var TeeCommand = &cobra.Command{
	Use:   "tee qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("tee a remote object from stdin"),
	Long: i18n.Sprintf(`qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Tee object: qsctl tee qs://prefix/a"),
	),
	Args:    cobra.ExactArgs(1),
	RunE:    teeRun,
	PreRunE: validateTeeFlag,
}

func teeRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewBetweenStorageTask(10)
	err = utils.ParseBetweenStorageInput(rootTask, "-", args[0])
	if err != nil {
		return
	}
	t := task.NewCopyStream(rootTask)
	t.SetPartSize(constants.DefaultPartSize)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return nil
}

func initTeeFlag() {
	TeeCommand.PersistentFlags().StringVar(&teeInput.ExpectSize,
		constants.ExpectSizeFlag,
		"",
		i18n.Sprintf("expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)"),
	)
	TeeCommand.PersistentFlags().StringVar(&teeInput.MaxMemory,
		constants.MaximumMemoryContentFlag,
		"",
		i18n.Sprintf("maximum content loaded in memory\n"+
			"(only used for input from stdin)"),
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
