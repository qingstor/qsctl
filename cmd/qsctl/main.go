package main

import (
	"context"
	"os"

	"github.com/qingstor/log"

	"github.com/qingstor/qsctl/v2/internal/pkg/ilog"
)

func main() {
	ctx := context.Background()

	logger := ilog.InitLoggerWithDebug(globalFlag.debug)
	ctx = log.ContextWithLogger(ctx, logger)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
