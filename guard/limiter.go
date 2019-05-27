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
	"net/http"

	"github.com/redmaner/MaguroHTTP/cache"
	"github.com/redmaner/MaguroHTTP/router"
	"golang.org/x/time/rate"
)

// Limiter is a type containing a MaguroHTTP Guard limiter
type Limiter struct {
	cache        *cache.SpearCache
	RatePerSec   rate.Limit
	RateBurst    int
	ErrorHandler router.ErrorHandler
	FilterOnIP   bool
}

// NewLimiter returns a new guard.Limiter
func NewLimiter(ratePerMin float64, rateBurst int, filterIP bool) *Limiter {
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
		FilterOnIP: filterIP,
	}
}

// LimitHTTP is a HTTP middleware function that can be used to add rate limiting
// to HTTP handlers.
func (l *Limiter) LimitHTTP(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var limit *rate.Limiter

		remoteAddr := router.StripHostPort(r.RemoteAddr)
		if !l.FilterOnIP {
			remoteAddr = remoteAddr + r.Header.Get("User-Agent")
		}
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
