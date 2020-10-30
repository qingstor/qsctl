package main

import (
	"fmt"
	"path/filepath"

	"github.com/aos-dev/go-storage/v2/types"
	tsk "github.com/qingstor/noah/pkg/task"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

type syncFlags struct {
	checkMD5       bool
	dryRun         bool
	existing       bool
	ignoreExisting bool
	recursive      bool
	update         bool
	multipartFlags
	inExcludeFlags
}

var syncFlag = syncFlags{}

// SyncCommand will handle sync command.
var SyncCommand = &cobra.Command{
	Use:   "sync <source-path> <dest-path>",
	Short: i18n.Sprintf("sync between local directory and QS-Directory"),
	Long: i18n.Sprintf(`qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Sync local directory to QS-Directory: qsctl sync . qs://bucket-name/dir/"),
		i18n.Sprintf("Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/"),
		i18n.Sprintf("Sync directory recursively: qsctl sync qs://bucket-name/test/ test_local/ -r"),
		i18n.Sprintf("Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name/dir/ --ignore-existing"),
		i18n.Sprintf("Only sync files that newer than files on receiver: qsctl sync . qs://bucket-name/dir/ --update"),
		i18n.Sprintf("Only sync files that already exist on receiver: qsctl sync . qs://bucket-name/dir/ --existing"),
		i18n.Sprintf("Show files that would sync (but not really do): qsctl sync . qs://bucket-name/dir/ --dry-run"),
	),
	Args: cobra.ExactArgs(2),
	PreRunE: func(c *cobra.Command, args []string) error {
		if err := validateSyncFlag(c, args); err != nil {
			return err
		}
		if err := parseSyncFlag(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "sync", err.Error())
		}
	},
	PostRun: func(_ *cobra.Command, _ []string) {
		syncFlag = syncFlags{}
	},
}

func syncRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask()
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() != types.ObjectTypeDir || rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("both source and destination should be directories"))
	}

	t := task.NewSync(rootTask)
	t.SetCheckMD5(syncFlag.checkMD5)
	t.SetRecursive(syncFlag.recursive)
	t.SetPartThreshold(syncFlag.partThreshold)
	if syncFlag.partSize != 0 {
		t.SetPartSize(syncFlag.partSize)
	}

	// set check functions
	var fn []func(task tsk.Task) tsk.Task
	if syncFlag.existing {
		fn = append(fn, task.NewIsDestinationObjectExistTask)
	}
	if syncFlag.ignoreExisting {
		fn = append(fn, task.NewIsDestinationObjectNotExistTask)
	}
	if syncFlag.update {
		fn = append(fn, task.NewIsUpdateAtGreaterTask)
	}
	if syncFlag.excludeRegx != nil {
		fn = append(fn, func(tt tsk.Task) tsk.Task {
			st := task.NewIsSourcePathExcludeInclude(tt)
			st.SetExcludeRegexp(syncFlag.excludeRegx)
			st.SetIncludeRegexp(syncFlag.includeRegx)
			return st
		})
	}
	t.SetCheckTasks(fn)

	if syncFlag.dryRun {
		t.SetDryRunFunc(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "%s\n", o.Name)
		})
	} else {
		t.SetHandleObjCallbackFunc(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "<%s> synced\n", o.Name)
		})
	}

	if err := t.Run(c.Context()); err != nil {
		return err
	}

	if h := taskutils.HandlerFromContext(c.Context()); h != nil {
		h.WaitProgress()
	}

	i18n.Fprintf(c.OutOrStdout(), "Dir <%s> and <%s> synced.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return nil

}

func initSyncFlag() {
	SyncCommand.Flags().BoolVarP(&syncFlag.dryRun, constants.DryRunFlag, "n", false,
		i18n.Sprintf(`show what would have been transferred`))
	SyncCommand.Flags().BoolVar(&syncFlag.existing, constants.ExistingFlag, false,
		i18n.Sprintf(`skip creating new files in dest dirs`))
	SyncCommand.Flags().BoolVar(&syncFlag.ignoreExisting, constants.IgnoreExistingFlag, false,
		i18n.Sprintf(`skip updating files in dest dirs, only copy those not exist`))
	SyncCommand.Flags().BoolVarP(&syncFlag.recursive, constants.RecursiveFlag, "r", false,
		i18n.Sprintf(`recurse into sub directories`))
	SyncCommand.Flags().BoolVarP(&syncFlag.update, constants.UpdateFlag, "u", false,
		i18n.Sprintf(`skip files that are newer in dest dirs`))
	SyncCommand.Flags().StringVar(&syncFlag.partThresholdStr, constants.PartThresholdFlag, "",
		i18n.Sprintf("set threshold to enable multipart upload"))
	SyncCommand.Flags().StringVar(&syncFlag.partSizeStr, constants.PartSizeFlag, "",
		i18n.Sprintf("set part size for multipart upload"))
	SyncCommand.Flags().StringVar(&syncFlag.excludeRegxStr, constants.ExcludeRegexp, "",
		i18n.Sprintf("regular expression for files to exclude"))
	SyncCommand.Flags().StringVar(&syncFlag.includeRegxStr, constants.IncludeRegexp, "",
		i18n.Sprintf("regular expression for files to include (not work if exclude-regx not set)"))
}

func validateSyncFlag(_ *cobra.Command, _ []string) (err error) {
	if syncFlag.existing && syncFlag.ignoreExisting {
		return fmt.Errorf(i18n.Sprintf("both --existing and --ignore-existing are set, no files would be synced"))
	}
	return nil
}

func parseSyncFlag() error {
	if err := syncFlag.multipartFlags.parse(); err != nil {
		return err
	}

	if err := syncFlag.inExcludeFlags.parse(); err != nil {
		return err
	}
	return nil
}
