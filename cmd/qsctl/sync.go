package main

import (
	"fmt"
	"path/filepath"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/task"
	"github.com/qingstor/qsctl/v2/utils"
)

var syncInput struct {
	IgnoreExisting bool
}

// SyncCommand will handle sync command.
var SyncCommand = &cobra.Command{
	Use:   "sync <source-path> <dest-path>",
	Short: i18n.Sprintf("sync between local directory and QS-Directory"),
	Long: i18n.Sprintf(`qsctl sync between local directory and QS-Directory. The first path argument
is the source directory and second the destination directory.

When a key(file) already exists in the destination directory, program will compare
the modified time of source file(key) and destination key(file). The destination
key(file) will be overwritten only if the source one newer than destination one.`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Sync local directory to QS-Directory: qsctl sync . qs://bucket-name"),
		i18n.Sprintf("Sync QS-Directory to local directory: qsctl sync qs://bucket-name/test/ test_local/"),
		i18n.Sprintf("Sync skip updating files that already exist on receiver: qsctl sync . qs://bucket-name --ignore-existing"),
	),
	Args: cobra.ExactArgs(2),
	RunE: syncRun,
}

func syncRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewBetweenStorageTask(10)
	err = utils.ParseBetweenStorageInput(rootTask, args[0], args[1])
	if err != nil {
		return
	}

	if rootTask.GetSourceType() != types.ObjectTypeDir || rootTask.GetDestinationType() != types.ObjectTypeDir {
		return fmt.Errorf("both source and destination should be directories")
	}

	if err = HandleSyncStorageBaseAndPath(rootTask); err != nil {
		return err
	}

	t := task.NewSync(rootTask)
	t.SetIgnoreExisting(syncInput.IgnoreExisting)
	t.Run()

	if t.GetFault().HasError() {
		return t.GetFault()
	}
	syncOutput(t)
	return nil

}

func initSyncFlag() {
	SyncCommand.Flags().BoolVar(&syncInput.IgnoreExisting, "ignore-existing", false,
		i18n.Sprintf(`skip creating new files in dest dirs, only copy newer by time`))
}

func syncOutput(t *task.SyncTask) {
	i18n.Printf("Dir <%s> and <%s> synced.\n", t.GetSourcePath(), t.GetDestinationPath())
}

// HandleSyncStorageBaseAndPath set work dir and path for sync cmd.
func HandleSyncStorageBaseAndPath(t *taskutils.BetweenStorageTask) error {
	srcPath, err := filepath.Abs(t.GetSourcePath())
	if err != nil {
		return err
	}
	if err := t.GetSourceStorage().Init(pairs.WithWorkDir(srcPath)); err != nil {
		return err
	}
	t.SetSourcePath("")

	dstPath, err := filepath.Abs(t.GetDestinationPath())
	if err != nil {
		return err
	}
	if err := t.GetDestinationStorage().Init(pairs.WithWorkDir(dstPath)); err != nil {
		return err
	}
	t.SetDestinationPath("")
	return nil
}
