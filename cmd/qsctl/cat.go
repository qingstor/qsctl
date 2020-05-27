package main

import (
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
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

func catRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	_, _, err = utils.ParseBetweenStorageInput(rootTask, args[0], "-")
	if err != nil {
		return
	}
	t := task.NewCopyFile(rootTask)
	t.SetCheckMD5(false)
	t.SetCheckTasks(nil)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return nil
}
