package main

import (
	"strings"

	"github.com/redmaner/smux"
)

func (m *micro) configureRouter() {

	m.router = smux.NewRouter()

	// Make routes for each vhost, if vhosts are enabled
	if m.config.Serve.VirtualHosting {

		// Loop over each Vhost
		for k, _ := range m.config.Serve.VirtualHosts {

			// Start with proxy
			if m.vhosts[k].Proxy.Enabled {
				for host, _ := range m.vhosts[k].Proxy.Rules {
					m.router.AddRoute(host, "/", true, "GET", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "PUT", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "POST", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "DELETE", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "HEAD", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "CONNECT", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "PATCH", "*", m.httpProxy())
					m.router.AddRoute(host, "/", true, "OPTIONS", "*", m.httpProxy())
				}

				// Default is serve
			} else {

				// Loop over each supported method
				for path, method := range m.vhosts[k].Methods {

					contentType := ""

					// Loop over each Content-Type for given path
					if content, ok := m.vhosts[k].ContentTypes.RequestTypes[path]; ok {
						contentType = content
					}

					if strings.IndexByte(method, ';') > -1 {
						for _, mtd := range strings.Split(method, ";") {
							m.router.AddRoute(k, path, true, mtd, contentType, m.httpServe())
						}
					} else {
						m.router.AddRoute(k, path, true, method, contentType, m.httpServe())
					}
				}
			}
		}
	}

	// Start with proxy
	if m.config.Proxy.Enabled {
		for host, _ := range m.config.Proxy.Rules {
			m.router.AddRoute(host, "/", true, "GET", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "PUT", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "POST", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "DELETE", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "HEAD", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "CONNECT", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "PATCH", "*", m.httpProxy())
			m.router.AddRoute(host, "/", true, "OPTIONS", "*", m.httpProxy())
		}

		// Default is serve

	} else {
		// Normal serve is enabled
		// Loop over each supported method
		for path, method := range m.config.Methods {

			contentType := ";"

			// Loop over each Content-Type for given path
			if content, ok := m.config.ContentTypes.RequestTypes[path]; ok {
				contentType = content
			}

			if strings.IndexByte(method, ';') > -1 {
				for _, mtd := range strings.Split(method, ";") {
					m.router.AddRoute(smux.DefaultHost, path, true, mtd, contentType, m.httpServe())
				}
			} else {
				m.router.AddRoute(smux.DefaultHost, path, true, method, contentType, m.httpServe())
			}
		}
	}
}
