package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var rbInput struct {
	force bool
}

// RbCommand will handle remove object command.
var RbCommand = &cobra.Command{
	Use:   "rb [qs://]<bucket_name> [--force/-f]",
	Short: "delete a bucket",
	Long:  "qsctl rb delete a qingstor bucket",
	Example: utils.AlignPrintWithColon(
		"delete an empty bucket: qsctl rb qs://bucket-name",
		"forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f",
	),
	Args: cobra.ExactArgs(1),
	RunE: rbRun,
}

func initRbFlag() {
	RbCommand.Flags().BoolVarP(&rbInput.force, constants.ForceFlag, "f", false,
		"Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket.",
	)
}

func rbParse(t *task.RemoveBucketTask, _ []string) (err error) {
	// Parse flags.
	t.SetForce(rbInput.force)
	return nil
}

func rbRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewRemoveBucketTask(func(t *task.RemoveBucketTask) {
		if err = rbParse(t, args); err != nil {
			t.TriggerFault(err)
			return
		}

		if err := utils.ParseAtServiceInput(t); err != nil {
			t.TriggerFault(err)
			return
		}

		_, bucketName, _, err := utils.ParseQsPath(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetBucketName(bucketName)

		// only rb -f use storage
		if t.GetForce() {
			if err := utils.ParseAtStorageInput(t, args[0]); err != nil {
				t.TriggerFault(err)
				return
			}
		}
	})

	t.Run()
	t.Wait()

	if t.ValidateFault() {
		return t.GetFault()
	}

	rbOutput(t)
	return nil
}

func rbOutput(t *task.RemoveBucketTask) {
	fmt.Printf("Bucket <%s> removed.\n", t.GetBucketName())
}
