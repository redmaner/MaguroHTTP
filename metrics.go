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

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Metrics type, part of MicroHTTP config
type metrics struct {
	Enabled  bool
	Address  string
	Path     string
	User     string
	Password string
	Out      string
}

// Metrics handler for root "/", takes a GET request only
// Metrics will show the login page to access the metrics. Login is a username password configuration
// that can be set in the configuration. Password and username is authenticated / validated using JWT token
func (m *micro) httpMetricsRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		m.httpSetHeaders(w, m.config.Headers)
		w.Header().Set("Content-Security-Policy", "")
		io.WriteString(w, htmlStart)
		io.WriteString(w, metricsHTMLLogin(m.config.Metrics.Path+"/retrieve"))
		io.WriteString(w, htmlEnd)
	}
}

// Metrics handler for admin "/admin", takes GET request with "application/json" Content-Type
// This handler validates the body. It expects a JWT token which is validated against the
// username and password that are set in the configuration. If the token cannot be validated a HTTP 404 error is returned.
func (m *micro) httpMetricsAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bb, err := ioutil.ReadAll(r.Body); err == nil {
			if ok, err := jwtValidateToken(string(bb), m.config.Metrics.Password, m.config.Metrics.User, "MicroMetrics"); ok && err == nil {
				w.Header().Set("Content-Type", "text/html")
				m.httpSetHeaders(w, m.config.Headers)
				w.Header().Set("Content-Security-Policy", "")
				io.WriteString(w, htmlStart)
				m.md.display(w)
				io.WriteString(w, htmlEnd)
				logNetwork(200, r)
			} else {
				logAction(logERROR, err)
				m.httpError(w, r, 404)
			}
		}
	}
}

// Metrics handler for retrieve "/retrieve", takes POST request
// This takes the username and password from the form of the metrics index page
// and adds it in a JWT token. The token is signed with the given password.
// This JWT token is passed to the admin page. The response is shown back to the user.
func (m *micro) httpMetricsRetrieve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		password := r.FormValue("password")

		// Retrieve JWT token
		if token, err := jwtSignToken(password, user, "MicroMetrics", 30*time.Second); err == nil {

			var addr string
			if m.config.TLS.Enabled {
				addr = "https://" + m.config.Metrics.Address + ":" + m.config.Port + m.config.Metrics.Path + "/admin"
			} else {
				addr = "http://" + m.config.Metrics.Address + ":" + m.config.Port + m.config.Metrics.Path + "/admin"
			}
			if req, err2 := http.NewRequest("GET", addr, bytes.NewReader(token)); err2 == nil {
				req.Header.Set("Content-Type", "application/json")
				if resp, err3 := m.client.Do(req); err3 == nil {
					w.WriteHeader(resp.StatusCode)
					w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
					w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
					w.Header().Set("Content-Security-Policy", "")
					io.Copy(w, resp.Body)
					resp.Body.Close()
					logNetwork(resp.StatusCode, r)
				} else {
					logAction(logERROR, err)
					m.httpError(w, r, 404)
				}
			} else {
				logAction(logERROR, err)
				m.httpError(w, r, 404)
			}
		} else {
			logAction(logERROR, err)
			m.httpError(w, r, 404)
		}
	}
}

// Function to retrieve the index metrics page
func metricsHTMLLogin(p string) string {
	return fmt.Sprintf(`<h1>MicroHTTP Metrics login</h1>
	<p><form action="%s" method="POST" accept-charset="UTF-8">
		<input type="text" name="user" placeholder="Username" autofocus autocomplete="off"><br><br>
		<input type="password" name="password" placeholder="Password" autofocus autocomplete="off"><br>
		<input type="submit" name="action" value="Login"><br>
	</form><br><br></p>
	`, p)
}
