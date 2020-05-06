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
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/redmaner/MaguroHTTP/debug"
)

// Function to set headers defined in configuration
func (s *Server) setHeaders(w http.ResponseWriter, h map[string]string, isProxy bool) {

	if !isProxy {
		// MaguroHTTP sets security headers at the most strict configuration
		// These headers can be overwritten with the headers element in the configratution file
		// Overwriting is possible in both the main configuration and the vhost configuration
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Feature-Policy", "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';")
		w.Header().Set("Server", "MaguroHTTP")
	}

	// If TLS is enabled, the Strict-Transport-Security header is set
	// These settings can be set in the configuration
	if s.Cfg.Core.TLS.Enabled {
		hstr := fmt.Sprintf("max-age=%d;", s.Cfg.Core.TLS.HSTS.MaxAge)
		if s.Cfg.Core.TLS.HSTS.IncludeSubdomains {
			hstr = hstr + " includeSubdomains;"
		}
		if s.Cfg.Core.TLS.HSTS.Preload {
			hstr = hstr + " preload"
		}
		w.Header().Set("Strict-Transport-Security", hstr)
	}

	if h != nil {
		// All headers set in the configuration are set
		for k, v := range h {
			w.Header().Set(k, v)
		}
	}
}

// Function to set MIME type depending on the file extension
func getMIMEType(p string, cts MIMETypes) string {
	ext := filepath.Ext(p)
	switch ext {
	case ".aac":
		return "audio/aac"
	case ".avi":
		return "video/x-msvideo"
	case ".bmp":
		return "image/bmp"
	case ".css":
		return "text/css"
	case ".csv":
		return "text/csv"
	case ".gif":
		return "image/gif"
	case ".html", ".htm":
		return "text/html"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".js":
		return "text/javascript"
	case ".json":
		return "text/json"
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
		return "application/xml"
	case ".zip":
		return "application/zip"
	}

	// Load custom content type if it exists
	if val, ok := cts.ResponseTypes[p]; ok {
		return val
	}

	// Default Content-Type is octet-stream
	return "application/octet-stream"
}

// Copy HTTP header to an existing HTTP header
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Clone a HTTP header
func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// WriteString is a wrapper function to write strings with incoperated error handling
func (s *Server) WriteString(w io.Writer, str string) {

	var err error
	sw, ok := w.(io.StringWriter)

	switch {
	case ok:
		_, err = sw.WriteString(str)
	default:
		_, err = w.Write([]byte(str))
	}

	s.Log(debug.LogError, err)
}
