package main

import (
	"fmt"
	"path/filepath"

	"github.com/Xuanwo/storage/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type cpFlags struct {
	checkMD5             bool
	expectSize           string
	maximumMemoryContent string
	recursive            bool
}

var cpFlag = cpFlags{}

// CpCommand will handle copy command.
var CpCommand = &cobra.Command{
	Use:   "cp <source-path> <dest-path>",
	Short: i18n.Sprintf("copy from/to qingstor"),
	Long:  i18n.Sprintf("qsctl cp can copy file/folder/stdin to qingstor or copy qingstor objects to local/stdout"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Copy file: qsctl cp /path/to/file qs://prefix/a"),
		i18n.Sprintf("Copy folder: qsctl cp /path/to/folder qs://prefix/a/ -r"),
		i18n.Sprintf("Copy all files in folder: qsctl cp /path/to/folder/ qs://prefix/a/ -r"),
		i18n.Sprintf("Read from stdin: cat /path/to/file | qsctl cp - qs://prefix/stdin"),
		i18n.Sprintf("Write to stdout: qsctl cp qs://prefix/b - > /path/to/file"),
	),
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cpRun(cmd, args); err != nil {
			i18n.Printf("Execute %s command error: %s", "cp", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		cpFlag = cpFlags{}
	},
}

func initCpFlag() {
	CpCommand.PersistentFlags().StringVar(&cpFlag.expectSize,
		constants.ExpectSizeFlag,
		"",
		i18n.Sprintf(`expected size of the input file
accept: 100MB, 1.8G
(only used and required for input from stdin)`),
	)
	CpCommand.PersistentFlags().StringVar(&cpFlag.maximumMemoryContent,
		constants.MaximumMemoryContentFlag,
		"",
		i18n.Sprintf(`maximum content loaded in memory
(only used for input from stdin)`),
	)
	CpCommand.Flags().BoolVarP(&cpFlag.recursive,
		constants.RecursiveFlag,
		"r",
		false,
		i18n.Sprintf("copy directory recursively"),
	)
}

func cpRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !cpFlag.recursive {
		return fmt.Errorf(i18n.Sprintf("-r is required to copy a directory"))
	}

	if cpFlag.recursive && rootTask.GetSourceType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("src should be a directory while -r is set"))
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir &&
		rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("cannot copy a directory to a non-directory dest"))
	}

	if cpFlag.recursive {
		t := task.NewCopyDir(rootTask)
		t.SetCheckMD5(cpFlag.checkMD5)
		t.SetHandleObjCallback(func(o *types.Object) {
			fmt.Println(i18n.Sprintf("<%s> copied", o.Name))
		})
		t.SetCheckTasks(nil)
		t.Run()

		if t.GetFault().HasError() {
			return t.GetFault()
		}

		i18n.Printf("Dir <%s> copied to <%s>.\n",
			filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
		return nil
	}

	t := task.NewCopyFile(rootTask)
	t.SetCheckMD5(cpFlag.checkMD5)
	t.SetCheckTasks(nil)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	i18n.Printf("File <%s> copied to <%s>.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return
}
