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

package micro

import (
	"strings"

	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/guard"
	"github.com/redmaner/MicroHTTP/html"
	"github.com/redmaner/MicroHTTP/router"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) addRoutesFromConfig() {

	var limiter *guard.Limiter
	var firewall *guard.Firewall

	// Make routes for each vhost, if vhosts are enabled
	if s.Cfg.Core.VirtualHosting {

		// Loop over each Vhost
		for vhost := range s.Cfg.Core.VirtualHosts {

			// Each virtual host gets it's own limiter
			limiter = guard.NewLimiter(s.Vhosts[vhost].Guard.Rate, s.Vhosts[vhost].Guard.RateBurst)
			limiter.ErrorHandler = s.handleError

			if s.Vhosts[vhost].Guard.Firewall.Enabled {
				firewall = &guard.Firewall{
					Blacklisting: s.Vhosts[vhost].Guard.Firewall.Blacklisting,
					Subpath:      s.Vhosts[vhost].Guard.Firewall.Subpath,
					Rules:        s.Vhosts[vhost].Guard.Firewall.Rules,
					ErrorHandler: s.handleError,
				}
			}

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

		limiter = guard.NewLimiter(s.Cfg.Guard.Rate, s.Cfg.Guard.RateBurst)
		limiter.ErrorHandler = s.handleError

		if s.Cfg.Guard.Firewall.Enabled {
			firewall = &guard.Firewall{
				Blacklisting: s.Cfg.Guard.Firewall.Blacklisting,
				Subpath:      s.Cfg.Guard.Firewall.Subpath,
				Rules:        s.Cfg.Guard.Firewall.Rules,
				ErrorHandler: s.handleError,
			}
		}

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

	if s.Cfg.Metrics.Enabled {

		// Make a set of users
		usrs := make(map[string]guard.User)

		for k, v := range s.Cfg.Metrics.Users {
			bcryptPass, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
			if err != nil {
				s.Log(debug.LogError, err)
				continue
			}

			usrs[k] = guard.User{
				Username: k,
				Password: bcryptPass,
			}
		}

		// Make a authorizer struct
		authr := guard.Authorizer{
			Users:    usrs,
			Sessions: make(map[string]string),
			TLS:      s.Cfg.Core.TLS.Enabled,

			RedirectAuth:  s.Cfg.Metrics.Path + "/auth",
			RedirectRoot:  s.Cfg.Metrics.Path,
			RedirectLogin: s.Cfg.Metrics.Path + "/login",

			LogInstance:   s.logInterface,
			LoginTemplate: html.NewTemplate(s.Cfg.Core.FileDir+"templates/", "login.html"),
		}

		s.Router.AddRoute(router.DefaultHost, s.Cfg.Metrics.Path+"/login", false, "GET", "", authr.HandleLogin())
		s.Router.AddRoute(router.DefaultHost, s.Cfg.Metrics.Path+"/auth", false, "POST", "application/x-www-form-urlencoded", authr.HandleAuth())
		s.Router.AddRoute(router.DefaultHost, s.Cfg.Metrics.Path, false, "GET", "", s.handleMetrics())
		s.Router.UseMiddleware(router.DefaultHost, s.Cfg.Metrics.Path+"/login", router.MiddlewareHandlerFunc(authr.Auth))
		s.Router.UseMiddleware(router.DefaultHost, s.Cfg.Metrics.Path, router.MiddlewareHandlerFunc(authr.Auth))
	}

}
