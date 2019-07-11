package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

func init() {
	initCpCommandFlag()
	initMbCommandFlag()
	// initRmCommandFlag()
	initTeeCommandFlag()
}

// ParseFlagIntoContexts will executed before any commands to init the flags in contexts.
// And check required flags
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

	return checkRequiredFlags(cmd)
}

func checkRequiredFlags(cmd *cobra.Command) error {
	// iterate all required flags for current cmd
	for _, requireFlag := range cmdToFlagSet.GetRequiredFlags(cmd.Name()) {
		// get flag name
		flagName := requireFlag.GetName()
		log.Debugf("required flag <%s> value: <%v>", flagName, cmd.Flag(flagName).Value)
		// if CheckRequired not ok, return error
		if !requireFlag.CheckRequired(cmd.Flag(flagName).Value.String()) {
			log.Errorf("Flag <%s> is required", flagName)
			return constants.ErrorRequiredFlagsNotSet
		}
	}
	return nil
}
