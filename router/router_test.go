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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type test struct {
	requestHost  string
	host         string
	method       string
	path         string
	pathFallback bool
	content      string
}

var routes = []test{
	{
		host:         "127.0.0.1",
		method:       "GET",
		path:         "/",
		pathFallback: false,
		content:      ";",
	},
	{
		host:         DefaultHost,
		method:       "POST",
		path:         "/v1/report",
		pathFallback: false,
		content:      "*",
	},
	{
		host:         "127.0.0.1",
		method:       "POST",
		path:         "/v1/report/",
		pathFallback: false,
		content:      "application/json;text/html",
	},
	{
		host:         DefaultHost,
		method:       "GET",
		path:         "/specificA",
		pathFallback: false,
		content:      "",
	},
	{
		host:         "127.0.0.1",
		method:       "GET",
		path:         "/specificB",
		pathFallback: false,
		content:      "",
	},
	{
		host:         DefaultHost,
		method:       "GET",
		path:         "/v1/report/specificA",
		pathFallback: false,
		content:      "",
	},
	{
		host:         "127.0.0.1",
		method:       "GET",
		path:         "/v1/report/specificB",
		pathFallback: false,
		content:      "",
	},
}

var tests = []test{

	// Tests that test different scenario's for route matching.
	{
		requestHost:  DefaultHost,
		host:         DefaultHost,
		method:       "GET",
		path:         "/specificA",
		pathFallback: false,
		content:      "",
	},
	{
		requestHost:  "127.0.0.1",
		host:         "127.0.0.1",
		method:       "GET",
		path:         "/specificB",
		pathFallback: false,
		content:      "",
	},
	{
		requestHost:  "127.0.0.1",
		host:         DefaultHost,
		method:       "GET",
		path:         "/v1/report/specificA",
		pathFallback: false,
		content:      "",
	},
	{
		requestHost:  "127.0.0.1",
		host:         "127.0.0.1",
		method:       "GET",
		path:         "/v1/report/specificB",
		pathFallback: false,
		content:      "",
	},

	// Test no fallback
	{
		requestHost:  "127.0.0.1",
		host:         "127.0.0.1",
		method:       "POST",
		path:         "/v1/report/specificC",
		pathFallback: false,
		content:      "application/json",
	},
	// Test wrong method
	{
		requestHost:  "127.0.0.1",
		host:         "127.0.0.1",
		method:       "HEAD",
		path:         "/v1/report/specificC",
		pathFallback: false,
		content:      "application/json",
	},
	// Test wrong content
	{
		requestHost:  "127.0.0.1",
		host:         "127.0.0.1",
		method:       "POST",
		path:         "/v1/report/specificC",
		pathFallback: false,
		content:      `application/java; charset="UTF-8"`,
	},
	// Test wrong content
	{
		requestHost:  DefaultHost,
		host:         DefaultHost,
		method:       "POST",
		path:         "/v1/report/",
		pathFallback: false,
		content:      `application/java; charset="UTF-8"`,
	},
}

func TestRouter(t *testing.T) {

	fmt.Printf("\nTEST FOR SINGLE REQUEST\n")

	ro := NewRouter()

	for _, rt := range routes {
		ro.AddRoute(rt.host, rt.path, rt.pathFallback, rt.method, rt.content, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(test.method, test.path, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err)
		}
		req.Header.Set("Content-Type", test.content)
		req.Host = test.requestHost
		ro.ServeHTTP(w, req)
		fmt.Printf("Ran test with method %s for path %s with Content-Type %s. Result was %d\n", test.method, test.path, test.content, w.Code)
	}
}

func TestRouterCollision(t *testing.T) {
	ro := NewRouter()
	ro.AddRoute("127.0.0.1", "/", false, "GET", "*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ro.AddRoute(DefaultHost, "/Test", false, "GET", "*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/Test", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	req.Host = "127.0.0.1"
	req.Header.Set("Content-Type", `text/html; charset="UTF-8"`)
	ro.ServeHTTP(w, req)
}

func TestWebDAVMethods(t *testing.T) {
	ro := NewRouter()
	ro.WebDAV = true
	ro.AddRoute("127.0.0.1", "/", false, "PROPFIND", "*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
}
