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

package micro

import (
	"fmt"
	"net/http"

	"github.com/redmaner/MicroHTTP/debug"
)

// Log is a function to log messages using the debug.Logger type
func (s *Server) Log(logLevel int, err error) {
	s.logInterface.Log(logLevel, err)
}

// LogNetwork is a function to log network activity using the debug.Logger type
func (s *Server) LogNetwork(statusCode int, r *http.Request) {
	s.Log(debug.LogNet, fmt.Errorf("%d request=%s %s%s%s IP=%s User-Agent=%s", statusCode, r.Method, r.Host, r.URL.Path, r.URL.RawQuery, r.RemoteAddr, r.Header.Get("User-Agent")))
}
