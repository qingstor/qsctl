package main

import (
	"fmt"

	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var mbInput struct {
	Zone string
}

// MbCommand will handle make bucket command.
var MbCommand = &cobra.Command{
	Use:   "mb [qs://]<bucket-name>",
	Short: i18n.Sprintf("make a new bucket"),
	Long: i18n.Sprintf(`qsctl mb can make a new bucket with the specific name,

bucket name should follow DNS name rule with:
* length between 6 and 63;
* can only contains lowercase letters, numbers and hyphen -
* must start and end with lowercase letter or number
* must not be an available IP address
	`),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Make bucket: qsctl mb bucket-name"),
	),
	Args:    cobra.ExactArgs(1),
	RunE:    mbRun,
	PreRunE: validateMbFlag,
}

func mbRun(_ *cobra.Command, args []string) (err error) {
	rootTask := taskutils.NewAtServiceTask(10)
	err = utils.ParseAtServiceInput(rootTask)
	if err != nil {
		return
	}

	_, bucketName, _, err := utils.ParseQsPath(args[0])
	if err != nil {
		return
	}

	t := task.NewCreateStorage(rootTask)
	t.SetStorageName(bucketName)
	t.SetZone(mbInput.Zone)

	t.Run()
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	mbOutput(t)
	return
}

func mbOutput(t *task.CreateStorageTask) {
	i18n.Printf("Bucket <%s> created.\n", t.GetStorageName())
}

func initMbFlag() {
	MbCommand.Flags().StringVarP(&mbInput.Zone, constants.ZoneFlag, "z",
		"", i18n.Sprintf("in which zone to make the bucket (required)"))
}

func validateMbFlag(_ *cobra.Command, _ []string) error {
	// check zone flag (required)
	if mbInput.Zone == "" {
		// TODO: we need to return an error here.
		return fmt.Errorf("flag zone is required, but not found")
	}
	return nil
}
