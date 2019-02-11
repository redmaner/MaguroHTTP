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

package micro

import (
	"strings"

	"github.com/redmaner/MicroHTTP/router"
)

func (s *Server) addRoutesFromConfig() {

	// Make routes for each vhost, if vhosts are enabled
	if s.Cfg.Core.VirtualHosting {

		// Loop over each Vhost
		for vhost := range s.Cfg.Core.VirtualHosts {

			// Start with proxy
			if s.Vhosts[vhost].Proxy.Enabled {
				for host := range s.Vhosts[vhost].Proxy.Rules {
					s.Router.AddRoute(host, "/", true, "GET", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "PUT", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "POST", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "DELETE", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "HEAD", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "CONNECT", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "PATCH", "*", s.handleProxy())
					s.Router.AddRoute(host, "/", true, "OPTIONS", "*", s.handleProxy())
				}

			} else if s.Vhosts[vhost].Serve.Download.Enabled {
				s.Router.AddRoute(vhost, "/", true, "GET", "", s.handleDownload())

				// Default is serve
			} else {

				// Loop over each supported method
				for path, method := range s.Vhosts[vhost].Serve.Methods {

					var fallback bool
					contentType := ";"

					if path[len(path)-1] == '/' {
						fallback = true
					}

					// Loop over each Content-Type for given path
					if content, ok := s.Vhosts[vhost].Serve.MIMETypes.RequestTypes[path]; ok {
						contentType = content
					}

					if strings.IndexByte(method, ';') > -1 {
						for _, mtd := range strings.Split(method, ";") {
							s.Router.AddRoute(vhost, path, fallback, mtd, contentType, s.handleServe())
						}
					} else {
						s.Router.AddRoute(vhost, path, fallback, method, contentType, s.handleServe())
					}
				}
			}
		}
	} else {

		// Start with proxy
		if s.Cfg.Proxy.Enabled {
			for host := range s.Cfg.Proxy.Rules {
				s.Router.AddRoute(host, "/", true, "GET", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "PUT", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "POST", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "DELETE", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "HEAD", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "CONNECT", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "PATCH", "*", s.handleProxy())
				s.Router.AddRoute(host, "/", true, "OPTIONS", "*", s.handleProxy())
			}

		} else if s.Cfg.Serve.Download.Enabled {
			s.Router.AddRoute(router.DefaultHost, "/", true, "GET", "", s.handleDownload())

			// Default is serve
		} else {
			// Normal serve is enabled
			// Loop over each supported method
			for path, method := range s.Cfg.Serve.Methods {

				var fallback bool
				contentType := ";"

				if path[len(path)-1] == '/' {
					fallback = true
				}

				// Loop over each Content-Type for given path
				if content, ok := s.Cfg.Serve.MIMETypes.RequestTypes[path]; ok {
					contentType = content
				}

				if strings.IndexByte(method, ';') > -1 {
					for _, mtd := range strings.Split(method, ";") {
						s.Router.AddRoute(router.DefaultHost, path, fallback, mtd, contentType, s.handleServe())
					}
				} else {
					s.Router.AddRoute(router.DefaultHost, path, fallback, method, contentType, s.handleServe())
				}
			}
		}
	}
}
