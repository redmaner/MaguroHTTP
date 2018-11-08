package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

	if mCfg.TLS {
		hstr := fmt.Sprintf("max-age=%d;", mCfg.HSTS.MaxAge)
		if mCfg.HSTS.IncludeSubdomains {
			hstr = hstr + " includeSubdomains;"
		}
		if mCfg.HSTS.Preload {
			hstr = hstr + " preload"
		}
		w.Header().Set("Strict-Transport-Security", hstr)
	}

	for k, v := range h {
		w.Header().Set(k, v)
	}
}

// Function to write HTTP error to ResponseWriter
func httpThrowError(w http.ResponseWriter, r *http.Request, e int) {
	if val, ok := mCfg.Errors[strconv.Itoa(e)]; ok {
		if _, err := os.Stat(val); err == nil {
			http.ServeFile(w, r, val)
			return
		}
	}
	w.WriteHeader(e)
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, htmlStart)
	switch e {
	case 404:
		io.WriteString(w, "<h3>Error 404 - Page not found</h3>")
	case 405:
		io.WriteString(w, "<h3>Error 405 - Method not allowed</h3>")
	default:
		io.WriteString(w, fmt.Sprintf("<h3>Error %d</h3>", e))
	}
	io.WriteString(w, htmlEnd)
	logNetwork(e, r)
}

// Function to set Content-Type depending on the file that is served
func httpSetContentType(w http.ResponseWriter, p string) {
	ext := filepath.Ext(p)
	switch ext {
	case ".aac":
		w.Header().Set("Content-Type", "audio/aac")
	case ".avi":
		w.Header().Set("Content-Type", "video/x-msvideo")
	case ".bmp":
		w.Header().Set("Content-Type", "image/bmp")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".html", ".htm":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".jpeg", ".jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".mpeg":
		w.Header().Set("Content-Type", "video/mpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain")
	case ".xhtml":
		w.Header().Set("Content-Type", "application/xhtml-xml")
	case ".xml":
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	case ".zip":
		w.Header().Set("Content-Type", "application/zip")
	}

	// Load custom content type if it exists
	if val, ok := mCfg.ContentTypes.ResponseTypes[p]; ok {
		w.Header().Set("Content-Type", val)
	}

}

// Function to validate a request Content-Type
func httpValidateRequestContentType(rct string, cts contentTypes) bool {
	if len(cts.RequestTypes) != 0 {
		for _, v := range cts.RequestTypes {
			if v == rct {
				return true
			}
		}
	}
	return false
}

// Function to trim port of an address
func httpTrimPort(s string) string {
	if match, err := regexp.MatchString(":", s); match && err == nil {
		hs := strings.Split(s, ":")
		return hs[0]
	}
	return s
}
