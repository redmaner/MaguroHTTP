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
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/redmaner/MicroHTTP/html"
	"github.com/redmaner/MicroHTTP/router"
)

// Function to write HTTP error to ResponseWriter
func (s *Server) handleError(w http.ResponseWriter, r *http.Request, e int) {

	s.LogNetwork(e, r)
	s.setHeaders(w, map[string]string{})

	host := router.StripHostPort(r.Host)
	cfg := s.Cfg

	// If virtual hosting is enabled, the configuration is switched to the
	// configuration of the vhost
	if cfg.Core.VirtualHosting {
		if _, ok := cfg.Core.VirtualHosts[host]; ok {
			cfg = s.Vhosts[host]
		}
	}

	// Custom error pages can be set in the configuration.
	if val, ok := cfg.Errors[strconv.Itoa(e)]; ok {
		if _, err := os.Stat(val); err == nil {
			http.ServeFile(w, r, val)
			return
		}
	}

	// If custom error pages aren't set the default error message is shown.
	// This is a very basic HTTP error code page without any technical information.
	w.WriteHeader(e)
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, html.PageTemplateStart)
	switch e {
	case 403:
		io.WriteString(w, "<h3>Error 403 - Forbidden</h3>")
	case 404:
		io.WriteString(w, "<h3>Error 404 - Page not found</h3>")
	case 405:
		io.WriteString(w, "<h3>Error 405 - Method not allowed</h3>")
	case 406:
		io.WriteString(w, "<h3>Error 406 - Unacceptable</h3>")
	case 429:
		io.WriteString(w, "<h3>Error 429 - Too many requests</h3>")
	case 502:
		io.WriteString(w, "<h3>Error 502 - Bad gateway</h3>")
	default:
		io.WriteString(w, fmt.Sprintf("<h3>Error %d</h3>", e))
	}
	io.WriteString(w, html.PageTemplateEnd)
}
