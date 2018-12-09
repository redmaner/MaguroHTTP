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
	"encoding/json"
	"fmt"
	"os"
)

// MicroHTTP config type
// The main configuration and the configuration of vhosts use this type
type microConfig struct {
	Core     coreConfig
	Serve    serveConfig
	Errors   map[string]string
	Proxy    proxy
	Firewall firewall
	Metrics  metrics
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
	TLS            tlsConfig
}

// This loads a configuration with type microConfig from a file
// The expected input is a json file with elements that match the microConfig spec
func loadConfigFromFile(p string, c *microConfig) {

	// check if config exists
	if _, err := os.Stat(p); err != nil {
		logAction(logERROR, err)
		os.Exit(1)
	}

	// load config
	file, err := os.Open(p)
	if err != nil {
		logAction(logERROR, err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		logAction(logERROR, err)
		os.Exit(1)
	}
}

// Function to validate the main config loaded with loadConfigFromFile(p, c)
// This function validates for the exsistence of several elements that are necessary
// to start the server
func validateConfig(p string, c *microConfig) (bool, error) {

	// First validate the coreConfig
	// Address and port need to be defined
	if c.Core.Address == "" || c.Core.Port == "" {
		return false, fmt.Errorf("%s: The server configuration has missing elements: check Address and Port", p)
	}

	// LogOut needs to be defined
	if c.Core.LogOut == "" {
		return false, fmt.Errorf("%s: LogOut is undefined", p)
	}

	// LogLevel cannot be lower than zero
	if c.Core.LogLevel < 0 {
		return false, fmt.Errorf("%s: LogLevel must be higher than 0", p)
	}

	// Test TLS
	if c.Core.TLS.Enabled {

		// Test autocert
		if c.Core.TLS.AutoCert.Enabled {

			// Certificates need to be defined
			if len(c.Core.TLS.AutoCert.Certificates) == 0 {
				return false, fmt.Errorf("%s: TLS autocert is enabled but certificates are not defined", p)
			}

			// Autocert only works in combination with https port (443)
			if c.Core.Port != "443" {
				return false, fmt.Errorf("%s: TLS autocert is enabled and cannot be used with a port different than 443 (HTTPS)", p)
			}

			// Certificates will be saved locally, so it requires an explicit directory
			if c.Core.TLS.AutoCert.CertDir == "" {
				return false, fmt.Errorf("%s: TLS autocert is enabled but CertDir is empty or not defined", p)
			}
		} else {

			// Autocert is disabled, so make sure custom certificate / key combination is defined
			if c.Core.TLS.TLSCert == "" || c.Core.TLS.TLSKey == "" {
				return false, fmt.Errorf("%s: TLS is enabled but certificates are not defined", p)
			}
		}
	}

	// Test virtual hosts
	if c.Core.VirtualHosting {
		if len(c.Core.VirtualHosts) == 0 {
			return false, fmt.Errorf("%s: VirtualHosting is enabled but VirtualHosts is empty", p)
		}
		for k, v := range c.Core.VirtualHosts {
			if v == "" {
				return false, fmt.Errorf("%s: Virtual host configuration not defined. Check reference for %s", p, k)
			}
		}
	}

	// Test serve
	if !c.Core.VirtualHosting && !c.Proxy.Enabled && c.Serve.Download.Enabled {
		if c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
			return false, fmt.Errorf("%s: The server configuration has missing elements: check ServeDir and ServeIndex", p)
		}

		// We automatically fix ServeDir that doesn't end with a slash
		if c.Serve.ServeDir[len(c.Serve.ServeDir)-1] != '/' {
			c.Serve.ServeDir = c.Serve.ServeDir + "/"
		}
	}

	// Test proxy
	if !c.Core.VirtualHosting && c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			return false, fmt.Errorf("%s: Proxy is enabled but no rules are defined", p)
		}
	}

	// Test firewall
	if !c.Core.VirtualHosting && c.Firewall.Enabled {
		if len(c.Firewall.Rules) == 0 {
			return false, fmt.Errorf("%s: Firewall is enabled but rules are not defined", p)
		}
	}

	if c.Metrics.Enabled {
		if c.Metrics.Path == "" || c.Metrics.Path == "/" {
			return false, fmt.Errorf("%s: Metrics path cannot be empty or /", p)
		}
		if c.Metrics.User == "" || c.Metrics.Password == "" || c.Metrics.Address == "" {
			return false, fmt.Errorf("%s: Metrics user, password and address cannot be empty", p)
		}
		if c.Metrics.Out == "" {
			return false, fmt.Errorf("%s: Metrics out cannot be empty", p)
		}
	}

	return true, nil

}

// Function to validate a vhost config loaded with loadConfigFromFile(p, c)
// This function validates for the exsistence of several elements that are necessary
// to use a defined vhost
func validateConfigVhost(p string, c *microConfig) (bool, error) {

	// Test virtual hosts
	if c.Core.VirtualHosting {
		return false, fmt.Errorf("%s: VirtualHosting cannot be enabled in Vhost configuration", p)
	}

	if c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			return false, fmt.Errorf("%s: Proxy is enabled but no rules are defined", p)
		}
	} else if !c.Serve.Download.Enabled && !c.Proxy.Enabled {
		if c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
			return false, fmt.Errorf("%s: The server configuration has missing elements: check ServeDir and ServeIndex", p)
		}

		// We automatically fix ServeDir that doesn't end with a slash
		if c.Serve.ServeDir[len(c.Serve.ServeDir)-1] != '/' {
			c.Serve.ServeDir = c.Serve.ServeDir + "/"
		}
	}

	if c.Firewall.Enabled {
		if len(c.Firewall.Rules) == 0 {
			return false, fmt.Errorf("%s: Firewall is enabled but rules are not defined", p)
		}
	}

	return true, nil
}
