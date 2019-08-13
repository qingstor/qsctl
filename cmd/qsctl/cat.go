package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qsctl/v2/action"
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
	// Package handler
	catHandler := &action.CopyHandler{}
	return catHandler.
		WithBench(bench).
		WithDest("-").
		WithSrc(args[0]).
		Copy()
}
