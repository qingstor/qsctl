package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/action"
)

// CatCommand will handle cat command.
var CatCommand = &cobra.Command{
	Use:   "cat <remote-path>",
	Short: "cat a remote object to stdout",
	Long:  "qsctl cat can cat a remote object to stdout",
	Example: `Cat object:       qsctl cat qs://prefix/a
`,
	Args: cobra.ExactArgs(1),
	RunE: catRun,
}

func catRun(_ *cobra.Command, args []string) (err error) {
	err = action.Copy(args[0], "-")
	if err != nil {
		panic(err)
	}
	return
}
