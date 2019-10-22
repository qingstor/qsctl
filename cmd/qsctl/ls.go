package main

import (
	"fmt"
	"strconv"

	"github.com/Xuanwo/storage/types"
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var lsInput struct {
	HumanReadable bool
	LongFormat    bool
	Recursive     bool
	Zone          string
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
	t.SetHumanReadable(lsInput.HumanReadable)
	t.SetLongFormat(lsInput.LongFormat)
	t.SetRecursive(lsInput.Recursive)
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
		_, bucketName, key, err := utils.ParseQsPath(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetDestinationPath(key)

		store, err := srv.Get(bucketName, types.WithLocation(t.GetZone()))
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetDestinationStorage(store)
		t.SetBucketName(bucketName)

		// init object channel, then stream output by goroutine
		oc := make(chan *types.Object)
		t.SetObjectChannel(oc)
	})

	t.Run()

	// list bucket output here
	if t.GetListType() == constants.ListTypeBucket {
		t.Wait()
		if t.ValidateFault() {
			return t.GetFault()
		}
		listBucketOutput(t)
		return
	}

	// list objects sync with channel, so do not need wait here
	listObjectOutput(t)
	// but we have to get fault after output, otherwise fault will not be triggered
	if t.ValidateFault() {
		return t.GetFault()
	}
	return
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&lsInput.HumanReadable, constants.HumanReadableFlag, "h", false,
		"print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,"+
			" in order to reduce the number of digits to three or less using base 2 for sizes")
	LsCommand.Flags().BoolVarP(&lsInput.LongFormat, constants.LongFormatFlag, "l", false,
		"list in long format and a total sum for all the file sizes is"+
			" output on a line before the long listing")
	LsCommand.Flags().BoolVarP(&lsInput.Recursive, constants.RecursiveFlag, "R", false,
		"recursively list subdirectories encountered")
	// LsCommand.Flags().BoolVarP(&reverse, constants.ReverseFlag, "r", false,
	// 	"reverse the order of the sort to get reverse lexicographical order")
	LsCommand.Flags().StringVarP(&lsInput.Zone, constants.ZoneFlag, "z", "",
		"in which zone to do the operation")
}

// listBucketOutput list buckets with normal slice
func listBucketOutput(t *task.ListTask) {
	for _, v := range t.GetBucketList() {
		fmt.Println(v)
	}
}

// listObjectOutput get object from channel asynchronously, and pack them into output format
func listObjectOutput(t *task.ListTask) {
	if !t.GetLongFormat() {
		for v := range t.GetObjectChannel() {
			fmt.Println(v.Name)
		}
		return
	}

	var err error
	for v := range t.GetObjectChannel() {
		objACL := constants.ACLObject
		if v.Type == types.ObjectTypeDir {
			objACL = constants.ACLDirectory
		}

		size, ok := v.Metadata.GetSize()
		if !ok {
			// if size not exists (like dir), set size to 0
			size = 0
		}

		// default print size by bytes
		readableSize := strconv.FormatInt(size, 10)
		if t.GetHumanReadable() {
			// if human readable flag true, print size as human readable format
			readableSize, err = utils.UnixReadableSize(datasize.ByteSize(size).HR())
			if err != nil {
				t.TriggerFault(err)
				log.Debugf("parse size <%v> failed [%v], key: <%s>", size, err, v.Name)
			}
		}

		// if modified not exists (like dir), init str with blank
		modifiedStr := ""
		if modified, ok := v.Metadata.GetUpdatedAt(); ok {
			modifiedStr = modified.Format(constants.LsDefaultFormat)
		}
		// output order: acl  size  lastModified  key
		// join with two space
		fmt.Printf("%s  %s  %s  %s\n", objACL, readableSize, modifiedStr, v.Name)
	}
}
