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

package guard

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/html"
	"golang.org/x/crypto/bcrypt"
)

const templateLogin = `
<form method="POST" action="%s">
	<input type="text" name="username" placeholder="Username" autofocus autocomplete="off"/>
	<input type="password" name="password" placeholder="Password" autofocus autocomplete="off"/>
	<input type="submit" name= "action" value="Login" class="primary"/>
</form>
`

// Authorizer is a type that provides HTTP middleware to add authentication
// and authorization to HTTP handlers.
type Authorizer struct {
	lock sync.Mutex

	Users    map[string]User
	Sessions map[string]string
	TLS      bool

	RedirectAuth  string
	RedirectLogin string
	RedirectRoot  string
	LogInstance   *debug.Logger

	LoginTemplate *html.TemplateHandler
}

// User holds a username and password, used by Authorizer type
type User struct {
	Username string
	Password []byte
}

type loginTemplate struct {
	LoginPane  template.HTML
	LoginError template.HTML
}

// Log is a function to log messages to debug.Logger instance
func (a *Authorizer) Log(level int, err error) {
	a.LogInstance.Log(level, err)
}

// Auth is a middleware parser, which handles authorization of users using session id / cookies
func (a *Authorizer) Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Make sure we aren't already on the login page
		path := r.URL.Path

		// concurrency safety
		a.lock.Lock()
		defer a.lock.Unlock()

		// Check for a cookie
		ck, err := r.Cookie("session-id")
		switch {

		case err == http.ErrNoCookie:
			if path == a.RedirectLogin {
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, a.RedirectLogin, 302)
			return

		default:

			if len(a.Sessions) == 0 || ck == nil {
				if path == a.RedirectLogin {
					handler.ServeHTTP(w, r)
					return
				}
				http.Redirect(w, r, a.RedirectLogin, 302)
				return

			}

			// We found a cookie so we check if the session exists
			if v, ok := a.Sessions[ck.Value]; ok {

				// We check if the session is of a user
				if _, exists := a.Users[v]; exists {
					if path == a.RedirectLogin {
						http.Redirect(w, r, a.RedirectRoot, 302)
						return
					}
					handler.ServeHTTP(w, r)
					return
				}
			}

			// The session no longer exists so we remove the cookie
			ck.MaxAge = -1
			http.SetCookie(w, ck)

			// We redirect back to login
			if path == a.RedirectLogin {
				handler.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, a.RedirectLogin, 302)
			return
		}
	}
}

// HandleLogin handles the login page using the login template
func (a *Authorizer) HandleLogin() http.HandlerFunc {
	a.LoginTemplate.Init()
	return func(w http.ResponseWriter, r *http.Request) {

		a.LoginTemplate.Execute(w, loginTemplate{
			LoginPane: a.templateData(),
		})
	}
}

// HandleAuth handles authentication it expects a POST request
func (a *Authorizer) HandleAuth() http.HandlerFunc {
	a.LoginTemplate.Init()
	return func(w http.ResponseWriter, r *http.Request) {

		// Concurrency safety
		a.lock.Lock()
		defer a.lock.Unlock()

		username := r.FormValue("username")
		password := r.FormValue("password")
		remember := r.FormValue("remember")

		// we check if the user is known
		us, exists := a.Users[username]

		if !exists {
			a.LoginTemplate.Execute(w, loginTemplate{
				LoginPane:  a.templateData(),
				LoginError: template.HTML("<br>* Username or password is invalid"),
			})
			return
		}

		// The user is known, so we check if the password matches
		// Passwords are hashed with bycrypt hashing algorithm
		err := bcrypt.CompareHashAndPassword(us.Password, []byte(password))
		if err != nil {
			a.LoginTemplate.Execute(w, loginTemplate{
				LoginPane:  a.templateData(),
				LoginError: template.HTML("<br>* Username or password is invalid"),
			})
		}

		// The password matches so we generate a session
		u, err := uuid.NewV4()
		a.Log(debug.LogError, err)

		// Session is added to sessions
		a.Sessions[u.String()] = username

		// Session is set to cookie
		ck := &http.Cookie{
			Name:     "session-id",
			Value:    u.String(),
			HttpOnly: true,
			Secure:   a.TLS,
		}

		if remember == "remember" {
			ck.MaxAge = 604800
		}

		http.SetCookie(w, ck)

		// We redirect login back to index
		http.Redirect(w, r, a.RedirectRoot, 302)
		return
	}
}

func (a *Authorizer) templateData() template.HTML {
	return template.HTML(fmt.Sprintf(templateLogin, a.RedirectAuth))
}
