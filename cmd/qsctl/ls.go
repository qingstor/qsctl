package main

import (
	"fmt"

	"github.com/Xuanwo/storage/types"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var lsInput struct {
	Zone string
}

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

func lsParse(t *task.ListTask, _ []string) (err error) {
	// Parse flags.
	t.SetZone(lsInput.Zone)
	return nil
}

func lsRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewListTask(func(t *task.ListTask) {
		err = lsParse(t, args)
		if err != nil {
			t.TriggerFault(err)
			return
		}

		srv, err := NewQingStorService()
		if err != nil {
			t.TriggerFault(err)
			return
		}

		// if no args, handle cmd as list buckets, otherwise list objects.
		if len(args) == 0 {
			t.SetListType(constants.ListTypeBucket)
			t.SetDestinationService(srv)
			return
		}

		t.SetListType(constants.ListTypeKey)
		_, bucketName, _, err := utils.ParseKey(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}

		store, err := srv.Get(bucketName, types.WithLocation(t.GetZone()))
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetDestinationStorage(store)
		t.SetBucketName(bucketName)
	})

	t.Run()
	t.Wait()
	if t.ValidateFault() {
		return t.GetFault()
	}

	lsOutput(t)
	return
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&humanReadable, constants.HumanReadableFlag, "h", false,
		"print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,"+
			" in order to reduce the number of digits to three or less using base 2 for sizes")
	LsCommand.Flags().BoolVarP(&longFormat, constants.LongFormatFlag, "l", false,
		"list in long format and a total sum for all the file sizes is"+
			" output on a line before the long listing")
	LsCommand.Flags().BoolVarP(&recursive, constants.RecursiveFlag, "R", false,
		"recursively list subdirectories encountered")
	LsCommand.Flags().BoolVarP(&reverse, constants.ReverseFlag, "r", false,
		"reverse the order of the sort to get reverse lexicographical order")
	LsCommand.Flags().StringVarP(&lsInput.Zone, constants.ZoneFlag, "z", "",
		"in which zone to do the operation")
}

func lsOutput(t *task.ListTask) {
	if t.GetListType() == constants.ListTypeBucket {
		for _, v := range t.GetBucketList() {
			fmt.Println(v)
		}
		return
	}
}
