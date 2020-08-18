package main

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/aos-dev/go-storage/v2"
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/c2h5oh/datasize"
	"github.com/jedib0t/go-pretty/text"
	"github.com/qingstor/log"
	"github.com/qingstor/noah/pkg/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type lsFlags struct {
	humanReadable bool
	longFormat    bool
	recursive     bool
}

var lsFlag = lsFlags{}

// LsCommand will handle list command.
var LsCommand = &cobra.Command{
	Use:   "ls [qs://<bucket-name/prefix>]",
	Short: i18n.Sprintf("list objects or buckets"),
	Long:  i18n.Sprintf(`qsctl ls can list all qingstor buckets or qingstor keys under a prefix.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("List buckets: qsctl ls"),
		i18n.Sprintf("List buckets by long format: qsctl ls -l"),
		i18n.Sprintf("List bucket's all objects: qsctl ls qs://bucket-name -R"),
		i18n.Sprintf("List objects with prefix: qsctl ls qs://bucket-name/prefix"),
		i18n.Sprintf("List objects with prefix recursively: qsctl ls qs://bucket-name/prefix -R"),
		i18n.Sprintf("List objects by long format: qsctl ls qs://bucket-name -l"),
	),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := lsRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "ls", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		lsFlag = lsFlags{}
	},
}

func lsRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	// if no args, handle as list buckets
	if len(args) == 0 {
		rootTask := taskutils.NewAtServiceTask(10)
		err = utils.ParseAtServiceInput(rootTask)
		if err != nil {
			return
		}

		t := task.NewListStorage(rootTask)
		t.SetZone(globalFlag.zone)
		if lsFlag.longFormat {
			t.SetStoragerFunc(func(s storage.Storager) {
				listBucketLongOutput(c.Context(), c.OutOrStdout(), s, t)
			})
		} else {
			t.SetStoragerFunc(func(s storage.Storager) {
				listBucketOutput(c.Context(), c.OutOrStdout(), s)
			})
		}

		t.Run(c.Context())
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

	t.SetFileFunc(func(o *typ.Object) {
		listFileOutput(c.Context(), c.OutOrStdout(), o)
	})
	if lsFlag.recursive {
		t.SetDirFunc(func(o *typ.Object) {
			listDirFunc(c.Context(), c.OutOrStdout(), t, o)
		})
	} else {
		t.SetDirFunc(func(o *typ.Object) {
			listFileOutput(c.Context(), c.OutOrStdout(), o)
		})
	}
	t.Run(c.Context())

	// but we have to get fault after output, otherwise fault will not be triggered
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	return
}

func listDirFunc(ctx context.Context, w io.Writer, t *task.ListDirTask, o *typ.Object) {
	listFileOutput(ctx, w, o)
	st := task.NewListDir(t)
	st.SetPath(o.Name)
	st.SetFileFunc(func(oo *typ.Object) {
		listFileOutput(ctx, w, oo)
	})
	st.SetDirFunc(func(oo *typ.Object) {
		listDirFunc(ctx, w, st, oo)
	})
	t.GetScheduler().Sync(ctx, st)
}

func initLsFlag() {
	LsCommand.Flags().BoolVarP(&lsFlag.humanReadable, constants.HumanReadableFlag, "h", false,
		i18n.Sprintf(`print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,
in order to reduce the number of digits to three or less using base 2 for sizes`))
	LsCommand.Flags().BoolVarP(&lsFlag.longFormat, constants.LongFormatFlag, "l", false,
		i18n.Sprintf(`list in long format and a total sum for all the file sizes is
output on a line before the long listing`))
	LsCommand.Flags().BoolVarP(&lsFlag.recursive, constants.RecursiveFlag, "R", false,
		i18n.Sprintf("recursively list subdirectories encountered"))
	// LsCommand.Flags().BoolVarP(&reverse, constants.ReverseFlag, "r", false,
	// 	"reverse the order of the sort to get reverse lexicographical order")
	// LsCommand.Flags().StringVarP(&lsFlag.Zone, constants.ZoneFlag, "z", "",
	// 	i18n.Sprintf("in which zone to do the operation"))
}

// listBucketOutput list buckets with normal slice
func listBucketOutput(ctx context.Context, w io.Writer, s storage.Storager) {
	logger := log.FromContext(ctx)
	m, err := s.Metadata()
	if err != nil {
		logger.Debug(
			log.String("action", "get metadata in listBucketOutput"),
			log.String("err", err.Error()),
		)
		i18n.Fprintf(w, "get metadata failed: %v\n", err)
		return
	}
	i18n.Fprintf(w, "%s\n", m.Name)
}

// listBucketLongOutput list buckets with long format
func listBucketLongOutput(ctx context.Context, w io.Writer, s storage.Storager, t types.SchedulerGetter) {
	logger := log.FromContext(ctx)
	rt := taskutils.NewAtStorageTask(10)
	rt.SetStorage(s)
	sst := task.NewStatStorage(rt)
	t.GetScheduler().Sync(ctx, sst)
	m, err := s.Metadata()
	if err != nil {
		logger.Debug(
			log.String("action", "get metadata in listBucketLongOutput"),
			log.String("err", err.Error()),
		)
		i18n.Fprintf(w, "get metadata failed: %v\n", err)
		return
	}

	// handle size separately from stat output for -h
	var size string
	if v, ok := sst.GetStorageInfo().GetSize(); ok {
		if lsFlag.humanReadable {
			size, err = utils.UnixReadableSize(datasize.ByteSize(v).HR())
			if err != nil {
				logger.Debug(
					log.String("action", "parse size in listBucketLongOutput"),
					log.Int("size", v),
					log.String("err", err.Error()),
				)
				i18n.Fprintf(w, "parse size <%v> failed [%v]\n", v, err)
				return
			}
		} else {
			size = datasize.ByteSize(v).String()
		}
	}
	statStorageOutput(w, m, sst.GetStorageInfo(), fmt.Sprintf("%%n %%l %s %%c", size))
}

func listFileOutput(ctx context.Context, w io.Writer, o *typ.Object) {
	logger := log.FromContext(ctx)
	if !lsFlag.longFormat {
		i18n.Fprintf(w, "%s\n", o.Name)
		return
	}

	var err error

	objACL := constants.ACLObject
	if o.Type == typ.ObjectTypeDir {
		objACL = constants.ACLDirectory
	}

	// default print size by bytes
	readableSize := strconv.FormatInt(o.Size, 10)
	if lsFlag.humanReadable {
		// if human readable flag true, print size as human readable format
		readableSize, err = utils.UnixReadableSize(datasize.ByteSize(o.Size).HR())
		if err != nil {
			logger.Debug(
				log.String("action", "parse size in listFileOutput"),
				log.Int("size", o.Size),
				log.String("key", o.Name),
				log.String("err", err.Error()),
			)
			i18n.Fprintf(w, "parse size <%v> failed [%v], key: <%s>\n", o.Size, err, o.Name)
			return
		}
		// 7 is the widest size of readable-size, like 1023.9K
		readableSize = text.AlignRight.Apply(readableSize, 7)
	}

	// if modified not exists (like dir), init str with blank
	modifiedStr := o.UpdatedAt.Format(constants.LsDefaultFormat)
	// output order: acl  size  lastModified  key
	// join with two space
	i18n.Fprintf(w, "%s  %s  %s  %s\n", objACL, readableSize, modifiedStr, o.Name)
}
