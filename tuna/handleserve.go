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

package tuna

import (
	"net/http"
	"os"

	"github.com/redmaner/MaguroHTTP/router"
)

// Function to handle HTTP requests to MaguroHTTP server
// This can be further configurated in the configuration file
// MaguroHTTP is capable to host multiple websites on one server using virtual hosts
func (s *Server) handleServe() http.HandlerFunc {
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

		path := r.URL.Path

		// If path ends with a slash, add ServeIndex
		if path[len(path)-1] == '/' {
			path = path + cfg.Serve.ServeIndex
		}

		// Serve the file that is requested by path if it esists in ServeDir.
		// If the requested path doesn't exist, return a 404 error
		if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			s.setHeaders(w, cfg.Serve.Headers, false)
			w.Header().Set("Content-Type", getMIMEType(path, cfg.Serve.MIMETypes))
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			s.LogNetwork(200, r)
		} else {

			// Path wasn't found, so we return a 404 not found error
			s.HandleError(w, r, 404)
			return
		}
	}
}
