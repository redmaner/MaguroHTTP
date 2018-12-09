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
	"io"
	"net/http"
)

// Proxy type, part of MicroHTTP config
type proxy struct {
	Enabled bool
	Rules   map[string]string
}

// Function to proxy. The proxy can be configurated in configuration
// MicroHTTP is capable to serve HTTP and to proxy along side each other using virtual hosts
func (m *micro) httpProxy() http.HandlerFunc {
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

		if val, ok := cfg.Proxy.Rules[host]; ok {

			if block := firewallProxy(&cfg, remote, host); block {
				m.httpError(w, r, 403)
				return
			}

			req, err := http.NewRequest(r.Method, val, r.Body)
			if err != nil {
				logAction(logERROR, err)
				m.httpError(w, r, 502)
				return
			}
			req.URL.Path = r.URL.Path
			req.URL.RawPath = r.URL.RawPath
			req.URL.RawQuery = r.URL.RawQuery
			req.RemoteAddr = r.RemoteAddr

			for k, v := range r.Header {
				req.Header[k] = v
			}

			if resp, err := m.client.Do(req); err == nil {

				w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
				w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

				// Proxy back all response headers
				for k, v := range resp.Header {
					w.Header().Set(k, v[0])
				}

				// Write header last. If header is written, headers can no longer be set
				w.WriteHeader(resp.StatusCode)

				io.Copy(w, resp.Body)
				resp.Body.Close()
				logNetwork(resp.StatusCode, r)
				m.md.concat(resp.StatusCode, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
			} else {
				logAction(logERROR, err)
				m.httpError(w, r, 502)
			}
		}
	}
}
