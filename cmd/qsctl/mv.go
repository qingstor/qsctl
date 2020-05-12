package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Xuanwo/storage/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var mvInput struct {
	Recursive bool
}

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
	RunE: mvRun,
}

func initMvFlag() {
	MvCommand.Flags().BoolVarP(&mvInput.Recursive,
		constants.RecursiveFlag,
		"r",
		false,
		i18n.Sprintf("move directory recursively"))
}

func mvRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !mvInput.Recursive {
		return fmt.Errorf("-r is required to move a directory")
	}

	if cpInput.Recursive && rootTask.GetSourceType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("src should be a directory while -r is set"))
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir &&
		rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("cannot move a directory to a non-directory dest"))
	}

	// only show progress bar without no-progress flag set
	if !noProgress {
		go func() {
			taskutils.StartProgress(time.Second)
		}()
		defer taskutils.FinishProgress()
	}

	if mvInput.Recursive {
		t := task.NewMoveDir(rootTask)
		t.SetHandleObjCallback(func(o *types.Object) {
			fmt.Println(i18n.Sprintf("<%s> moved", o.Name))
		})
		t.Run()

		if t.GetFault().HasError() {
			return t.GetFault()
		}

		taskutils.WaitProgress()
		i18n.Printf("Dir <%s> moved to <%s>.\n",
			filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
		return nil
	}

	t := task.NewMoveFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	taskutils.WaitProgress()
	i18n.Printf("File <%s> moved to <%s>.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return
}
