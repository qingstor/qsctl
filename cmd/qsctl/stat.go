package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/utils"
)

// StatCommand will handle stat command.
var StatCommand = &cobra.Command{
	Use:   "stat qs://<bucket_name>/<object_key>",
	Short: "stat a remote object",
	Long:  "qsctl stat show the detailed info of this object",
	Example: utils.AlignPrintWithColon(
		"Stat object: qsctl stat qs://prefix/a",
	),
	Args: cobra.ExactArgs(1),
	RunE: statRun,
}

func statRun(_ *cobra.Command, args []string) (err error) {
	sh := action.StatHandler{}
	return sh.WithZone(zone).WithFormat(format).WithRemote(args[0]).Stat()
}

func initStatFlag() {
	StatCommand.Flags().StringVar(&format, constants.FormatFlag, "",
		`use the specified FORMAT instead of the default;
output a newline after each use of FORMAT

The valid format sequences for files:

  %F   file type
  %h   content md5 of the file
  %n   file name
  %s   total size, in bytes
  %y   time of last data modification, human-readable
  %Y   time of last data modification, seconds since Epoch
	`,
	)
}
