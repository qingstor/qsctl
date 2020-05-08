package main

import (
	"fmt"

	"github.com/Xuanwo/storage/pkg/segment"
	typ "github.com/Xuanwo/storage/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var rbInput struct {
	force bool
}

// RbCommand will handle remove object command.
var RbCommand = &cobra.Command{
	Use:   "rb [qs://]<bucket_name> [--force/-f]",
	Short: i18n.Sprintf("delete a bucket"),
	Long:  i18n.Sprintf("qsctl rb delete a qingstor bucket"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("delete an empty bucket: qsctl rb qs://bucket-name"),
		i18n.Sprintf("forcely delete a nonempty bucket: qsctl rb qs://bucket-name -f"),
	),
	Args: cobra.ExactArgs(1),
	RunE: rbRun,
}

func initRbFlag() {
	RbCommand.Flags().BoolVarP(&rbInput.force, constants.ForceFlag, "f", false,
		i18n.Sprintf("Delete an empty qingstor bucket or forcely delete nonempty qingstor bucket."),
	)
}

func rbRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtServiceTask(10)
	err = utils.ParseAtServiceInput(rootTask)
	if err != nil {
		return
	}

	_, bucketName, _, err := utils.ParseQsPath(args[0])
	if err != nil {
		return
	}

	if rbInput.force {
		var match bool
		match, err = utils.DoubleCheckString(bucketName,
			i18n.Sprintf(`This operation will delete all data (including segments) in your bucket <%s>, which cannot be recovered.
Please input the bucket name to confirm:`, bucketName))
		if err != nil {
			return
		}
		if !match {
			return fmt.Errorf(i18n.Sprintf("The bucket name you just input is not match. Bucket <%s> not removed.", bucketName))
		}
	}

	t := task.NewDeleteStorage(rootTask)
	t.SetStorageName(bucketName)
	t.SetForce(rbInput.force)
	t.SetHandleObjCallback(func(o *typ.Object) {
		fmt.Println(i18n.Sprintf("<%s> removed", o.Name))
	})
	t.SetHandleSegmentCallback(func(seg segment.Segment) {
		fmt.Println(i18n.Sprintf("segment id <%s>, path <%s> removed", seg.ID(), seg.Path()))
	})

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
