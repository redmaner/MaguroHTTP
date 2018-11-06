package main

import (
	"net/http"
	"os"
)

const defaultMethods = "GET;"

// Function to handle HTTP requests to MicroHTTP server
func handleHTTP(w http.ResponseWriter, r *http.Request) {

	remote := httpTrimPort(r.RemoteAddr)

	// Validate request content type
	rct := r.Header.Get("Content-Type")
	if !httpValidateRequestContentType(rct) {
		httpThrowError(w, r, 406)
		return
	}

	path := r.URL.Path

	if block := firewallHTTP(remote, path); block {
		httpThrowError(w, r, 403)
		return
	}

	// Determine allowed methods
	var methods string
	if val, ok := mCfg.Methods["/"]; ok {
		methods = val
	} else {
		methods = defaultMethods
	}

	// If the url path is root, serve the ServeIndex file.
	if path == "/" {
		if httpMethodAllowed(r.Method, methods) {
			if _, err := os.Stat(mCfg.ServeDir + mCfg.ServeIndex); err == nil {
				httpSetContentType(w, mCfg.ServeIndex)
				httpSetHeaders(w, mCfg.Headers)
				http.ServeFile(w, r, mCfg.ServeDir+mCfg.ServeIndex)
			} else if path != "" {
				httpThrowError(w, r, 404)
				return
			}
		} else {
			httpThrowError(w, r, 405)
			return
		}

		// If path is not root, serve the file that is requested by path if it esists
		// in ServeDir. If the requested path doesn't exist, return a 404 error
	} else if _, err := os.Stat(mCfg.ServeDir + path); err == nil {

		if val, ok := mCfg.Methods[path]; ok {
			methods = val
		}

		if httpMethodAllowed(r.Method, methods) {
			httpSetContentType(w, path)
			httpSetHeaders(w, mCfg.Headers)
			http.ServeFile(w, r, mCfg.ServeDir+path)
		} else {
			httpThrowError(w, r, 405)
			return
		}
	} else {
		httpThrowError(w, r, 404)
		return
	}
}
