package main

import (
	"fmt"
	"path/filepath"

	"github.com/aos-dev/go-storage/v2/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type cpFlags struct {
	checkMD5             bool
	expectSize           string
	maximumMemoryContent string
	recursive            bool
	multipartFlags
}

var cpFlag = cpFlags{}

// CpCommand will handle copy command.
var CpCommand = &cobra.Command{
	Use:   "cp <source-path> <dest-path>",
	Short: i18n.Sprintf("copy from/to qingstor"),
	Long:  i18n.Sprintf("qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Copy file: qsctl cp /path/to/file qs://prefix/a"),
		i18n.Sprintf("Copy folder: qsctl cp /path/to/folder qs://prefix/a/ -r"),
		i18n.Sprintf("Copy all files in folder: qsctl cp /path/to/folder/ qs://prefix/a/ -r"),
		i18n.Sprintf("Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin"),
		i18n.Sprintf("Write to stdout: qsctl cp qs://prefix/b - > /path/to/file"),
	),
	Args: cobra.ExactArgs(2),
	PreRunE: func(c *cobra.Command, args []string) error {
		if err := parseCpFlag(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cpRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "cp", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		cpFlag = cpFlags{}
	},
}

func initCpFlag() {
	CpCommand.PersistentFlags().StringVar(&cpFlag.expectSize,
		constants.ExpectSizeFlag,
		"",
		i18n.Sprintf(`expected size of the input file
accept: 100MB, 1.8G
(only used and required for input from stdin)`),
	)
	CpCommand.PersistentFlags().StringVar(&cpFlag.maximumMemoryContent,
		constants.MaximumMemoryContentFlag,
		"",
		i18n.Sprintf(`maximum content loaded in memory
(only used for input from stdin)`),
	)
	CpCommand.Flags().BoolVarP(&cpFlag.recursive,
		constants.RecursiveFlag,
		"r",
		false,
		i18n.Sprintf("copy directory recursively"),
	)
	CpCommand.Flags().StringVar(&cpFlag.partThresholdStr,
		constants.PartThresholdFlag,
		"",
		i18n.Sprintf("set threshold to enable multipart upload"),
	)
	CpCommand.Flags().StringVar(&cpFlag.partSizeStr,
		constants.PartSizeFlag,
		"",
		i18n.Sprintf("set part size for multipart upload"),
	)
}

func cpRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !cpFlag.recursive {
		return fmt.Errorf(i18n.Sprintf("-r is required to copy a directory"))
	}

	if cpFlag.recursive && rootTask.GetSourceType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("src should be a directory while -r is set"))
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir &&
		rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("cannot copy a directory to a non-directory dest"))
	}

	if cpFlag.recursive {
		t := task.NewCopyDir(rootTask)
		t.SetCheckMD5(cpFlag.checkMD5)
		t.SetPartThreshold(cpFlag.partThreshold)
		if cpFlag.partSize != 0 {
			t.SetPartSize(cpFlag.partSize)
		}
		t.SetHandleObjCallback(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "<%s> copied\n", o.Name)
		})
		t.SetCheckTasks(nil)
		t.Run(c.Context())

		if t.GetFault().HasError() {
			return t.GetFault()
		}

		if h := taskutils.HandlerFromContext(c.Context()); h != nil {
			h.WaitProgress()
		}

		i18n.Fprintf(c.OutOrStdout(), "Dir <%s> copied to <%s>.\n",
			filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
		return nil
	}

	t := task.NewCopyFile(rootTask)
	t.SetCheckMD5(cpFlag.checkMD5)
	t.SetPartThreshold(cpFlag.partThreshold)
	if cpFlag.partSize != 0 {
		t.SetPartSize(cpFlag.partSize)
	}
	t.SetCheckTasks(nil)
	t.Run(c.Context())

	if t.GetFault().HasError() {
		return t.GetFault()
	}

	if h := taskutils.HandlerFromContext(c.Context()); h != nil {
		h.WaitProgress()
	}

	i18n.Fprintf(c.OutOrStdout(), "File <%s> copied to <%s>.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return
}

func parseCpFlag() (err error) {
	// parse multipart chunk size
	if cpFlag.partSizeStr != "" {
		// do not set chunk size default value, we need to check it when task init
		cpFlag.partSize, err = utils.ParseByteSize(cpFlag.partSizeStr)
		if err != nil {
			return err
		}
	}

	// parse multipart partThreshold
	if cpFlag.partThresholdStr == "" {
		cpFlag.partThreshold = constants.MaximumAutoMultipartSize
	} else {
		cpFlag.partThreshold, err = utils.ParseByteSize(cpFlag.partThresholdStr)
		if err != nil {
			return err
		}
	}
	return nil
}
