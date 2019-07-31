package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
