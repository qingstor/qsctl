package main

import (
	"context"
	"errors"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/cosiner/argv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
	PreRun: func(_ *cobra.Command, _ []string) {
		log.SetOutput(os.Stdout)
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

	if err = checkShellCmd(args); err != nil {
		i18n.Printf("check command failed: %s\n", err)
		return
	}

	rootCmd.SetArgs(args)
	if err = rootCmd.ExecuteContext(context.Background()); err != nil {
		i18n.Printf("execute command failed: %s\n", err)
		return
	}
	return
}

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
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	// if not a qingstor path, try to suggest local files
	if !utils.IsQsPath(curWord) {
		s = getFileSuggests(d)
		return s
	}
	return
}

func shellRun(_ *cobra.Command, _ []string) {
	p := prompt.New(executor, completeFunc,
		prompt.OptionPrefix(constants.Name+"> "),
		prompt.OptionTitle(constants.Name),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)

	p.Run()
	return
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

func isTransferCmd(c string) bool {
	switch c {
	case CpCommand.Name(), SyncCommand.Name(), MvCommand.Name():
		return true
	default:
		return false
	}
}

func checkShellCmd(args []string) error {
	cmdName := args[0]
	switch cmdName {
	case RbCommand.Name():
		err := RbCommand.Flags().Parse(args[1:])
		if err != nil {
			return err
		}
		if rbInput.force {
			_, bucketName, _, _ := utils.ParseQsPath(RbCommand.Flags().Args()[0])
			for {
				confirm, err := utils.DoubleCheckString(bucketName, i18n.Sprintf("input bucket name <%s> to confirm:", bucketName))
				if err != nil {
					if errors.Is(err, terminal.InterruptErr) {
						continue
					}
					return err
				}
				if !confirm {
					return errors.New("not confirmed")
				}
			}

		}
	case RmCommand.Name():
		_, _, key, _ := utils.ParseQsPath(args[1])
		for {
			confirm, err := utils.CheckConfirm(i18n.Sprintf("confirm to remove <%s>?", key))
			if err != nil {
				if errors.Is(err, terminal.InterruptErr) {
					continue
				}
				return err
			}
			if !confirm {
				return errors.New("not confirmed")
			}
		}

	case CatCommand.Name(), CpCommand.Name(), LsCommand.Name(),
		MbCommand.Name(), MvCommand.Name(), PresignCommand.Name(),
		StatCommand.Name(), SyncCommand.Name(), TeeCommand.Name():
		break
	default:
		return constants.ErrCmdNotSupport
	}

	return nil
}

func getCmdSuggests() (s []prompt.Suggest) {
	for _, command := range basicCommands() {
		s = append(s,
			prompt.Suggest{Text: command.Name(), Description: command.Short},
		)
	}
	return s
}

func getFlagSuggests(d prompt.Document) (s []prompt.Suggest) {
	for _, command := range basicCommands() {
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

func getFileSuggests(d prompt.Document) []prompt.Suggest {
	var fileCompleter = completer.FilePathCompleter{
		IgnoreCase: true,
		Filter: func(fi os.FileInfo) bool {
			return true
		},
	}
	return fileCompleter.Complete(d)
}
