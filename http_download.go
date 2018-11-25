package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Download type, part of the MicroHTTP config
type download struct {
	Enabled bool
	Exts    []string
}

// Type fileinfo to gather information about files
type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

// Function to handle HTTP requests to MicroHTTP server
// This can be further configurated in the configuration file
// MicroHTTP is capable to host multiple websites on one server using virtual hosts
func (m *micro) httpServeDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host := httpTrimPort(r.Host)
		remote := httpTrimPort(r.RemoteAddr)

		var dlurls []fileInfo

		cfg := m.config

		// If virtual hosting is enabled, the configuration is switched to the
		// configuration of the vhost
		if cfg.Serve.VirtualHosting {
			if _, ok := cfg.Serve.VirtualHosts[host]; ok {
				cfg = m.vhosts[host]
			}
		}

		// Collect downloadable files
		if cfg.Download.Enabled {
			for _, v := range cfg.Download.Exts {
				filepath.Walk(cfg.Serve.ServeDir, func(path string, f os.FileInfo, _ error) error {
					if !f.IsDir() {
						if filepath.Ext(f.Name()) == v {
							dlurls = append(dlurls, fileInfo{f.Name(), f.Size(), f.ModTime()})
						}
					}
					return nil
				})
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
		if path == cfg.Serve.ServeIndex {
			w.Header().Set("Content-Type", "text/html")
			io.WriteString(w, htmlStart)
			io.WriteString(w, "<h1>Downloads</h1>")
			io.WriteString(w, fmt.Sprintln(`<table border="0" cellpadding="0" cellspacing="0">`))
			io.WriteString(w, fmt.Sprintln(`<tr><td height="auto" width="200px"><span><b>Name</b></span><td height="auto" width="120px"><span><b>Size</b></span></td><td height="auto" width="auto"><span><b>Modification date</b></span></td></tr>`))
			for _, v := range dlurls {
				io.WriteString(w, fmt.Sprintln(`<tr><td height="auto" width="200px"><span><a href="/`, v.name, `">`, v.name, `</a><br></span><td height="auto" width="120px"><span >`, v.size, `</b></span></td><td height="auto" width="auto"><span>`, v.modTime, `</b></span></td></tr>`))
			}
			io.WriteString(w, fmt.Sprintln("</table><br>"))
			io.WriteString(w, htmlEnd)
			logNetwork(200, r)
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		} else if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			w.Header().Set("Content-Type", httpGetContentType(&path, &cfg.ContentTypes))
			if cfg.Download.Enabled {
				w.Header().Set("Content-Disposition", "attachement")
			}
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
