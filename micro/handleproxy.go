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
	"io"
	"net/http"

	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/router"
)

// Function to proxy. The proxy can be configurated in configuration
// MicroHTTP is capable to serve HTTP and to proxy along side each other using virtual hosts
func (s *Server) handleProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host := router.StripHostPort(r.Host)

		cfg := s.Cfg

		// If virtual hosting is enabled, the configuration is switched to the
		// configuration of the vhost
		if cfg.Core.VirtualHosting {
			if _, ok := cfg.Core.VirtualHosts[host]; ok {
				cfg = s.Vhosts[host]
			}
		}

		if val, ok := cfg.Proxy.Rules[host]; ok {

			req, err := http.NewRequest(r.Method, val+r.RequestURI, r.Body)
			if err != nil {
				s.Log(debug.LogError, err)
				s.handleError(w, r, 502)
				return
			}

			req.Header = cloneHeader(r.Header)

			if resp, err := s.transport.RoundTrip(req); err == nil {

				// Proxy back all response headers
				copyHeader(w.Header(), resp.Header)

				// Write header last. If header is written, headers can no longer be set
				w.WriteHeader(resp.StatusCode)

				io.Copy(w, resp.Body)
				resp.Body.Close()
				s.LogNetwork(resp.StatusCode, r)
			} else {
				s.Log(debug.LogError, err)
				s.handleError(w, r, 502)
			}
		}
	}
}
