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

package router

import "net/http"

// Middleware is an inteface used to execute middleware functions by the router
type Middleware interface {
	MiddlewareHTTP(http.Handler) http.Handler
}

// MiddlewareHandler is a function that implements the Middleware interface
// using http.Handler
type MiddlewareHandler func(handler http.Handler) http.Handler

// MiddlewareHTTP implements the Middleware interface for MiddlewareHandler
func (mh MiddlewareHandler) MiddlewareHTTP(handler http.Handler) http.Handler {
	return mh(handler)
}

// MiddlewareHandlerFunc is a function that implements the Middleware interface
// using http.HandlerFunc
type MiddlewareHandlerFunc func(handler http.HandlerFunc) http.HandlerFunc

// MiddlewareHTTP implements the Middleware interface for MiddlewareHandler
func (mhf MiddlewareHandlerFunc) MiddlewareHTTP(handler http.Handler) http.Handler {
	return mhf(handler.ServeHTTP)
}

func (sr *SRouter) UseMiddleware(host, path string, handler Middleware) {

	sr.mu.Lock()
	defer sr.mu.Unlock()

	// We don't want empty parameters
	if host == "" || path == "" {
		panic("smux: found illegal blank parameters")
	}

	// If the path is not the root, we don't want paths ending with a "/"
	// smux.SRouter uses pathFallback to configure fallback, and no weird slashes like http.ServeMux
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Handler cannot be nil. This is rare, but we check anyway.
	if handler == nil {
		panic("smux: nil handler")
	}

	if sr.routes == nil {
		sr.routes = make(map[string]pathRoute)
	}

	if _, ok := sr.routes[host+path]; !ok {
		sr.routes[host+path] = pathRoute{
			subRoutes:  make(map[string]methodRoute),
			middleware: []Middleware{},
		}
	}

	pr := sr.routes[host+path]
	pr.middleware = append(pr.middleware, handler)
	sr.routes[host+path] = pr

}
