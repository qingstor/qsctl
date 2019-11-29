package main

import (
	"fmt"

	"github.com/Xuanwo/storage/types"
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/pkg/i18n"

	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var mvInput struct {
	Recursive bool
}

// MvCommand will handle move command.
var MvCommand = &cobra.Command{
	Use:   "mv <source-path> <dest-path>",
	Short: i18n.Sprint("move from/to qingstor"),
	Long:  i18n.Sprint("qsctl mv can move file/folder to qingstor or move qingstor objects to local"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprint("Move file: qsctl mv /path/to/file qs://prefix/a"),
		i18n.Sprint("Move folder: qsctl mv qs://prefix/a /path/to/folder -r"),
	),
	Args: cobra.ExactArgs(2),
	RunE: mvRun,
}

func initMvFlag() {
	MvCommand.Flags().BoolVarP(&mvInput.Recursive,
		constants.RecursiveFlag,
		"r",
		false,
		i18n.Sprint("move directory recursively"))
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
		mvOutput(args[0])
		return nil
	}

	t := task.NewMoveFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	mvOutput(args[0])
	return
}

func mvOutput(path string) {
	i18n.Printf("Key <%s> moved.\n", path)
}
