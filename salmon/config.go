package salmon

import "golang.org/x/time/rate"

type Config struct {
	Address             string
	Port                string
	KeepAliveInSeconds  int
	MaxConnections      rate.Limit
	MaxConnectionsBurst int
	Firewall            FirewallConfig
}

type FirewallConfig struct {
	Enabled      bool
	Blacklisting bool
	Rules        []string
}
