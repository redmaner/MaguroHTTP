package micro

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/redmaner/MicroHTTP/debug"
	"golang.org/x/crypto/acme"
)

// Functio to check if defined TLS certificates exist
func (s *Server) httpCheckTLS() bool {
	if s.Cfg.Core.TLS.AutoCert.Enabled {
		return true
	}
	if s.Cfg.Core.TLS.TLSCert != "" && s.Cfg.Core.TLS.TLSKey != "" {
		if _, err := os.Stat(s.Cfg.Core.TLS.TLSCert); err != nil {
			s.Log(debug.LogError, err)
			return false
		}
		if _, err := os.Stat(s.Cfg.Core.TLS.TLSKey); err != nil {
			s.Log(debug.LogError, err)
			return false
		}
		return true
	}
	return false
}

// Function to create a TLS configuration based on server configuration
func (s *Server) httpCreateTLSConfig() *tls.Config {
	tlsc := tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	if s.Cfg.Core.TLS.AutoCert.Enabled {
		tlsc.NextProtos = []string{
			"h2", "http/1.1", // enable HTTP/2
			acme.ALPNProto, // enable tls-alpn ACME challenges
		}
	}

	// Check if the config uses a PrivateCA and load it
	if len(s.Cfg.Core.TLS.PrivateCA) > 0 {
		var capool *x509.CertPool
		if runtime.GOOS == "windows" {
			capool = x509.NewCertPool()
		} else {
			capool, _ = x509.SystemCertPool()
		}
		for _, v := range s.Cfg.Core.TLS.PrivateCA {
			if cacert, err := ioutil.ReadFile(v); err == nil {
				if ok := capool.AppendCertsFromPEM(cacert); ok {
					tlsc.RootCAs = capool
				}
			} else {
				s.Log(debug.LogError, err)
			}
		}
		tlsc.RootCAs = capool
	}
	return &tlsc
}
