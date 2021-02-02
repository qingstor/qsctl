package main

import (
	"fmt"

	"github.com/qingstor/noah/task"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/shellutils"
	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

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
		i18n.Sprintf("Make bucket: qsctl mb bucket-name --zone=zone-name"),
	),
	Args:    cobra.ExactArgs(1),
	PreRunE: validateMbFlag,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true // handle runtime errors with i18n, do not show error
		if err := mbRun(cmd, args); err != nil {
			i18n.Fprintf(cmd.OutOrStderr(), "Execute %s command error: %s\n", "mb", err.Error())
			return err
		}
		return nil
	},
}

func mbRun(c *cobra.Command, args []string) (err error) {
	silenceUsage(c) // silence usage when handled error returns
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
	t.SetZone(globalFlag.zone)

	t.Run(c.Context())
	if t.GetFault().HasError() {
		return t.GetFault()
	}

	i18n.Fprintf(c.OutOrStdout(), "Bucket <%s> created.\n", t.GetStorageName())
	return
}

func initMbFlag() {}

func validateMbFlag(_ *cobra.Command, _ []string) error {
	// check zone flag (required)
	if globalFlag.zone == "" {
		return fmt.Errorf(i18n.Sprintf("flag zone is required, but not found"))
	}
	return nil
}

type mbShellHandler struct {
	bucketName string
}

// preRunE do bucket name parse before mb run in shell
func (h *mbShellHandler) preRunE(args []string) error {
	err := MbCommand.Flags().Parse(args)
	if err != nil {
		return err
	}
	if len(MbCommand.Flags().Args()) <= 0 {
		return fmt.Errorf(i18n.Sprintf("Error: at least one arg is needed for %s", "mb"))
	}
	_, bucketName, _, err := utils.ParseQsPath(MbCommand.Flags().Args()[0])
	if err != nil {
		return err
	}
	h.bucketName = bucketName
	return nil
}

// postRun add bucket name into cache list if no error while run
func (h mbShellHandler) postRun(err error) {
	if err == nil {
		shellutils.AddBucketIntoList(h.bucketName)
	}
}
