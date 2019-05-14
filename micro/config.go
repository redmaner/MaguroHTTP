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

package micro

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/hashicorp/hcl"
)

// Config is type holding the main configurtion
type Config struct {
	Core    CoreConfig
	Serve   serveConfig
	Errors  map[string]string
	Proxy   proxyConfig
	Guard   guardConfig
	Metrics metricsConfig
}

// CoreConfig is part of the main configuration.
// coreConfig is not used by vhosts
type CoreConfig struct {
	Address  string
	Port     string
	FileDir  string
	LogLevel int
	LogOut   string

	ReadTimeout       int
	ReadHeaderTimeout int
	WriteTimeout      int

	WebDAV         bool
	VirtualHosting bool
	VirtualHosts   map[string]string
	TLS            tlsConfig
}

// TLSConfig holds information about TLS and is part of MicroHTTP core config
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
	Certificates []string
}

// HSTS type, part of MicroHTTP core/tls config
type hstsConfig struct {
	MaxAge            int
	Preload           bool
	IncludeSubdomains bool
}

// Serve type, part of the MicroHTTP config
type serveConfig struct {
	ServeDir   string
	ServeIndex string
	Headers    map[string]string
	Methods    map[string]string
	MIMETypes  MIMETypes
	Download   download
}

// Download type, part of the MicroHTTP config
type download struct {
	Enabled bool
	Exts    []string
}

// FileInfo to gather information about files
type fileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}

// MIMETypes type, part of MicroHTTP serveConfig
type MIMETypes struct {
	ResponseTypes map[string]string
	RequestTypes  map[string]string
}

// Proxy type, part of MicroHTTP config
type proxyConfig struct {
	Enabled bool
	Rules   map[string]string
	Methods []string
}

// guardConfig
type guardConfig struct {
	Rate       float64
	RateBurst  int
	FilterOnIP bool

	Firewall firewallConfig
}

// Firewall type, part of MicroHTTP config
type firewallConfig struct {
	Enabled      bool
	Blacklisting bool
	Subpath      bool
	Rules        map[string][]string
}

// Metrics type, part of MicroHTTP config
type metricsConfig struct {
	Enabled bool
	Path    string
	Out     string
	Users   map[string]string
}

// LoadConfigFromFile is a function which a loads the Config type microConfig from a json file
func LoadConfigFromFile(p string, c *Config) {

	// check if config exists
	if _, err := os.Stat(p); err != nil {
		log.Fatal(err)
	}

	// load config
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("%s: %v", p, err)
	}

	err = hcl.Unmarshal(data, c)
	if err != nil {
		log.Fatalf("%s: %v", p, err)
	}
}

// Validate can be used to validate a Config type
func (c *Config) Validate(p string, isVhost bool) {

	if !isVhost {
		if c.Core.Address == "" || c.Core.Port == "" {
			log.Fatalf("%s: The server configuration has missing elements: check Address and Port", p)
		}

		// LogOut needs to be defined
		if c.Core.LogOut == "" {
			log.Fatalf("%s: LogOut is undefined", p)
		}

		// LogLevel cannot be lower than zero
		if c.Core.LogLevel < 0 {
			log.Fatalf("%s: LogLevel must be higher than 0", p)
		}

		// FileDir must be defined
		if c.Core.FileDir == "" || c.Core.FileDir == "/" {
			log.Fatalf("%s: FileDir is not defined or is pointing to root", p)
		}

		// We automatically fix FileDir if it doesn't end with a slash
		if c.Core.FileDir[len(c.Core.FileDir)-1] != '/' {
			c.Core.FileDir = c.Core.FileDir + "/"
		}

		// Test timeouts and set defaults if they are zero
		if c.Core.ReadTimeout <= 0 {
			c.Core.ReadTimeout = 8
		}

		if c.Core.ReadHeaderTimeout <= 0 {
			c.Core.ReadHeaderTimeout = 4
		}

		if c.Core.WriteTimeout <= 0 {
			c.Core.WriteTimeout = 16
		}

		// Test TLS
		if c.Core.TLS.Enabled {

			// Test autocert
			if c.Core.TLS.AutoCert.Enabled {

				// Certificates need to be defined
				if len(c.Core.TLS.AutoCert.Certificates) == 0 {
					log.Fatalf("%s: TLS autocert is enabled but certificates are not defined", p)
				}

				// Autocert only works in combination with https port (443)
				if c.Core.Port != "443" {
					log.Fatalf("%s: TLS autocert is enabled and cannot be used with a port different than 443 (HTTPS)", p)
				}
			} else {

				// Autocert is disabled, so make sure custom certificate / key combination is defined
				if c.Core.TLS.TLSCert == "" || c.Core.TLS.TLSKey == "" {
					log.Fatalf("%s: TLS is enabled but certificates are not defined", p)
				}
			}
		}
	}

	// Test virtual hosts
	if !isVhost && c.Core.VirtualHosting {
		if len(c.Core.VirtualHosts) == 0 {
			log.Fatalf("%s: VirtualHosting is enabled but VirtualHosts is empty", p)
		}
		for k, v := range c.Core.VirtualHosts {
			if v == "" {
				log.Fatalf("%s: Virtual host configuration not defined. Check reference for %s", p, k)
			}
		}
		return

	} else if c.Core.VirtualHosting {
		log.Fatalf("%s: Virtual hosting cannot be enabled in a vhost configuration", p)
	}

	// Test serve
	if !c.Proxy.Enabled && c.Serve.Download.Enabled {
		if c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
			log.Fatalf("%s: The server configuration has missing elements: check ServeDir and ServeIndex", p)
		}

		// We automatically fix ServeDir that doesn't end with a slash
		if c.Serve.ServeDir[len(c.Serve.ServeDir)-1] != '/' {
			c.Serve.ServeDir = c.Serve.ServeDir + "/"
		}
	}

	// Test proxy
	if c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			log.Fatalf("%s: Proxy is enabled but no rules are defined", p)
		}
	}
}
