package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/action"
	"github.com/yunify/qsctl/utils"
)

// StatCommand will handle stat command.
var StatCommand = &cobra.Command{
	Use:   "stat <remote-path>",
	Short: "stat a remote object",
	Long:  "qsctl stat show the detailed info of this object",
	Example: utils.AlignPrintWithColon(
		"Stat object: qsctl stat qs://prefix/a",
	),
	Args: cobra.ExactArgs(1),
	RunE: statRun,
}

func statRun(_ *cobra.Command, args []string) (err error) {
	return action.Stat(args[0])
}
