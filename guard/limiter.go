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
	ErrorHandler router.ErrorHandler
}

// NewLimiter returns a new guard.Limiter
func NewLimiter(ratePerMin float64) *Limiter {
	return &Limiter{
		cache:        cache.NewCache(),
		RatePerSec:   rate.Limit(ratePerMin / 60.00),
		ErrorHandler: handleError(),
	}
}

func handleError() router.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, code int) {
		switch code {
		case 429:
			http.Error(w, "Too many requests", 429)
		}
	}
}
