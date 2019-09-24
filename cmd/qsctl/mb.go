package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task"
	taskUtils "github.com/yunify/qsctl/v2/task/utils"
	"github.com/yunify/qsctl/v2/utils"
)

var mbInput struct {
	Zone string
}

// MbCommand will handle make bucket command.
var MbCommand = &cobra.Command{
	Use:   "mb [qs://]<bucket-name>",
	Short: "make a new bucket",
	Long: `qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`,
	Example: utils.AlignPrintWithColon(
		"Make bucket: qsctl mb bucket-name",
	),
	Args:    cobra.ExactArgs(1),
	RunE:    mbRun,
	PreRunE: validateMbFlag,
}

func mbParse(t *task.MakeBucketTask, args []string) (err error) {
	// Parse flags.
	t.SetZone(mbInput.Zone)
	return nil
}

func mbRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewMakeBucketTask(func(t *task.MakeBucketTask) {
		if err = mbParse(t, args); err != nil {
			return
		}
		keyType, bucketName, _, e := taskUtils.ParseKey(args[0])
		if e != nil {
			err = e
			return
		}
		if keyType != constants.KeyTypeBucket {
			err = constants.ErrorQsPathInvalid
			return
		}
		t.SetBucketName(bucketName)

		stor, e := storage.NewQingStorObjectStorage()
		if e != nil {
			err = e
			return
		}
		t.SetStorage(stor)

		if err = stor.SetupBucket(t.GetBucketName(), t.GetZone()); err != nil {
			return
		}
	})

	t.Run()
	t.Wait()
	return
}

func initMbFlag() {
	MbCommand.Flags().StringVarP(&mbInput.Zone, constants.ZoneFlag, "z",
		"", "in which zone to make the bucket (required)")
}

func validateMbFlag(_ *cobra.Command, _ []string) error {
	// check zone flag (required)
	if mbInput.Zone == "" {
		return constants.ErrorZoneRequired
	}
	return nil
}
