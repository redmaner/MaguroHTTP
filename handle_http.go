package main

import (
	"fmt"
	"net/http"
	"os"
)

const defaultMethods = "GET;"

// Function to handle HTTP requests to MicroHTTP server
func handleHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	logAction(logVERBOSE, fmt.Errorf(path))

	// Determine allowed methods
	var methods string
	if val, ok := mCfg.Methods["/"]; ok {
		methods = val
	} else {
		methods = defaultMethods
	}

	// If the url path is root, serve the ServeIndex file.
	if path == "/" {
		if methodAllowed(r.Method, methods) {
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

		// If path is not root, serve the file that is requested by path if it esists
		// in ServeDir. If the requested path doesn't exist, return a 404 error
	} else if _, err := os.Stat(mCfg.ServeDir + path); err == nil {

		if val, ok := mCfg.Methods[path]; ok {
			methods = val
		}

		if methodAllowed(r.Method, methods) {
			//w.Header().Set("Content-Type", "text/html")
			setHeaders(w, mCfg.Headers)
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
