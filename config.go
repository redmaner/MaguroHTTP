package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	VirtualHosts   map[string]vhost
}

type vhost struct {
	ServeDir     string
	ServeIndex   string
	Headers      map[string]string
	Methods      map[string]string
	ContentTypes contentTypes
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
	ProxyRules   map[string][]string
	HTTPRules    map[string][]string
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

func validateConfig(c *microConfig) (bool, error) {

	// Test for empty elements that cannot be empty
	if c.Address == "" || c.Port == "" || c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
		return false, fmt.Errorf("The server configuration has missing elements: check Address, Port, ServeDir and ServeIndex")
	}

	// Test virtual hosts
	if c.Serve.VirtualHosting {
		if len(c.Serve.VirtualHosts) == 0 {
			return false, fmt.Errorf("VirtualHosting is enabled but VirtualHosts is empty")
		} else {
			for k, v := range c.Serve.VirtualHosts {
				if v.ServeDir == "" || v.ServeIndex == "" {
					return false, fmt.Errorf("Virtual host %s has missing elements: check ServeDir and ServeIndex for %s", k, k)
				}
			}
		}
	}

	if c.Proxy.Enabled {
		if len(c.Proxy.Rules) == 0 {
			return false, fmt.Errorf("Proxy is enabled but no rules are defined")
		}
	}

	if c.TLS {
		if c.TLSCert == "" || c.TLSKey == "" {
			return false, fmt.Errorf("TLS is enabled but certificates are not defined")
		}
	}

	if c.Firewall.Enabled {
		if c.Proxy.Enabled && len(c.Firewall.ProxyRules) == 0 {
			return false, fmt.Errorf("Firewall is enabled but rules are not defined")
		} else if !c.Proxy.Enabled && len(c.Firewall.HTTPRules) == 0 {
			return false, fmt.Errorf("Firewall is enabled but rules are not defined")
		}
	}

	return true, nil

}
