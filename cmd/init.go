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
	formatFlag               = "format"
	humanReadableFlag        = "human-readable"
	longFormatFlag           = "long-format"
	maximumMemoryContentFlag = "maximum-memory-content"
	recursiveFlag            = "recursive"
	zoneFlag                 = "zone"
)

var (
	// register available flag vars here
	expectSize           string
	format               string
	humanReadable        bool
	longFormat           bool
	maximumMemoryContent string
	recursive            bool
	zone                 string
)

// initFlags will init all available flags.
func initFlags() {
	flagSet = pflag.NewFlagSet("", pflag.ExitOnError)

	flagSet.StringVar(&expectSize,
		expectSizeFlag,
		"",
		`expected size of the input file
accept: 100MB, 1.8G
(only used for input from stdin)`,
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

	flagSet.BoolVarP(&humanReadable,
		humanReadableFlag,
		"h",
		false,
		"print size by using unit suffixes: Byte, Kilobyte, Megabyte, Gigabyte, Terabyte and Petabyte,"+
			" in order to reduce the number of digits to three or less using base 2 for sizes",
	)

	flagSet.BoolVarP(&longFormat,
		longFormatFlag,
		"l",
		false,
		"list in long format and a total sum for all the file sizes is output on a line before the long listing",
	)

	flagSet.StringVar(&maximumMemoryContent,
		maximumMemoryContentFlag,
		"",
		"maximum content loaded in memory \n (only used for input from stdin)",
	)

	flagSet.BoolVarP(&recursive,
		recursiveFlag,
		"r",
		false,
		"execute recursively",
	)

	flagSet.StringVarP(&zone,
		zoneFlag,
		"z",
		"",
		"in which zone to do the operation",
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

	if format != "" {
		contexts.Format = format
	}

	contexts.HumanReadable = humanReadable
	contexts.LongFormat = longFormat

	if maximumMemoryContent != "" {
		contexts.MaximumMemoryContent, err = utils.ParseByteSize(maximumMemoryContent)
		if err != nil {
			return
		}
	}

	contexts.Recursive = recursive

	if zone != "" {
		contexts.Zone = zone
	}

	return nil
}

func init() {
	initFlags()

	// Flags for cp.
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(expectSizeFlag))
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(maximumMemoryContentFlag))

	// Flags for ls.
	LsCommand.Flags().AddFlag(flagSet.Lookup(humanReadableFlag))
	LsCommand.Flags().AddFlag(flagSet.Lookup(longFormatFlag))
	LsCommand.Flags().AddFlag(flagSet.Lookup(recursiveFlag))
	LsCommand.Flags().AddFlag(flagSet.Lookup(zoneFlag))

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
