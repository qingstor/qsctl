package main

import (
	"fmt"

	"github.com/spf13/cobra"

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

func rmParse(t *task.RemoveObjectTask, _ []string) (err error) {
	// Parse flags.
	t.SetRecursive(rmInput.recursive)
	return nil
}

func rmRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewRemoveObjectTask(func(t *task.RemoveObjectTask) {
		if err = rmParse(t, args); err != nil {
			t.TriggerFault(err)
			return
		}

		keyType, bucketName, objectKey, err := utils.ParseKey(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}

		if (keyType == constants.KeyTypePseudoDir || keyType == constants.KeyTypeBucket) && !t.GetRecursive() {
			t.TriggerFault(fmt.Errorf("-r is required for removing dir operation"))
			return
		}

		t.SetKey(objectKey)
		srv, err := NewQingStorService()
		if err != nil {
			t.TriggerFault(err)
			return
		}

		stor, err := srv.Get(bucketName)
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetDestinationStorage(stor)
	})

	t.Run()
	t.Wait()

	if t.ValidateFault() {
		return t.GetFault()
	}

	rmOutput(t)
	return nil
}

func rmOutput(t *task.RemoveObjectTask) {
	fmt.Printf("Object <%s> removed.\n", t.GetKey())
}
