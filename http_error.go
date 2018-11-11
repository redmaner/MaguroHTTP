package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Function to write HTTP error to ResponseWriter
func (m *micro) httpError(w http.ResponseWriter, r *http.Request, e int) {

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
	m.md.concat(e, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
}
