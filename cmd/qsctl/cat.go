package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/utils"
)

// CatCommand will handle cat command.
var CatCommand = &cobra.Command{
	Use:   "cat qs://<bucket_name>/<object_key>",
	Short: "cat a remote object to stdout",
	Long:  "qsctl cat can cat a remote object to stdout",
	Example: utils.AlignPrintWithColon(
		"Cat object: qsctl cat qs://prefix/a",
	),
	Args: cobra.ExactArgs(1),
	RunE: catRun,
}

func catRun(_ *cobra.Command, args []string) (err error) {
	// Package context
	ctx = contexts.SetContext(ctx, "src", args[0])
	ctx = contexts.SetContext(ctx, "dest", "-")
	return action.Copy(ctx)
}
