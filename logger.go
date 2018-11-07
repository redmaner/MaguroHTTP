package main

import (
	"log"
	"os"
)

const (
	logNONE    = 0
	logNET     = 1
	logERROR   = 2
	logDEBUG   = 3
	logTRACE   = 4
	logVERBOSE = 5
)

var debug = logTRACE
var logger *log.Logger

// Function to initialize logger
func initLogger(s string) {
	logger = log.New(os.Stdout, s, log.Ldate|log.Ltime)
}

// General function to write errors and messages to log
func logAction(l int, err error) {
	if err == nil {
		return
	}

	switch {
	case debug >= logNET && l == logNET:
		logger.Println("ERROR:", err)
	case debug >= logERROR && l == logERROR:
		logger.Println("ERROR:", err)
	case debug >= logDEBUG && l == logDEBUG:
		logger.Println("DEBUG:", err)
	case debug >= logTRACE && l == logTRACE:
		logger.Println("TRACE:", err)
	case debug >= logVERBOSE && l == logVERBOSE:
		logger.Println("VERBOSE:", err)
	}
}
