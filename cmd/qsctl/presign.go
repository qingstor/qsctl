package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/utils"
)

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

func presignRun(_ *cobra.Command, args []string) error {
	return nil
}

func initPresignFlag() {
	PresignCommand.Flags().IntVarP(&expire, constants.ExpireFlag, "e", 0,
		"the number of seconds until the pre-signed URL expires. Default is 300 seconds")
	PresignCommand.Flags().StringVarP(&zone, constants.ZoneFlag, "z",
		"", "specify the zone of the bucket which contains the object")
}

func validatePresignFlag(_ *cobra.Command, _ []string) {
	// set expire default to DefaultPresignExpire
	if expire <= 0 {
		expire = constants.DefaultPresignExpire
	}
}
