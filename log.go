package main

import (
	"io/ioutil"
	"log"
	"os"
)

func getLogger(verbose, cui *bool) *log.Logger {
	if *cui || !*verbose {
		return log.New(ioutil.Discard, "", 0)
	}
	return log.New(os.Stdout, "", 0)
}
