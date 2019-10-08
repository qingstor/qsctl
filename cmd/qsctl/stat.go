package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	storageType "github.com/yunify/qsctl/v2/pkg/types/storage"
	"github.com/yunify/qsctl/v2/storage"
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
		keyType, bucketName, objectKey, err := utils.ParseKey(args[0])
		if err != nil {
			t.TriggerFault(err)
			return
		}
		// for now, only support stat object
		if keyType != constants.KeyTypeObject {
			t.TriggerFault(fmt.Errorf("key type is not match"))
			return
		}
		t.SetKey(objectKey)

		stor, err := storage.NewQingStorObjectStorage()
		if err != nil {
			t.TriggerFault(err)
			return
		}
		t.SetStorage(stor)

		if err = stor.SetupBucket(bucketName, ""); err != nil {
			t.TriggerFault(err)
			return
		}
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
  %y   time of last data modification, human-readable
  %Y   time of last data modification, seconds since Epoch
	`,
	)
}

func statFormat(input string, om *storageType.ObjectMeta) string {
	input = strings.ReplaceAll(input, "%F", om.ContentType)
	input = strings.ReplaceAll(input, "%h", om.ETag)
	input = strings.ReplaceAll(input, "%n", om.Key)
	input = strings.ReplaceAll(input, "%s", strconv.FormatInt(om.ContentLength, 10))
	input = strings.ReplaceAll(input, "%y", om.LastModified.String())
	input = strings.ReplaceAll(input, "%Y", strconv.FormatInt(om.LastModified.Unix(), 10))

	return input
}

func statOutput(t *task.StatTask, format string) {
	// if format string was set, print result as format string
	if format != "" {
		fmt.Println(statFormat(format, t.GetObjectMeta()))
		return
	}

	om := t.GetObjectMeta()
	content := []string{
		"Key: " + om.Key,
		"Size: " + datasize.ByteSize(om.ContentLength).String(),
		"Type: " + om.ContentType,
		"Modify: " + om.LastModified.String(),
		"StorageClass: " + om.StorageClass,
	}

	if om.ETag != "" {
		content = append(content, "MD5: "+om.ETag)
	}

	fmt.Println(utils.AlignPrintWithColon(content...))
}
