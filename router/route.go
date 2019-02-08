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

import (
	"net/http"
	"path"
	"strings"
)

type pathRoute struct {
	subRoutes map[string]methodRoute
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
func (sr *SRouter) getRoute(h, p, m, c string) (methodRoute, int) {

	// For concurrency safety, lock mutex
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	// Define an empty pathRoute.
	var pr pathRoute
	var em bool

	// Match a route for host and DefaultHost
	hostFound, hostExact, hostMatch := sr.matchRoute(h, p)
	defaultFound, defaultExact, defaultMatch := sr.matchRoute(DefaultHost, p)

	// We check which route we use, by checking several cases.
	// To increase readability of the code all cases have been separated.
	switch {

	// case 1: We found a route for host, but not for DefaultHost
	// Selected route: host
	case hostFound && !defaultFound:
		pr = hostMatch
		em = hostExact

	// case 2: We found a route for DefaultHost,  but not for Host
	// Selected route: DefaultHost
	case defaultFound && !hostFound:
		pr = defaultMatch
		em = defaultExact

	// case 3: We found a route for both hosts, but host route was an exact match
	// Selected route: Host
	case defaultFound && !defaultExact && hostFound && hostExact:
		pr = hostMatch
		em = hostExact

	// case 4: We found a route for both hosts, but DefaultHost route was an exact match
	// Selected route: DefaultHost
	case defaultFound && defaultExact && hostFound && !hostExact:
		pr = defaultMatch
		em = defaultExact

	// case 5: we found a route for both hosts, and both were an exact match
	// Selected route: Host
	case defaultFound && defaultExact && hostFound && hostExact:
		pr = hostMatch
		em = hostExact

	// case 6: we found a route for both hosts, and both were not an exact match
	// Selected route: Host
	case defaultFound && !defaultExact && hostFound && !hostExact:
		pr = hostMatch
		em = hostExact

	// case 7: we didn't found any route
	// Selected route: none, we return a 404 HTTP Not Found error
	default:
		return methodRoute{}, 404
	}

	// We have found a pathRoute. We now search for a methodRoute that matches the
	// method of the request.
	if mr, ok := pr.subRoutes[m]; ok {

		// We have found a route that matches the host, path and method. We now determine
		// if the request Content-Type matches the route.
		if mr.contentAllowed(c) {

			// All criteria have matched: host, path, method and Content-Type. If the
			// pathRoute wasn't an exact match we now determine if fallback is allowed to
			// the subpath. We do this on this level because we allow different fallback rules
			// for each methodRoute
			if !em && !mr.pathFallback {
				return methodRoute{}, 404
			}

			// We got a winner, return the found methodRoute with a 200 OK status code
			return mr, 200
		}

		// We have found a route with matching host, path and method. The request
		// Content-Type is not allowed. We return an empty method route with a
		// 406 Media not allowed status code.
		return methodRoute{}, 406
	}

	// We have found a route with matching host and path, but the method wasn't found.
	// we return an empty method route with a 405 Method  not allowed status code.
	return methodRoute{}, 405
}

// Function to match a route for a given host + path combination
func (sr *SRouter) matchRoute(host, p string) (bool, bool, pathRoute) {

	var pr pathRoute

	// Set exact match to false
	var em bool
	var match bool

	for {

		// Search for an exact host+path match
		if rt, ok := sr.routes[host+p]; ok {
			match = true
			em = true
			pr = rt
			break
		}

		// We haven't found an exact match, so we search for a subpath. For example:
		// request has path /foo/bar. The path /foo/bar doesn't exist. So we search if
		// the path /foo exists. We don't bother with fallback allowance just yet.
		for pa := p; pa != "/"; pa = path.Dir(pa) {
			if rt, ok := sr.routes[host+pa]; ok {
				match = true
				em = false
				pr = rt
				break
			}
		}

		if match {
			break
		}

		// We haven't found a sub path, so as a last resort we check if a root path exists.
		if rt, ok := sr.routes[host+"/"]; ok {
			match = true
			em = false
			pr = rt
			break
		}

		// We haven't found anything so we break
		break
	}

	return match, em, pr
}

func (mr *methodRoute) contentAllowed(c string) bool {

	if mr.content == "*" {
		return true
	}

	if strings.IndexByte(c, ';') > -1 {
		c = strings.Split(c, ";")[0]
	}

	if strings.IndexByte(mr.content, ';') > -1 {
		for _, v := range strings.Split(mr.content, ";") {
			if v == c {
				return true
			}
		}
	}
	return mr.content == c
}
