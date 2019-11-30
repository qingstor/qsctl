package main

import (
	"fmt"
	"strconv"
	"strings"

	typ "github.com/Xuanwo/storage/types"
	"github.com/c2h5oh/datasize"
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/yunify/qsctl/v2/pkg/i18n"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var statInput struct {
	format string
}

// StatCommand will handle stat command.
var StatCommand = &cobra.Command{
	Use:   "stat qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("stat a remote object"),
	Long:  i18n.Sprintf("qsctl stat show the detailed info of this object"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Stat object: qsctl stat qs://prefix/a"),
	),
	Args: cobra.ExactArgs(1),
	RunE: statRun,
}

func statRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtStorageTask(10)
	err = utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	t := task.NewStatFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	statOutput(t, statInput.format)
	return
}

func initStatFlag() {
	StatCommand.Flags().StringVar(&statInput.format, constants.FormatFlag, "",
		i18n.Sprintf(`use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`),
	)
}

func statFormat(input string, om *typ.Object) string {
	input = strings.ReplaceAll(input, "%n", om.Name)

	if v, ok := om.GetType(); ok {
		input = strings.ReplaceAll(input, "%F", v)
	}
	if v, ok := om.GetChecksum(); ok {
		input = strings.ReplaceAll(input, "%h", v)
	}
	if v, ok := om.GetSize(); ok {
		input = strings.ReplaceAll(input, "%s", strconv.FormatInt(v, 10))
	}
	input = strings.ReplaceAll(input, "%y", om.UpdatedAt.String())
	input = strings.ReplaceAll(input, "%Y", strconv.FormatInt(om.UpdatedAt.Unix(), 10))

	return input
}

func statOutput(t *task.StatFileTask, format string) {
	// if format string was set, print result as format string
	if format != "" {
		fmt.Println(statFormat(format, t.GetObject()))
		return
	}

	om := t.GetObject()
	var content []string

	content = append(content, "Key: "+om.Name)
	if v, ok := om.GetSize(); ok {
		content = append(content, "Size: "+datasize.ByteSize(v).String())
	}
	if v, ok := om.GetType(); ok {
		content = append(content, "Type: "+v)
	}
	if v, ok := om.GetClass(); ok {
		content = append(content, "StorageClass: "+v)
	}
	if v, ok := om.GetChecksum(); ok {
		content = append(content, "MD5: "+v)
	}
	content = append(content, "UpdatedAt: "+om.UpdatedAt.String())

	fmt.Println(utils.AlignPrintWithColon(content...))
}
