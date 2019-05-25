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
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/redmaner/MicroHTTP/router"
)

// BasicAuth provides HTTP middleware for protecting URIs with HTTP Basic Authentication
// as per RFC 2617. The server authenticates a user:password combination provided in the
// "Authorization" HTTP header.
//
// Note: HTTP Basic Authentication credentials are sent in plain text, and therefore it does
// not make for a wholly secure authentication mechanism. You should serve your content over
// HTTPS to mitigate this, noting that "Basic Authentication" is meant to be just that: basic!
type BasicAuth struct {
	Realm               string
	Users               map[string]AuthUser
	AuthFunc            func(string, string, *http.Request) bool
	UnauthorizedHandler router.ErrorHandler
}

type AuthUser struct {
	User     string
	Password string
}

// Authenticate is a middleware function that adds HTTP BasicAuth to a http handler
func (b BasicAuth) Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check if we have a user-provided error handler, else set a default
		if b.UnauthorizedHandler == nil {
			b.UnauthorizedHandler = router.ErrorHandler(func(w http.ResponseWriter, r *http.Request, code int) {
				switch code {
				case 403:
					http.Error(w, "Forbidden", 403)
				}
			})
		}

		// Check that the provided details match
		if !b.authenticate(r) {
			b.requestAuth(w, r)
			return
		}

		// Call the next handler on success.
		handler.ServeHTTP(w, r)
	}
}

// authenticate retrieves and then validates the user:password combination provided in
// the request header. Returns 'false' if the user has not successfully authenticated.
func (b *BasicAuth) authenticate(r *http.Request) bool {
	const basicScheme string = "Basic "

	if r == nil {
		return false
	}

	// In simple mode, prevent authentication with empty credentials if User is
	// not set. Allow empty passwords to support non-password use-cases.
	if b.AuthFunc == nil && len(b.Users) == 0 {
		return false
	}

	// Confirm the request is sending Basic Authentication credentials.
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		return false
	}

	// Get the plain-text username and password from the request.
	// The first six characters are skipped - e.g. "Basic ".
	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return false
	}

	// Split on the first ":" character only, with any subsequent colons assumed to be part
	// of the password. Note that the RFC2617 standard does not place any limitations on
	// allowable characters in the password.
	creds := bytes.SplitN(str, []byte(":"), 2)

	if len(creds) != 2 {
		return false
	}

	givenUser := string(creds[0])
	givenPass := string(creds[1])

	// Default to Simple mode if no AuthFunc is defined.
	if b.AuthFunc == nil {
		b.AuthFunc = b.simpleBasicAuthFunc
	}

	return b.AuthFunc(givenUser, givenPass, r)
}

// simpleBasicAuthFunc authenticates the supplied username and password against
// the User and Password set in the Options struct.
func (b *BasicAuth) simpleBasicAuthFunc(user, pass string, r *http.Request) bool {

	// Equalize lengths of supplied and required credentials by hashing them
	givenUser := sha256.Sum256([]byte(user))
	givenPass := sha256.Sum256([]byte(pass))
	requiredUser := sha256.Sum256([]byte(b.Users[user].User))
	requiredPass := sha256.Sum256([]byte(b.Users[user].Password))

	// Compare the supplied credentials to those set in our options
	if subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 &&
		subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
		return true
	}

	return false
}

// Require authentication, and serve our error handler otherwise.
func (b *BasicAuth) requestAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, b.Realm))
	b.UnauthorizedHandler(w, r, 401)
}

func SimpleBasicAuth(users map[string]string) *BasicAuth {

	aUsers := make(map[string]AuthUser)
	for k, v := range users {
		aUsers[k] = AuthUser{
			User:     k,
			Password: v,
		}
	}

	return &BasicAuth{
		Realm: "Restricted",
		Users: aUsers,
	}

}
