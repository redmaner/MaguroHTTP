package guard

import (
	"net/http"

	"github.com/redmaner/MicroHTTP/cache"
	"github.com/redmaner/MicroHTTP/router"
	"golang.org/x/time/rate"
)

// Limiter is a type containing a MicroHTTP Guard limiter
type Limiter struct {
	cache        *cache.SpearCache
	RatePerSec   rate.Limit
	RateBurst    int
	ErrorHandler router.ErrorHandler
}

// NewLimiter returns a new guard.Limiter
func NewLimiter(ratePerMin float64, rateBurst int) *Limiter {
	return &Limiter{
		cache:      cache.NewCache(),
		RatePerSec: rate.Limit(ratePerMin / 60.00),
		RateBurst:  rateBurst,
		ErrorHandler: router.ErrorHandler(func(w http.ResponseWriter, r *http.Request, code int) {
			switch code {
			case 429:
				http.Error(w, "Too many requests", 429)
			}
		}),
	}
}

// LimitHTTP is a HTTP middleware function that can be used to add rate limiting
// to HTTP handlers.
func (l *Limiter) LimitHTTP(h http.HandlerFunc) http.HandlerFunc {
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
			limit = rate.NewLimiter(l.RatePerSec, l.RateBurst)
		}

		l.cache.Set(remoteAddr, limit)

		if !limit.Allow() {
			l.ErrorHandler(w, r, 429)
			return
		}

		h.ServeHTTP(w, r)
	}
}
