package main

import (
	"path"
)

func (m *micro) firewallHTTP(h, p string) bool {

	rules := m.config.Firewall.HTTPRules

	if m.config.Firewall.Enabled {
		for pt := p; pt != "/"; pt = path.Dir(pt) {
			if val, ok := rules[pt]; ok {
				for _, v := range val {
					if v == h || v == "*" {
						return m.config.Firewall.Blacklisting
					}
				}
			}
		}
		if val, ok := rules["/"]; ok {
			for _, v := range val {
				if v == h || v == "*" {
					return m.config.Firewall.Blacklisting
				}
			}
		}
		return !m.config.Firewall.Blacklisting
	}
	return false
}

func (m *micro) firewallProxy(h, p string) bool {

	rules := m.config.Firewall.ProxyRules

	if m.config.Firewall.Enabled {
		if val, ok := rules[p]; ok {
			for _, v := range val {
				if v == h || v == "*" {
					return m.config.Firewall.Blacklisting
				}
			}
		}
	}

	return false

}
