package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/cosiner/argv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/qingstor/qsctl/v2/constants"
	"github.com/qingstor/qsctl/v2/pkg/i18n"
	"github.com/qingstor/qsctl/v2/utils"
)

var yesRx = regexp.MustCompile("^(?i:y(?:es)?)$")

// ShellCommand will handle qsctl shell command.
var ShellCommand = &cobra.Command{
	Use:   "shell",
	Short: i18n.Sprintf("start an interactive shell of qsctl"),
	Long:  i18n.Sprintf("qsctl shell can execute command interactively, input exit to quit"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprintf("Start shell: qsctl shell"),
	),
	Args: cobra.NoArgs,
	RunE: shellRun,
	PreRun: func(_ *cobra.Command, _ []string) {
		log.SetOutput(os.Stdout)
	},
}

func shellRun(_ *cobra.Command, _ []string) (err error) {
	// sigFin used to pass signal into sigChan when shellHandler execute without context cancel
	var sigFin syscall.Signal

	sh := shellHandler{
		reader:         os.Stdin,
		sigChan:        make(chan os.Signal),
		inputChan:      make(chan string),
		readyInputChan: make(chan struct{}),
	}

	signal.Notify(sh.sigChan, os.Interrupt)
	go sh.initReader()

	for {
		sh.ReadyToInput()
		input := sh.loopInput(fmt.Sprintf("%s> ", constants.Name))
		// exit if user input exit
		if isExit(input) {
			os.Exit(0)
		}
		// if input blank, continue to wait for next input
		if input == "" {
			continue
		}

		args, err := parseArgs(input)
		if err != nil {
			i18n.Printf("get args failed: %s\n", err)
			continue
		}

		if err = sh.checkShellCmd(args); err != nil {
			i18n.Printf("check command failed: %s\n", err)
			continue
		}

		// execute command with args if check passed
		// get new ctx for each execution
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			sh.waitSignal()
			cancel()
		}()

		err = sh.ExecuteWithContext(ctx, args)
		select {
		case <-ctx.Done():
			break
		default:
			// if not canceled, we need to trigger signal to avoid goroutine leak, which waits signal to cancel ctx
			sh.triggerSignal(sigFin)
			// cancel func is idempotent
			cancel()
		}
		if err != nil {
			i18n.Printf("execute %v failed: %s\n", args, err)
			continue
		}
	}
}

func isExit(input string) bool {
	return input == constants.ExitCmd
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

// shellHandler is the struct to handle shellHandler command's execution
type shellHandler struct {
	// reader is where to get input
	reader io.Reader
	// sigChan get signal for some response
	sigChan chan os.Signal
	// inputChan used to transfer input
	inputChan chan string
	// readyInputChan used to mark ready to get input
	readyInputChan chan struct{}
}

func (s shellHandler) initReader() {
	reader := bufio.NewReader(s.reader)
	for {
		// wait until ready to input
		<-s.readyInputChan
		// try to read from stdin
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			i18n.Printf("\nread string failed: %s\n", err)
			fmt.Printf("%s> ", constants.Name)
			continue
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		// put input cmd out
		s.inputChan <- cmdString
	}
}

func (s shellHandler) ReadyToInput() {
	s.readyInputChan <- struct{}{}
}

func (s shellHandler) Execute(args []string) error {
	return s.ExecuteWithContext(context.Background(), args)
}

func (s shellHandler) ExecuteWithContext(ctx context.Context, args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.ExecuteContext(ctx)
}

func (s shellHandler) checkShellCmd(args []string) error {
	cmdName := args[0]
	switch cmdName {
	case RbCommand.Name():
		err := RbCommand.Flags().Parse(args[1:])
		if err != nil {
			return err
		}
		if rbInput.force {
			_, bucketName, _, _ := utils.ParseQsPath(args[1])
			s.ReadyToInput()
			input := s.loopInput(fmt.Sprintf("input bucket name <%s> to confirm: ", bucketName))
			if bucketName != input {
				return errors.New("not confirmed")
			}
		}
	case RmCommand.Name():
		_, _, key, _ := utils.ParseQsPath(args[1])
		var confirm bool
		s.ReadyToInput()
		input := s.loopInput(fmt.Sprintf("confirm to remove <%s>?: (y/N) ", key))
		confirm = yesRx.MatchString(input)
		if !confirm {
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

func (s shellHandler) loopInput(tip string) (res string) {
	for {
		fmt.Print(tip)
		select {
		case <-s.sigChan:
			// if interrupt signal caught, print tip and continue
			i18n.Printf("\ninterrupted by input\n")
			continue
		case res = <-s.inputChan:
			return
		}
	}
}

func (s shellHandler) waitSignal() {
	<-s.sigChan
}

func (s shellHandler) triggerSignal(sig os.Signal) {
	s.sigChan <- sig
}
