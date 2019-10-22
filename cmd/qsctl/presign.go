package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var presignInput struct {
	expire int
}

// PresignCommand will handle list command.
var PresignCommand = &cobra.Command{
	Use:   "presign qs://<bucket_name>/<object_key>",
	Short: "get the pre-signed URL by the object key",
	Long: `qsctl presign can generate a pre-signed URL for the object. 
Within the given expire time, anyone who receives this URL can retrieve
the object with an HTTP GET request. If an object belongs to a public bucket, 
generate a URL spliced by bucket name, zone and its name, anyone who receives 
this URL can always retrieve the object with an HTTP GET request.`,
	Example: utils.AlignPrintWithColon(
		"Presign object: qsctl qs://bucket-name/object-name",
	),
	Args:   cobra.ExactArgs(1),
	RunE:   presignRun,
	PreRun: validatePresignFlag,
}

func presignParse(t *task.PresignTask, _ []string) (err error) {
	// Parse flags.
	t.SetExpire(presignInput.expire)
	return nil
}

func presignRun(_ *cobra.Command, args []string) error {
	t := task.NewPresignTask(func(t *task.PresignTask) {
		if err := presignParse(t, args); err != nil {
			t.TriggerFault(err)
			return
		}

		err := utils.ParseAtStorageInput(t, args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
	})

	t.Run()
	t.Wait()

	if t.ValidateFault() {
		return t.GetFault()
	}

	presignOutput(t)
	return nil
}

func presignOutput(t *task.PresignTask) {
	fmt.Println(t.GetURL())
}

func initPresignFlag() {
	PresignCommand.Flags().IntVarP(&presignInput.expire, constants.ExpireFlag, "e", 0,
		"the number of seconds until the pre-signed URL expires. Default is 300 seconds")
}

func validatePresignFlag(_ *cobra.Command, _ []string) {
	// set expire default to DefaultPresignExpire
	if presignInput.expire <= 0 {
		presignInput.expire = constants.DefaultPresignExpire
	}
}
