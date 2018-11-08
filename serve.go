package main

import (
	"net/http"
	"os"
)

const defaultMethods = "GET;"

// Function to handle HTTP requests to MicroHTTP server
func (m *micro) httpServe(w http.ResponseWriter, r *http.Request) {

	host := httpTrimPort(r.Host)
	remote := httpTrimPort(r.RemoteAddr)

	cfg := m.config
	if cfg.Serve.VirtualHosting {
		if _, ok := cfg.Serve.VirtualHosts[host]; ok {
			cfg = m.vhosts[host]
		}
	}

	// Check for proxy
	if cfg.Proxy.Enabled {
		m.httpProxy(w, r, &cfg)
		return
	}

	// Validate request content type
	rct := r.Header.Get("Content-Type")
	if !httpValidateRequestContentType(&rct, &cfg.ContentTypes) {
		m.httpThrowError(w, r, 406)
		return
	}

	path := r.URL.Path
	if block := firewallHTTP(&cfg, remote, path); block {
		m.httpThrowError(w, r, 403)
		return
	}

	// Determine allowed methods
	var methods string
	if val, ok := cfg.Methods[path]; ok {
		methods = val
	} else {
		methods = defaultMethods
	}

	// Correct path to ServeIndex when path is root
	if path == "/" {
		path = cfg.Serve.ServeIndex
		if val, ok := cfg.Methods["/"]; ok {
			methods = val
		}
	}

	// Serve the file that is requested by path if it esists in ServeDir.
	// If the requested path doesn't exist, return a 404 error
	if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {

		if httpMethodAllowed(&r.Method, &methods) {
			w.Header().Set("Content-Type", httpGetContentType(&path, &cfg.ContentTypes))
			m.httpSetHeaders(w, cfg.Headers)
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			logNetwork(200, r)
		} else {
			m.httpThrowError(w, r, 405)
			return
		}
	} else {
		m.httpThrowError(w, r, 404)
		return
	}
}
