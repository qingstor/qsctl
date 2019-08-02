package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

// LsCommand will handle list command.
var LsCommand = &cobra.Command{
	Use:   "ls [qs://<bucket-name/prefix>]",
	Short: "list objects or buckets",
	Long:  `qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`,
	Example: utils.AlignPrintWithColon(
		"List buckets: qsctl ls",
		"List bucket's all objects: qsctl ls qs://bucket-name -r",
		"List objects with prefix: qsctl ls qs://bucket-name/prefix",
		"List objects by long format: qsctl ls qs://bucket-name -l",
	),
	Args: cobra.MaximumNArgs(1),
	RunE: lsRun,
}

func lsRun(_ *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return action.ListBuckets(contexts.Zone)
	}
	return action.ListObjects(args[0])
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&contexts.HumanReadable, "human-readable", "h",
		false, "print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,"+
			" in order to reduce the number of digits to three or less using base 2 for sizes")
	LsCommand.Flags().BoolVarP(&contexts.LongFormat, "long-format", "l",
		false, "list in long format and a total sum for all the file sizes is"+
			" output on a line before the long listing")
	LsCommand.Flags().StringVarP(&contexts.Zone, "zone", "z",
		"", "in which zone to do the operation")
	LsCommand.Flags().BoolVarP(&contexts.Recursive, "recursive", "R",
		false, "recursively list subdirectories encountered")
	LsCommand.Flags().BoolVarP(&contexts.Reverse, "reverse", "r",
		false, "reverse the order of the sort to get reverse lexicographical order")
}
