package main

import (
	"os"

	"github.com/Xuanwo/go-locale"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/i18n"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	tag, err := locale.Detect()
	if err != nil {
		os.Exit(1)
	}
	i18n.Setup(tag)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
