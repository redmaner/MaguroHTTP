// Copyright 2018 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const version = "1.0 beta2"

// Main function
func main() {

	initLogger("MicroHTTP-", "stderr")

	args := os.Args
	if len(args) == 1 {
		showHelp()
	}

	// Handle arguments
	// To start MicroHTTP you need to define the path to the main configuration file
	if _, err := os.Stat(args[1]); err == nil {
		var mCfg microConfig
		loadConfigFromFile(args[1], &mCfg)
		if valid, err := validateConfig(args[1], &mCfg); valid && err == nil {
			startServer(&mCfg)
		} else {
			logAction(logERROR, err)
			os.Exit(1)
		}

	} else {
		showHelp()
	}
}

// Function to start Server
func startServer(mCfg *microConfig) {

	// Set micro struct
	m := micro{
		config: *mCfg,
		vhosts: make(map[string]microConfig),
	}

	// Setup metrics
	m.loadMetrics()

	// If virtual hosting is enabled, all the configurations of the vhosts are loaded
	if m.config.Serve.VirtualHosting {
		for k, v := range m.config.Serve.VirtualHosts {
			var cfg microConfig
			loadConfigFromFile(v, &cfg)
			if valid, err := validateConfigVhost(v, &cfg); !valid || err != nil {
				logAction(logERROR, err)
				os.Exit(1)
			}
			m.vhosts[k] = cfg
		}
	}

	// Configure logging
	debug = m.config.LogLevel
	initLogger("MicroHTTP-", m.config.LogOut)

	// Configure router
	m.configureRouter()

	// Get client
	m.client = serverClient(m.config.TLS)

	// If TLS is enabled the server will start in TLS
	if m.config.TLS.Enabled && httpCheckTLS(m.config.TLS) {
		logAction(logNONE, fmt.Errorf("MicroHTTP %s is listening on port %s with TLS", version, mCfg.Port))
		tlsc := httpCreateTLSConfig(m.config.TLS)
		ms := http.Server{
			Addr:      mCfg.Address + ":" + mCfg.Port,
			Handler:   m.router,
			TLSConfig: tlsc,
		}

		// This is meant to listen for signals. A signal will stop MicroHTTP
		done := make(chan bool)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			logAction(logNONE, fmt.Errorf("server is shutting down"))
			m.flushMDToFile(m.config.Metrics.Out)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ms.SetKeepAlivesEnabled(false)
			if err := ms.Shutdown(ctx); err != nil {
				logAction(logNONE, fmt.Errorf("could not gracefully shutdown the server: %v", err))
			}
			close(done)
		}()

		// Start the server
		err := ms.ListenAndServeTLS(mCfg.TLS.TLSCert, mCfg.TLS.TLSKey)
		if err != nil && err != http.ErrServerClosed {
			logAction(logERROR, fmt.Errorf("Starting server failed: %s", err))
			return
		}

		<-done
		logAction(logNONE, fmt.Errorf("MicroHTTP stopped"))

		// IF TLS is disabled the server is started without TLS
		// Never run non TLS servers in production!
	} else {
		logAction(logNONE, fmt.Errorf("MicroHTTP %s is listening on port %s", version, mCfg.Port))
		http.ListenAndServe(mCfg.Address+":"+mCfg.Port, m.router)
	}
}

// Function to show help
func showHelp() {
	fmt.Printf("MicroHTTP version %s\n\nUsage: microhttp </path/to/config.json>\n\n", version)
	os.Exit(1)
}
