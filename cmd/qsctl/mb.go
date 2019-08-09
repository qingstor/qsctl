package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/utils"
)

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

func mbRun(_ *cobra.Command, args []string) (err error) {
	bh := &action.BucketHandler{}
	return bh.WithZone(zone).WithRemote(args[0]).MakeBucket()
}

func initMbFlag() {
	MbCommand.Flags().StringVarP(&zone, constants.ZoneFlag, "z",
		"", "in which zone to make the bucket (required)")
}

func validateMbFlag(_ *cobra.Command, _ []string) error {
	// check zone flag (required)
	if zone == "" {
		return constants.ErrorZoneRequired
	}
	return nil
}
