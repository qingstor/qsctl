package main

import (
	"fmt"

	"github.com/Xuanwo/storage/types"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"

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
	Short: "remove a remote object",
	Long:  "qsctl rm remove the object with given object key",
	Example: utils.AlignPrintWithColon(
		"Remove a single object: qsctl rm qs://bucket-name/object-key",
	),
	Args: cobra.ExactArgs(1),
	RunE: rmRun,
}

func initRmFlag() {
	RmCommand.Flags().BoolVarP(&rmInput.recursive, constants.RecursiveFlag, "r",
		false, "recursively delete keys under a specific prefix")
}

func rmRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtStorageTask(10)
	err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	// TODO: handle remove dir.
	if rootTask.GetType() == types.ObjectTypeDir && !rmInput.recursive {
		return fmt.Errorf("-r is required while remove a directory")
	}

	if rmInput.recursive {
		t := task.NewDeleteDir(rootTask)
		t.Run()

		if t.GetFault().HasError() {
			return t.GetFault()
		}
		return rmDirOutput(t)
	}

	t := task.NewDeleteFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	rmFileOutput(t)
	return nil
}

func rmFileOutput(t *task.DeleteFileTask) {
	fmt.Printf("Object <%s> removed.\n", t.GetPath())
}

func rmDirOutput(t *task.DeleteDirTask) error {
	if t.GetPath() == "" || t.GetPath() == "/" {
		md, err := t.GetStorage().Metadata()
		if err != nil {
			return err
		}
		bucketName, ok := md.GetName()
		if ok {
			fmt.Printf("Objects in bucket <%s> removed.\n", bucketName)
			return nil
		}
		fmt.Printf("Objects removed.\n")
		return nil
	}
	fmt.Printf("Objects in <%s> removed.\n", t.GetPath())
	return nil
}
