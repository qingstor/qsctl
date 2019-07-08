package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/yunify/qsctl/contexts"
	"github.com/yunify/qsctl/utils"
)

const (
	expectSizeFlag           = "expect-size"
	maximumMemoryContentFlag = "maximum-memory-content"
)

var (
	expectSize           string
	maximumMemoryContent string
)

// initFlags will init all available flags.
func initFlags() {
	pflag.StringVar(&expectSize,
		expectSizeFlag,
		"",
		`expected size of the input file
accept: 100MB, 1.8G
(only used for input from stdin)`)

	pflag.StringVar(&maximumMemoryContent,
		maximumMemoryContentFlag,
		"",
		"maximum content loaded in memory \n (only used for input from stdin)")
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

	return nil
}

func init() {
	initFlags()

	// Flags for cp.
	CpCommand.PersistentFlags().AddFlag(pflag.Lookup(expectSizeFlag))
	CpCommand.PersistentFlags().AddFlag(pflag.Lookup(maximumMemoryContentFlag))

	// Flags for tee.
	TeeCommand.PersistentFlags().AddFlag(pflag.Lookup(expectSizeFlag))
	TeeCommand.PersistentFlags().AddFlag(pflag.Lookup(maximumMemoryContentFlag))
}
