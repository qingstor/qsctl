package main

import (
	"fmt"
	"strconv"
	"strings"

	typ "github.com/Xuanwo/storage/types"
	"github.com/c2h5oh/datasize"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
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

func statRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtStorageTask(10)
	_, err = utils.ParseAtStorageInput(rootTask, args[0])
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
	input = strings.ReplaceAll(input, "%n", om.ID)

	if v, ok := om.GetContentType(); ok {
		input = strings.ReplaceAll(input, "%F", v)
	}
	if v, ok := om.GetContentMD5(); ok {
		input = strings.ReplaceAll(input, "%h", v)
	}
	input = strings.ReplaceAll(input, "%s", strconv.FormatInt(om.Size, 10))
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

	content = append(content, i18n.Sprintf("Key: %s", om.ID))
	content = append(content, i18n.Sprintf("Size: %s", datasize.ByteSize(om.Size).String()))
	if v, ok := om.GetContentType(); ok {
		content = append(content, i18n.Sprintf("Type: %s", v))
	}
	if v, ok := om.GetStorageClass(); ok {
		content = append(content, i18n.Sprintf("StorageClass: %s", v))
	}
	if v, ok := om.GetContentMD5(); ok {
		content = append(content, i18n.Sprintf("MD5: %s", v))
	}
	content = append(content, i18n.Sprintf("UpdatedAt: %s", om.UpdatedAt.String()))

	fmt.Println(utils.AlignPrintWithColon(content...))
}
