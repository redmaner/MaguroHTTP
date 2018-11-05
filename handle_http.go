package main

import (
	"net/http"
	"os"
)

// Function to handle HTTP requests to MicroHTTP server
func handleHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	// If the url path is root, serve the ServeIndex file.
	if path == "/" {
		if methodAllowed(r.Method, mCfg.Methods) {
			if _, err := os.Stat(mCfg.ServeDir + mCfg.ServeIndex); err == nil {
				w.Header().Set("Content-Type", "text/html")
				setHeaders(w, mCfg.Headers)
				http.ServeFile(w, r, mCfg.ServeDir+mCfg.ServeIndex)
			} else if path != "" {
				throwError(w, r, "404")
			}
		} else {
			throwError(w, r, "405")
		}

		// Serve path if it matches with the serve directory
	} else if _, err := os.Stat(mCfg.ServeDir + path); err == nil {

		var fCfg = mCfg

		if _, err := os.Stat(mCfg.ServeDir + path + ".json"); err == nil {
			loadConfigFromFile(mCfg.ServeDir+path+".json", &fCfg)
		}

		if methodAllowed(r.Method, fCfg.Methods) {
			w.Header().Set("Content-Type", "text/html")
			setHeaders(w, fCfg.Headers)
			http.ServeFile(w, r, mCfg.ServeDir+path)
		} else {
			throwError(w, r, "405")
		}
	} else {
		throwError(w, r, "404")
	}
}

// Function to set headers defined in configuration
func setHeaders(w http.ResponseWriter, h map[string]string) {
	setSecurityHeaders(w)
	w.Header().Set("Server", "MicroHTTP")
	for k, v := range h {
		w.Header().Set(k, v)
	}
}

// Function to set security headers
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("Feature-Policy", "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';")
}
