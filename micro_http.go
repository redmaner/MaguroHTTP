package main

import (
	"fmt"
	"net/http"
	"os"
)

const version = "v0.5"

var mCfg = microConfig{}

// Main function
func main() {

	initLogger("MicroHTTP-")

	args := os.Args
	if len(args) == 1 {
		showHelp()
	}

	if _, err := os.Stat(args[1]); err == nil {
		loadConfigFromFile(args[1], &mCfg)
		startServer()
	} else {
		showHelp()
	}
}

// Function to start Server
func startServer() {

	mux := http.NewServeMux()

	// If ProxyMode is enabled, use proxy handler
	if mCfg.Proxy.Enabled {
		mux.HandleFunc("/", handleProxy)
	} else {
		mux.HandleFunc("/", handleHTTP)
	}

	logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s", mCfg.Port))
	http.ListenAndServe(mCfg.Address+":"+mCfg.Port, mux)

}

// Function to show help
func showHelp() {
	fmt.Printf("MicroHTTP version %s\n\nUsage: microhttp </path/to/config.json>\n\n", version)
	os.Exit(1)
}
