package main

import (
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type presignFlags struct {
	expire int
}

var presignFlag = presignFlags{}

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
	PreRun: validatePresignFlag,
	Run: func(cmd *cobra.Command, args []string) {
		if err := presignRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "presign", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		presignFlag = presignFlags{}
	},
}

func presignRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtStorageTask()
	_, err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	t := task.NewReachFile(rootTask)
	t.SetReacher(rootTask.GetStorage().(typ.Reacher))
	t.SetExpire(presignFlag.expire)

	if err := t.Run(c.Context()); err != nil {
		return err
	}

	i18n.Fprintf(c.OutOrStdout(), "%s\n", t.GetURL())
	return nil
}

func initPresignFlag() {
	PresignCommand.Flags().IntVarP(&presignFlag.expire, constants.ExpireFlag, "e", 0,
		i18n.Sprintf("the number of seconds until the pre-signed URL expires. Default is 300 seconds"))
}

func validatePresignFlag(_ *cobra.Command, _ []string) {
	// set expire default to DefaultPresignExpire
	if presignFlag.expire <= 0 {
		presignFlag.expire = constants.DefaultPresignExpire
	}
}
