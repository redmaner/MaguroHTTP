package main

import (
	"fmt"
	"net/http"
	"os"
)

// Serve type, part of the MicroHTTP config
type serve struct {
	ServeDir       string
	ServeIndex     string
	VirtualHosting bool
	VirtualHosts   map[string]string
}

// Function to handle HTTP requests to MicroHTTP server
// This can be further configurated in the configuration file
// MicroHTTP is capable to host multiple websites on one server using virtual hosts
func (m *micro) httpServe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		path := r.URL.Path

		// Check firewall for path
		if block := firewallHTTP(&cfg, remote, path); block {
			m.httpError(w, r, 403)
			return
		}

		// Correct path to ServeIndex when path is root
		if path == "/" {
			path = cfg.Serve.ServeIndex
		}

		// Serve the file that is requested by path if it esists in ServeDir.
		// If the requested path doesn't exist, return a 404 error
		if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			w.Header().Set("Content-Type", httpGetContentType(&path, &cfg.ContentTypes))
			m.httpSetHeaders(w, cfg.Headers)
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			logNetwork(200, r)
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		} else {
			m.httpError(w, r, 404)
			return
		}
	}
}
