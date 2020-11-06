package main

import (
	"fmt"
	"path/filepath"

	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/c-bata/go-prompt"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/shellutils"
	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type rmFlags struct {
	recursive bool
}

var rmFlag = rmFlags{}

// RmCommand will handle remove object command.
var RmCommand = &cobra.Command{
	Use:   "rm qs://<bucket_name>/<object_key>",
	Short: i18n.Sprintf("remove a remote object"),
	Long:  i18n.Sprintf("qsctl rm remove the object with given object key"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Remove a single object: qsctl rm qs://bucket-name/object-key"),
		i18n.Sprintf("Remove objects with prefix: qsctl rm qs://bucket-name/prefix -r"),
	),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := rmRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "rm", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		rmFlag = rmFlags{}
	},
}

func initRmFlag() {
	RmCommand.Flags().BoolVarP(&rmFlag.recursive, constants.RecursiveFlag, "r",
		false, i18n.Sprintf("recursively delete keys under a specific prefix"))
}

func rmRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewAtStorageTask()
	workDir, err := utils.ParseAtStorageInput(rootTask, args[0])
	if err != nil {
		return
	}

	if rootTask.GetType() == typ.ObjectTypeDir && !rmFlag.recursive {
		return fmt.Errorf(i18n.Sprintf("-r is required to remove a directory"))
	}

	if rmFlag.recursive && rootTask.GetType() != typ.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("path should be a directory while -r is set"))
	}

	key := filepath.Join(workDir, rootTask.GetPath())
	if rmFlag.recursive {
		t := task.NewDeleteDir(rootTask)
		t.SetHandleObjCallbackFunc(func(o *typ.Object) {
			i18n.Fprintf(c.OutOrStdout(), "<%s> removed\n", o.Name)
		})
		if err := t.Run(c.Context()); err != nil {
			return err
		}

		i18n.Fprintf(c.OutOrStdout(), "Dir <%s> removed.\n", key)
		return nil
	}

	t := task.NewDeleteFile(rootTask)
	if err := t.Run(c.Context()); err != nil {
		return err
	}

	i18n.Fprintf(c.OutOrStdout(), "File <%s> removed.\n", key)
	return nil
}

// rmShellHandler handle rm in shell
type rmShellHandler struct{}

func (r rmShellHandler) preRunE(args []string) error {
	err := RmCommand.Flags().Parse(args)
	if err != nil {
		return err
	}
	_, _, key, err := utils.ParseQsPath(RmCommand.Flags().Args()[0])
	if err != nil {
		return err
	}
	// try to get confirm input
	input := prompt.Input(
		i18n.Sprintf("confirm to remove <%s>? [y/N] ", key),
		noSuggests)
	if !shellutils.CheckYes(input) {
		return fmt.Errorf(i18n.Sprintf("not confirmed"))
	}
	return nil
}

func (r rmShellHandler) postRun(_ error) {
	return
}
