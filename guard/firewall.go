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

package guard

import (
	"net/http"
	"path"

	"github.com/redmaner/MicroHTTP/router"
)

// Firewall is type that holds multiple middleware functions to add a firewall
// to http handlers
type Firewall struct {
	Blacklisting bool
	Subpath      bool
	Rules        map[string][]string
	ErrorHandler router.ErrorHandler
}

// NewFirewall returns a *Firewall type
func NewFirewall() *Firewall {
	return &Firewall{
		Blacklisting: true,
		ErrorHandler: router.ErrorHandler(func(w http.ResponseWriter, r *http.Request, code int) {
			switch code {
			case 403:
				http.Error(w, "Forbidden", 403)
			}
		}),
	}
}

// BlockHTTP is a middleware function to add a firewall to HTTP handlers
func (f *Firewall) BlockHTTP(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := r.URL.Path
		host := router.StripHostPort(r.RemoteAddr)

		for pt := p; pt != "/"; pt = path.Dir(pt) {
			if val, ok := f.Rules[pt]; ok {
				for _, v := range val {
					if v == host || v == "*" {
						if f.Blacklisting {
							f.ErrorHandler(w, r, 403)
							return
						}
						handler.ServeHTTP(w, r)
						return
					}
				}
			}
		}

		// The firewall subpath element allows blocking on specific subpaths of a website
		// This is only when you want to be extremely specific when configuring the firewall.
		// Subpath blocking is disabled by default and can be enabled in the configuration.
		if val, ok := f.Rules["/"]; ok && p == "/" || ok && !f.Subpath {
			for _, v := range val {
				if v == host || v == "*" {
					if f.Blacklisting {
						f.ErrorHandler(w, r, 403)
						return
					}
					handler.ServeHTTP(w, r)
					return
				}
			}
		}

		if f.Blacklisting {
			f.ErrorHandler(w, r, 403)
			return
		}
		handler.ServeHTTP(w, r)
		return
	}
}

// BlockProxy is a middleware function to add a firewall to HTTP proxy
func (f *Firewall) BlockProxy(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := r.URL.Path
		host := router.StripHostPort(r.RemoteAddr)

		for pt := p; pt != "/"; pt = path.Dir(pt) {
			if val, ok := f.Rules[pt]; ok {
				for _, v := range val {
					if v == host || v == "*" {
						if f.Blacklisting {
							f.ErrorHandler(w, r, 403)
							return
						}
						handler.ServeHTTP(w, r)
						return
					}
				}
			}
		}

		// The firewall subpath element allows blocking on specific subpaths of a website
		// This is only when you want to be extremely specific when configuring the firewall.
		// Subpath blocking is disabled by default and can be enabled in the configuration.
		if val, ok := f.Rules["/"]; ok && p == "/" || ok && !f.Subpath {
			for _, v := range val {
				if v == host || v == "*" {
					if f.Blacklisting {
						f.ErrorHandler(w, r, 403)
						return
					}
					handler.ServeHTTP(w, r)
					return
				}
			}
		}

		if f.Blacklisting {
			handler.ServeHTTP(w, r)
			return
		}

		f.ErrorHandler(w, r, 403)
		return
	}
}
