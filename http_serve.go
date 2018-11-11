package main

import (
	"fmt"
	"net/http"
	"os"
)

const defaultMethods = "GET;"

// Function to handle HTTP requests to MicroHTTP server
// This can be further configurated in the configuration file
// MicroHTTP is capable to host multiple websites on one server using virtual hosts
func (m *micro) httpServe(w http.ResponseWriter, r *http.Request) {

	host := httpTrimPort(r.Host)
	remote := httpTrimPort(r.RemoteAddr)

	cfg := m.config

	// If virtual hosting is enabled, the configuration is switched to the
	// configuration of the vhost
	if cfg.Serve.VirtualHosting {
		if _, ok := cfg.Serve.VirtualHosts[host]; ok {
			cfg = m.vhosts[host]
		}
	}

	// Check for proxy. If proxy is enabled, httpProxy is called.
	if cfg.Proxy.Enabled {
		m.httpProxy(w, r, &cfg)
		return
	}

	// Validate request content type. If the request Content-Type is not supported
	// MicroHTTP will throw an HTTP 406 error. The supported request Content-Type can
	// be configurated in the configuration
	rct := r.Header.Get("Content-Type")
	if !httpValidateRequestContentType(&rct, &cfg.ContentTypes) {
		m.httpError(w, r, 406)
		return
	}

	path := r.URL.Path

	// Check firewall for path
	if block := firewallHTTP(&cfg, remote, path); block {
		m.httpError(w, r, 403)
		return
	}

	// Determine allowed methods
	// If the method is not allowed MicroHTTP will throw a HTTP 405 error
	// The allowed HTTP methods can be set in the configuration
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
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		} else {
			m.httpError(w, r, 405)
			return
		}
	} else {
		m.httpError(w, r, 404)
		return
	}
}
