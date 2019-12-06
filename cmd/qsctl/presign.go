package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/task"
	"github.com/qingstor/qsctl/v2/utils"
)

var presignInput struct {
	expire int
}

// PresignCommand will handle list command.
var PresignCommand = &cobra.Command{
	Use:   "presign qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("get the pre-signed URL by the object key"),
	Long: i18n.Sprintf(`qsctl presign can generate a pre-signed URL for the object.
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket,
generate a URL spliced by bucket name, zone and its name, anyone who receives
this URL can always retrieve the object with an HTTP GET request.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Presign object: qsctl qs://bucket-name/object-name"),
	),
	Args:   cobra.ExactArgs(1),
	RunE:   presignRun,
	PreRun: validatePresignFlag,
}

func presignRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtStorageTask(10)
	err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	t := task.NewReachFile(rootTask)
	t.SetExpire(presignInput.expire)

	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	presignOutput(t)
	return nil
}

func presignOutput(t *task.ReachFileTask) {
	fmt.Println(t.GetURL())
}

func initPresignFlag() {
	PresignCommand.Flags().IntVarP(&presignInput.expire, constants.ExpireFlag, "e", 0,
		i18n.Sprintf("the number of seconds until the pre-signed URL expires. Default is 300 seconds"))
}

func validatePresignFlag(_ *cobra.Command, _ []string) {
	// set expire default to DefaultPresignExpire
	if presignInput.expire <= 0 {
		presignInput.expire = constants.DefaultPresignExpire
	}
}
