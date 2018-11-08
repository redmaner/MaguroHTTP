package main

import (
	"crypto/tls"
	"os"
)

// Functio to check if defined TLS certificates exist
func httpCheckTLS(c *microConfig) bool {
	if c.TLSCert != "" && c.TLSKey != "" {
		if _, err := os.Stat(c.TLSCert); err != nil {
			logAction(logERROR, err)
			return false
		}
		if _, err := os.Stat(c.TLSKey); err != nil {
			logAction(logERROR, err)
			return false
		}
		return true
	}
	return false
}

// Function to create a TLS configuration based on server configuration
func httpCreateTLSConfig() *tls.Config {
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
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		},
	}
	return &tlsc
}
