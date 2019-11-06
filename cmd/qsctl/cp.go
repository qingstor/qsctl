package main

import (
	"fmt"
	"path/filepath"

	"github.com/Xuanwo/storage/types"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/cmd/qsctl/taskutils"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task"
	"github.com/yunify/qsctl/v2/utils"
)

var cpInput struct {
	ExpectSize           string
	MaximumMemoryContent string
	Recursive            bool
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
		constants.ExpectSizeFlag,
		"",
		"expected size of the input file"+
			"accept: 100MB, 1.8G\n"+
			"(only used and required for input from stdin)",
	)
	CpCommand.PersistentFlags().StringVar(&cpInput.MaximumMemoryContent,
		constants.MaximumMemoryContentFlag,
		"",
		"maximum content loaded in memory\n"+
			"(only used for input from stdin)",
	)
	CpCommand.Flags().BoolVarP(&cpInput.Recursive,
		constants.RecursiveFlag,
		"r",
		false,
		"copy directory recursively")
}

func cpRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewBetweenStorageTask(10)
	err = utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() == types.ObjectTypeDir && !cpInput.Recursive {
		return fmt.Errorf("-r is required to delete a directory")
	}

	if err = HandleCpStorageBaseAndPath(rootTask); err != nil {
		return err
	}

	if cpInput.Recursive {
		t := task.NewCopyDir(rootTask)
		t.Run()

		if t.GetFault().HasError() {
			return t.GetFault()
		}
		cpOutput(args[0])
		return nil
	}

	t := task.NewCopyFile(rootTask)
	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}
	cpOutput(args[0])
	return
}

func cpOutput(path string) {
	fmt.Printf("Key <%s> copied.\n", path)
}

// HandleCpStorageBaseAndPath set work dir and path for cp cmd.
func HandleCpStorageBaseAndPath(t *taskutils.BetweenStorageTask) error {
	// In operation cp, we set source storage to dir of the source path.
	srcPath, err := filepath.Abs(t.GetSourcePath())
	if err != nil {
		return err
	}
	if err = t.GetSourceStorage().Init(types.WithWorkDir(filepath.Dir(srcPath))); err != nil {
		return err
	}
	t.SetSourcePath(filepath.Base(srcPath))

	// Destination path depends on different condition.
	dstPath, err := filepath.Abs(t.GetDestinationPath())
	if err != nil {
		return err
	}
	// if copy dir
	if cpInput.Recursive {
		if err := t.GetDestinationStorage().Init(types.WithWorkDir(dstPath)); err != nil {
			return err
		}
		t.SetDestinationPath("")
		return nil
	}
	// NOT copy dir. Copy file to a dir, we need to get destination key from the source.
	if t.GetDestinationType() == types.ObjectTypeDir {
		if err := t.GetDestinationStorage().Init(types.WithWorkDir(dstPath)); err != nil {
			return err
		}
		t.SetDestinationPath(t.GetSourcePath())
		return nil
	}
	// Copy to a file, get destination directly.
	if err := t.GetDestinationStorage().Init(types.WithWorkDir(filepath.Dir(dstPath))); err != nil {
		return err
	}
	t.SetDestinationPath(filepath.Base(dstPath))
	return nil
}
