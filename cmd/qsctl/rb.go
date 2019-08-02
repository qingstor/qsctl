package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/utils"
)

// RbCommand will handle remove object command.
var RbCommand = &cobra.Command{
	Use:   "rb [qs://]<bucket_name>",
	Short: "delete a bucket",
	Long:  "qsctl rb delete an empty qingstor bucket",
	Example: utils.AlignPrintWithColon(
		"delete an empty bucket: qsctl rb qs://bucket-name",
	),
	Args: cobra.ExactArgs(1),
	RunE: rbRun,
}

func rbRun(_ *cobra.Command, args []string) (err error) {
	return action.RemoveBucket(args[0])
}
