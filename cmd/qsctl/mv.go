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

type mvFlags struct {
	checkMD5  bool
	recursive bool
	multipartFlags
}

var mvFlag = mvFlags{}

// MvCommand will handle move command.
var MvCommand = &cobra.Command{
	Use:   "mv <source-path> <dest-path>",
	Short: i18n.Sprintf("move from/to qingstor"),
	Long:  i18n.Sprintf("qsctl mv can move file/folder to qingstor or move qingstor objects to local"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Move file: qsctl mv /path/to/file qs://prefix/a"),
		i18n.Sprintf("Move folder: qsctl mv /path/to/folder qs://prefix/a/ -r"),
		i18n.Sprintf("Move all files in folder: qsctl mv /path/to/folder/ qs://prefix/a/ -r"),
	),
	Args: cobra.ExactArgs(2),
	PreRunE: func(c *cobra.Command, args []string) error {
		if err := parseMvFlag(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := mvRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "mv", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		mvFlag = mvFlags{}
	},
}

func initMvFlag() {
	MvCommand.Flags().BoolVarP(&mvFlag.recursive,
		constants.RecursiveFlag,
		"r",
		false,
		i18n.Sprintf("move directory recursively"))
	MvCommand.Flags().StringVar(&mvFlag.multipartThresholdStr,
		constants.MultipartThresholdFlag,
		"",
		i18n.Sprintf("set threshold to enable multipart upload"),
	)
	MvCommand.Flags().StringVar(&mvFlag.partSizeStr,
		constants.PartsizeFlag,
		"",
		i18n.Sprintf("set chunk size of multipart upload"),
	)
}

func mvRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !mvFlag.recursive {
		return fmt.Errorf(i18n.Sprintf("-r is required to move a directory"))
	}

	if mvFlag.recursive && rootTask.GetSourceType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("src should be a directory while -r is set"))
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir &&
		rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("cannot move a directory to a non-directory dest"))
	}

	if mvFlag.recursive {
		t := task.NewMoveDir(rootTask)
		t.SetHandleObjCallback(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "<%s> moved\n", o.Name)
		})
		t.SetCheckMD5(mvFlag.checkMD5)
		t.SetPartThreshold(mvFlag.multipartThreshold)
		if mvFlag.partSize != 0 {
			t.SetPartSize(mvFlag.partSize)
		}
		t.Run(c.Context())

		if t.GetFault().HasError() {
			return t.GetFault()
		}

		if h := taskutils.HandlerFromContext(c.Context()); h != nil {
			h.WaitProgress()
		}

		i18n.Fprintf(c.OutOrStdout(), "Dir <%s> moved to <%s>.\n",
			filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
		return nil
	}

	t := task.NewMoveFile(rootTask)
	t.SetCheckMD5(mvFlag.checkMD5)
	t.SetPartThreshold(mvFlag.multipartThreshold)
	if mvFlag.partSize != 0 {
		t.SetPartSize(mvFlag.partSize)
	}
	t.Run(c.Context())
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	if h := taskutils.HandlerFromContext(c.Context()); h != nil {
		h.WaitProgress()
	}

	i18n.Fprintf(c.OutOrStdout(), "File <%s> moved to <%s>.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return
}

func parseMvFlag() (err error) {
	// parse multipart chunk size
	if mvFlag.partSizeStr != "" {
		// do not set chunk size default value, we need to check it when task init
		mvFlag.partSize, err = utils.ParseByteSize(mvFlag.partSizeStr)
		if err != nil {
			return err
		}
	}

	// parse multipart threshold
	if mvFlag.multipartThresholdStr == "" {
		mvFlag.multipartThreshold = constants.MaximumAutoMultipartSize
	} else {
		mvFlag.multipartThreshold, err = utils.ParseByteSize(mvFlag.multipartThresholdStr)
		if err != nil {
			return err
		}
	}
	return nil
}
