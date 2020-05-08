package main

import (
	"fmt"
	"path/filepath"

	typ "github.com/Xuanwo/storage/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var rmInput struct {
	recursive bool
}

// RmCommand will handle remove object command.
var RmCommand = &cobra.Command{
	Use:   "rm qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("remove a remote object"),
	Long:  i18n.Sprintf("qsctl rm remove the object with given object key"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Remove a single object: qsctl rm qs://bucket-name/object-key"),
		i18n.Sprintf("Remove objects with prefix: qsctl rm qs://bucket-name/prefix -r"),
	),
	Args: cobra.ExactArgs(1),
	RunE: rmRun,
}

func initRmFlag() {
	RmCommand.Flags().BoolVarP(&rmInput.recursive, constants.RecursiveFlag, "r",
		false, i18n.Sprintf("recursively delete keys under a specific prefix"))
}

func rmRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtStorageTask(10)
	workDir, err := utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	if rootTask.GetType() == typ.ObjectTypeDir && !rmInput.recursive {
		return fmt.Errorf(i18n.Sprintf("-r is required to remove a directory"))
	}

	if rmInput.recursive && rootTask.GetType() != typ.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("path should be a directory while -r is set"))
	}

	key := filepath.Join(workDir, rootTask.GetPath())
	confirm, err := utils.CheckConfirm(i18n.Sprintf(`This operation will delete <%s>, which cannot be recovered.
Confirm?:`, key))
	if err != nil {
		return
	}
	if !confirm {
		return fmt.Errorf(i18n.Sprintf("Not confirmed. Object <%s> not removed.", key))
	}

	if rmInput.recursive {
		t := task.NewDeleteDir(rootTask)
		t.SetHandleObjCallback(func(o *typ.Object) {
			fmt.Println(i18n.Sprintf("<%s> removed", o.Name))
		})
		t.Run()
		if t.GetFault().HasError() {
			return t.GetFault()
		}

		i18n.Printf("Dir <%s> removed.\n", key)
		return nil
	}

	t := task.NewDeleteFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	i18n.Printf("File <%s> removed.\n", key)
	return nil
}
