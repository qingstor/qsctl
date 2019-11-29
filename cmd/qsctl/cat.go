package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qsctl/v2/pkg/i18n"

	"github.com/yunify/qsctl/v2/utils"
)

// CatCommand will handle cat command.
var CatCommand = &cobra.Command{
	Use:   i18n.Sprint("cat qs://<bucket_name>/<object_key>"),
	Short: i18n.Sprint("cat a remote object to stdout"),
	Long:  i18n.Sprint("qsctl cat can cat a remote object to stdout"),
	Example: utils.AlignPrintWithColon(
		i18n.Sprint("Cat object: qsctl cat qs://prefix/a"),
	),
	Args: cobra.ExactArgs(1),
	RunE: catRun,
}

func catRun(_ *cobra.Command, args []string) (err error) {
	// Package handler
	return nil
}
