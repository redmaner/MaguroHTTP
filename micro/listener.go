package micro

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/redmaner/MicroHTTP/cache"
	"github.com/redmaner/MicroHTTP/router"
	"golang.org/x/time/rate"
)

type TCPSecListener struct {
	listener      *net.TCPListener
	cache         *cache.SpearCache
	maxConns      rate.Limit
	maxConnsBurst int
}

func NewTCPSecListener(addr string, maxConnsPerMin float32, maxConnsBurst int) *TCPSecListener {
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}
	return &TCPSecListener{
		listener:      ln.(*net.TCPListener),
		cache:         cache.NewCache(),
		maxConns:      rate.Limit(maxConnsPerMin / 60),
		maxConnsBurst: maxConnsBurst,
	}
}

func (ln *TCPSecListener) Accept() (net.Conn, error) {

	tc, err := ln.listener.AcceptTCP()

	var limit *rate.Limiter

	if err != nil {
		return nil, err
	}

	remoteAddr := router.StripHostPort(tc.RemoteAddr().String())
	ok, lmt := ln.cache.Get(remoteAddr, 900000000000)

	switch {
	case ok:
		if as, ok := lmt.(*rate.Limiter); ok {
			limit = as
		}
	default:
		limit = rate.NewLimiter(ln.maxConns, ln.maxConnsBurst)
	}

	ln.cache.Set(remoteAddr, limit)

	if !limit.Allow() {
		tc.Close()
		return nil, fmt.Errorf("Host %s reached max connections", remoteAddr)
	}

	tc.SetKeepAlive(true)

	tc.SetKeepAlivePeriod(3 * time.Minute)

	return tc, nil
}

func (ln *TCPSecListener) Close() error {
	return ln.listener.Close()
}

func (ln *TCPSecListener) Addr() net.Addr {
	return ln.listener.Addr()
}
