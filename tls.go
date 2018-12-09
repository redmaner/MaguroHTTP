// Copyright 2018 Jake van der Putten.
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

package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"golang.org/x/crypto/acme"
)

// TLSConfig, part of MicroHTTP core config
type tlsConfig struct {
	Enabled   bool
	TLSCert   string
	TLSKey    string
	PrivateCA []string
	AutoCert  autocertConfig
	HSTS      hstsConfig
}

// autocertConfig, part of MicroHTTP core/tls configuration
type autocertConfig struct {
	Enabled      bool
	CertDir      string
	Certificates []string
}

// HSTS type, part of MicroHTTP core/tls config
type hstsConfig struct {
	MaxAge            int
	Preload           bool
	IncludeSubdomains bool
}

// Functio to check if defined TLS certificates exist
func httpCheckTLS(c tlsConfig) bool {
	if c.AutoCert.Enabled {
		return true
	}
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
func httpCreateTLSConfig(c tlsConfig) *tls.Config {
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
		NextProtos: []string{
			"h2", "http/1.1", // enable HTTP/2
			acme.ALPNProto, // enable tls-alpn ACME challenges
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
func serverClient(c tlsConfig) *http.Client {
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
