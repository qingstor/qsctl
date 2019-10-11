package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
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

func mbParse(t *task.MakeBucketTask, _ []string) (err error) {
	// Parse flags.
	t.SetZone(mbInput.Zone)
	return nil
}

func mbRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewMakeBucketTask(func(t *task.MakeBucketTask) {
		if err = mbParse(t, args); err != nil {
			t.TriggerFault(err)
			return
		}
		keyType, bucketName, _, err := utils.ParseKey(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
		if keyType != constants.KeyTypeBucket {
			t.TriggerFault(fmt.Errorf("key type is not match"))
			return
		}
		t.SetBucketName(bucketName)

		cfg := NewQingstorConfig(
			WriteBase(),
			WriteBucketName(t.GetBucketName()),
			WriteZone(t.GetZone()),
		)

		stor, err := cfg.New()
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

	mbOutput(t)
	return
}

func mbOutput(t *task.MakeBucketTask) {
	fmt.Printf("Bucket <%s> created.\n", t.GetBucketName())
}

func initMbFlag() {
	MbCommand.Flags().StringVarP(&mbInput.Zone, constants.ZoneFlag, "z",
		"", "in which zone to make the bucket (required)")
}

func validateMbFlag(_ *cobra.Command, _ []string) error {
	// check zone flag (required)
	if mbInput.Zone == "" {
		// TODO: we need to return an error here.
		return fmt.Errorf("flag zone is required, but not found")
	}
	return nil
}
