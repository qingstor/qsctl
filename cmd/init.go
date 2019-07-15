package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

// flagSet stores all flags in itself
var flagSet *pflag.FlagSet

const (
	// all flags' input here
	expectSizeFlag           = "expect-size"
	maximumMemoryContentFlag = "maximum-memory-content"
	zoneFlag                 = "zone"
	formatFlag               = "format"
)

var (
	// register available flag vars here
	expectSize           string
	maximumMemoryContent string
	zone                 string
	format               string
)

// initFlags will init all available flags.
func initFlags() {
	flagSet = pflag.NewFlagSet("", pflag.ExitOnError)

	flagSet.StringVar(&expectSize,
		expectSizeFlag,
		"",
		`expected size of the input file
accept: 100MB, 1.8G
(only used for input from stdin)`)

	flagSet.StringVar(&maximumMemoryContent,
		maximumMemoryContentFlag,
		"",
		"maximum content loaded in memory \n (only used for input from stdin)")

	flagSet.StringVarP(&zone,
		zoneFlag,
		"z",
		"",
		"In which zone to do the operation",
	)

	flagSet.StringVarP(&format,
		formatFlag,
		"",
		"",
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

// ParseFlagIntoContexts will executed before any commands to init the flags in contexts.
func ParseFlagIntoContexts(cmd *cobra.Command, args []string) (err error) {
	if expectSize != "" {
		contexts.ExpectSize, err = utils.ParseByteSize(expectSize)
		if err != nil {
			return
		}
	}

	if maximumMemoryContent != "" {
		contexts.MaximumMemoryContent, err = utils.ParseByteSize(maximumMemoryContent)
		if err != nil {
			return
		}
	}

	if zone != "" {
		contexts.Zone = zone
	}

	if format != "" {
		contexts.Format = format
	}

	return nil
}

func init() {
	initFlags()

	// Flags for cp.
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(expectSizeFlag))
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(maximumMemoryContentFlag))

	// Flags for mb.
	MbCommand.Flags().AddFlag(flagSet.Lookup(zoneFlag))

	// Flags for rm.
	RmCommand.Flags().AddFlag(flagSet.Lookup(zoneFlag))

	// Flags for stat.
	StatCommand.Flags().AddFlag(flagSet.Lookup(formatFlag))

	// Flags for tee.
	TeeCommand.PersistentFlags().AddFlag(flagSet.Lookup(expectSizeFlag))
	TeeCommand.PersistentFlags().AddFlag(flagSet.Lookup(maximumMemoryContentFlag))
}
