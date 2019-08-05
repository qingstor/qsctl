package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
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
	ctx = contexts.SetContext(ctx, "remote", args[0])
	return action.Delete(ctx)
}

func initRmFlag() {
	RmCommand.Flags().StringVarP(&zone, constants.ZoneFlag, "z",
		"", "in which zone to remove object")
}
