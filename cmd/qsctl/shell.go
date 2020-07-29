package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/cosiner/argv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/shellutils"
	cutils "github.com/qingstor/qsctl/v2/cmd/utils"
	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

// ShellCommand will handle qsctl shell command.
var ShellCommand = &cobra.Command{
	Use:   "shell",
	Short: i18n.Sprintf("start an interactive shell of qsctl"),
	Long:  i18n.Sprintf("qsctl shell can execute command interactively, input exit to quit"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Start shell: qsctl shell"),
	),
	Args: cobra.NoArgs,
	Run:  shellRun,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		if !cutils.IsInteractiveEnable() {
			return fmt.Errorf(i18n.Sprintf("not interactive shell, cannot call shell"))
		}
		log.SetOutput(os.Stdout)
		return nil
	},
}

func executor(t string) {
	if t == "" {
		return
	}
	args, err := parseArgs(t)
	if err != nil {
		i18n.Printf("get args failed: %s\n", err)
		return
	}

	subCmdName := args[0]
	sh, err := shellHandlerFactory(subCmdName)
	if err != nil {
		i18n.Printf("%s\n", err)
		return
	}
	if err = sh.preRunE(args[1:]); err != nil {
		i18n.Printf("%s\n", err)
		return
	}

	rootCmd.SetArgs(args)
	rootCmd.SetOut(os.Stdout)
	if err = rootCmd.ExecuteContext(context.Background()); err != nil {
		return
	}

	sh.postRun(err)
	return
}

// completeFunc handle auto-completion by suggests while input
func completeFunc(d prompt.Document) (s []prompt.Suggest) {
	// if first word inputting, try to suggest commands
	if !strings.Contains(d.CurrentLineBeforeCursor(), " ") {
		s := getCmdSuggests()
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	curWord := d.GetWordBeforeCursor()
	// if not start to input, do not suggest
	if curWord == "" {
		return
	}

	// if start to input flags, which starts with "-", try to suggest flags
	if strings.HasPrefix(curWord, "-") {
		s = getFlagSuggests(d)
		return prompt.FilterHasPrefix(s, curWord, true)
	}

	// if not a qingstor path, try to suggest local files
	if utils.IsQsPath(curWord) {
		s = getBucketSuggests()
		return prompt.FilterHasPrefix(s, strings.TrimPrefix(curWord, "qs://"), true)
	}

	s = getFileSuggests(d)
	return s
}

func shellRun(_ *cobra.Command, _ []string) {
	shellutils.InitBucketList()

	p := prompt.New(executor, completeFunc,
		prompt.OptionPrefix(constants.Name+"> "),
		prompt.OptionTitle(constants.Name),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)

	p.Run()
	return
}

// shellHandler contains preRunE and postRun methods
type shellHandler interface {
	preRunE(args []string) error
	postRun(err error)
}

// shellHandlerFactory create shellHandler by factory pattern
func shellHandlerFactory(cmd string) (shellHandler, error) {
	switch cmd {
	case MbCommand.Name():
		return &mbShellHandler{}, nil
	case RbCommand.Name():
		return &rbShellHandler{}, nil
	case RmCommand.Name():
		return rmShellHandler{}, nil
	// remove cat and tee command support
	case CpCommand.Name(), LsCommand.Name(), MvCommand.Name(),
		PresignCommand.Name(), StatCommand.Name(), SyncCommand.Name():
		return blankShellHandler{}, nil
	default:
		return nil, constants.ErrCmdNotSupport
	}
}

// parseArgs parse input string into string slice like os.Args
func parseArgs(input string) ([]string, error) {
	// handle args as os.Args
	args, err := argv.Argv(input, func(bq string) (string, error) {
		return bq, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(args) > 1 {
		log.Warnf(i18n.Sprint("pipe not supported in shell, input after %v would be abandoned"), args[0])
	}
	return args[0], nil
}

// getCmdSuggests add all sub commands into suggest list
func getCmdSuggests() (s []prompt.Suggest) {
	for _, command := range shellSubCommands() {
		s = append(s,
			prompt.Suggest{Text: command.Name(), Description: command.Short},
		)
	}
	return s
}

// getFlagSuggests returns flag suggest list
func getFlagSuggests(d prompt.Document) (s []prompt.Suggest) {
	// get the specific sub commands' flags into suggest
	for _, command := range shellSubCommands() {
		if input := d.TextBeforeCursor(); strings.HasPrefix(input, command.Name()) {
			command.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				s = append(s,
					prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage},
				)
			})
			break
		}
	}
	// add global flags to suggest
	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		s = append(s,
			prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage},
		)
	})
	return s
}

// getFileSuggests returns local files suggest list
func getFileSuggests(d prompt.Document) []prompt.Suggest {
	var fileCompleter = completer.FilePathCompleter{
		IgnoreCase: true,
		Filter: func(fi os.FileInfo) bool {
			return true
		},
	}
	return fileCompleter.Complete(d)
}

// getBucketSuggests return suggest list which contains all buckets
func getBucketSuggests() (s []prompt.Suggest) {
	for _, b := range shellutils.GetBucketList() {
		s = append(s, prompt.Suggest{Text: b})
	}
	return s
}

// shellSubCommands return all available sub command in shell
func shellSubCommands() []*cobra.Command {
	return []*cobra.Command{
		CpCommand,
		LsCommand,
		MbCommand,
		MvCommand,
		PresignCommand,
		RbCommand,
		RmCommand,
		StatCommand,
		SyncCommand,
	}
}

// blankShellHandler implements shellHandler and do nothing
type blankShellHandler struct{}

func (b blankShellHandler) preRunE(_ []string) error {
	return nil
}

func (b blankShellHandler) postRun(_ error) {
	return
}

// noSuggests is the func that return empty prompt.Suggest
func noSuggests(_ prompt.Document) []prompt.Suggest {
	return nil
}
