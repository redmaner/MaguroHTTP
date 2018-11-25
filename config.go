package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// MicroHTTP config type
// The main configuration and the configuration of vhosts use this type
type microConfig struct {
	Address      string
	Port         string
	LogLevel     int
	LogOut       string
	Serve        serve
	Download     download
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
	Metrics      metrics
}

// Proxy type, part of MicroHTTP config
type proxy struct {
	Enabled bool
	Rules   map[string]string
}

// contentTypes type, part of MicroHTTP config
type contentTypes struct {
	ResponseTypes map[string]string
	RequestTypes  map[string]string
}

// HSTS type, part of MicroHTTP config
type hsts struct {
	MaxAge            int
	Preload           bool
	IncludeSubdomains bool
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

	// Test for empty elements that cannot be empty
	if c.Address == "" || c.Port == "" || c.Serve.ServeDir == "" || c.Serve.ServeIndex == "" {
		return false, fmt.Errorf("%s: The server configuration has missing elements: check Address, Port, ServeDir and ServeIndex", p)
	}

	if c.LogOut =="" {
		return false, fmt.Errorf("%s: LogOut is undefined", p)
	}

	if c.LogLevel < 0 {
		return false, fmt.Errorf("%s: LogLevel must be higher than 0", p)
	}

	// Test virtual hosts
	if c.Serve.VirtualHosting {
		if len(c.Serve.VirtualHosts) == 0 {
			return false, fmt.Errorf("%s: VirtualHosting is enabled but VirtualHosts is empty", p)
		}
		for k, v := range c.Serve.VirtualHosts {
			if v == "" {
				return false, fmt.Errorf("%s: Virtual host configuration not defined. Check reference for %s", p, k)
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

	if c.Metrics.Enabled {
		if c.Metrics.Path == "" || c.Metrics.Path == "/" {
			return false, fmt.Errorf("%s: Metrics path cannot be empty or /", p)
		}
		if c.Metrics.User == "" || c.Metrics.Password == "" || c.Metrics.Address == "" {
			return false, fmt.Errorf("%s: Metrics user, password and address cannot be empty", p)
		}
	}

	return true, nil

}

// Function to validate a vhost config loaded with loadConfigFromFile(p, c)
// This function validates for the exsistence of several elements that are necessary
// to use a defined vhost
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
