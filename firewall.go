package main

import (
	"path"
)

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
		if val, ok := rules["/"]; ok && p == "/" {
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
