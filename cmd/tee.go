package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/action"
)

// TeeCommand will handle tee command.
var TeeCommand = &cobra.Command{
	Use:   "tee <remote-path>",
	Short: "tee a remote object from stdin",
	Long: `qsctl tee can tee a remote object from stdin.

NOTICE: qsctl will not tee the content to stdout like linux tee command does.
`,
	Example: `Tee object:       qsctl tee qs://prefix/a
`,
	Args: cobra.ExactArgs(1),
	RunE: teeRun,
}

func teeRun(_ *cobra.Command, args []string) (err error) {
	err = action.Copy("-", args[0])
	if err != nil {
		panic(err)
	}
	return
}
