
# Core configuration
Core {
	Address = "0.0.0.0"
	Port = "80"

	LogLevel = "3"
	LogOut = "stdout"

	FileDir = "/usr/lib/microhttp/"

	# Virtual host configuration
	VirtualHosting = false
	VirtualHosts {
		"localhost" = "/path/to/vhost1.hcl"
		"127.0.0.1" = "/path/to/vhost2.hcl"
	}

	# Metrics settings
	Metrics {
		Enabled = true
		Path = "/MicroMetrics"
		Out = "/usr/lib/microhttp/metrics.json"
		Users {
			"Admin" = "Your amazing passphrase goes here, because passphrases are the way to go"
		}
	}

	# TLS configuration
	TLS {
		Enabled = false
		TLSCert = "/path/to/tls_certificate"
		TLSKey = "/path/to/tls_key"

		# Autocert will automatically retrieve certificates for you
		AutoCert {
			Enabled = false
			Certificates = [
				"example.com",
				"subdomain.example.com"
			]
		}

		# HSTS settings (HTTP Strict Transport Security)
		HSTS {
			MaxAge = 63072000
			Preload = true
			IncludeSubdomains = true
		}
	}
}

# Server example configuration
Serve {
	ServeDir = "/usr/lib/microhttp/www/"
	ServeIndex = "index.html"

	# Custom HTTP headers
	Headers {
		Content-Security-Policy		= "default-src 'self'",
		Feature-Policy						= "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';",
		Referrer-Policy						= "no-referrer",
		X-Content-Type-Options 		= "nosniff",
		X-Frame-Options						= "SAMEORIGIN",
		X-Xss-Protection					= "1; mode=block"
	}

	Methods {
		"/" = "GET"
	}

	Download {
	  Enabled = false
		Exts = [ ".zip" ]
	}
}

# Proxy settings
Proxy {
	Enabled = false
	Rules {
		"localhost" = "https://proxy-to.example.com"
		"127.0.0.1" = "https://proxy-to.example.eu"
	}
}

# Guard settings
# WARNING: rate and rateburst must be set higher than 0, or your site will be unreachable
Guard {
	Rate = 100
	RateBurst = 10
}
