// Copyright 2018-2019 Jake van der Putten.
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

package micro

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redmaner/MicroHTTP/debug"
	"golang.org/x/crypto/acme/autocert"
)

// Serve is used to serve a MicroHTTP instance
func (s *Server) Serve() {

	// TLS
	var tlsCert string
	var tlsKey string

	// Define server struct
	server := http.Server{
		Addr:              s.Cfg.Core.Address + ":" + s.Cfg.Core.Port,
		Handler:           s.Router,
		ReadTimeout:       time.Duration(s.Cfg.Core.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.Cfg.Core.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.Cfg.Core.WriteTimeout) * time.Second,
		IdleTimeout:       30 * time.Second,
		ErrorLog:          s.logInterface.Instance,
	}

	go func() {
		// Gracefully stop the server in case of a signal
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case signal := <-sig:
				fmt.Printf("Signal (%d) received, stopping\n", signal)

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				// Flush metrics on server stop
				if s.Cfg.Core.Metrics.Enabled {
					s.flushMetrics()
				}

				server.SetKeepAlivesEnabled(false)
				if err := server.Shutdown(ctx); err != nil {
					s.Log(debug.LogNone, fmt.Errorf("could not gracefully shutdown the server: %v", err))
				}
				close(sig)
				os.Exit(1)
			}
		}
	}()

	switch {

	// If TLS is enabled the server will start in TLS
	case s.Cfg.Core.TLS.Enabled && s.httpCheckTLS():
		s.Log(debug.LogNone, fmt.Errorf("MicroHTTP is listening on port %s with TLS", s.Cfg.Core.Port))
		tlsc := s.httpCreateTLSConfig()

		// Handle autocert
		if s.Cfg.Core.TLS.AutoCert.Enabled && len(s.Cfg.Core.TLS.AutoCert.Certificates) > 0 {
			acm := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.Cfg.Core.TLS.AutoCert.Certificates...),
				Cache:      autocert.DirCache(s.Cfg.Core.FileDir + "certs/"),
			}
			tlsc.GetCertificate = acm.GetCertificate
		} else {
			tlsCert = s.Cfg.Core.TLS.TLSCert
			tlsKey = s.Cfg.Core.TLS.TLSKey
		}
		server.TLSConfig = tlsc

		err := server.ListenAndServeTLS(tlsCert, tlsKey)
		if err != nil {
			panic(err)
		}

	// if TLS is not enabled HTTP will be served
	default:
		s.Log(debug.LogNone, fmt.Errorf("MicroHTTP is listening on port %s", s.Cfg.Core.Port))
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}
}
