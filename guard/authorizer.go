package guard

import (
	"html/template"
	"net/http"
	"sync"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/html"
	"golang.org/x/crypto/bcrypt"
)

type Authorizer struct {
	lock sync.Mutex

	Users    map[string]User
	sessions map[string]string
	TLS      bool

	RedirectLogin string
	RedirectRoot  string
	LogInstance   *debug.Logger

	LoginTemplate *html.TemplateHandler
}

type User struct {
	Username string
	Password []byte
}

func (a *Authorizer) Log(level int, err error) {
	a.LogInstance.Log(level, err)
}

// Auth is a middleware parser, which handles authorization of users using session id / cookies
func (a *Authorizer) Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// concurrency safety
		a.lock.Lock()
		defer a.lock.Unlock()

		// Check for a cookie
		ck, err := r.Cookie("session-id")
		switch {

		case err == http.ErrNoCookie:
			http.Redirect(w, r, a.RedirectLogin, 302)
			return

		default:

			if len(a.sessions) == 0 || ck == nil {
				return
			}

			// We found a cookie so we check if the session exists
			if v, ok := a.sessions[ck.Value]; ok {

				// We check if the session is of a user
				if _, exists := a.Users[v]; exists {
					handler.ServeHTTP(w, r)
					return
				}
			}

			// The session no longer exists so we remove the cookie
			ck.MaxAge = -1
			http.SetCookie(w, ck)

			// We redirect back to login
			http.Redirect(w, r, a.RedirectLogin, 302)
			return
		}
	}
}

// HandleLogin handles the login page using the login template
func (a *Authorizer) HandleLogin() http.HandlerFunc {
	a.LoginTemplate.Init()
	return func(w http.ResponseWriter, r *http.Request) {

		a.LoginTemplate.Execute(w, nil)
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
			a.LoginTemplate.Execute(w, template.HTML("<br>* Gebruikernaam of wachtwoord onjuist"))
			return
		}

		// The user is known, so we check if the password matches
		// Passwords are hashed with bycrypt hashing algorithm
		err := bcrypt.CompareHashAndPassword(us.Password, []byte(password))
		if err != nil {
			a.LoginTemplate.Execute(w, template.HTML("<br>* Gebruikernaam of wachtwoord onjuist"))
		}

		// The password matches so we generate a session
		u, err := uuid.NewV4()
		a.Log(debug.LogError, err)

		// Session is added to sessions
		a.sessions[u.String()] = username

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
