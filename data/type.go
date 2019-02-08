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

package data

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Config is type holding the main configurtion
type Config struct {
	Core   coreConfig
	Serve  serveConfig
	Errors map[string]string
	Proxy  proxy
}

// coreConfig is part of the main configuration.
// coreConfig is not used by vhosts
type coreConfig struct {
	Address        string
	Port           string
	LogLevel       int
	LogOut         string
	VirtualHosting bool
	VirtualHosts   map[string]string
	TLS            TLSConfig
}

// TLSConfig holds information about TLS and is part of MicroHTTP core config
type TLSConfig struct {
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
type FileInfo struct {
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
type proxy struct {
	Enabled bool
	Rules   map[string]string
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

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
	}
}
