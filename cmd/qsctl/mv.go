package main

import (
	"fmt"

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
		i18n.Sprintf("Move folder: qsctl mv qs://prefix/a /path/to/folder -r"),
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

func mvRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewBetweenStorageTask(10)
	err = utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !mvInput.Recursive {
		return fmt.Errorf("-r is required to move a directory")
	}

	// mv cmd set storage wd and path reuse cp cmd
	cpInput.Recursive = mvInput.Recursive
	if err = HandleBetweenStorageWdAndPath(rootTask, mvInput.Recursive); err != nil {
		return err
	}

	if mvInput.Recursive {
		t := task.NewMoveDir(rootTask)
		t.Run()

		if t.GetFault().HasError() {
			return t.GetFault()
		}
		i18n.Printf("Dir <%s> moved to <%s>.\n", t.GetSourcePath(), t.GetDestinationPath())
		return nil
	}

	t := task.NewMoveFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	i18n.Printf("File <%s> moved to <%s>.\n", t.GetSourcePath(), t.GetDestinationPath())
	return
}
