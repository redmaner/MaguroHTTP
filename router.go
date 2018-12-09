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
	"strings"

	"github.com/redmaner/smux"
)

func (m *micro) configureRouter() {

	m.router = &smux.SRouter{
		ErrorHandler: m.httpError,
	}

	// Make routes for each vhost, if vhosts are enabled
	if m.config.Core.VirtualHosting {

		// Loop over each Vhost
		for vhost := range m.config.Core.VirtualHosts {

			// Start with proxy
			if m.vhosts[vhost].Proxy.Enabled {
				for host := range m.vhosts[vhost].Proxy.Rules {
					m.router.AddRoute(host, "/", true, "GET", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "PUT", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "POST", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "DELETE", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "HEAD", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "CONNECT", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "PATCH", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "OPTIONS", "*", m.httpProxy())
				}

			} else if m.vhosts[vhost].Serve.Download.Enabled {
				m.router.AddRoute(vhost, "/", true, "GET", "", m.httpServeDownload())

				// Default is serve
			} else {

				// Loop over each supported method
				for path, method := range m.vhosts[vhost].Serve.Methods {

					var fallback bool
					contentType := ";"

					if path[len(path)-1] == '/' {
						fallback = true
					}

					// Loop over each Content-Type for given path
					if content, ok := m.vhosts[vhost].Serve.ContentTypes.RequestTypes[path]; ok {
						contentType = content
					}

					if strings.IndexByte(method, ';') > -1 {
						for _, mtd := range strings.Split(method, ";") {
							m.router.AddRoute(vhost, path, fallback, mtd, contentType, m.httpServe())
						}
					} else {
						m.router.AddRoute(vhost, path, fallback, method, contentType, m.httpServe())
					}
				}
			}
		}
	} else {

		// Start with proxy
		if m.config.Proxy.Enabled {
			for host := range m.config.Proxy.Rules {
				m.router.AddRoute(host, "/", true, "GET", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "PUT", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "POST", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "DELETE", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "HEAD", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "CONNECT", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "PATCH", "*", m.httpProxy())
				m.router.AddRoute(host, "/", true, "OPTIONS", "*", m.httpProxy())
			}

		} else if m.config.Serve.Download.Enabled {
			m.router.AddRoute(smux.DefaultHost, "/", true, "GET", "", m.httpServeDownload())

			// Default is serve
		} else {
			// Normal serve is enabled
			// Loop over each supported method
			for path, method := range m.config.Serve.Methods {

				var fallback bool
				contentType := ";"

				if path[len(path)-1] == '/' {
					fallback = true
				}

				// Loop over each Content-Type for given path
				if content, ok := m.config.Serve.ContentTypes.RequestTypes[path]; ok {
					contentType = content
				}

				if strings.IndexByte(method, ';') > -1 {
					for _, mtd := range strings.Split(method, ";") {
						m.router.AddRoute(smux.DefaultHost, path, fallback, mtd, contentType, m.httpServe())
					}
				} else {
					m.router.AddRoute(smux.DefaultHost, path, fallback, method, contentType, m.httpServe())
				}
			}
		}
	}

	// MicroMetrics
	if m.config.Metrics.Enabled {
		m.router.AddRoute(smux.DefaultHost, m.config.Metrics.Path, true, "GET", "", m.httpMetricsRoot())
		m.router.AddRoute(smux.DefaultHost, m.config.Metrics.Path+"/admin", false, "GET", "application/json", m.httpMetricsAdmin())
		m.router.AddRoute(smux.DefaultHost, m.config.Metrics.Path+"/retrieve", false, "POST", ";application/x-www-form-urlencoded", m.httpMetricsRetrieve())
	}
}
