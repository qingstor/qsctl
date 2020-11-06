package main

import (
	"io"
	"strconv"
	"strings"

	"github.com/aos-dev/go-service-qingstor"
	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/c2h5oh/datasize"
	"github.com/qingstor/noah/pkg/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type statFlags struct {
	format string
}

var statFlag = statFlags{}

// StatCommand will handle stat command.
var StatCommand = &cobra.Command{
	Use:   "stat qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("stat a remote object"),
	Long:  i18n.Sprintf("qsctl stat show the detailed info of this object"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Stat object: qsctl stat qs://prefix/a"),
		i18n.Sprintf("Stat bucket: qsctl stat qs://bucket-name"),
	),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := statRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "stat", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		statFlag = statFlags{}
	},
}

func statRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtStorageTask()
	workDir, err := utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	// work dir is root path and path blank, handle it as stat bucket
	if workDir == "/" && rootTask.GetPath() == "" {
		t := task.NewStatStorage(rootTask)
		if err := t.Run(c.Context()); err != nil {
			return err
		}
		sm, err := t.GetStorage().Metadata()
		if err != nil {
			return types.NewErrUnhandled(err)
		}

		statStorageOutput(c.OutOrStdout(), sm, t.GetStorageInfo(), statFlag.format)
		return nil
	}

	t := task.NewStatFile(rootTask)
	if err := t.Run(c.Context()); err != nil {
		return err
	}

	statFileOutput(c.OutOrStdout(), t.GetObject(), statFlag.format)
	return
}

func initStatFlag() {
	StatCommand.Flags().StringVar(&statFlag.format, constants.FormatFlag, "",
		i18n.Sprintf(`use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content etag of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch

The valid format sequences for buckets:

  %n   bucket name
  %l   bucket location
  %s   total size, in bytes
  %c   count of files in this bucket
	`),
	)
}

// statFileFormat replace format placeholder into relevant object's attr
func statFileFormat(input string, om *typ.Object) string {
	input = strings.ReplaceAll(input, "%n", om.ID)

	if v, ok := om.GetContentType(); ok {
		input = strings.ReplaceAll(input, "%F", v)
	}
	if v, ok := om.GetETag(); ok {
		input = strings.ReplaceAll(input, "%h", v)
	}
	if v, ok := om.GetSize(); ok {
		input = strings.ReplaceAll(input, "%s", strconv.FormatInt(v, 10))
	}
	if v, ok := om.GetUpdatedAt(); ok {
		input = strings.ReplaceAll(input, "%y", v.String())
		input = strings.ReplaceAll(input, "%Y", strconv.FormatInt(v.Unix(), 10))
	}
	return input
}

// statStorageFormat replace format placeholder into relevant storage's attr
func statStorageFormat(input string, sm typ.StorageMeta, ss typ.StorageStatistic) string {
	input = strings.ReplaceAll(input, "%n", sm.Name)

	if v, ok := sm.GetLocation(); ok {
		input = strings.ReplaceAll(input, "%l", v)
	}
	if v, ok := ss.GetSize(); ok {
		input = strings.ReplaceAll(input, "%s", strconv.FormatInt(v, 10))
	}
	if v, ok := ss.GetCount(); ok {
		input = strings.ReplaceAll(input, "%c", strconv.FormatInt(v, 10))
	}
	return input
}

// statFileOutput print object info to writer with given format
func statFileOutput(w io.Writer, om *typ.Object, format string) {
	// if format string was set, print result as format string
	if format != "" {
		i18n.Fprintf(w, "%s\n", statFileFormat(format, om))
		return
	}

	var content []string

	content = append(content, i18n.Sprintf("Key: %s", om.ID))
	if v, ok := om.GetSize(); ok {
		content = append(content, i18n.Sprintf("Size: %s", datasize.ByteSize(v).String()))
	}
	if v, ok := om.GetContentType(); ok {
		content = append(content, i18n.Sprintf("Type: %s", v))
	}
	if v, ok := om.GetETag(); ok {
		content = append(content, i18n.Sprintf("ETag: %s", v))
	}
	if v, ok := om.ObjectMeta.Get(qingstor.InfoObjectMetaStorageClass); ok {
		content = append(content, i18n.Sprintf("StorageClass: %s", v))
	}
	if v, ok := om.GetUpdatedAt(); ok {
		content = append(content, i18n.Sprintf("UpdatedAt: %s", v.String()))
	}

	i18n.Fprintf(w, "%s\n", utils.AlignPrintWithColon(content...))
}

// statStorageOutput print stat storage info to writer with given StorageMeta, StorageStatistic and format
func statStorageOutput(w io.Writer, sm typ.StorageMeta, ss typ.StorageStatistic, format string) {
	if format != "" {
		i18n.Fprintf(w, "%s\n", statStorageFormat(format, sm, ss))
		return
	}

	var content []string
	content = append(content, i18n.Sprintf("Name: %s", sm.Name))
	if v, ok := sm.GetLocation(); ok {
		content = append(content, i18n.Sprintf("Location: %s", v))
	}
	if v, ok := ss.GetSize(); ok {
		content = append(content, i18n.Sprintf("Size: %s", datasize.ByteSize(v).String()))
	}
	if v, ok := ss.GetCount(); ok {
		content = append(content, i18n.Sprintf("Count: %s", strconv.FormatInt(v, 10)))
	}

	i18n.Fprintf(w, "%s\n", utils.AlignPrintWithColon(content...))
}
