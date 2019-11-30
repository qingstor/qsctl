package main

import (
	"fmt"
	"strconv"

	"github.com/Xuanwo/storage"
	typ "github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/pkg/i18n"

	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"
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
	Use:   i18n.Sprintf("ls [qs://<bucket-name/prefix>]"),
	Short: i18n.Sprintf("list objects or buckets"),
	Long:  i18n.Sprintf(`qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("List buckets: qsctl ls"),
		i18n.Sprintf("List bucket's all objects: qsctl ls qs://bucket-name -R"),
		i18n.Sprintf("List objects with prefix: qsctl ls qs://bucket-name/prefix"),
		i18n.Sprintf("List objects by long format: qsctl ls qs://bucket-name -l"),
	),
	Args: cobra.MaximumNArgs(1),
	RunE: lsRun,
}

func lsRun(_ *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		rootTask := taskutils.NewAtServiceTask(10)
		err = utils.ParseAtServiceInput(rootTask)
		if err != nil {
			return
		}

		t := task.NewListStorage(rootTask)
		t.SetZone(lsInput.Zone)
		t.SetStoragerFunc(listBucketOutput)

		t.Run()
		if t.GetFault().HasError() {
			return t.GetFault()
		}
		return
	}

	rootTask := taskutils.NewAtStorageTask(10)
	err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	if err = HandleLsStorageWdAndPath(rootTask); err != nil {
		return err
	}
	t := task.NewListDir(rootTask)
	t.SetFileFunc(listFileOutput)
	t.SetDirFunc(func(object *typ.Object) {})

	t.Run()

	// but we have to get fault after output, otherwise fault will not be triggered
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&lsInput.HumanReadable, constants.HumanReadableFlag, "h", false,
		i18n.Sprintf("print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,"+
			" in order to reduce the number of digits to three or less using base 2 for sizes"))
	LsCommand.Flags().BoolVarP(&lsInput.LongFormat, constants.LongFormatFlag, "l", false,
		i18n.Sprintf("list in long format and a total sum for all the file sizes is"+
			" output on a line before the long listing"))
	LsCommand.Flags().BoolVarP(&lsInput.Recursive, constants.RecursiveFlag, "R", false,
		i18n.Sprintf("recursively list subdirectories encountered"))
	// LsCommand.Flags().BoolVarP(&reverse, constants.ReverseFlag, "r", false,
	// 	"reverse the order of the sort to get reverse lexicographical order")
	LsCommand.Flags().StringVarP(&lsInput.Zone, constants.ZoneFlag, "z", "",
		i18n.Sprintf("in which zone to do the operation"))
}

// listBucketOutput list buckets with normal slice
func listBucketOutput(s storage.Storager) {
	m, err := s.Metadata()
	if err != nil {
		log.Debugf("listBucketOutput: %v", err)
	}
	name, ok := m.GetName()
	if !ok {
		log.Warnf("bucket failed to get name")
	}
	fmt.Println(name)
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

	size, ok := o.Metadata.GetSize()
	if !ok {
		// if size not exists (like dir), set size to 0
		size = 0
	}

	// default print size by bytes
	readableSize := strconv.FormatInt(size, 10)
	if lsInput.HumanReadable {
		// if human readable flag true, print size as human readable format
		readableSize, err = utils.UnixReadableSize(datasize.ByteSize(size).HR())
		if err != nil {
			log.Debugf("parse size <%o> failed [%o], key: <%s>", size, err, o.Name)
		}
	}

	// if modified not exists (like dir), init str with blank
	modifiedStr := o.UpdatedAt.Format(constants.LsDefaultFormat)
	// output order: acl  size  lastModified  key
	// join with two space
	i18n.Printf("%s  %s  %s  %s\n", objACL, readableSize, modifiedStr, o.Name)
}

// HandleLsStorageWdAndPath set work dir and path for ls cmd.
func HandleLsStorageWdAndPath(t *taskutils.AtStorageTask) error {
	if err := t.GetStorage().Init(pairs.WithWorkDir(t.GetPath())); err != nil {
		return err
	}
	t.SetPath("")
	return nil
}
