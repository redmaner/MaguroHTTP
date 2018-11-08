package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type micro struct {
	config microConfig
	vhosts map[string]microConfig
}

// MicroHTTP config type
type microConfig struct {
	Address      string
	Port         string
	Serve        serve
	Errors       map[string]string
	Headers      map[string]string
	Methods      map[string]string
	ContentTypes contentTypes
	Proxy        proxy
	TLS          bool
	TLSCert      string
	TLSKey       string
	HSTS         hsts
	Firewall     firewall
}

type serve struct {
	ServeDir       string
	ServeIndex     string
	VirtualHosting bool
	VirtualHosts   map[string]string
}

// Proxy type, part of MicroHTTP config
type proxy struct {
	Enabled bool
	Rules   map[string]string
}

// contentTypes type, part of MicroHTTP config
type contentTypes struct {
	ResponseTypes map[string]string
	RequestTypes  []string
}

// HSTS type, part of MicroHTTP config
type hsts struct {
	MaxAge            int
	Preload           bool
	IncludeSubdomains bool
}

// Firewall type, part of MicroHTTP config
type firewall struct {
	Enabled      bool
	Blacklisting bool
	Subpath      bool
	Rules        map[string][]string
}

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

func validateConfig(p string, c *microConfig) (bool, error) {

	// Test for empty elements that cannot be empty
	if c.Address == "" || c.Port == "" || c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
		return false, fmt.Errorf("%s: The server configuration has missing elements: check Address, Port, ServeDir and ServeIndex", p)
	}

	// Test virtual hosts
	if c.Serve.VirtualHosting {
		if len(c.Serve.VirtualHosts) == 0 {
			return false, fmt.Errorf("%s: VirtualHosting is enabled but VirtualHosts is empty", p)
		} else {
			for k, v := range c.Serve.VirtualHosts {
				if v == "" {
					return false, fmt.Errorf("%s: Virtual host configuration not defined. Check reference for %s", p, k)
				}
			}
		}
	}

	if c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			return false, fmt.Errorf("%s: Proxy is enabled but no rules are defined", p)
		}
	}

	if c.TLS {
		if c.TLSCert == "" || c.TLSKey == "" {
			return false, fmt.Errorf("%s: TLS is enabled but certificates are not defined", p)
		}
	}

	if c.Firewall.Enabled {
		if len(c.Firewall.Rules) == 0 {
			return false, fmt.Errorf("%s: Firewall is enabled but rules are not defined", p)
		}
	}

	return true, nil

}

func validateConfigVhost(p string, c *microConfig) (bool, error) {

	// Test virtual hosts
	if c.Serve.VirtualHosting {
		return false, fmt.Errorf("%s: VirtualHosting cannot be enabled in Vhost configuration", p)
	}

	if c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			return false, fmt.Errorf("%s: Proxy is enabled but no rules are defined", p)
		}
	} else if c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
		return false, fmt.Errorf("%s: The server configuration has missing elements: check ServeDir and ServeIndex", p)
	}

	if c.Firewall.Enabled {
		if len(c.Firewall.Rules) == 0 {
			return false, fmt.Errorf("%s: Firewall is enabled but rules are not defined", p)
		}
	}

	return true, nil
}
