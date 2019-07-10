package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
	"github.com/yunify/qsctl/utils"
)

// RmCommand will handle remove object command.
var RmCommand = &cobra.Command{
	Use:   "rm qs://<bucket_name>/<object_key>",
	Short: "remove a remote object",
	Long:  "qsctl rm remove the object with given object key",
	Example: utils.AlignPrintWithColon(
		"Remove a single object: qsctl rm qs://bucket-name/object-key",
	),
	Args: cobra.ExactArgs(1),
	RunE: rmRun,
}

func rmRun(_ *cobra.Command, args []string) (err error) {
	return action.DeleteObject(args[0])
}
