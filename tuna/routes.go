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
	"strings"

	"github.com/redmaner/MaguroHTTP/guard"
	"github.com/redmaner/MaguroHTTP/router"
)

func (s *Server) addRoutesFromConfig() {

	var limiter *guard.Limiter
	var firewall *guard.Firewall

	// Make routes for each vhost, if vhosts are enabled
	if s.Cfg.Core.VirtualHosting {

		// Loop over each Vhost
		for vhost := range s.Cfg.Core.VirtualHosts {

			// Each virtual host gets it's own limiter
			limiter = guard.NewLimiter(s.Vhosts[vhost].Guard.Rate, s.Vhosts[vhost].Guard.RateBurst, s.Vhosts[vhost].Guard.FilterOnIP)
			limiter.ErrorHandler = s.HandleError

			if s.Vhosts[vhost].Guard.Firewall.Enabled {
				firewall = &guard.Firewall{
					Blacklisting: s.Vhosts[vhost].Guard.Firewall.Blacklisting,
					Subpath:      s.Vhosts[vhost].Guard.Firewall.Subpath,
					Rules:        s.Vhosts[vhost].Guard.Firewall.Rules,
					ErrorHandler: s.HandleError,
				}
			}

			// Start with proxy
			if s.Vhosts[vhost].Proxy.Enabled {
				for host := range s.Vhosts[vhost].Proxy.Rules {
					for _, mtd := range s.Vhosts[vhost].Proxy.Methods {
						s.Router.AddRoute(host, "/", true, mtd, "*", s.handleProxy())
					}

					// Add firewall as middleware if enabled
					if s.Vhosts[vhost].Guard.Firewall.Enabled {
						s.Router.UseMiddleware(host, "/", router.MiddlewareHandlerFunc(firewall.BlockProxy))
					}

					// Add limiter as middleware
					s.Router.UseMiddleware(host, "/", router.MiddlewareHandlerFunc(limiter.LimitHTTP))
				}

			} else if s.Vhosts[vhost].Serve.Download.Enabled {
				s.Router.AddRoute(vhost, "/", true, "GET", "", s.handleDownload())

				// Add firewall as middleware if enabled
				if s.Vhosts[vhost].Guard.Firewall.Enabled {
					s.Router.UseMiddleware(vhost, "/", router.MiddlewareHandlerFunc(firewall.BlockHTTP))
				}

				// Add limiter as middleware
				s.Router.UseMiddleware(vhost, "/", router.MiddlewareHandlerFunc(limiter.LimitHTTP))

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

					// Add firewall as middleware if enabled
					if s.Vhosts[vhost].Guard.Firewall.Enabled {
						s.Router.UseMiddleware(vhost, path, router.MiddlewareHandlerFunc(firewall.BlockHTTP))
					}

					// Add limiter as middleware
					s.Router.UseMiddleware(vhost, path, router.MiddlewareHandlerFunc(limiter.LimitHTTP))
				}
			}
		}
	} else {

		limiter = guard.NewLimiter(s.Cfg.Guard.Rate, s.Cfg.Guard.RateBurst, s.Cfg.Guard.FilterOnIP)
		limiter.ErrorHandler = s.HandleError

		if s.Cfg.Guard.Firewall.Enabled {
			firewall = &guard.Firewall{
				Blacklisting: s.Cfg.Guard.Firewall.Blacklisting,
				Subpath:      s.Cfg.Guard.Firewall.Subpath,
				Rules:        s.Cfg.Guard.Firewall.Rules,
				ErrorHandler: s.HandleError,
			}
		}

		// Start with proxy
		if s.Cfg.Proxy.Enabled {
			for host := range s.Cfg.Proxy.Rules {
				for _, mtd := range s.Cfg.Proxy.Methods {
					s.Router.AddRoute(host, "/", true, mtd, "*", s.handleProxy())
				}

				// Add firewall as middleware if enabled
				if s.Cfg.Guard.Firewall.Enabled {
					s.Router.UseMiddleware(host, "/", router.MiddlewareHandlerFunc(firewall.BlockProxy))
				}

				// Add limiter as middleware
				s.Router.UseMiddleware(host, "/", router.MiddlewareHandlerFunc(limiter.LimitHTTP))
			}

		} else if s.Cfg.Serve.Download.Enabled {
			s.Router.AddRoute(router.DefaultHost, "/", true, "GET", "", s.handleDownload())

			// Add firewall as middleware if enabled
			if s.Cfg.Guard.Firewall.Enabled {
				s.Router.UseMiddleware(router.DefaultHost, "/", router.MiddlewareHandlerFunc(firewall.BlockHTTP))
			}

			// Add limiter as middleware
			s.Router.UseMiddleware(router.DefaultHost, "/", router.MiddlewareHandlerFunc(limiter.LimitHTTP))

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

				// Add firewall as middleware if enabled
				if s.Cfg.Guard.Firewall.Enabled {
					s.Router.UseMiddleware(router.DefaultHost, path, router.MiddlewareHandlerFunc(firewall.BlockHTTP))
				}

				// Add limiter as middleware
				s.Router.UseMiddleware(router.DefaultHost, path, router.MiddlewareHandlerFunc(limiter.LimitHTTP))
			}
		}
	}

	if s.Cfg.Core.Metrics.Enabled {

		ba := guard.SimpleBasicAuth(s.Cfg.Core.Metrics.Users)
		ba.UnauthorizedHandler = s.HandleError

		s.Router.AddRoute(router.DefaultHost, s.Cfg.Core.Metrics.Path, false, "GET", "", s.handleMetrics())
		s.Router.UseMiddleware(router.DefaultHost, s.Cfg.Core.Metrics.Path, router.MiddlewareHandlerFunc(ba.Authenticate))
	}
}
