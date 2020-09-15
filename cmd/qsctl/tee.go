package main

import (
	"path/filepath"

	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type teeFlags struct {
	expectSize string
	maxMemory  string
	multipartFlags
}

var teeFlag = teeFlags{}

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
	Args: cobra.ExactArgs(1),
	PreRunE: func(c *cobra.Command, args []string) error {
		if err := validateTeeFlag(c, args); err != nil {
			return err
		}
		if err := parseTeeFlag(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := teeRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "tee", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		teeFlag = teeFlags{}
	},
}

func teeRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	_, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, "-", args[0])
	if err != nil {
		return
	}

	t := task.NewCopyStream(rootTask)
	t.SetCheckMD5(false)
	t.SetPartSize(teeFlag.partSize)
	t.Run(c.Context())

	if t.GetFault().HasError() {
		return t.GetFault()
	}
	i18n.Fprintf(c.OutOrStdout(), "Stdin copied to <%s>.\n", filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return nil
}

func initTeeFlag() {
	TeeCommand.PersistentFlags().StringVar(&teeFlag.expectSize,
		constants.ExpectSizeFlag,
		"",
		i18n.Sprintf("expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)"),
	)
	TeeCommand.PersistentFlags().StringVar(&teeFlag.maxMemory,
		constants.MaximumMemoryContentFlag,
		"",
		i18n.Sprintf("maximum content loaded in memory\n"+
			"(only used for input from stdin)"),
	)
	TeeCommand.Flags().StringVar(&teeFlag.partSizeStr,
		constants.PartsizeFlag,
		"",
		i18n.Sprintf("set part size for multipart upload"),
	)
}

func validateTeeFlag(_ *cobra.Command, _ []string) (err error) {
	// TODO: parse should be moved into teeParse func
	if teeFlag.expectSize != "" {
		teeExpectSize, err := utils.ParseByteSize(teeFlag.expectSize)
		_ = teeExpectSize
		if err != nil {
			return err
		}
	}

	if teeFlag.maxMemory != "" {
		teeMaxMemory, err := utils.ParseByteSize(teeFlag.maxMemory)
		_ = teeMaxMemory
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTeeFlag() (err error) {
	// parse multipart chunk size
	if teeFlag.partSizeStr != "" {
		teeFlag.partSize, err = utils.ParseByteSize(teeFlag.partSizeStr)
		if err != nil {
			return err
		}
	} else {
		teeFlag.partSize = constants.DefaultPartSize
	}
	return nil
}
