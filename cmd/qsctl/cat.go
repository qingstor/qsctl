package main

import (
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/task"
	"github.com/qingstor/qsctl/v2/utils"
)

// CatCommand will handle cat command.
var CatCommand = &cobra.Command{
	Use:   "cat qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("cat a remote object to stdout"),
	Long:  i18n.Sprintf("qsctl cat can cat a remote object to stdout"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Cat object: qsctl cat qs://prefix/a"),
	),
	Args: cobra.ExactArgs(1),
	RunE: catRun,
}

func catRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewBetweenStorageTask(10)
	err = utils.ParseBetweenStorageInput(rootTask, args[0], "-")
	if err != nil {
		return
	}
	t := task.NewCopyFile(rootTask)
	t.SetCheckTasks(nil)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return nil
}
