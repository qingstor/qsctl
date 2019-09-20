package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var cpInput struct {
	ExpectSize           string
	MaximumMemoryContent string
}

// CpCommand will handle copy command.
var CpCommand = &cobra.Command{
	Use:   "cp <source-path> <dest-path>",
	Short: "copy from/to qingstor",
	Long:  "qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout",
	Example: utils.AlignPrintWithColon(
		"Copy file: qsctl cp /path/to/file qs://prefix/a",
		"Copy folder: qsctl cp qs://prefix/a /path/to/folder -r",
		"Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin",
		"Write to stdout: qsctl cp qs://prefix/b - > /path/to/file",
	),
	Args: cobra.ExactArgs(2),
	RunE: cpRun,
}

func initCpFlag() {
	CpCommand.PersistentFlags().StringVar(&cpInput.ExpectSize,
		"expect-size",
		"",
		"expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)",
	)
	CpCommand.PersistentFlags().StringVar(&cpInput.MaximumMemoryContent,
		"maximum-memory-content",
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
}

func cpParse(t *task.CopyTask, args []string) (err error) {
	// Setup storage.
	stor, err := storage.NewQingStorObjectStorage()
	if err != nil {
		return err
	}
	t.SetStorage(stor)

	// Parse flow.
	src, dst := args[0], args[1]
	flow := utils.ParseFlow(src, dst)
	t.SetFlowType(flow)

	switch flow {
	case constants.FlowToRemote:
		pathType, err := utils.ParsePath(src)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		t.SetPath(src)

		keyType, bucketName, objectKey, err := utils.ParseKey(dst)
		if err != nil {
			return err
		}
		t.SetKeyType(keyType)
		t.SetObjectKey(objectKey)
		t.SetBucketName(bucketName)
	case constants.FlowToLocal, constants.FlowAtRemote:
		pathType, err := utils.ParsePath(dst)
		if err != nil {
			return err
		}
		t.SetPathType(pathType)
		t.SetPath(dst)

		keyType, bucketName, objectKey, err := utils.ParseKey(src)
		if err != nil {
			return err
		}
		t.SetKeyType(keyType)
		t.SetObjectKey(objectKey)
		t.SetBucketName(bucketName)
	default:
		panic("this case should never be switched")
	}
	err = stor.SetupBucket(t.GetBucketName(), "")
	if err != nil {
		return
	}

	return nil
}

func cpRun(cmd *cobra.Command, args []string) (err error) {
	t := task.NewCopyTask(func(t *task.CopyTask) {
		err = cpParse(t, args)
		if err != nil {
			return
		}
	})

	t.Run()
	t.Wait()
	return
}
