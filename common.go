package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// HTML start constant
const htmlStart = `<!doctype html>
<html class="no-js" lang="">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
<title></title>
<meta name="description" content="">
<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
`

// HTML end constant
const htmlEnd = `</body>
</html>
`

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
	if m.config.TLS {
		hstr := fmt.Sprintf("max-age=%d;", m.config.HSTS.MaxAge)
		if m.config.HSTS.IncludeSubdomains {
			hstr = hstr + " includeSubdomains;"
		}
		if m.config.HSTS.Preload {
			hstr = hstr + " preload"
		}
		w.Header().Set("Strict-Transport-Security", hstr)
	}

	// All headers set in the configuration are set
	for k, v := range h {
		w.Header().Set(k, v)
	}
}

// Function to write HTTP error to ResponseWriter
func (m *micro) httpThrowError(w http.ResponseWriter, r *http.Request, e int) {

	// Custom error pages can be set in the configuration.
	if val, ok := m.config.Errors[strconv.Itoa(e)]; ok {
		if _, err := os.Stat(val); err == nil {
			http.ServeFile(w, r, val)
			return
		}
	}

	// If custom error pages aren't set the default error message is shown.
	// This is a very basic HTTP error code page without any technical information.
	w.WriteHeader(e)
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, htmlStart)
	switch e {
	case 403:
		io.WriteString(w, "<h3>Error 403 - Forbidden</h3>")
	case 404:
		io.WriteString(w, "<h3>Error 404 - Page not found</h3>")
	case 405:
		io.WriteString(w, "<h3>Error 405 - Method not allowed</h3>")
	case 406:
		io.WriteString(w, "<h3>Error 406 - Unacceptable</h3>")
	case 502:
		io.WriteString(w, "<h3>Error 502 - Bad gateway</h3>")
	default:
		io.WriteString(w, fmt.Sprintf("<h3>Error %d</h3>", e))
	}
	io.WriteString(w, htmlEnd)
	logNetwork(e, r)
}
