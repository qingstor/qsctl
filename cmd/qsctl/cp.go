package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/task/types"

	"github.com/yunify/qsctl/v2/utils"
)

var cpInput struct {
	ExpectSize           string
	MaximumMemoryContent string
}

var cpOutput struct {
	ExpectSize           int64
	MaximumMemoryContent int64

	Flow     constants.FlowType
	Path     string
	PathType constants.PathType
	Key      string
	KeyType  constants.KeyType

	Storage storage.ObjectStorage
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

func cpParse(_ *cobra.Command, args []string) (err error) {
	// Parse flags.
	if cpInput.ExpectSize != "" {
		cpOutput.ExpectSize, err = utils.ParseByteSize(expectSize)
		if err != nil {
			return err
		}
	}

	if cpInput.MaximumMemoryContent != "" {
		cpOutput.MaximumMemoryContent, err = utils.ParseByteSize(maximumMemoryContent)
		if err != nil {
			return err
		}
	}

	// Setup storage.
	cpOutput.Storage, err = storage.NewQingStorObjectStorage()
	if err != nil {
		return err
	}

	// Parse flow.
	src, dst := args[0], args[1]
	cpOutput.Flow = utils.ParseFlow(src, dst)

	var bucketName, objectKey string

	switch cpOutput.Flow {
	case constants.FlowToRemote:
		cpOutput.PathType, err = utils.ParsePath(src)
		if err != nil {
			return
		}
		cpOutput.Path = src

		cpOutput.KeyType, bucketName, objectKey, err = utils.ParseKey(dst)
		if err != nil {
			return
		}
		cpOutput.Key = objectKey
	case constants.FlowToLocal, constants.FlowAtRemote:
		cpOutput.PathType, err = utils.ParsePath(dst)
		if err != nil {
			return
		}
		cpOutput.Path = dst

		cpOutput.KeyType, bucketName, objectKey, err = utils.ParseKey(src)
		if err != nil {
			return
		}
		cpOutput.Key = objectKey
	default:
		panic("this case should never be switched")
	}
	err = cpOutput.Storage.SetupBucket(bucketName, "")
	if err != nil {
		return
	}

	return nil
}

func cpRun(cmd *cobra.Command, args []string) (err error) {
	err = cpParse(cmd, args)
	if err != nil {
		return
	}

	var t types.Tasker

	switch cpOutput.Flow {
	case constants.FlowToLocal:
		t, err = cpToLocal()
	case constants.FlowToRemote:
		t, err = cpToRemote()
	default:
		panic("this case should never be switched")
	}
	if err != nil {
		return err
	}

	t.Run()
	t.GetPool().Wait()
	return
}

func cpToLocal() (t types.Tasker, err error) {
	// TODO: support -r

	switch cpOutput.PathType {
	case constants.PathTypeLocalDir:
		return nil, constants.ErrorActionNotImplemented
	case constants.PathTypeStream:
		// TODO: RUN xxxTask
	case constants.PathTypeFile:
		// TODO: Run XX task
	}

	return nil, nil
}

func cpToRemote() (t types.Tasker, err error) {
	switch cpOutput.PathType {
	case constants.PathTypeLocalDir:
		return nil, constants.ErrorActionNotImplemented
	case constants.PathTypeStream:
		t = task.NewCopyStreamTask(cpOutput.Key, cpOutput.Storage)
	case constants.PathTypeFile:
		t = task.NewCopyFileTask(cpOutput.Path, cpOutput.Key, cpOutput.Storage)
	default:
		panic("invalid path type")
	}
	return t, nil
}
