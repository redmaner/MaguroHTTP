package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Function that is used for MicroHTTP metrics
func (m *micro) httpMetrics(w http.ResponseWriter, r *http.Request) {

	remote := httpTrimPort(r.RemoteAddr)

	path := r.URL.Path

	// Check firewall for path
	if block := firewallHTTP(&m.config, remote, path); block {
		m.httpError(w, r, 403)
		return
	}

	// Metrics only accepts GET and POST requests, other methods are rejected with
	// HTTP 405 code
	switch r.Method {
	case "GET":
		switch path {

		// This is the base path. Metrics will show the login page to access the metrics
		// Login is a username password configuration that can be set in the configuration.
		// Password and username is authenticated / validated using JWT token
		case m.config.Metrics.Path + "/":
			w.Header().Set("Content-Type", "text/html")
			m.httpSetHeaders(w, m.config.Headers)
			w.Header().Set("Content-Security-Policy", "")
			io.WriteString(w, htmlStart)
			io.WriteString(w, metricsHtmlLogin(m.config.Metrics.Path+"/retrieve"))
			io.WriteString(w, htmlEnd)

		// This is the admin path that validates the body. It expects a JWT token which
		// is validated against the username and password that are set in the configuration.
		// If the token cannot be validated a HTTP 404 error is returned.
		case m.config.Metrics.Path + "/admin":

			// Only accept application/json as Content-Type. Other Content-Types are rejected
			if r.Header.Get("Content-Type") != "application/json" {
				m.httpError(w, r, 405)
				return
			}

			// Convert request body to string
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
			} else {
				logAction(logERROR, err)
				m.httpError(w, r, 404)
			}
		default:
			m.httpError(w, r, 404)
			return
		}
	case "POST":
		switch path {

		// This is the retrieve path. This takes the username and password from the form
		// of the metrics index page and adds it in a JWT token. The token is signed with
		// the given password. This JWT token is passed to the admin page. The response
		// is shown back to the user.
		case m.config.Metrics.Path + "/retrieve":
			user := r.FormValue("user")
			password := r.FormValue("password")

			// Retrieve JWT token
			if token, err := jwtSignToken(password, user, "MicroMetrics", 30*time.Second); err == nil {

				var addr string
				if m.config.TLS {
					addr = "https://" + m.config.Metrics.Address + ":" + m.config.Port + m.config.Metrics.Path + "/admin"
				} else {
					addr = "http://" + m.config.Metrics.Address + ":" + m.config.Port + m.config.Metrics.Path + "/admin"
				}
				if req, err := http.NewRequest("GET", addr, bytes.NewReader(token)); err == nil {
					req.Header.Set("Content-Type", "application/json")
					if resp, err := http.DefaultClient.Do(req); err == nil {
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
		default:
			m.httpError(w, r, 404)
			return
		}
	default:
		m.httpError(w, r, 405)
		return
	}
}

// Function to retrieve the index metrics page
func metricsHtmlLogin(p string) string {
	return fmt.Sprintf(`<h1>MicroHTTP Metrics login</h1>
	<p><form action="%s" method="POST" accept-charset="UTF-8">
		<input type="text" name="user" placeholder="Username" autofocus autocomplete="off"><br><br>
		<input type="password" name="password" placeholder="Password" autofocus autocomplete="off"><br>
		<input type="submit" name="action" value="Login"><br>
	</form><br><br></p>
	`, p)
}
