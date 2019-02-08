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

import "log"

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

				// Certificates will be saved locally, so it requires an explicit directory
				if c.Core.TLS.AutoCert.CertDir == "" {
					log.Fatalf("%s: TLS autocert is enabled but CertDir is empty or not defined", p)
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
