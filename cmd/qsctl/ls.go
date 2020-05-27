package main

import (
	"fmt"
	"strconv"

	"github.com/Xuanwo/storage"
	typ "github.com/Xuanwo/storage/types"
	"github.com/c2h5oh/datasize"
	"github.com/jedib0t/go-pretty/text"
	"github.com/qingstor/noah/pkg/types"
	"github.com/qingstor/noah/task"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var lsInput struct {
	HumanReadable bool
	LongFormat    bool
	Recursive     bool
}

// LsCommand will handle list command.
var LsCommand = &cobra.Command{
	Use:   "ls [qs://<bucket-name/prefix>]",
	Short: i18n.Sprintf("list objects or buckets"),
	Long:  i18n.Sprintf(`qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("List buckets: qsctl ls"),
		i18n.Sprintf("List bucket's all objects: qsctl ls qs://bucket-name -R"),
		i18n.Sprintf("List objects with prefix: qsctl ls qs://bucket-name/prefix"),
		i18n.Sprintf("List objects with prefix recursively: qsctl ls qs://bucket-name/prefix -R"),
		i18n.Sprintf("List objects by long format: qsctl ls qs://bucket-name -l"),
	),
	Args: cobra.MaximumNArgs(1),
	RunE: lsRun,
}

func lsRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	if len(args) == 0 {
		rootTask := taskutils.NewAtServiceTask(10)
		err = utils.ParseAtServiceInput(rootTask)
		if err != nil {
			return
		}

		t := task.NewListStorage(rootTask)
		t.SetZone(zone)
		t.SetStoragerFunc(listBucketOutput)

		t.Run()
		if t.GetFault().HasError() {
			return t.GetFault()
		}
		return
	}

	rootTask := taskutils.NewAtStorageTask(10)
	_, err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	t := task.NewListDir(rootTask)
	lister, ok := rootTask.GetStorage().(storage.DirLister)
	if !ok {
		return types.NewErrStorageInsufficientAbility(nil)
	}
	t.SetDirLister(lister)

	t.SetFileFunc(listFileOutput)
	if lsInput.Recursive {
		t.SetDirFunc(func(o *typ.Object) {
			listDirFunc(t, o)
		})
	} else {
		t.SetDirFunc(listFileOutput)
	}
	t.Run()

	// but we have to get fault after output, otherwise fault will not be triggered
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return
}

func listDirFunc(t *task.ListDirTask, o *typ.Object) {
	listFileOutput(o)
	sf := task.NewListDir(t)
	sf.SetPath(o.Name)
	sf.SetFileFunc(listFileOutput)
	sf.SetDirFunc(func(oo *typ.Object) {
		listDirFunc(sf, oo)
	})
	t.GetScheduler().Sync(sf)
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&lsInput.HumanReadable, constants.HumanReadableFlag, "h", false,
		i18n.Sprintf(`print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,
in order to reduce the number of digits to three or less using base 2 for sizes`))
	LsCommand.Flags().BoolVarP(&lsInput.LongFormat, constants.LongFormatFlag, "l", false,
		i18n.Sprintf(`list in long format and a total sum for all the file sizes is
output on a line before the long listing`))
	LsCommand.Flags().BoolVarP(&lsInput.Recursive, constants.RecursiveFlag, "R", false,
		i18n.Sprintf("recursively list subdirectories encountered"))
	// LsCommand.Flags().BoolVarP(&reverse, constants.ReverseFlag, "r", false,
	// 	"reverse the order of the sort to get reverse lexicographical order")
	// LsCommand.Flags().StringVarP(&lsInput.Zone, constants.ZoneFlag, "z", "",
	// 	i18n.Sprintf("in which zone to do the operation"))
}

// listBucketOutput list buckets with normal slice
func listBucketOutput(s storage.Storager) {
	m, err := s.Metadata()
	if err != nil {
		log.Debugf("listBucketOutput: %v", err)
	}
	fmt.Println(m.Name)
}

func listFileOutput(o *typ.Object) {
	if !lsInput.LongFormat {
		fmt.Println(o.Name)
		return
	}

	var err error

	objACL := constants.ACLObject
	if o.Type == typ.ObjectTypeDir {
		objACL = constants.ACLDirectory
	}

	// default print size by bytes
	readableSize := strconv.FormatInt(o.Size, 10)
	if lsInput.HumanReadable {
		// if human readable flag true, print size as human readable format
		readableSize, err = utils.UnixReadableSize(datasize.ByteSize(o.Size).HR())
		if err != nil {
			log.Debugf("parse size <%o> failed [%o], key: <%s>", o.Size, err, o.Name)
		}
		// 7 is the widest size of readable-size, like 1023.9K
		readableSize = text.AlignRight.Apply(readableSize, 7)
	}

	// if modified not exists (like dir), init str with blank
	modifiedStr := o.UpdatedAt.Format(constants.LsDefaultFormat)
	// output order: acl  size  lastModified  key
	// join with two space
	i18n.Printf("%s  %s  %s  %s\n", objACL, readableSize, modifiedStr, o.Name)
}
