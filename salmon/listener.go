package salmon

import (
	"log"
	"net"
	"time"

	"github.com/redmaner/MaguroHTTP/cache"
	"github.com/redmaner/MaguroHTTP/router"
	"golang.org/x/time/rate"
)

// Listener provides a TCP listener with rate limting integration
type Listener struct {
	listener      *net.TCPListener
	cache         *cache.SpearCache
	maxConns      rate.Limit
	maxConnsBurst int
}

func NewListener(addr string, maxConnsPerMin float32, maxConnsBurst int) *Listener {
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}
	return &Listener{
		listener:      ln.(*net.TCPListener),
		cache:         cache.NewCache(),
		maxConns:      rate.Limit(maxConnsPerMin / 60),
		maxConnsBurst: maxConnsBurst,
	}
}

func (ln *Listener) Accept() (net.Conn, error) {

	for {
		tc, err := ln.listener.AcceptTCP()
		if err != nil {
			continue
		}

		remoteAddr := router.StripHostPort(tc.RemoteAddr().String())

		// If connection isn't allowed, close the connection
		if !ln.allowConnection(remoteAddr) {
			tc.Close()
			continue
		}

		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(1 * time.Minute)
		return tc, nil
	}
}

func (ln *Listener) Close() error {
	return ln.listener.Close()
}

func (ln *Listener) Addr() net.Addr {
	return ln.listener.Addr()
}
