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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func exampleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is an example function\n")
	}
}

func exampleMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is an example middleware function\n")
		handler.ServeHTTP(w, r)
	}
}

func exampleMiddlewareTwo(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is an example middleware function, number two\n")
		handler.ServeHTTP(w, r)
	}
}

func TestMiddleware(t *testing.T) {

	ro := NewRouter()

	ro.AddRoute(DefaultHost, "/", true, "GET", "", exampleHandler())
	ro.UseMiddleware(DefaultHost, "/", MiddlewareHandlerFunc(exampleMiddleware))
	ro.UseMiddleware(DefaultHost, "/", MiddlewareHandlerFunc(exampleMiddlewareTwo))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	ro.ServeHTTP(w, req)

	resp, _ := ioutil.ReadAll(w.Result().Body)
	fmt.Println(string(resp))

}
