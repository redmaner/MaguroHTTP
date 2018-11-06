package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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
const htmlEnd = `</body>
</html>
`

// Test whether a HTTP method is allowed
func httpMethodAllowed(m string, a string) bool {
	am := make(map[string]int)
	if match, err := regexp.MatchString(";", a); match && err == nil {
		sc := strings.Split(a, ";")
		for _, k := range sc {
			if k != "" {
				am[k] = 0
			}
		}
	} else {
		if a != "" {
			am[a] = 0
		}
	}
	if _, ok := am[m]; ok {
		return true
	}
	return false
}

// Function to set headers defined in configuration
func httpSetHeaders(w http.ResponseWriter, h map[string]string) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Feature-Policy", "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';")
	w.Header().Set("Server", "MicroHTTP")
	for k, v := range h {
		w.Header().Set(k, v)
	}
}

// Function to write error
func httpThrowError(w http.ResponseWriter, r *http.Request, e string) {
	if val, ok := mCfg.Errors[e]; ok {
		if _, err := os.Stat(val); err == nil {
			http.ServeFile(w, r, val)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, htmlStart)
	switch e {
	case "404":
		io.WriteString(w, "<h3>Error 404 - Page not found</h3>")
	case "405":
		io.WriteString(w, "<h3>Error 405 - Method not allowed</h3>")
	default:
		io.WriteString(w, fmt.Sprintf("<h3>Error %s</h3>", e))
	}
	io.WriteString(w, htmlEnd)
}
