// Copyright 2018-2019 Jake van der Putten.
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

package tuna

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/redmaner/MaguroHTTP/router"
)

// Function to handle HTTP requests to MaguroHTTP download server
// This can be further configurated in the configuration file
// MaguroHTTP download server generates a table of downloadable files based on extensions
func (s *Server) handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		host := router.StripHostPort(r.Host)

		var dlurls []fileInfo

		cfg := s.Cfg

		// If virtual hosting is enabled, the configuration is switched to the
		// configuration of the vhost
		if cfg.Core.VirtualHosting {
			if _, ok := cfg.Core.VirtualHosts[host]; ok {
				cfg = s.Vhosts[host]
			}
		}

		// Collect downloadable files
		if cfg.Serve.Download.Enabled {
			for _, v := range cfg.Serve.Download.Exts {
				err := filepath.Walk(cfg.Serve.ServeDir, func(path string, f os.FileInfo, _ error) error {
					if !f.IsDir() {
						if filepath.Ext(f.Name()) == v {
							dlurls = append(dlurls, fileInfo{
								Name:    f.Name(),
								Size:    f.Size(),
								ModTime: f.ModTime(),
							})
						}
					}
					return nil
				})
				if err != nil {
					s.handleError(w, r, 500)
					return
				}
			}
		}

		path := r.URL.Path

		// Correct path to ServeIndex when path is root
		if path == "/" {
			path = cfg.Serve.ServeIndex
		}

		buf := bytes.NewBufferString("")

		// If the request path is ServeIndex, generate the index page with downloadable files
		if path == cfg.Serve.ServeIndex {
			w.Header().Set("Content-Type", "text/html")
			s.setHeaders(w, cfg.Serve.Headers, false)
			s.WriteString(buf, "<h1>Downloads</h1>")
			s.WriteString(buf, fmt.Sprintln(`<table border="0" cellpadding="0" cellspacing="0">`))
			s.WriteString(buf, fmt.Sprintln(`<tr><td height="auto" width="200px"><span><b>Name</b></span><td height="auto" width="120px"><span><b>Size</b></span></td><td height="auto" width="auto"><span><b>Modification date</b></span></td></tr>`))
			for _, v := range dlurls {
				s.WriteString(buf, fmt.Sprint(`<tr><td height="auto" width="200px"><span><a href="/`, v.Name, `">`, v.Name, `</a><br></span><td height="auto" width="120px"><span >`, v.Size, `</b></span></td><td height="auto" width="auto"><span>`, v.ModTime, `</b></span></td></tr>`))
			}
			s.WriteString(buf, fmt.Sprintln("</table><br>"))

			data := struct {
				DownloadTable template.HTML
			}{
				DownloadTable: template.HTML(buf.String()),
			}

			if err := s.templates.download.Execute(w, data); err != nil {
				s.handleError(w, r, 500)
				return
			}
			s.LogNetwork(200, r)

			// If the request path is not the index, and the file does exist in ServeDir
			// the file is served and forced to be downloaded by the recipient.
			// If the file doesn't exist, a 404 error is returned.
		} else if _, err := os.Stat(cfg.Serve.ServeDir + path); err == nil {
			w.Header().Set("Content-Type", getMIMEType(path, cfg.Serve.MIMETypes))
			w.Header().Set("Content-Disposition", "attachement")
			s.setHeaders(w, cfg.Serve.Headers, false)
			http.ServeFile(w, r, cfg.Serve.ServeDir+path)
			s.LogNetwork(200, r)
		} else {

			// Path wasn't found, so we return a 404 not found error.
			s.handleError(w, r, 404)
			return
		}
	}
}
