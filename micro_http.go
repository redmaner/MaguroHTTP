package main

import (
	"fmt"
	"net/http"
	"os"
)

const version = "v0.7"

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

	if mCfg.TLS && httpCheckTLS() {
		logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s with TLS", mCfg.Port))
		tlsc := httpCreateTLSConfig()
		ms := http.Server{
			Addr:      mCfg.Address + ":" + mCfg.Port,
			Handler:   mux,
			TLSConfig: tlsc,
		}
		err := ms.ListenAndServeTLS(mCfg.TLSCert, mCfg.TLSKey)
		logAction(logERROR, err)
	} else {
		logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s", mCfg.Port))
		http.ListenAndServe(mCfg.Address+":"+mCfg.Port, mux)
	}
}

// Function to show help
func showHelp() {
	fmt.Printf("MicroHTTP version %s\n\nUsage: microhttp </path/to/config.json>\n\n", version)
	os.Exit(1)
}
