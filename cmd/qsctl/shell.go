package main

import (
	"bufio"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"

	"github.com/cosiner/argv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

const shellName = "shell"

// ShellCommand will handle qsctl shell command.
var ShellCommand = &cobra.Command{
	Use:   shellName,
	Short: i18n.Sprintf("start an interactive shell of qsctl"),
	Long:  i18n.Sprintf("qsctl shell can execute command interactively, input exit to quit"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Start shell: qsctl shell"),
	),
	Args: cobra.ExactArgs(0),
	RunE: shellRun,
	PreRun: func(_ *cobra.Command, _ []string) {
		log.SetOutput(os.Stdout)
	},
}

func shellRun(_ *cobra.Command, _ []string) (err error) {
	sig, cmdChan, inputReadyChan := make(chan os.Signal), make(chan string), make(chan struct{})
	signal.Notify(sig, os.Interrupt)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			// try to read from stdin
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				i18n.Printf("\nread string failed: %s\n", err)
				fmt.Printf("%s> ", constants.Name)
				continue
			}
			cmdString = strings.TrimSuffix(cmdString, "\n")
			// put input cmd out
			cmdChan <- cmdString
			// wait until next input
			<-inputReadyChan
		}
	}()

	for {
		fmt.Printf("%s> ", constants.Name)
		select {
		case <-sig:
			// if interrupt signal caught, print tip and continue
			i18n.Printf("\ninterrupted by input\n")
			continue
		case input := <-cmdChan:
			// exit if user input exit
			if isExit(input) {
				os.Exit(0)
			}
			// if input blank, continue to wait for next input
			if input == "" {
				break
			}

			args, err := getArgs(input)
			if err != nil {
				i18n.Printf("get args failed: %s\n", err)
				break
			}

			if err := handleArgs(args); err != nil {
				i18n.Printf("execute %v failed: %s\n", args, err)
				break
			}
		}
		inputReadyChan <- struct{}{}
	}
}

func isExit(input string) bool {
	return input == constants.ExitCmd
}

func getArgs(input string) ([]string, error) {
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

func handleArgs(args []string) error {
	fmt.Println(args)

	if err := checkShellCmd(args); err != nil {
		return err
	}

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
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
			fmt.Print("input bucket name to confirm: ")
			_, bucketName, _, _ := utils.ParseQsPath(args[1])
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			input = strings.TrimSuffix(input, "\n")
			if bucketName != input {
				return errors.New("not confirmed")
			}
			fmt.Print("get input", input)
		}
	case RmCommand.Name():
		err := RmCommand.Flags().Parse(args[1:])
		if err != nil {
			return err
		}
		_, _, key, _ := utils.ParseQsPath(args[1])
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSuffix(input, "\n")
		if key != input {
			return errors.New("not confirmed")
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
