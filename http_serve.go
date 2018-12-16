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
	"fmt"
	"net/http"
	"os"
)

// Serve type, part of the MicroHTTP config
type serveConfig struct {
	ServeDir   string
	ServeIndex string
	Headers    map[string]string
	Methods    map[string]string
	MIMETypes  MIMETypes
	Download   download
}

// MIMETypes type, part of MicroHTTP serveConfig
type MIMETypes struct {
	ResponseTypes map[string]string
	RequestTypes  map[string]string
}

// Function to handle HTTP requests to MicroHTTP server
// This can be further configurated in the configuration file
// MicroHTTP is capable to host multiple websites on one server using virtual hosts
func (m *micro) httpServe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host := httpTrimPort(r.Host)
		remote := httpTrimPort(r.RemoteAddr)

		cfg := m.config

		// If virtual hosting is enabled, the configuration is switched to the
		// configuration of the vhost
		if cfg.Core.VirtualHosting {
			if _, ok := cfg.Core.VirtualHosts[host]; ok {
				cfg = m.vhosts[host]
			}
		}

		path := r.URL.Path

		// Check firewall for path
		if block := firewallHTTP(&cfg, remote, path); block {
			m.httpError(w, r, 403)
			return
		}

		// Correct path to ServeIndex when path is root
		if path == "/" {
			path = cfg.Serve.ServeIndex
		}

		// Serve the file that is requested by path if it esists in ServeDir.
		// If the requested path doesn't exist, return a 404 error
		if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			w.Header().Set("Content-Type", httpGetMIMEType(path, cfg.Serve.MIMETypes))
			m.httpSetHeaders(w, cfg.Serve.Headers)
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			logNetwork(200, r)
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		} else {

			// Path wasn't found, so we return a 404 not found error
			m.httpError(w, r, 404)
			return
		}
	}
}
