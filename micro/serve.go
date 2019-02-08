package micro

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redmaner/MicroHTTP/debug"
	"gitlab.com/EDSN/griffin/logger"
	"golang.org/x/crypto/acme/autocert"
)

// Serve is used to serve a MicroHTTP instance
func (s *Server) Serve() {

	// TLS
	var tlsCert string
	var tlsKey string

	// Define server struct
	server := http.Server{
		Addr:              s.Cfg.Core.Address + ":" + s.Cfg.Core.Port,
		Handler:           s.Router,
		ReadTimeout:       4 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       12 * time.Second,
		ErrorLog:          s.logInterface.Instance,
	}

	switch {

	// If TLS is enabled the server will start in TLS
	case s.Cfg.Core.TLS.Enabled && s.httpCheckTLS():
		s.Log(logger.LogNone, fmt.Errorf("MicroHTTP is listening on port %s with TLS", s.Cfg.Core.Port))
		tlsc := s.httpCreateTLSConfig()

		// Handle autocert
		if s.Cfg.Core.TLS.AutoCert.Enabled && len(s.Cfg.Core.TLS.AutoCert.Certificates) > 0 {
			acm := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.Cfg.Core.TLS.AutoCert.Certificates...),
				Cache:      autocert.DirCache(s.Cfg.Core.TLS.AutoCert.CertDir),
			}
			tlsc.GetCertificate = acm.GetCertificate
		} else {
			tlsCert = s.Cfg.Core.TLS.TLSCert
			tlsKey = s.Cfg.Core.TLS.TLSKey
		}
		server.TLSConfig = tlsc

		err := server.ListenAndServeTLS(tlsCert, tlsKey)
		if err != nil {
			panic(err)
		}

	// if TLS is not enabled HTTP will be served
	default:
		s.Log(logger.LogNone, fmt.Errorf("MicroHTTP is listening on port %s", s.Cfg.Core.Port))
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}

	// Gracefully stop the server in case of a signal
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case signal := <-sig:
			fmt.Printf("Signal (%d) received, stopping\n", signal)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.SetKeepAlivesEnabled(false)
			if err := server.Shutdown(ctx); err != nil {
				s.Log(debug.LogNone, fmt.Errorf("could not gracefully shutdown the server: %v", err))
			}
			close(sig)
			os.Exit(1)
		}
	}
}
