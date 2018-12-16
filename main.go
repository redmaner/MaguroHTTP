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

	"golang.org/x/crypto/acme/autocert"
)

const version = "r2"

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

	// Empty strings for TLS key and certificate
	var tlsCert string
	var tlsKey string

	// Set micro struct
	m := micro{
		config: *mCfg,
		vhosts: make(map[string]microConfig),
	}

	// Setup metrics
	m.loadMetrics()

	// If virtual hosting is enabled, all the configurations of the vhosts are loaded
	if m.config.Core.VirtualHosting {
		for k, v := range m.config.Core.VirtualHosts {
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
	debug = m.config.Core.LogLevel
	initLogger("MicroHTTP-", m.config.Core.LogOut)

	// Configure router
	m.configureRouter()

	// Get client
	m.client = serverClient(m.config.Core.TLS)

	// Setup the MicroHTTP server
	ms := http.Server{
		Addr:              mCfg.Core.Address + ":" + mCfg.Core.Port,
		Handler:           m.router,
		ReadTimeout:       4 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       12 * time.Second,
	}

	// If TLS is enabled the server will start in TLS
	if m.config.Core.TLS.Enabled && httpCheckTLS(m.config.Core.TLS) {
		logAction(logNONE, fmt.Errorf("MicroHTTP %s is listening on port %s with TLS", version, mCfg.Core.Port))
		tlsc := httpCreateTLSConfig(m.config.Core.TLS)

		// Handle autocert
		if m.config.Core.TLS.AutoCert.Enabled && len(m.config.Core.TLS.AutoCert.Certificates) > 0 {
			acm := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(m.config.Core.TLS.AutoCert.Certificates...),
				Cache:      autocert.DirCache(m.config.Core.TLS.AutoCert.CertDir),
			}
			tlsc.GetCertificate = acm.GetCertificate
		} else {
			tlsCert = m.config.Core.TLS.TLSCert
			tlsKey = m.config.Core.TLS.TLSKey
		}
		ms.TLSConfig = tlsc
	}

	// This is meant to listen for signals. A signal will stop MicroHTTP
	done := make(chan bool)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logAction(logNONE, fmt.Errorf("server is shutting down"))
		if m.config.Metrics.Enabled {
			m.flushMDToFile(m.config.Metrics.Out)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ms.SetKeepAlivesEnabled(false)
		if err := ms.Shutdown(ctx); err != nil {
			logAction(logNONE, fmt.Errorf("could not gracefully shutdown the server: %v", err))
		}
		close(done)
	}()

	// Start the server
	// If TLS is enabled in the configuration, TLS is used.
	// Otherwise a non TLS server is used.
	if m.config.Core.TLS.Enabled {
		logAction(logNONE, fmt.Errorf("MicroHTTP %s is listening on port %s with TLS", version, mCfg.Core.Port))

		// Start the server with TLS
		err := ms.ListenAndServeTLS(tlsCert, tlsKey)
		if err != nil && err != http.ErrServerClosed {
			logAction(logERROR, fmt.Errorf("Starting server failed: %s", err))
			return
		}
	} else {

		// start the server without TLS
		logAction(logNONE, fmt.Errorf("MicroHTTP %s is listening on port %s", version, mCfg.Core.Port))
		err := ms.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logAction(logERROR, fmt.Errorf("Starting server failed: %s", err))
			return
		}
	}

	<-done
	logAction(logNONE, fmt.Errorf("MicroHTTP stopped"))
}

// Function to show help
func showHelp() {
	fmt.Printf("MicroHTTP version %s\n\nUsage: microhttp </path/to/config.json>\n\n", version)
	os.Exit(1)
}
