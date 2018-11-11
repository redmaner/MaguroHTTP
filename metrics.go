package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func (m *micro) httpMetrics(w http.ResponseWriter, r *http.Request) {

	remote := httpTrimPort(r.RemoteAddr)

	// Validate request Content-Type
	rct := r.Header.Get("Content-Type")
	act := contentTypes{
		RequestTypes: []string{"", "text/html", "application/x-www-form-urlencoded", "application/json"},
	}
	if !httpValidateRequestContentType(&rct, &act) {
		m.httpError(w, r, 406)
		return
	}

	path := r.URL.Path

	// Check firewall for path
	if block := firewallHTTP(&m.config, remote, path); block {
		m.httpError(w, r, 403)
		return
	}

	switch r.Method {
	case "GET":
		switch path {
		case m.config.Metrics.Path + "/":
			w.Header().Set("Content-Type", "text/html")
			m.httpSetHeaders(w, m.config.Headers)
			w.Header().Set("Content-Security-Policy", "")
			io.WriteString(w, htmlStart)
			io.WriteString(w, metricsHtmlLogin(m.config.Metrics.Path+"/retrieve"))
			io.WriteString(w, htmlEnd)
		case m.config.Metrics.Path + "/admin":

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
		case m.config.Metrics.Path + "/retrieve":
			user := r.FormValue("user")
			password := r.FormValue("password")

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

func metricsHtmlLogin(p string) string {
	return fmt.Sprintf(`<h1>MicroHTTP Metrics login</h1>
	<p><form action="%s" method="POST" accept-charset="UTF-8">
		<input type="text" name="user" placeholder="Username" autofocus autocomplete="off"><br><br>
		<input type="text" name="password" placeholder="Password" autofocus autocomplete="off"><br>
		<input type="submit" name="action" value="Login"><br>
	</form><br><br></p>
	`, p)
}
