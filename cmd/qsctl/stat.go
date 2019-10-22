package main

import (
	"fmt"
	"strconv"
	"strings"

	storageType "github.com/Xuanwo/storage/types"
	"github.com/c2h5oh/datasize"
	"github.com/spf13/cobra"

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
	Short: "stat a remote object",
	Long:  "qsctl stat show the detailed info of this object",
	Example: utils.AlignPrintWithColon(
		"Stat object: qsctl stat qs://prefix/a",
	),
	Args: cobra.ExactArgs(1),
	RunE: statRun,
}

func statRun(_ *cobra.Command, args []string) (err error) {
	t := task.NewStatTask(func(t *task.StatTask) {
		_, bucketName, objectKey, err := utils.ParseQsPath(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
		// for now, only support stat object
		// if keyType != constants.KeyTypeObject {
		// 	t.TriggerFault(fmt.Errorf("key type is not match"))
		// 	return
		// }
		t.SetDestinationPath(objectKey)

		srv, err := NewQingStorService()
		if err != nil {
			t.TriggerFault(err)
			return
		}

		stor, err := srv.Get(bucketName)
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetDestinationStorage(stor)
	})

	t.Run()
	t.Wait()

	if t.ValidateFault() {
		return t.GetFault()
	}

	statOutput(t, statInput.format)
	return
}

func initStatFlag() {
	StatCommand.Flags().StringVar(&statInput.format, constants.FormatFlag, "",
		`use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable, e.g: 2006-01-02 15:04:05 +0000 UTC
  %Y   time of last data modification, seconds since Epoch
	`,
	)
}

func statFormat(input string, om *storageType.Object) string {
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
	if v, ok := om.GetUpdatedAt(); ok {
		input = strings.ReplaceAll(input, "%y", v.String())
		input = strings.ReplaceAll(input, "%Y", strconv.FormatInt(v.Unix(), 10))
	}

	return input
}

func statOutput(t *task.StatTask, format string) {
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
	if v, ok := om.GetStorageClass(); ok {
		content = append(content, "StorageClass: "+v)
	}
	if v, ok := om.GetChecksum(); ok {
		content = append(content, "MD5: "+v)
	}
	if v, ok := om.GetUpdatedAt(); ok {
		content = append(content, "UpdatedAt: "+v.String())
	}

	fmt.Println(utils.AlignPrintWithColon(content...))
}
