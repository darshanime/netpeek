package main

import (
	"io/ioutil"
	"log"
	"os"
)

func getLogger(quiet, cui *bool) *log.Logger {
	if *cui {
		return log.New(ioutil.Discard, "", 0)
	}
	if *quiet {
		return log.New(os.Stderr, "", 0)
	}
	return log.New(os.Stdout, "", 0)
}
