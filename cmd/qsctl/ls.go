package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/c2h5oh/datasize"
	"github.com/jedib0t/go-pretty/text"
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
		rootTask := taskutils.NewAtServiceTask()
		err = utils.ParseAtServiceInput(rootTask)
		if err != nil {
			return
		}

		t := task.NewListStorage(rootTask)
		t.SetZone(globalFlag.zone)
		if err := t.Run(c.Context()); err != nil {
			return err
		}

		if err := listStorageFunc(c.Context(), t, c.OutOrStdout()); err != nil {
			return err
		}
		return nil
	}

	rootTask := taskutils.NewAtStorageTask()
	_, err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	// conduct the init listDirTask
	t := task.NewListDir(rootTask)
	lister, ok := rootTask.GetStorage().(typ.DirLister)
	if !ok {
		return types.NewErrStorageInsufficientAbility(nil)
	}
	t.SetDirLister(lister)
	// list dir recursively
	if err := listDirFunc(c.Context(), c.OutOrStdout(), t, t.GetPath()); err != nil {
		return err
	}

	return nil
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

// listDirFunc handle list dir task recursively
func listDirFunc(ctx context.Context, w io.Writer, t *task.ListDirTask, path string) error {
	st := task.NewListDir(t)
	st.SetPath(path)
	if err := t.Sync(ctx, st); err != nil {
		return err
	}

	it := st.GetObjectIter()
	for {
		s, err := it.Next()
		if err != nil {
			if errors.Is(err, typ.IterateDone) {
				break
			}
			return err
		}

		switch s.Type {
		case typ.ObjectTypeDir:
			// always print dir as file first
			if err := listFileOutput(ctx, w, s); err != nil {
				return err
			}
			// if recursive flag not set, do not list dir recursively
			if !lsFlag.recursive {
				continue
			}
			// else list dir recursively
			if err := listDirFunc(ctx, w, st, s.Name); err != nil {
				return err
			}
		case typ.ObjectTypeFile:
			if err := listFileOutput(ctx, w, s); err != nil {
				return err
			}
		default: // print tip for other type object
			i18n.Fprintf(w, "invalid object <%s> type: %v\n", s.Name, s.Type)
			continue
		}
	}
	return nil
}

// listStorageFunc handle storage iter with different output func
func listStorageFunc(ctx context.Context, t types.StorageIterGetter, w io.Writer) error {
	it := t.GetStorageIter()
	for {
		s, err := it.Next()
		if err != nil {
			if errors.Is(err, typ.IterateDone) {
				break
			}
			return err
		}

		if lsFlag.longFormat {
			err = listBucketLongOutput(ctx, w, s)
		} else {
			err = listBucketOutput(ctx, w, s)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// listBucketOutput list buckets with normal slice
func listBucketOutput(_ context.Context, w io.Writer, s typ.Storager) error {
	m, err := s.Metadata()
	if err != nil {
		i18n.Fprintf(w, "get metadata failed: %v\n", err)
		return err
	}
	i18n.Fprintf(w, "%s\n", m.Name)
	return nil
}

// listBucketLongOutput list buckets with long format
func listBucketLongOutput(ctx context.Context, w io.Writer, s typ.Storager) error {
	t := taskutils.NewAtStorageTask()
	st := task.NewStatStorage(t)
	st.SetStorage(s)
	if err := st.Run(ctx); err != nil {
		return err
	}
	m, err := s.MetadataWithContext(ctx)
	if err != nil {
		i18n.Fprintf(w, "get metadata failed: %v\n", err)
		return err
	}

	// handle size separately from stat output for -h
	// because we want to reuse statStorageOutput(), which cannot handle size with humanReadable,
	// so transfer size into string here and conduct the format to call statStorageOutput()
	var size string
	if v, ok := st.GetStorageInfo().GetSize(); ok {
		if lsFlag.humanReadable {
			size, err = utils.UnixReadableSize(datasize.ByteSize(v).HR())
			if err != nil {
				i18n.Fprintf(w, "parse size <%v> failed [%v]\n", v, err)
				return err
			}
		} else {
			size = datasize.ByteSize(v).String()
		}
	}

	// conduct format for stat storage output
	format := fmt.Sprintf("%%n %%l %s %%c", size) // %n %l size %c
	statStorageOutput(w, m, st.GetStorageInfo(), format)
	return nil
}

// listFileOutput print object with given flags, such as long format, human readable
func listFileOutput(_ context.Context, w io.Writer, o *typ.Object) (err error) {
	if !lsFlag.longFormat {
		i18n.Fprintf(w, "%s\n", o.Name)
		return nil
	}

	objACL := constants.ACLObject
	if o.Type == typ.ObjectTypeDir {
		objACL = constants.ACLDirectory
	}

	// default print size by bytes, if object size not valid, set size to 0
	var size int64
	if v, ok := o.GetSize(); ok {
		size = v
	}
	readableSize := strconv.FormatInt(size, 10)
	// if human readable flag true, print size as human readable format
	if lsFlag.humanReadable {
		readableSize, err = utils.UnixReadableSize(datasize.ByteSize(size).HR())
		if err != nil {
			i18n.Fprintf(w, "parse size <%v> failed [%v], key: <%s>\n", size, err, o.Name)
			return err
		}
		// 7 is the widest size of readable-size, like 1023.9K
		readableSize = text.AlignRight.Apply(readableSize, 7)
	}

	// if updatedAt not exists (like dir), init with blank time (Jan 01 00:00)
	var updatedAt time.Time
	if v, ok := o.GetUpdatedAt(); ok {
		updatedAt = v
	}
	modifiedStr := updatedAt.Format(constants.LsDefaultFormat)
	// output order: acl  size  lastModified  key
	// join with two space
	i18n.Fprintf(w, "%s  %s  %s  %s\n", objACL, readableSize, modifiedStr, o.Name)
	return nil
}
