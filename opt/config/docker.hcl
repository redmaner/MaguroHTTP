# Core configuration
core {
	address = "0.0.0.0"
	port = "80"

	log_level = "3"
	log_out = "stdout"

	file_directory = "/usr/lib/magurohttp/"

	# Virtual host configuration
	virtual_hosting = false
	virtual_hosts {
		"localhost" = "/path/to/vhost1.hcl"
		"127.0.0.1" = "/path/to/vhost2.hcl"
	}

	# Metrics settings
	metrics {
		enabled = true
		path = "/MicroMetrics"
		out = "/usr/lib/magurohttp/metrics.json"
		users {
			"Admin" = "Your amazing passphrase goes here, because passphrases are the way to go"
		}
	}

	# TLS configuration
  tls {
		enabled = false
		tls_cert = "/path/to/tls_certificate"
		tls_key = "/path/to/tls_key"

		# Autocert will automatically retrieve certificates for you
		auto_cert {
			enabled = false
			certificates = [
				"example.com",
				"subdomain.example.com"
			]
		}

		# HSTS settings (HTTP Strict Transport Security)
		hsts {
			max_age = 63072000
			preload = true
			include_subdomains = true
		}
	}
}

# Server example configuration
serve {
	serve_directory = "/usr/lib/magurohttp/www/"
	serve_index = "index.html"

	# Custom HTTP headers
	headers {
		Content-Security-Policy		= "default-src 'self'",
		Feature-Policy						= "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';",
		Referrer-Policy						= "no-referrer",
		X-Content-Type-Options 		= "nosniff",
		X-Frame-Options						= "SAMEORIGIN",
		X-Xss-Protection					= "1; mode=block"
	}

	methods {
		"/" = "GET"
	}

	download {
	  enabled = false
		extensions = [ ".zip" ]
	}
}

# Proxy settings
proxy {
	enabled = false
	rules {
		"localhost" = "https://proxy-to.example.com"
		"127.0.0.1" = "https://proxy-to.example.eu"
	}
}

# Guard settings
# WARNING: rate and rateburst must be set higher than 0, or your site will be unreachable
guard {
	rate = 100
	rate_burst = 10
}
