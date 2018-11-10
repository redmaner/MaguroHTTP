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
	case debug >= logNONE && l == logNONE:
		logger.Println(err)
	case debug >= logNET && l == logNET:
		logger.Println("NET:", err)
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

func logNetwork(sc int, r *http.Request) {
	logAction(logNET, fmt.Errorf("%d request=%s %s%s%s IP=%s User-Agent=%s", sc, r.Method, r.Host, r.URL.Path, r.URL.RawQuery, r.RemoteAddr, r.Header.Get("User-Agent")))
}
