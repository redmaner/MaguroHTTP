package guard

import (
	"net/http"

	"github.com/redmaner/MicroHTTP/router"
	"golang.org/x/time/rate"
)

// GuardHTTP is a HTTP middleware function that can be used to add rate limiting
// to HTTP handlers.
func (l *Limiter) GuardHTTP(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var limit *rate.Limiter

		remoteAddr := router.StripHostPort(r.RemoteAddr)
		ok, lmt := l.cache.Get(remoteAddr, 900000000000)

		switch {
		case ok:
			if as, ok := lmt.(*rate.Limiter); ok {
				limit = as
			}
		default:
			limit = rate.NewLimiter(l.RatePerSec, 10)
		}

		l.cache.Set(remoteAddr, limit)

		if !limit.Allow() {
			l.ErrorHandler(w, r, 429)
			return
		}

		h.ServeHTTP(w, r)

	}
}
