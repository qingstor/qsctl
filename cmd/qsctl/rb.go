package main

import (
	"fmt"

	"github.com/Xuanwo/storage/pkg/segment"
	typ "github.com/Xuanwo/storage/types"
	"github.com/c-bata/go-prompt"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/shellutils"
	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type rbFlags struct {
	force bool
}

var rbFlag = rbFlags{}

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
	Run: func(cmd *cobra.Command, args []string) {
		if err := rbRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "rb", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		rbFlag = rbFlags{}
	},
}

func initRbFlag() {
	RbCommand.Flags().BoolVarP(&rbFlag.force, constants.ForceFlag, "f", false,
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

	t := task.NewDeleteStorage(rootTask)
	t.SetStorageName(bucketName)
	t.SetForce(rbFlag.force)
	t.SetHandleObjCallback(func(o *typ.Object) {
		i18n.Fprintf(c.OutOrStdout(), "<%s> removed\n", o.Name)
	})
	t.SetHandleSegmentCallback(func(seg segment.Segment) {
		i18n.Fprintf(c.OutOrStdout(), "segment id <%s>, path <%s> removed\n", seg.ID(), seg.Path())
	})

	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	i18n.Fprintf(c.OutOrStdout(), "Bucket <%s> removed.\n", t.GetStorageName())
	return nil
}

type rbShellHandler struct {
	bucketName string
}

// preRunE do pre-run check before rb in shell
func (r *rbShellHandler) preRunE(args []string) error {
	err := RbCommand.Flags().Parse(args)
	if err != nil {
		return err
	}
	_, bucketName, _, err := utils.ParseQsPath(RbCommand.Flags().Args()[0])
	if err != nil {
		return err
	}
	if rbFlag.force {
		input := prompt.Input(
			i18n.Sprintf("input bucket name <%s> to confirm: ", bucketName),
			noSuggests)
		if input != bucketName {
			return fmt.Errorf(i18n.Sprintf("not confirmed"))
		}
	}
	r.bucketName = bucketName
	return nil
}

// postRun remove bucket from cache list if no error while run
func (r rbShellHandler) postRun(err error) {
	if err == nil {
		shellutils.RemoveBucketFromList(r.bucketName)
	}
}
