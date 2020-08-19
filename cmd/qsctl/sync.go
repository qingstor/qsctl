package main

import (
	"fmt"
	"path/filepath"

	"github.com/aos-dev/go-storage/v2/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
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
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() != types.ObjectTypeDir || rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf(i18n.Sprintf("both source and destination should be directories"))
	}

	if syncFlag.existing && syncFlag.ignoreExisting {
		return fmt.Errorf(i18n.Sprintf("both --existing and --ignore-existing are set, no files would be synced"))
	}

	t := task.NewSync(rootTask)
	t.SetDryRun(syncFlag.dryRun)
	t.SetExisting(syncFlag.existing)
	t.SetIgnoreExisting(syncFlag.ignoreExisting)
	t.SetCheckMD5(syncFlag.checkMD5)
	t.SetRecursive(syncFlag.recursive)
	t.SetUpdate(syncFlag.update)
	if syncFlag.dryRun {
		t.SetDryRunFunc(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "%s\n", o.Name)
		})
	} else {
		t.SetDryRunFunc(nil)
		t.SetHandleObjCallback(func(o *types.Object) {
			i18n.Fprintf(c.OutOrStdout(), "<%s> synced\n", o.Name)
		})
	}

	t.Run(c.Context())

	if t.GetFault().HasError() {
		return t.GetFault()
	}

	if h := taskutils.HandlerFromContext(c.Context()); h != nil {
		h.WaitProgress()
	}

	i18n.Fprintf(c.OutOrStdout(), "Dir <%s> and <%s> synced.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return nil

}

func initSyncFlag() {
	SyncCommand.Flags().BoolVarP(&syncFlag.dryRun, "dry-run", "n", false,
		i18n.Sprintf(`show what would have been transferred`))
	SyncCommand.Flags().BoolVar(&syncFlag.existing, "existing", false,
		i18n.Sprintf(`skip creating new files in dest dirs`))
	SyncCommand.Flags().BoolVar(&syncFlag.ignoreExisting, "ignore-existing", false,
		i18n.Sprintf(`skip updating files in dest dirs, only copy those not exist`))
	SyncCommand.Flags().BoolVarP(&syncFlag.recursive, "recursive", "r", false,
		i18n.Sprintf(`recurse into sub directories`))
	SyncCommand.Flags().BoolVarP(&syncFlag.update, "update", "u", false,
		i18n.Sprintf(`skip files that are newer in dest dirs`))

}
