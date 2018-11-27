package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

type tlsconfig struct {
	Enabled   bool
	TLSCert   string
	TLSKey    string
	PrivateCA []string
}

// Functio to check if defined TLS certificates exist
func httpCheckTLS(c tlsconfig) bool {
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
func httpCreateTLSConfig(c tlsconfig) *tls.Config {
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

	// Check if the config uses a PrivateCA and load it
	if len(c.PrivateCA) > 0 {
		var capool *x509.CertPool
		if runtime.GOOS == "windows" {
			capool = x509.NewCertPool()
		} else {
			capool, _ = x509.SystemCertPool()
		}
		for _, v := range c.PrivateCA {
			if cacert, err := ioutil.ReadFile(v); err == nil {
				if ok := capool.AppendCertsFromPEM(cacert); ok {
					tlsc.RootCAs = capool
				}
			} else {
				logAction(logERROR, err)
			}
		}
		tlsc.RootCAs = capool
	}
	return &tlsc
}

// Function that returns an http client. If TLS is enabled on the server, a
// TLS client is returned
func serverClient(c tlsconfig) *http.Client {
	if c.Enabled {

		tlsc := httpCreateTLSConfig(c)

		t := &http.Transport{
			TLSClientConfig: tlsc,
		}

		return &http.Client{
			Timeout:   time.Minute,
			Transport: t,
		}
	}
	return http.DefaultClient
}
