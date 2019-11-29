package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/yunify/qsctl/v2/pkg/i18n"

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
	Short: i18n.Sprint("delete a bucket"),
	Long:  i18n.Sprint("qsctl rb delete a qingstor bucket"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprint("delete an empty bucket: qsctl rb qs://bucket-name"),
		i18n.Sprint("forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f"),
	),
	Args: cobra.ExactArgs(1),
	RunE: rbRun,
}

func initRbFlag() {
	RbCommand.Flags().BoolVarP(&rbInput.force, constants.ForceFlag, "f", false,
		i18n.Sprint("Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket."),
	)
}

func rbRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtServiceTask(10)
	err = utils.ParseAtServiceInput(rootTask)
	if err != nil {
		return
	}

	_, bucketName, _, err := utils.ParseQsPath(args[0])
	if err != nil {
		return
	}

	t := task.NewDeleteStorage(rootTask)
	t.SetStorageName(bucketName)
	t.SetForce(rbInput.force)

	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	rbOutput(t)
	return nil
}

func rbOutput(t *task.DeleteStorageTask) {
	i18n.Printf("Bucket <%s> removed.\n", t.GetStorageName())
}
