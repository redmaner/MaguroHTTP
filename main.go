package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const version = "v0.11"

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
		if valid, err := validateConfig(args[1], &mCfg); valid && err == nil {
			go startServer()
			sig := make(chan os.Signal)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			for {
				select {
				case s := <-sig:
					logAction(logDEBUG, fmt.Errorf("Signal (%d) received, stopping\n", s))
				}
			}
		} else {
			logAction(logERROR, err)
			os.Exit(1)
		}

	} else {
		showHelp()
	}
}

// Function to start Server
func startServer() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpServe)

	if mCfg.TLS && httpCheckTLS() {
		logAction(logDEBUG, fmt.Errorf("MicroHTTP is listening on port %s with TLS", mCfg.Port))
		tlsc := httpCreateTLSConfig()
		ms := http.Server{
			Addr:      mCfg.Address + ":" + mCfg.Port,
			Handler:   mux,
			TLSConfig: tlsc,
		}

		err := ms.ListenAndServeTLS(mCfg.TLSCert, mCfg.TLSKey)
		if err != nil {
			logAction(logERROR, fmt.Errorf("Starting server failed: %s", err))
			return
		}

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
