package main

import (
	"log"
	"os"
)

func getLogger(quiet, cui *bool) *log.Logger {
	if *quiet || *cui {
		return log.New(os.Stderr, "", 0)
	}
	return log.New(os.Stdout, "", 0)
}
