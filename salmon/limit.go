package salmon

import "golang.org/x/time/rate"

func (ln *Listener) allowConnection(ip string) bool {

	var limit *rate.Limiter
	ok, lmt := ln.cache.Get(ip, 900000000000)
	switch {
	case ok:
		if as, ok := lmt.(*rate.Limiter); ok {
			limit = as
		}
	default:
		limit = rate.NewLimiter(ln.maxConns, ln.maxConnsBurst)
	}
	ln.cache.Set(ip, limit)

	return limit.Allow()
}
