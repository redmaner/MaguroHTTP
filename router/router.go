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

// Package router provides a fast and simple, security orientated HTTP router for GO (golang)
package router

import (
	"log"
	"net/http"
	"sync"
)

const (
	// DefaultHost contains the default host the router uses to match a route to a HTTP request.
	DefaultHost = "DEFAULT"

	// DefaultFallback defines the default behavior of the router whether to fallback to a subpath of a request path.
	// Example: the request contains a request for /foo/bar, but /foo/bar is not registered as a route.
	// DefaultFallback determines whether the routers is allowed to fallback to /foo if it is registered.
	// pathFallback can be set for each individual method route.
	DefaultFallback = false
)

var (

	// Allowed methods for default HTTP
	allowedMethodsHTTP = map[string]struct{}{
		"GET": struct{}{}, "PUT": struct{}{}, "POST": struct{}{}, "DELETE": struct{}{},
		"PATCH": struct{}{}, "CONNECT": struct{}{}, "OPTIONS": struct{}{}, "HEAD": struct{}{},
	}

	// Allowed methods for WebDAV extension (disabled by default)
	allowedMethodsHTTPWebDAV = map[string]struct{}{
		"PROPFIND": struct{}{}, "PROPPATCH": struct{}{}, "MKCOL": struct{}{}, "COPY": struct{}{},
		"MOVE": struct{}{}, "LOCK": struct{}{}, "UNLOCK": struct{}{},
	}
)

// SRouter dispatches HTTP requests to a defined handler. This router implements
// the http.Handler and http.HandlerFunc to dispatch requests. Requests will be dispatched to defined routes.
// This router dispatches requests that match the request's: host, path, method, and Content-Type.
//
// Requests that don't match the host or path will receive a standard HTTP 404 error.
// Requests that don't match the method for the host or path will receive a standard HTTP 405 error.
// Requests that don't match the Content-Type for the host or path will receive a standard HTTP 406 error.
//
// By design router.SRouter requires explicit defition of routes. It does however support
// fallback to subpaths, if the request path cannot be found.
// Example: the request contains a request for /foo/bar, but /foo/bar is not registered as a route.
// router.SRouter will dispatch that request to /foo route if that route is registered and supports fallback.
//
type SRouter struct {
	mu     sync.RWMutex
	routes map[string]pathRoute

	// ErrorHandler allows to define a custom handler for errors. It takes ErrorHandler as type,
	// which implements the http.Error function (w http.ResponseWriter, error string, code int).
	ErrorHandler ErrorHandler

	// WebDAV. If enabled, the router allows WebDAV methods to be registered as routes
	WebDAV bool
}

// ErrorHandler is a type of func(w http.ResponseWriter, r *http.Request, code int) where code
// corresponds to a HTTP status code (200, 404, 405 etc.). By default SRouter uses the http.Error to return errors.
// Th ErrorHandler type allows users to define their own implementation of an ErrorHandler to handle HTTP errors from the Router.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, code int)

// NewRouter returns a default router.SRouter
func NewRouter() *SRouter {
	return &SRouter{
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, code int) {
			switch code {
			case 404:
				http.Error(w, "Not Found", 404)
			case 405:
				http.Error(w, "Method Not Allowed", 405)
			case 406:
				http.Error(w, "Media Not Supported", 406)
			}
		},
	}
}

// AddRoute can be used to add a route to the Router
// Parameters:
//
// 1. Host as string. You should use router.DefaultHost if you do not want to use a custom host.
//
// 2. path as string. Every path should start with a /
//
// 3. path fallback as bool.
//
// 4. method as string. This can only be one single HTTP method
//
// 5. Content-Type. This allows multiple entries, separated by semicolon. For example "text/html;application/json"
// An empty Content-Type can also be valid, use a single semicolon to do so. For example ";"
// Character sets are not checked by router.SRouter so you do not have to define these explicitly
//
// 6. handler of type http.Handler or http.HandlerFunc. The http.HandlerFunc does implement the http.Handler interface
// and can therefore be passed into AddRoute as well.
func (sr *SRouter) AddRoute(host, path string, fallback bool, method, content string, handler http.Handler) {

	sr.mu.Lock()
	defer sr.mu.Unlock()

	// We don't want empty parameters
	if host == "" || path == "" || method == "" {
		panic("router: found illegal blank parameters")
	}

	// If the path is not the root, we don't want paths ending with a "/"
	// router.SRouter uses pathFallback to configure fallback, and no weird slashes like http.ServeMux
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Test validity of HTTP methods
	_, methodHTTPAllowed := allowedMethodsHTTP[method]
	_, methodWebDAVAllowed := allowedMethodsHTTPWebDAV[method]

	switch {
	case sr.WebDAV:
		if !methodWebDAVAllowed && !methodHTTPAllowed {
			log.Fatalf("router: method %s is not allowed", method)
		}
	default:
		if !methodHTTPAllowed {
			log.Fatalf("router: method %s is not allowed", method)
		}
	}

	// Handler cannot be nil. This is rare, but we check anyway.
	if handler == nil {
		log.Fatalf("router: nil handler")
	}

	if sr.routes == nil {
		sr.routes = make(map[string]pathRoute)
	}

	pathRouteAdd := pathRoute{
		subRoutes:  make(map[string]methodRoute),
		middleware: []Middleware{},
	}

	if _, ok := sr.routes[host+path]; !ok {
		sr.routes[host+path] = pathRouteAdd
	}

	methodRouteAdd := methodRoute{
		handler:      handler,
		host:         host,
		path:         path,
		pathFallback: fallback,
		method:       method,
		content:      content,
	}

	sr.routes[host+path].subRoutes[method] = methodRouteAdd
}

func (sr *SRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	host := StripHostPort(r.Host)
	path := cleanPath(r.URL.Path)
	method := r.Method
	content := r.Header.Get("Content-Type")

	r.URL.Path = path

	methodRouteMatch, middleWare, status := sr.getRoute(host, path, method, content)

	if status != 200 {
		sr.ErrorHandler(w, r, status)
		return
	}

	// Handle middleware
	lenMw := len(middleWare)
	mwHandler := methodRouteMatch.handler

	switch {

	// We do some middleware execution magic right here
	case lenMw > 0:
		for i := lenMw - 1; i >= 0; i-- {
			if i == 0 {
				middleWare[i].MiddlewareHTTP(mwHandler).ServeHTTP(w, r)
				return
			}
			mwHandler = middleWare[i].MiddlewareHTTP(mwHandler)
		}

	// Request can be handled by handler, so dispatch to defined handler
	default:
		methodRouteMatch.handler.ServeHTTP(w, r)
	}
}
