package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/yunify/qsctl/v2/pkg/i18n"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var rmInput struct {
	recursive bool
}

// RmCommand will handle remove object command.
var RmCommand = &cobra.Command{
	Use:   "rm qs://<bucket_name>/<object_key>",
	Short: i18n.Sprint("remove a remote object"),
	Long:  i18n.Sprint("qsctl rm remove the object with given object key"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprint("Remove a single object: qsctl rm qs://bucket-name/object-key"),
	),
	Args: cobra.ExactArgs(1),
	RunE: rmRun,
}

func initRmFlag() {
	RmCommand.Flags().BoolVarP(&rmInput.recursive, constants.RecursiveFlag, "r",
		false, i18n.Sprint("recursively delete keys under a specific prefix"))
}

func rmRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtStorageTask(10)
	err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	if rmInput.recursive {
		t := task.NewDeleteDir(rootTask)
		t.Run()
		if t.GetFault().HasError() {
			return t.GetFault()
		}

		rmDirOutput(t)
		return nil
	}

	t := task.NewDeleteFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	rmOutput(t)
	return nil
}

func rmOutput(t *task.DeleteFileTask) {
	i18n.Printf("Object <%s> removed.\n", t.GetPath())
}

func rmDirOutput(t *task.DeleteDirTask) {
	i18n.Printf("Dir <%s> removed.\n", t.GetPath())
}
