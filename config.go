package main

import (
	"encoding/json"
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
	ServeDir   string
	ServeIndex string
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
