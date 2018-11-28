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
