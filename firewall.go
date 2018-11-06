package main

import (
	"path"
)

func firewallHttp(h, p string) bool {

	rules := mCfg.Firewall.HttpRules

	if mCfg.Firewall.Enabled {
		for pt := p; pt != "/"; pt = path.Dir(pt) {
			if val, ok := rules[pt]; ok {
				for _, v := range val {
					if v == h || v == "*" {
						return mCfg.Firewall.Blacklisting
					}
				}
			}
		}
		if val, ok := rules["/"]; ok {
			for _, v := range val {
				if v == h || v == "*" {
					return mCfg.Firewall.Blacklisting
				}
			}
		}
	}

	return !mCfg.Firewall.Blacklisting

}

func firewallProxy(h, p string) bool {

	rules := mCfg.Firewall.ProxyRules

	if mCfg.Firewall.Enabled {
		if val, ok := rules[p]; ok {
			for _, v := range val {
				if v == h || v == "*" {
					return mCfg.Firewall.Blacklisting
				}
			}
		}
	}

	return false

}
