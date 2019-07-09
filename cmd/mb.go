package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
	"github.com/yunify/qsctl/utils"
)

// MbCommand will handle make bucket command.
var MbCommand = &cobra.Command{
	Use:   "mb <bucket-name>",
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
	Args: cobra.ExactArgs(1),
	RunE: mbRun,
}

func mbRun(_ *cobra.Command, args []string) (err error) {
	return action.MakeBucket(args[0])
}
