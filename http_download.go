// Copyright 2018 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Function to handle HTTP requests to MicroHTTP download server
// This can be further configurated in the configuration file
// MicroHTTP download server generates a table of downloadable files based on extensions
func (m *micro) httpServeDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host := httpTrimPort(r.Host)
		remote := httpTrimPort(r.RemoteAddr)

		var dlurls []fileInfo

		cfg := m.config

		// If virtual hosting is enabled, the configuration is switched to the
		// configuration of the vhost
		if cfg.Core.VirtualHosting {
			if _, ok := cfg.Core.VirtualHosts[host]; ok {
				cfg = m.vhosts[host]
			}
		}

		// Collect downloadable files
		if cfg.Serve.Download.Enabled {
			for _, v := range cfg.Serve.Download.Exts {
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

		// If the request path is ServeIndex, generate the index page with downloadable files
		if path == cfg.Serve.ServeIndex {
			w.Header().Set("Content-Type", "text/html")
			m.httpSetHeaders(w, cfg.Serve.Headers)
			io.WriteString(w, htmlStart)
			io.WriteString(w, "<h1>Downloads</h1>")
			io.WriteString(w, fmt.Sprintln(`<table border="0" cellpadding="0" cellspacing="0">`))
			io.WriteString(w, fmt.Sprintln(`<tr><td height="auto" width="200px"><span><b>Name</b></span><td height="auto" width="120px"><span><b>Size</b></span></td><td height="auto" width="auto"><span><b>Modification date</b></span></td></tr>`))
			for _, v := range dlurls {
				io.WriteString(w, fmt.Sprint(`<tr><td height="auto" width="200px"><span><a href="/`, v.name, `">`, v.name, `</a><br></span><td height="auto" width="120px"><span >`, v.size, `</b></span></td><td height="auto" width="auto"><span>`, v.modTime, `</b></span></td></tr>`))
			}
			io.WriteString(w, fmt.Sprintln("</table><br>"))
			io.WriteString(w, htmlEnd)
			logNetwork(200, r)
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))

			// If the request path is not the index, and the file does exist in ServeDir
			// the file is served and forced to be downloaded by the recipient.
			// If the file doesn't exist, a 404 error is returned.
		} else if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			w.Header().Set("Content-Type", httpGetContentType(&path, &cfg.Serve.ContentTypes))
			if cfg.Serve.Download.Enabled {
				w.Header().Set("Content-Disposition", "attachement")
			}
			m.httpSetHeaders(w, cfg.Serve.Headers)
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			logNetwork(200, r)
			m.md.concat(200, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		} else {

			// Path wasn't found, so we return a 404 not found error.
			m.httpError(w, r, 404)
			return
		}
	}
}
