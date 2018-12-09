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
	"path/filepath"
	"regexp"
	"strings"
)

// Function to set headers defined in configuration
func (m *micro) httpSetHeaders(w http.ResponseWriter, h map[string]string) {

	// MicroHTTP sets security headers at the most strict configuration
	// These headers can be overwritten with the headers element in the configratution file
	// Overwriting is possible in both the main configuration and the vhost configuration
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Feature-Policy", "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';")
	w.Header().Set("Server", "MicroHTTP")

	// If TLS is enabled, the Strict-Transport-Security header is set
	// These settings can be set in the configuration
	if m.config.Core.TLS.Enabled {
		hstr := fmt.Sprintf("max-age=%d;", m.config.Core.TLS.HSTS.MaxAge)
		if m.config.Core.TLS.HSTS.IncludeSubdomains {
			hstr = hstr + " includeSubdomains;"
		}
		if m.config.Core.TLS.HSTS.Preload {
			hstr = hstr + " preload"
		}
		w.Header().Set("Strict-Transport-Security", hstr)
	}

	// All headers set in the configuration are set
	for k, v := range h {
		w.Header().Set(k, v)
	}
}

// Function to set Content-Type depending on the file that is served
func httpGetContentType(p *string, cts *contentTypes) string {
	ext := filepath.Ext(*p)
	switch ext {
	case ".aac":
		return "audio/aac"
	case ".avi":
		return "video/x-msvideo"
	case ".bmp":
		return "image/bmp"
	case ".css":
		return "text/css; charset=utf-8"
	case ".csv":
		return "text/csv; charset=utf-8"
	case ".gif":
		return "image/gif"
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".mpeg":
		return "video/mpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".svg":
		return "image/svg+xml"
	case ".txt":
		return "text/plain"
	case ".xhtml":
		return "application/xhtml-xml"
	case ".xml":
		return "application/xml; charset=utf-8"
	case ".zip":
		return "application/zip"
	}

	// Load custom content type if it exists
	if val, ok := cts.ResponseTypes[*p]; ok {
		return val
	}

	// Default Content-Type is text/html
	return "application/x-unknown"
}

// Function to trim port of an address
func httpTrimPort(s string) string {
	if match, err := regexp.MatchString(":", s); match && err == nil {
		hs := strings.Split(s, ":")
		return hs[0]
	}
	return s
}
