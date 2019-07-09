package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/utils"
)

// flagSet stores all flags in itself
var flagSet *pflag.FlagSet

const (
	// all flags' input here
	expectSizeFlag           = "expect-size"
	maximumMemoryContentFlag = "maximum-memory-content"
	zoneFlag                 = "zone"
)

var (
	// register available flag vars here
	expectSize           string
	maximumMemoryContent string
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
(only used for input from stdin)`)

	flagSet.StringVar(&maximumMemoryContent,
		maximumMemoryContentFlag,
		"",
		"maximum content loaded in memory \n (only used for input from stdin)")

	flagSet.StringVarP(&zone,
		zoneFlag,
		"z",
		"",
		"In which zone to do the operation (required)",
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

	return nil
}

func init() {
	initFlags()

	// Flags for cp.
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(expectSizeFlag))
	CpCommand.PersistentFlags().AddFlag(flagSet.Lookup(maximumMemoryContentFlag))

	// Flags for tee.
	TeeCommand.PersistentFlags().AddFlag(flagSet.Lookup(expectSizeFlag))
	TeeCommand.PersistentFlags().AddFlag(flagSet.Lookup(maximumMemoryContentFlag))

	// Flags for mb.
	MbCommand.Flags().AddFlag(flagSet.Lookup(zoneFlag))
	// Mark flag "zone" required
	if err := MbCommand.MarkFlagRequired(zoneFlag); err != nil {
		log.Errorf("cmd mb: Mark flag zone required failed [%v]", err)
	}
}
