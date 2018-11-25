package main

import (
	"path"
)

// Firewall type, part of MicroHTTP config
type firewall struct {
	Enabled      bool
	Blacklisting bool
	Subpath      bool
	Rules        map[string][]string
}

// Function to determine whether to block a host when serving HTTP
// The firewall can be set in the configuration and is disabled by default.
// By default the firewall is a whitelist firewall, which only allows hosts
// that are explicitly set in the configuration file
func firewallHTTP(cfg *microConfig, h, p string) bool {

	rules := cfg.Firewall.Rules

	if cfg.Firewall.Enabled {
		for pt := p; pt != "/"; pt = path.Dir(pt) {
			if val, ok := rules[pt]; ok {
				for _, v := range val {
					if v == h || v == "*" {
						return cfg.Firewall.Blacklisting
					}
				}
			}
		}

		// The firewall subpath element allows blocking on specific subpaths of a website
		// This is only when you want to be extremely specific when configuring the firewall.
		// Subpath blocking is disabled by default and can be enabled in the configuration.
		if val, ok := rules["/"]; ok && p == "/" || ok && !cfg.Firewall.Subpath {
			for _, v := range val {
				if v == h || v == "*" {
					return cfg.Firewall.Blacklisting
				}
			}
		}
		return !cfg.Firewall.Blacklisting
	}
	return false
}

// Function to determine whether to block a host when serving HTTP
// The firewall can be set in the configuration and is disabled by default.
// By default the firewall is a whitelist firewall, which only allows hosts
// that are explicitly set in the configuration file
func firewallProxy(cfg *microConfig, h, p string) bool {

	rules := cfg.Firewall.Rules

	if cfg.Firewall.Enabled {
		if val, ok := rules[p]; ok {
			for _, v := range val {
				if v == h || v == "*" {
					return cfg.Firewall.Blacklisting
				}
			}
		}
	}
	return false
}
