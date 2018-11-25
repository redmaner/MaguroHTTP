package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Constants for logging levels
const (
	logNONE    = 0
	logNET     = 1
	logERROR   = 2
	logDEBUG   = 3
)

var debug = logERROR
var logger *log.Logger

// Function to initialize logger
func initLogger(s, o string) {
	switch o {
	case "stdout":
		logger = log.New(os.Stdout, s, log.Ldate|log.Ltime)
	case "stderr":
		logger = log.New(os.Stderr, s, log.Ldate|log.Ltime)
	default:
		var logFile *os.File
		var err error

		if _, err = os.Stat(o); err != nil {
			logFile, err = os.Create(o)
			if err != nil {
				fmt.Printf("An error occured creating %s\n", o)
				os.Exit(1)
			}
		}

		logFile, err = os.OpenFile(o, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("An error occured opening %s\n", o)
			os.Exit(1)
		}

		logger = log.New(logFile, s, log.Ldate|log.Ltime)
	}
}

// General function to write errors and messages to log
func logAction(l int, err error) {
	if err == nil {
		return
	}

	switch {
	case debug >= logNONE && l == logNONE:
		logger.Println(err)
	case debug >= logNET && l == logNET:
		logger.Println("NET:", err)
	case debug >= logERROR && l == logERROR:
		logger.Println("ERROR:", err)
	case debug >= logDEBUG && l == logDEBUG:
		logger.Println("DEBUG:", err)
	}
}

func logNetwork(sc int, r *http.Request) {
	logAction(logNET, fmt.Errorf("%d request=%s %s%s%s IP=%s User-Agent=%s", sc, r.Method, r.Host, r.URL.Path, r.URL.RawQuery, r.RemoteAddr, r.Header.Get("User-Agent")))
}
