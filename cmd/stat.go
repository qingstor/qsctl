package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/action"
)

// StatCommand will handle stat command.
var StatCommand = &cobra.Command{
	Use:   "stat <remote-path>",
	Short: "stat a remote object",
	Long:  "qsctl stat show the detailed info of this object",
	Example: `Stat object:       qsctl stat qs://prefix/a
`,
	Args: cobra.ExactArgs(1),
	RunE: statRun,
}

func statRun(_ *cobra.Command, args []string) (err error) {
	err = action.Stat(args[0])
	if err != nil {
		panic(err)
	}
	return
}
