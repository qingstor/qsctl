package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Xuanwo/storage/types"
	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var syncInput struct {
	DryRun         bool
	Existing       bool
	IgnoreExisting bool
	Recursive      bool
	Update         bool
}

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
	RunE: syncRun,
}

func syncRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
	rootTask := taskutils.NewBetweenStorageTask(10)
	srcWorkDir, dstWorkDir, err := utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() != types.ObjectTypeDir || rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf("both source and destination should be directories")
	}

	if syncInput.Existing && syncInput.IgnoreExisting {
		return fmt.Errorf("both --existing and --ignore-existing are set, no files would be synced")
	}

	go func() {
		taskutils.StartProgress(time.Second)
	}()
	defer taskutils.FinishProgress()

	t := task.NewSync(rootTask)
	t.SetDryRun(syncInput.DryRun)
	t.SetExisting(syncInput.Existing)
	t.SetIgnoreExisting(syncInput.IgnoreExisting)
	t.SetRecursive(syncInput.Recursive)
	t.SetUpdate(syncInput.Update)
	if syncInput.DryRun {
		t.SetDryRunFunc(func(o *types.Object) {
			fmt.Println(o.Name)
		})
	} else {
		t.SetDryRunFunc(nil)
		t.SetHandleObjCallback(func(o *types.Object) {
			fmt.Println(i18n.Sprintf("<%s> synced", o.Name))
		})
	}

	t.Run()

	if t.GetFault().HasError() {
		return t.GetFault()
	}

	taskutils.WaitProgress()
	i18n.Printf("Dir <%s> and <%s> synced.\n",
		filepath.Join(srcWorkDir, t.GetSourcePath()), filepath.Join(dstWorkDir, t.GetDestinationPath()))
	return nil

}

func initSyncFlag() {
	SyncCommand.Flags().BoolVarP(&syncInput.DryRun, "dry-run", "n", false,
		i18n.Sprintf(`show what would have been transferred`))
	SyncCommand.Flags().BoolVar(&syncInput.Existing, "existing", false,
		i18n.Sprintf(`skip creating new files in dest dirs`))
	SyncCommand.Flags().BoolVar(&syncInput.IgnoreExisting, "ignore-existing", false,
		i18n.Sprintf(`skip updating files in dest dirs, only copy those not exist`))
	SyncCommand.Flags().BoolVarP(&syncInput.Recursive, "recursive", "r", false,
		i18n.Sprintf(`recurse into sub directories`))
	SyncCommand.Flags().BoolVarP(&syncInput.Update, "update", "u", false,
		i18n.Sprintf(`skip files that are newer in dest dirs`))

}
