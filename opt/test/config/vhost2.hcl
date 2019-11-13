#Server settings
serve {
	serve_directory = "./opt/test/web"
	serve_index = "index.html"

	headers {
		"Content-Security-Policy" 	=	"default-src 'self'",
		"Feature-Policy"						= "geolocation 'none'; midi 'none'; notifications 'none'; push 'none'; sync-xhr 'none'; microphone 'none'; camera 'none'; magnetometer 'none'; gyroscope 'none'; speaker 'none'; vibrate 'none'; fullscreen 'none'; payment 'none';",
		"Referrer-Policy" 					= "no-referrer",
		"X-Content-Type-Options" 		= "nosniff",
		"X-Frame-Options" 					= "SAMEORIGIN",
		"X-Xss-Protection"					= "1; mode=block"
	}

	methods {
		"/" = "GET"
	}
}

# Custom error pages
errors {
	"429" = "./opt/test/web/27.html"
}

# Guard settings
guard {
	rate 			= 100
	rate_burst = 10

	# Firewall settings
	firewall {
		enabled = false
		blacklisting = true
		subpath = false
		rules {
			"/" = [ "127.0.0.1" ]
		}
	}
}
