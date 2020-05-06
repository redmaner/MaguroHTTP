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

package router

import (
	"net/http"
	"path"
	"strings"
)

type pathRoute struct {
	subRoutes  map[string]methodRoute
	middleware []Middleware
}

type methodRoute struct {
	handler      http.Handler
	host         string
	path         string
	pathFallback bool
	method       string
	content      string
}

// Function to retrieve a methodRoute for a HTTP request
func (sr *SRouter) getRoute(host, path, method, contentType string) (methodRoute, []Middleware, int) {

	// For concurrency safety, lock mutex
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	// Define an empty pathRoute.
	var pathRouteMatch pathRoute
	var exactMatch bool

	// Match a route for host and DefaultHost
	hostFound, hostExact, hostMatch := sr.matchRoute(host, path)
	defaultFound, defaultExact, defaultMatch := sr.matchRoute(DefaultHost, path)

	// We check which route we use, by checking several cases.
	// To increase readability of the code all cases have been separated.
	switch {

	// case 1: We found a route for host, but not for DefaultHost
	// Selected route: host
	case hostFound && !defaultFound:
		pathRouteMatch = hostMatch
		exactMatch = hostExact

	// case 2: We found a route for DefaultHost,  but not for Host
	// Selected route: DefaultHost
	case defaultFound && !hostFound:
		pathRouteMatch = defaultMatch
		exactMatch = defaultExact

	// case 3: We found a route for both hosts, but host route was an exact match
	// Selected route: Host
	case defaultFound && !defaultExact && hostFound && hostExact:
		pathRouteMatch = hostMatch
		exactMatch = hostExact

	// case 4: We found a route for both hosts, but DefaultHost route was an exact match
	// Selected route: DefaultHost
	case defaultFound && defaultExact && hostFound && !hostExact:
		pathRouteMatch = defaultMatch
		exactMatch = defaultExact

	// case 5: we found a route for both hosts, and both were an exact match
	// Selected route: Host
	case defaultFound && defaultExact && hostFound && hostExact:
		pathRouteMatch = hostMatch
		exactMatch = hostExact

	// case 6: we found a route for both hosts, and both were not an exact match
	// Selected route: Host
	case defaultFound && !defaultExact && hostFound && !hostExact:
		pathRouteMatch = hostMatch
		exactMatch = hostExact

	// case 7: we didn't found any route
	// Selected route: none, we return a 404 HTTP Not Found error
	default:
		return methodRoute{}, []Middleware{}, 404
	}

	// We have found a pathRoute. We now search for a methodRoute that matches the
	// method of the request.
	methodRouteMatch, ok := pathRouteMatch.subRoutes[method]

	// We have found a route with matching host and path, but the method wasn't found.
	// we return an empty method route with a 405 Method  not allowed status code.
	if !ok {
		return methodRoute{}, []Middleware{}, 405
	}

	// We have found a route with matching host, path and method. The request
	// Content-Type is not allowed. We return an empty method route with a
	// 406 Media not allowed status code.
	if !methodRouteMatch.contentAllowed(contentType) {
		return methodRoute{}, []Middleware{}, 406
	}

	// All criteria have matched: host, path, method and Content-Type. If the
	// pathRoute wasn't an exact match we now determine if fallback is allowed to
	// the subpath. We do this on this level because we allow different fallback rules
	// for each methodRoute
	if !exactMatch && !methodRouteMatch.pathFallback {
		return methodRoute{}, []Middleware{}, 404
	}

	// We got a winner, return the found methodRoute with a 200 OK status code
	return methodRouteMatch, pathRouteMatch.middleware, 200

}

// Function to match a route for a given host + path combination
func (sr *SRouter) matchRoute(host, urlPath string) (bool, bool, pathRoute) {

	var pathRouteMatch pathRoute

	// Set exact match to false
	var exactMatch bool
	var match bool

	for {

		// Search for an exact host+path match
		if route, ok := sr.routes[host+urlPath]; ok {
			match = true
			exactMatch = true
			pathRouteMatch = route
			break
		}

		// We haven't found an exact match, so we search for a subpath. For example:
		// request has path /foo/bar. The path /foo/bar doesn't exist. So we search if
		// the path /foo exists. We don't bother with fallback allowance just yet.
		for pa := urlPath; pa != "/"; pa = path.Dir(pa) {
			if route, ok := sr.routes[host+pa]; ok {
				match = true
				exactMatch = false
				pathRouteMatch = route
				break
			}
		}

		if match {
			break
		}

		// We haven't found a sub path, so as a last resort we check if a root path exists.
		if route, ok := sr.routes[host+"/"]; ok {
			match = true
			exactMatch = false
			pathRouteMatch = route
			break
		}

		// We haven't found anything so we break
		break
	}

	return match, exactMatch, pathRouteMatch
}

func (mr *methodRoute) contentAllowed(contentType string) bool {

	if mr.content == "*" {
		return true
	}

	if strings.IndexByte(contentType, ';') > -1 {
		contentType = strings.Split(contentType, ";")[0]
	}

	if strings.IndexByte(mr.content, ';') > -1 {
		for _, v := range strings.Split(mr.content, ";") {
			if v == contentType {
				return true
			}
		}
	}
	return mr.content == contentType
}
